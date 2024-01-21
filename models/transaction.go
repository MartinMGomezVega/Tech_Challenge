package models

import "time"

type Transaction struct {
	Cuil        string    `json:"cuil"`
	Date        time.Time `json:"date"`
	Transaction float64   `json:"transaction"`
}
