package server

import (
	"data-explorer/pkg/dataexplorer/conf"
	"data-explorer/pkg/dataexplorer/connection"
	"data-explorer/pkg/dataexplorer/controllers"
	"data-explorer/pkg/dataexplorer/models"
	"data-explorer/pkg/dataexplorer/repositories"
	"data-explorer/pkg/dataexplorer/services"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(connectionsConfiguration *conf.ConnectionsConfiguration) (*Server, error) {
	r := gin.Default()

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: gormLogger,
	})
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
	repository := repositories.NewRepository(db)
	mainController := controllers.NewMainController(repository, queryService)

	api := r.Group("/api")
	{
		api.POST("/query", queryController.Query)

		api.POST("/issues", mainController.CreateIssue)
		api.GET("/issues", mainController.ListIssues)
		api.GET("/issues/:issueId", mainController.GetIssue)
		api.DELETE("/issues/:issueId", mainController.DeleteIssue)
		api.PATCH("/issues/:issueId", mainController.PatchIssue)

		api.POST("/issues/:issueId/sections", mainController.CreateSection)
		api.GET("/issues/:issueId/sections", mainController.ListSections)
		api.GET("/issues/:issueId/sections/:sectionId", mainController.GetSection)
		api.DELETE("/sections/:sectionId", mainController.DeleteSection)
		api.PATCH("/sections/:sectionId", mainController.PatchSection)

		api.POST("/issues/:issueId/sections/:sectionId/queries", mainController.CreateQuery)
		api.GET("/issues/:issueId/sections/:sectionId/queries", mainController.ListQueries)
		api.GET("/issues/:issueId/sections/:sectionId/queries/:queryId", mainController.GetQuery)
		api.DELETE("/queries/:queryId", mainController.DeleteQuery)
		api.PATCH("/queries/:queryId", mainController.PatchQuery)
	}

	return &Server{
		engine: r,
	}, nil
}

func (server *Server) Run() error {
	return server.engine.Run()
}
