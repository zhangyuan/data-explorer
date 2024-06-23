package controllers

import (
	"data-explorer/pkg/dataexplorer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryRequest struct {
	ConnectionId string            `json:"connection_id" binding:"required"`
	Title        string            `json:"title"`
	Query        string            `json:"query"`
	Params       map[string]string `json:"params"`
}

type QueryController struct {
	queryService *services.QueryService
}

func NewQueryController(queryService *services.QueryService) *QueryController {
	return &QueryController{
		queryService: queryService,
	}
}

func (controller *QueryController) Query(c *gin.Context) {
	var request QueryRequest
	if err := c.Bind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sql := controller.queryService.CompileSQL(request.Query, request.Params)

	result, err := controller.queryService.Query(c, request.ConnectionId, sql)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"query":  request.Query,
		"params": request.Params,
		"result": result,
	})
}
