package porm

import "database/sql"

type Session struct {
	engine *Engine
	tx     *sql.Tx
}
