package models

import "time"

type Issue struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Sections []IssueSection `json:"sections"`
}

type IssueSection struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IssueID   uint64    `json:"issue_id"`

	Header string `json:"header"`
	Body   string `json:"body"`
	Footer string `json:"footer"`

	Queries []SQLQuery `json:"queries"`
}

type SQLQuery struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IssueSectionID uint64 `json:"section_id"`

	ConnectionId string `json:"connection_id"`
	Title        string `json:"title"`
	Query        string `json:"query"`
	Result       string `json:"result"`
}
