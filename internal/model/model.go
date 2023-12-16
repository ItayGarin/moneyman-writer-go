package model

import "time"

type Transaction struct {
	Identifier       string    `json:"identifier"`
	Type             string    `json:"type"`
	Status           string    `json:"status"`
	Date             time.Time `json:"date"`
	ProcessedDate    time.Time `json:"processedDate"`
	OriginalAmount   float64   `json:"originalAmount"`
	OriginalCurrency string    `json:"originalCurrency"`
	ChargedAmount    float64   `json:"chargedAmount"`
	ChargedCurrency  string    `json:"chargedCurrency"`
	Description      string    `json:"description"`
	Memo             string    `json:"memo"`
	Category         string    `json:"category"`
	Account          string    `json:"account"`
	CompanyId        string    `json:"companyId"`
	Hash             string    `json:"hash"`
}

type TransactionsFileUploadedEvent struct {
	Bucket      string    `json:"bucket"`
	Name        string    `json:"name"`
	TimeCreated time.Time `json:"timeCreated"`
}
