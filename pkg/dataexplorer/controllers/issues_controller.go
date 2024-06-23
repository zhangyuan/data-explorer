package controllers

import (
	"data-explorer/pkg/dataexplorer/models"
	"data-explorer/pkg/dataexplorer/repositories"
	"data-explorer/pkg/dataexplorer/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/datatypes"
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

func (controller *IssuesController) CreateIssue(c *gin.Context) {
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

func (controller *IssuesController) ListIssues(c *gin.Context) {
	page, err := GetIntOr(c.Query("page"), 1)
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	limit, err := GetIntOr(c.Query("page_size"), 20)
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	offset := (page - 1) * limit

	var issues []models.Issue

	if tx := controller.repository.DB.Limit(limit).Offset(offset).Find(&issues); tx.Error != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(tx.Error))
		return
	}

	c.JSON(200, issues)
}

func (controller *IssuesController) GetIssue(c *gin.Context) {
	issueId, err := GetUint(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	issue, err := controller.repository.FindIssueByID(issueId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, issue)
}

func (controller *IssuesController) DeleteIssue(c *gin.Context) {
	issueId, err := GetUint(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	if err := controller.repository.DeleteIssueByID(issueId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GetUint(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

func GetUint64Or(value string, defaultValue uint64) (uint64, error) {
	if value == "" {
		return defaultValue, nil
	}
	return strconv.ParseUint(value, 10, 64)
}

func GetIntOr(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
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

	issue, err := controller.repository.FindIssueByID(issueId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	var request CreateIssueSectionRequest

	if err := c.Bind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	section := models.Section{
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
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	issueId, err := GetUint(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	issue, err := controller.repository.FindIssueByID(issueId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	sectionId, err := GetUint(c.Param("sectionId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	section, err := controller.repository.FindSectionByID(sectionId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResponse(err))
		return
	}

	sqlQuery := models.SQLQuery{
		ConnectionId: request.ConnectionId,
		IssueID:      issue.ID,
		SectionID:    section.ID,
		Title:        request.Title,
		Query:        request.Query,
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

	sql := controller.queryService.CompileSQL(request.Query, request.Params)
	startTime := time.Now()
	queryResult, err := controller.queryService.Query(c, request.ConnectionId, sql)
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
	sqlQuery.Sql = sql

	if err := controller.repository.Save(&sqlQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewQueryResponse(&sqlQuery))
}

func (controller *IssuesController) ListQueries(c *gin.Context) {
	page, err := GetIntOr(c.Query("page"), 1)
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}
	limit, err := GetIntOr(c.Query("page_size"), 20)
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	offset := (page - 1) * limit

	issueId, err := GetUint(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	sectionId, err := GetUint(c.Param("sectionId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	var sqlQueries []models.SQLQuery

	if tx := controller.repository.DB.
		Limit(limit).Offset(offset).
		Where(&models.SQLQuery{IssueID: issueId, SectionID: sectionId}).
		Find(&sqlQueries); tx.Error != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(tx.Error))
		return
	}

	c.JSON(200, sqlQueries)
}

func NewQueryResponse(sqlQuery *models.SQLQuery) *QueryResponse {
	return &QueryResponse{
		ID:       sqlQuery.ID,
		Query:    sqlQuery.Query,
		Params:   sqlQuery.Params,
		SQL:      sqlQuery.Sql,
		Result:   sqlQuery.Result,
		Duration: sqlQuery.Duration,
	}
}

type QueryResponse struct {
	ID       uint64         `json:"id"`
	Query    string         `json:"query"`
	Params   datatypes.JSON `json:"params"`
	SQL      string         `json:"sql"`
	Result   datatypes.JSON `json:"result"`
	Duration int64          `json:"duration"`
}

func (controller *IssuesController) ListSections(c *gin.Context) {
	page, err := GetIntOr(c.Query("page"), 1)
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	limit, err := GetIntOr(c.Query("page_size"), 20)
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	offset := (page - 1) * limit

	issueId, err := GetUint(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	var sections []models.Section

	if tx := controller.repository.DB.Limit(limit).Offset(offset).Where(&models.Section{IssueID: issueId}).Find(&sections); tx.Error != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(tx.Error))
		return
	}

	c.JSON(200, sections)
}

func (controller *IssuesController) GetSection(c *gin.Context) {
	issueId, err := GetUint(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	sectionId, err := GetUint(c.Param("sectionId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	section, err := controller.repository.FindSection(sectionId, &models.Section{IssueID: issueId})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, section)
}

func (controller *IssuesController) DeleteSection(c *gin.Context) {
	sectionId, err := GetUint(c.Param("sectionId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	if err := controller.repository.DeleteSectionByID(sectionId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (controller *IssuesController) GetQuery(c *gin.Context) {
	issueId, err := GetUint(c.Param("issueId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	sectionId, err := GetUint(c.Param("sectionId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	queryId, err := GetUint(c.Param("queryId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	sqlQuery, err := controller.repository.FindQuery(queryId, &models.SQLQuery{
		SectionID: sectionId,
		IssueID:   issueId,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewQueryResponse(sqlQuery))
}

func (controller *IssuesController) DeleteQuery(c *gin.Context) {
	queryId, err := GetUint(c.Param("queryId"))
	if err != nil {
		c.AbortWithStatusJSON(400, NewErrorResponse(err))
		return
	}

	if err := controller.repository.DeleteQuery(queryId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
