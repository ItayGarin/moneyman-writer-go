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

type GoogleCloudStorageEvent struct {
	Kind                    string    `json:"kind"`
	ID                      string    `json:"id"`
	SelfLink                string    `json:"selfLink"`
	Name                    string    `json:"name"`
	Bucket                  string    `json:"bucket"`
	Generation              string    `json:"generation"`
	Metageneration          string    `json:"metageneration"`
	ContentType             string    `json:"contentType"`
	TimeCreated             time.Time `json:"timeCreated"`
	Updated                 time.Time `json:"updated"`
	StorageClass            string    `json:"storageClass"`
	TimeStorageClassUpdated time.Time `json:"timeStorageClassUpdated"`
	Size                    string    `json:"size"`
	Md5Hash                 string    `json:"md5Hash"`
	MediaLink               string    `json:"mediaLink"`
	Crc32c                  string    `json:"crc32c"`
	Etag                    string    `json:"etag"`
}
