package models

import "time"

type Issue struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Sections []IssueSection `json:"sections"`
}

type IssueSection struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IssueID   uint64    `json:"issueId"`

	Header string `json:"header"`
	Body   string `json:"body"`
	Footer string `json:"footer"`

	Queries []Query `json:"queries"`
}

type Query struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	IssueSectionID uint64 `json:"sectionId"`

	Title  string `json:"title"`
	Query  string `json:"query"`
	Result string `json:"result"`
}
