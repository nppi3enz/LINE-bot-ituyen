package models

import (
	"time"
)

type Expiry struct {
	ID         string                 `json:"id"`
	ExpireDate time.Time              `json:"expire_date"`
	Quantity   int64                  `json:"quantity"`
	Product    map[string]interface{} `json:"product"`
}
