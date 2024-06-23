package server

import (
	"data-explorer/pkg/dataexplorer/conf"
	"data-explorer/pkg/dataexplorer/connection"
	"data-explorer/pkg/dataexplorer/controllers"
	"data-explorer/pkg/dataexplorer/models"
	"data-explorer/pkg/dataexplorer/repositories"
	"data-explorer/pkg/dataexplorer/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(connectionsConfiguration *conf.ConnectionsConfiguration) (*Server, error) {
	r := gin.Default()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}

	if err := db.AutoMigrate(
		&models.Issue{},
		&models.IssueSection{},
		&models.SQLQuery{},
	); err != nil {
		return nil, err
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	connectionHolder := connection.NewConnectionHolder(connectionsConfiguration.Connections)
	queryService, err := services.NewQueryService(connectionHolder)
	if err != nil {
		return nil, err
	}

	queryController := controllers.NewQueryController(queryService)
	r.POST("/query", queryController.Query)

	repository := repositories.NewRepository(db)
	issueController := controllers.NewIssuesController(repository, queryService)
	r.POST("/issues", issueController.Create)
	r.POST("/issues/:issueId/sections", issueController.CreateSection)
	r.GET("/issues/:issueId/sections", issueController.ListSections)
	r.POST("/issues/:issueId/sections/:sectionId/queries", issueController.CreateQuery)

	return &Server{
		engine: r,
	}, nil
}

func (server *Server) Run() error {
	return server.engine.Run()
}
