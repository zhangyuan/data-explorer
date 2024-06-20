package controllers

import (
	"data-explorer/pkg/dataexplorer/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateIssueRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateIssueSectionRequest struct {
	Header string `json:"header"`
	Body   string `json:"body"`
	Footer string `json:"footer"`
}

type IssuesController struct {
	db *gorm.DB
}

func NewIssuesController(db *gorm.DB) *IssuesController {
	return &IssuesController{db: db}
}

func (controller *IssuesController) Create(c *gin.Context) {
	var request CreateIssueRequest

	if err := c.Bind(&request); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	issue := models.Issue{
		Title:       request.Title,
		Description: request.Description,
	}

	db := controller.db.Create(&issue)
	if db.Error != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": db.Error})
		return
	}

	c.JSON(200, issue)
}

func GetUint(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
	}
}

func (controller *IssuesController) CreateSection(c *gin.Context) {
	issueId, err := GetUint(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	var issue models.Issue
	if tx := controller.db.Select("id").First(&issue, issueId); tx.Error != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(tx.Error))
		return
	}

	var request CreateIssueSectionRequest

	if err := c.Bind(&request); err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	section := models.IssueSection{
		Header: request.Header,
		Body:   request.Body,
		Footer: request.Footer,
	}

	db := controller.db.Create(&section)
	if db.Error != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(db.Error))
		return
	}

	if err := db.Model(&issue).Association("Sections").Append(&section); err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	c.JSON(200, section)
}

func (controller *IssuesController) ListSections(c *gin.Context) {
	issueId, err := GetUint(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	var sections []models.IssueSection

	if tx := controller.db.Limit(100).Offset(0).Where(&models.IssueSection{IssueID: issueId}).Find(&sections); tx.Error != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(tx.Error))
		return
	}

	c.JSON(200, sections)
}
