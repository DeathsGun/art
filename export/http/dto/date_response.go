package dto

import "time"

type DateResponse struct {
	Date time.Time `json:"date"`
}
