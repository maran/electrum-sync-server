package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"path"
)

type SyncMaster struct {
	db     gorm.DB
	logger *lumber.ConsoleLogger
}

func newSyncMaster(opts DbOpts) SyncMaster {
	sync := SyncMaster{}
	sync.logger = lumber.NewConsoleLogger(lumber.DEBUG)
	sync.logger.Info("Starting Electrum Sync Server v%s", Version)

	var err error

	if opts.DbType == "sqlite3" {
		dbPath := opts.DbPath + "/sync.db"
		sync.logger.Info("Opening database at %s", dbPath)
		err = os.MkdirAll(opts.DbPath, 0700) // read, write and dir search for user
		if err != nil {
			log.Fatal("Error creating database folder", err)
		}
		newdb, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("Error opening sqlite3 database:", err)
		}
		sync.db = *newdb
	} else if opts.DbType == "postgres" {
		optss := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", opts.Host, opts.User, opts.Dbname, opts.Password)
		newdb, err := gorm.Open("postgres", optss)
		if err != nil {
			log.Fatal("Error connecting to Postgres database: ", err)
		}
		sync.db = *newdb
	} else {
		log.Fatal("Unknown database type. Please supply sqlite3 or postgres")
	}

	sync.db.AutoMigrate(&Label{})
	return sync
}

func (self *SyncMaster) makeLabel(labelRequest LabelRequest, w rest.ResponseWriter) Label {
	var label Label
	search := Label{WalletId: labelRequest.WalletId, ExternalId: labelRequest.ExternalId}

	self.db.Where(search).FirstOrInit(&label)

	self.logger.Debug("current label nonce: %d got a request nonce to overwrite it with %d", label.Nonce, labelRequest.WalletNonce)

	if label.Nonce > labelRequest.WalletNonce {
		rest.Error(w, "serverNonce is larger then walletNonce please sync first.", 400)
	}

	label.EncryptedLabel = labelRequest.EncryptedLabel
	label.Nonce = labelRequest.WalletNonce
	self.db.Save(&label)

	return label
}

func (self *SyncMaster) CreateLabels(w rest.ResponseWriter, r *rest.Request) {
	labelsRequest := LabelsRequest{}
	err := r.DecodeJsonPayload(&labelsRequest)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var labelsResponse LabelsResponse
	for _, batchLabel := range labelsRequest.Labels {
		labelRequest := LabelRequest{batchLabel.EncryptedLabel, batchLabel.ExternalId, labelsRequest.WalletId, labelsRequest.WalletNonce}
		labelsResponse.Labels = append(labelsResponse.Labels, self.makeLabel(labelRequest, w))
	}
	labelsResponse.Nonce = highestNonce(labelsResponse.Labels)
	w.WriteJson(labelsResponse)
}

func (self *SyncMaster) CreateLabel(w rest.ResponseWriter, r *rest.Request) {
	labelRequest := LabelRequest{}
	err := r.DecodeJsonPayload(&labelRequest)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if labelRequest.WalletId == "" {
		rest.Error(w, "walletId required", 400)
	}
	if labelRequest.EncryptedLabel == "" {
		rest.Error(w, "encryptedLabel required", 400)
	}
	if labelRequest.ExternalId == "" {
		rest.Error(w, "externalId required", 400)
	}
	if labelRequest.WalletNonce == 0 {
		rest.Error(w, "walletNonce required", 400)
	}
	self.logger.Debug("Got request:", labelRequest)
	label := self.makeLabel(labelRequest, w)
	w.WriteJson(label)
}

func (self *SyncMaster) GetLabels(w rest.ResponseWriter, r *rest.Request) {
	var labels []Label
	mpk := r.PathParam("mpk")
	nonce := r.PathParam("nonce")

	if mpk == "" {
		rest.Error(w, "walletId required", 400)
	}
	if nonce == "" {
		rest.Error(w, "nonce required", 400)
	}
	self.db.Where("wallet_id = ? AND nonce > ?", mpk, nonce).Find(&labels)
	highestNonce := highestNonce(labels)
	w.WriteJson(LabelsResponse{Nonce: highestNonce, Labels: labels})
}

func highestNonce(labels []Label) int {
	var highestNonce int
	for _, label := range labels {
		if label.Nonce > highestNonce {
			highestNonce = label.Nonce
		}
	}
	return highestNonce
}
func defaultDbDir() string {
	return path.Join(os.Getenv("HOME"), "/.electrum-sync-server/")
}
