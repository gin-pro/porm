package porm

import (
	"context"
	"database/sql"
	"github.com/gin-pro/porm/core"
)

type Engine struct {
	db *core.DB

	ctx context.Context
}

func NewEngine(driverName string, dataSourceName string) (*Engine, error) {
	db, err := core.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Engine{
		db:  db,
		ctx: context.Background(),
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

func (e *Engine) NewSession() *Session {
	return newSession(e)
}

func (e *Engine) Table(table string) *Session {
	session := e.NewSession()
	return session.Table(table)
}

func (e *Engine) Where(query interface{}, args ...interface{}) *Session {
	session := e.NewSession()
	return session.Where(query, args...)
}

func (e *Engine) Get(bean interface{}) (bool, error) {
	session := e.NewSession()
	return session.Get(bean)
}

func (e *Engine) Find(bean interface{}) (bool, error) {
	session := e.NewSession()
	return session.Find(bean)
}

func (e *Engine) Insert(bean interface{}) (int64, error) {
	session := e.NewSession()
	return session.Insert(bean)
}

func (e *Engine) Update(bean interface{}) (int64, error) {
	session := e.NewSession()
	return session.Update(bean)
}

func (e *Engine) Delete(bean interface{}) (int64, error) {
	session := e.NewSession()
	return session.Delete(bean)
}
