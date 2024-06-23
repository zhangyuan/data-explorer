package controllers

import (
	"data-explorer/pkg/dataexplorer/models"
	"data-explorer/pkg/dataexplorer/repositories"
	"data-explorer/pkg/dataexplorer/services"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateIssueRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateQueryRequest struct {
	ConnectionId string            `json:"connection_id" binding:"required"`
	Title        string            `json:"title"`
	Query        string            `json:"query" binding:"required"`
	Params       map[string]string `json:"params"`
}

type CreateIssueSectionRequest struct {
	Header string `json:"header"`
	Body   string `json:"body"`
	Footer string `json:"footer"`
}

type IssuesController struct {
	db           *gorm.DB
	repository   *repositories.Repository
	queryService *services.QueryService
}

func NewIssuesController(
	issueRepository *repositories.Repository,
	queryService *services.QueryService,
) *IssuesController {
	return &IssuesController{
		repository:   issueRepository,
		queryService: queryService,
	}
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

	if err := controller.repository.CreateIssue(&issue); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err})
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
	issue, err := controller.repository.FindIssueByStringID(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	var request CreateIssueSectionRequest

	if err := c.Bind(&request); err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	section := models.IssueSection{
		Header:  request.Header,
		Body:    request.Body,
		Footer:  request.Footer,
		IssueID: issue.ID,
	}

	if err := controller.repository.CreateSection(&section); err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	c.JSON(200, section)
}

func (controller *IssuesController) CreateQuery(c *gin.Context) {
	var request CreateQueryRequest

	if err := c.Bind(&request); err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	_, err := controller.repository.FindIssueByStringID(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	section, err := controller.repository.FindSectionByStringID(c.Param("sectionId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	sqlQuery := models.SQLQuery{
		IssueSectionID: section.ID,
		Title:          request.Title,
		Query:          request.Query,
	}

	if err := controller.repository.CreateQuery(&sqlQuery); err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	result, err := controller.queryService.QueryWithParams(c, request.ConnectionId, request.Query, request.Params)
	if err != nil {
		c.AbortWithStatusJSON(500, NewErrorResponse(err))
		return
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		c.AbortWithStatusJSON(500, NewErrorResponse(err))
		return
	}

	sqlQuery.Result = string(jsonBytes)

	if err := controller.repository.Save(&sqlQuery); err != nil {
		c.AbortWithStatusJSON(500, NewErrorResponse(err))
		return
	}

	c.JSON(200, result)
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
