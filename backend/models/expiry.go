package models

import (
	"time"
)

type Expiry struct {
	ExpireDate time.Time              `json:"expire_date"`
	Quantity   int64                  `json:"quantity"`
	Product    map[string]interface{} `json:"product"`
}
