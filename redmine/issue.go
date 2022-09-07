package redmine

import "time"

type IssueResponse struct {
	Issues     []Issue `json:"issues"`
	TotalCount int     `json:"total_count"`
	Offset     int     `json:"offset"`
	Limit      int     `json:"limit"`
}

type Issue struct {
	Id             int       `json:"id"`
	Project        NamedItem `json:"project"`
	Tracker        NamedItem `json:"tracker"`
	Status         NamedItem `json:"status"`
	Priority       NamedItem `json:"priority"`
	Author         NamedItem `json:"author"`
	AssignedTo     NamedItem `json:"assigned_to"`
	Category       NamedItem `json:"category"`
	Subject        string    `json:"subject"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	DueDate        time.Time `json:"due_date"`
	DoneRatio      int       `json:"done_ratio"`
	IsPrivate      bool      `json:"is_private"`
	EstimatedHours float64   `json:"estimated_hours"`
	CreatedOn      time.Time `json:"created_on"`
	UpdatedOn      time.Time `json:"updated_on"`
	ClosedOn       time.Time `json:"closed_on"`
}
