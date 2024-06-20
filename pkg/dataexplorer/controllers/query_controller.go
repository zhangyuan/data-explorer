package controllers

import (
	"data-explorer/pkg/dataexplorer/db"
	"data-explorer/pkg/dataexplorer/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type QueryRequest struct {
	ConnectionId string            `json:"connectionId"`
	Title        string            `json:"title"`
	Query        string            `json:"query"`
	Params       map[string]string `json:"params"`
}

type QueryController struct {
}

func NewQueryController() *QueryController {
	return &QueryController{}
}

func (controller *QueryController) Query(c *gin.Context) {
	var request QueryRequest
	if err := c.Bind(&request); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	dsn := os.Getenv("DSN")

	var query string
	if request.Params != nil {
		query = template.SimpleCompile(request.Query, request.Params)
	} else {
		query = request.Query
	}

	result, err := db.Query(c, dsn, query)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
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
