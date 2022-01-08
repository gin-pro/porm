package porm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-pro/porm/core"
	"github.com/gin-pro/porm/statements"
	"reflect"
	"strings"
)

type Session struct {
	engine    *Engine
	tx        *sql.Tx
	statement *statements.Statement

	ctx context.Context
}

func newSession(engine *Engine) *Session {
	return &Session{engine, nil, &statements.Statement{}, engine.ctx}
}

func (session *Session) Where(query interface{}, args ...interface{}) *Session {
	session.statement.Where(query, args...)
	return session
}

func (session *Session) Get(bean interface{}) (bool, error) {
	return session.get(bean)
}

func (session *Session) Find(bean interface{}) (bool, error) {
	return session.find(bean)
}

func (session *Session) resetStatement() {
	session.statement.Reset()
}

func (session *Session) get(bean interface{}) (bool, error) {
	defer session.resetStatement()

	if session.statement.LastError != nil {
		return false, session.statement.LastError
	}

	vv := reflect.ValueOf(bean)
	if vv.Kind() != reflect.Ptr {
		return false, errors.New("bean is not ptr")
	} else if vv.Elem().Kind() == reflect.Ptr {
		return false, errors.New("a ptr to a ptr is allow")
	} else if vv.IsNil() {
		return false, errors.New("bean is nil")
	}

	rows, err := session.queryRows()
	if err != nil {
		return false, err
	}
	defer rows.Close()

	switch vv.Elem().Kind() {
	case reflect.Struct:
		return session.getStruct(rows, bean)
	case reflect.Map:
		return session.getMap(rows, bean)
	}
	return session.getStruct(rows, bean)
}

func (session *Session) find(bean interface{}) (bool, error) {
	defer session.resetStatement()

	if session.statement.LastError != nil {
		return false, session.statement.LastError
	}

	vv := reflect.ValueOf(bean)
	if vv.Kind() != reflect.Ptr {
		return false, errors.New("bean is not ptr")
	} else if vv.Elem().Kind() == reflect.Ptr {
		return false, errors.New("a ptr to a ptr is allow")
	} else if vv.IsNil() {
		return false, errors.New("bean is nil")
	}

	rows, err := session.queryRows()
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return session.getSlice(rows, bean)
}

func (session *Session) queryRows() (*core.Rows, error) {
	builer := strings.Builder{}
	builer.WriteString("SELECT * FROM ")
	builer.WriteString("`" + session.statement.TableName + "` ")
	if len(session.statement.QueryStr) > 0 {
		builer.WriteString("WHERE ")
		builer.WriteString(session.statement.QueryStr)
	}
	builer.WriteString(";")
	sqls := builer.String()
	fmt.Println("Query sql  ", sqls)
	row, err := session.engine.DB().QueryContext(session.ctx, sqls, session.statement.Args...)
	if err != nil {
		return nil, err
	}
	return &core.Rows{
		Rows: row,
	}, nil
}

func (session *Session) Table(table string) *Session {
	session.statement.TableName = table
	return session
}

func (session *Session) getStruct(rows *core.Rows, bean interface{}) (bool, error) {
	if rows.Next() {
		err := rows.ScanStructByIndex(bean)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (session *Session) getSlice(rows *core.Rows, bean interface{}) (bool, error) {
	vvv := reflect.ValueOf(bean)
	tyo := reflect.TypeOf(bean).Elem().Elem()

	newArr := make([]reflect.Value, 0)
	for rows.Next() {
		var value reflect.Value
		if tyo.Kind() == reflect.Ptr {
			value = reflect.New(tyo.Elem())
		} else {
			value = reflect.New(tyo)
		}
		err := rows.ScanStructByIndex(value.Interface())
		if err != nil {
			return false, err
		}
		if tyo.Kind() != reflect.Ptr {
			value = value.Elem()
		}
		newArr = append(newArr, value)
	}
	resArr := reflect.Append(vvv.Elem(), newArr...)
	vvv.Elem().Set(resArr)
	return true, nil
}

func (session *Session) getMap(rows *core.Rows, bean interface{}) (bool, error) {
	if rows.Next() {
		err := rows.ScanStructByIndex(bean)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
