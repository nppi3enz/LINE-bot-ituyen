package models

type Expiry struct {
	ExpireDate string `json:"expire_date"`
	Quantity   int    `json:"quantity"`
	Product    string
}
