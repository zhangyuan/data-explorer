package services

import (
	"context"
	"data-explorer/pkg/dataexplorer/connection"
	"data-explorer/pkg/dataexplorer/template"
)

type QueryService struct {
	connectionHolder *connection.ConnectionHolder
}

func NewQueryService(connectionHolder *connection.ConnectionHolder) (*QueryService, error) {
	return &QueryService{
		connectionHolder: connectionHolder,
	}, nil
}

func (s *QueryService) Query(
	ctx context.Context,
	connectionId string,
	sqlQuery string,
) (*connection.QueryResult, error) {
	db, err := s.connectionHolder.GetDB(connectionId)
	if err != nil {
		return nil, err
	}
	return connection.Query(ctx, db, sqlQuery)
}

func (s *QueryService) CompileSQL(
	sqlQuery string,
	params map[string]string,
) string {
	if params != nil {
		sqlQuery = template.SimpleCompile(sqlQuery, params)
	}

	return sqlQuery
}
