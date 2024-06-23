package repositories

import (
	"data-explorer/pkg/dataexplorer/models"
	"strconv"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) FindIssueByID(issueId uint64) (*models.Issue, error) {
	var issue models.Issue
	if tx := r.DB.Select("id").First(&issue, issueId); tx.Error != nil {
		return nil, tx.Error
	}
	return &issue, nil
}

func (r *Repository) FindSectionByStringID(sectionId uint64) (*models.Section, error) {
	var section models.Section
	if tx := r.DB.Select("id").First(&section, sectionId); tx.Error != nil {
		return nil, tx.Error
	}
	return &section, nil
}

func (r *Repository) CreateIssue(issue *models.Issue) error {
	tx := r.DB.Create(&issue)
	return tx.Error
}

func (r *Repository) Save(value interface{}) error {
	tx := r.DB.Save(value)
	return tx.Error
}

func (r *Repository) CreateQuery(query *models.SQLQuery) error {
	return r.DB.Create(&query).Error
}

func (r *Repository) CreateSection(section *models.Section) error {
	return r.DB.Create(section).Error
}

func GetUint(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

func (r *Repository) FindQuery(issueId uint64, sectionId uint64, queryId uint64) (*models.SQLQuery, error) {
	var query models.SQLQuery
	tx := r.DB.Where(&models.SQLQuery{IssueID: issueId, SectionID: sectionId, ID: queryId}).First(&query)
	return &query, tx.Error
}
