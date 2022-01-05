package porm

import (
	"database/sql"
	"github.com/gin-pro/porm/core"
)

type Engine struct {
	db *core.DB
}

func NewEngine(driverName string, dataSourceName string) (*Engine, error) {
	db, err := core.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Engine{
		db: db,
	}, nil
}

func (e *Engine) NewDB(driverName, dataSourceName string) (*core.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &core.DB{
		DB: db,
	}, nil
}

func (e *Engine) DB() *core.DB {
	return e.db
}
