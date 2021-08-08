package models

type ProductHasExpiry struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Barcode    string `json:"barcode"`
	ExpireDate string `json:"expire_date"`
	Quantity   int    `json:"quantity"`
}
