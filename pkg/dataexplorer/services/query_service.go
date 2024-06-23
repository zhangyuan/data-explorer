package services

import (
	"context"
	"data-explorer/pkg/dataexplorer/conf"
	"data-explorer/pkg/dataexplorer/db"
	"data-explorer/pkg/dataexplorer/template"
	"errors"
)

type QueryService struct {
	Conf *conf.ConnectionsConfiguration
}

func NewQueryService(conf *conf.ConnectionsConfiguration) (*QueryService, error) {
	return &QueryService{
		Conf: conf,
	}, nil
}

func (s *QueryService) FindConnection(name string) *conf.Connection {
	for idx := range s.Conf.Connections {
		connection := s.Conf.Connections[idx]
		if connection.Id == name {
			return &connection
		}
	}
	return nil
}

func (s *QueryService) QueryWithParams(ctx context.Context, connectionName string, sqlQuery string, params map[string]string) (*db.QueryResult, error) {
	connection := s.FindConnection(connectionName)
	if connection == nil {
		return nil, errors.New("invalid connection name")
	}
	dsn := connection.DSN

	if params != nil {
		sqlQuery = template.SimpleCompile(sqlQuery, params)
	}

	return db.Query(ctx, dsn, sqlQuery)
}
