package connection

import (
	"data-explorer/pkg/dataexplorer/conf"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
)

type Connection struct {
	Id string
	DB *sqlx.DB
}

func NewConnection(id string, db *sqlx.DB) *Connection {
	return &Connection{
		Id: id,
		DB: db,
	}
}

type ConnectionHolder struct {
	mu            sync.Mutex
	Connections   []*Connection
	Configuration []conf.Connection
}

func NewConnectionHolder(configuration []conf.Connection) *ConnectionHolder {
	return &ConnectionHolder{
		Configuration: configuration,
	}
}

func (holder *ConnectionHolder) GetDB(id string) (*sqlx.DB, error) {
	holder.mu.Lock()
	defer holder.mu.Unlock()

	for idx := range holder.Connections {
		if holder.Connections[idx].Id == id {
			return holder.Connections[idx].DB, nil
		}
	}

	for idx := range holder.Configuration {
		conf := holder.Configuration[idx]
		if conf.Id == id {
			db, err := NewDB(conf.DSN)
			if err != nil {
				return nil, err
			}
			connection := NewConnection(id, db)
			holder.Connections = append(holder.Connections, connection)
			return db, nil
		}
	}

	return nil, fmt.Errorf("connection id is invalid: %s", id)
}
