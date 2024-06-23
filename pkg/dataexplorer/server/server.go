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
		&models.Section{},
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
	mainController := controllers.NewMainController(repository, queryService)
	r.POST("/issues", mainController.CreateIssue)
	r.GET("/issues", mainController.ListIssues)
	r.GET("/issues/:issueId", mainController.GetIssue)
	r.DELETE("/issues/:issueId", mainController.DeleteIssue)
	r.PATCH("/issues/:issueId", mainController.PatchIssue)

	r.POST("/issues/:issueId/sections", mainController.CreateSection)
	r.GET("/issues/:issueId/sections", mainController.ListSections)
	r.GET("/issues/:issueId/sections/:sectionId", mainController.GetSection)
	r.DELETE("/sections/:sectionId", mainController.DeleteSection)
	r.PATCH("/sections/:sectionId", mainController.PatchSection)

	r.POST("/issues/:issueId/sections/:sectionId/queries", mainController.CreateQuery)
	r.GET("/issues/:issueId/sections/:sectionId/queries", mainController.ListQueries)
	r.GET("/issues/:issueId/sections/:sectionId/queries/:queryId", mainController.GetQuery)
	r.DELETE("/queries/:queryId", mainController.DeleteQuery)
	r.PATCH("/queries/:queryId", mainController.PatchQuery)

	return &Server{
		engine: r,
	}, nil
}

func (server *Server) Run() error {
	return server.engine.Run()
}
