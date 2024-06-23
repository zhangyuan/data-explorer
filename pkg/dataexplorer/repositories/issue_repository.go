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

func (r *Repository) FindIssueByStringID(stringId string) (*models.Issue, error) {
	issueId, err := GetUint(stringId)
	if err != nil {
		return nil, err
	}

	var issue models.Issue
	if tx := r.DB.Select("id").First(&issue, issueId); tx.Error != nil {
		return nil, err
	}
	return &issue, nil
}

func (r *Repository) FindSectionByStringID(stringId string) (*models.IssueSection, error) {
	id, err := GetUint(stringId)
	if err != nil {
		return nil, err
	}

	var section models.IssueSection
	if tx := r.DB.Select("id").First(&section, id); tx.Error != nil {
		return nil, err
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
	db := r.DB.Create(&query)
	return db.Error
}

func (r *Repository) CreateSection(section *models.IssueSection) error {
	return r.DB.Create(section).Error
}

func GetUint(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}
