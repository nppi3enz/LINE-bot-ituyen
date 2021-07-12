package models

type Expire struct {
	ExpireDate string `json:"expire_date"`
	Quantity   int    `json:"quantity"`
	Product    string
}
