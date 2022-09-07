package redmine

import "time"

type TimeEntriesResponse struct {
	TimeEntries []TimeEntry `json:"time_entries"`
	TotalCount  int         `json:"total_count"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
}

type TimeEntry struct {
	Id           int            `json:"id"`
	Project      NamedItem      `json:"project"`
	Issue        IssueReference `json:"issue"`
	User         NamedItem      `json:"user"`
	Activity     NamedItem      `json:"activity"`
	Hours        float64        `json:"hours"`
	Comments     string         `json:"comments"`
	SpentOn      string         `json:"spent_on"`
	CreatedOn    time.Time      `json:"created_on"`
	UpdatedOn    time.Time      `json:"updated_on"`
	CustomFields []CustomField  `json:"custom_fields"`
}
