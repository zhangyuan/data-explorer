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

func (r *Repository) FindSectionByID(sectionId uint64) (*models.Section, error) {
	var section models.Section
	if tx := r.DB.Select("id").First(&section, sectionId); tx.Error != nil {
		return nil, tx.Error
	}
	return &section, nil
}

func (r *Repository) FindSection(sectionId uint64, condition *models.Section) (*models.Section, error) {
	var section models.Section
	if tx := r.DB.Where(condition).First(&section, sectionId); tx.Error != nil {
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

func (r *Repository) FindQuery(queryId uint64, where *models.SQLQuery) (*models.SQLQuery, error) {
	var query models.SQLQuery
	tx := r.DB.Where(where).First(&query, queryId)
	return &query, tx.Error
}

func (r *Repository) DeleteQuery(queryId uint64) error {
	return r.DB.Delete(&models.SQLQuery{}, queryId).Error
}

func (r *Repository) DeleteSectionByID(sectionId uint64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := r.DB.Where("section_id = ?", sectionId).Delete(&models.SQLQuery{}).Error; err != nil {
			return err
		}

		if err := r.DB.Delete(&models.Section{}, sectionId).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) DeleteIssueByID(issueId uint64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := r.DB.Where("issue_id = ?", issueId).Delete(&models.SQLQuery{}).Error; err != nil {
			return err
		}

		if err := r.DB.Where("issue_id = ?", issueId).Delete(&models.Section{}).Error; err != nil {
			return err
		}

		if err := r.DB.Delete(&models.Issue{}, issueId).Error; err != nil {
			return err
		}

		return nil
	})
}
