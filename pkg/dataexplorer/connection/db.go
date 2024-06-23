package connection

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	_ "github.com/aliyun/aliyun-odps-go-sdk/sqldriver"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
)

func NewDB(dsn string) (*sqlx.DB, error) {
	if strings.HasPrefix(dsn, "postgres") {
		return sqlx.Connect("postgres", strings.TrimPrefix(dsn, "postgres://"))
	} else if strings.HasPrefix(dsn, "mysql") {
		return sqlx.Connect("mysql", strings.TrimPrefix(dsn, "mysql://"))
	} else if strings.Contains(dsn, "maxcompute.aliyun.com/api") {
		return sqlx.Connect("odps", dsn)
	}
	return nil, errors.New("Invalid DSN")
}

type QueryResult struct {
	SQL         string        `json:"sql"`
	ColumnNames []string      `json:"column_names"`
	ColumnTypes []string      `json:"column_types"`
	Records     []interface{} `json:"records"`
}

func Query(ctx context.Context, db *sqlx.DB, query string) (*QueryResult, error) {
	rows, err := db.QueryxContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	queryResult := QueryResult{
		SQL: query,
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	queryResult.ColumnNames = columns

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	queryResult.ColumnTypes = lo.Map(columnTypes, func(item *sql.ColumnType, index int) string {
		return item.DatabaseTypeName()
	})

	for rows.Next() {
		if record, err := rows.SliceScan(); err != nil {
			return nil, err
		} else {
			queryResult.Records = append(queryResult.Records, record)
		}
	}

	return &queryResult, nil
}
