package models

type ProductHasExpiry struct {
	Name       string `json:"name"`
	Barcode    string `json:"barcode"`
	ExpireDate string `json:"expire_date"`
	Quantity   int    `json:"quantity"`
}
