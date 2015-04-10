package main

import (
	"fmt"
)

// SQL Object
type Label struct {
	Id             int    `sql:"index" json:"id"`
	ExternalId     string `json:"externalId"`
	EncryptedLabel string `json:"encryptedLabel"`
	Nonce          int    `json:"nonce"`
	WalletId       string `json:"walletId"`
}

// Rest response
type LabelsResponse struct {
	Nonce  int     `json:"nonce"`
	Labels []Label `json:"labels"`
}

// Rest request
type LabelRequest struct {
	EncryptedLabel string `json:"encryptedLabel"`
	ExternalId     string `json:"externalId"`
	WalletId       string `json:"walletId"`
	WalletNonce    int    `json:"walletNonce"`
}

func (self LabelRequest) String() string {
	return fmt.Sprintf(`
Request information:
encryptedLabel: %s
externalId: %s
walletId: %s
walletNonce: %d
	`, self.EncryptedLabel, self.ExternalId, self.WalletId, self.WalletNonce)
}

type LabelsRequest struct {
	WalletNonce int          `json:"walletNonce"`
	WalletId    string       `json:"walletId"`
	Labels      []BatchLabel `json:"labels"`
}

type BatchLabel struct {
	EncryptedLabel string `json:"encryptedLabel"`
	ExternalId     string `json:"externalId"`
}
