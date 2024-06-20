package models

import "time"

type Issue struct {
	ID          uint64         `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Sections    []IssueSection `json:"sections"`
}

type IssueSection struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IssueID   uint64    `json:"issueId"`

	Header string `json:"header"`
	Body   string `json:"body"`
	Footer string `json:"footer"`
}
