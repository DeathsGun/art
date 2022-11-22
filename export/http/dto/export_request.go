package dto

import "time"

type ExportRequest struct {
	Date     time.Time `json:"date" form:"date"`
	Provider string    `json:"provider" form:"provider"`
}
