package server

import (
	"data-explorer/pkg/dataexplorer/controllers"
	"data-explorer/pkg/dataexplorer/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	engine *gin.Engine
}

func NewServer() (*Server, error) {
	r := gin.Default()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}

	if err := db.AutoMigrate(&models.Issue{}, &models.IssueSection{}); err != nil {
		return nil, err
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	queryController := controllers.NewQueryController()
	r.POST("/query", queryController.Query)

	issueController := controllers.NewIssuesController(db)
	r.POST("/issues", issueController.Create)
	r.POST("/issues/:issueId/sections", issueController.CreateSection)
	r.GET("/issues/:issueId/sections", issueController.ListSections)

	return &Server{
		engine: r,
	}, nil
}

func (server *Server) Run() error {
	return server.engine.Run()
}
