package server

import (
	"data-explorer/pkg/dataexplorer/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	queryController := controllers.NewQueryController()

	r.POST("/query", queryController.Query)

	return &Server{
		engine: r,
	}
}

func (server *Server) Run() error {
	return server.engine.Run()
}
