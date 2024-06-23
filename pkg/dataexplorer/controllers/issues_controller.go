package controllers

import (
	"data-explorer/pkg/dataexplorer/connection"
	"data-explorer/pkg/dataexplorer/models"
	"data-explorer/pkg/dataexplorer/repositories"
	"data-explorer/pkg/dataexplorer/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/datatypes"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	var request CreateIssueSectionRequest

	if err := c.Bind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	section := models.IssueSection{
		Header:  request.Header,
		Body:    request.Body,
		Footer:  request.Footer,
		IssueID: issue.ID,
	}

	if err := controller.repository.CreateSection(&section); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
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
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	section, err := controller.repository.FindSectionByStringID(c.Param("sectionId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	sqlQuery := models.SQLQuery{
		IssueSectionID: section.ID,
		Title:          request.Title,
		Query:          request.Query,
	}

	if request.Params != nil {
		paramsBytes, err := jsoniter.Marshal(request.Params)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
			return
		}
		sqlQuery.Params = datatypes.JSON(paramsBytes)
	}

	if err := controller.repository.CreateQuery(&sqlQuery); err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	startTime := time.Now()
	queryResult, err := controller.queryService.QueryWithParams(c, request.ConnectionId, request.Query, request.Params)
	finishTime := time.Now()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	resultBytes, err := jsoniter.Marshal(queryResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	sqlQuery.Result = datatypes.JSON(resultBytes)
	sqlQuery.Duration = finishTime.Sub(startTime).Milliseconds()

	if err := controller.repository.Save(&sqlQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, QueryResponse{
		Query:    request.Query,
		Params:   request.Params,
		Result:   queryResult,
		Duration: sqlQuery.Duration,
	})
}

type QueryResponse struct {
	Query    string                  `json:"query"`
	Params   map[string]string       `json:"params"`
	Result   *connection.QueryResult `json:"result"`
	Duration int64                   `json:"duration"`
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
