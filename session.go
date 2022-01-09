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

	rows, err := session.queryRows(bean)
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

	rows, err := session.queryRows(bean)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return session.getSlice(rows, bean)
}

func (session *Session) queryRows(bean interface{}) (*core.Rows, error) {
	session.statement.QueryStructToField(bean)

	builder := strings.Builder{}
	builder.WriteString("SELECT ")
	field := ""
	for _, f := range session.statement.Fields {
		field += fmt.Sprintf("`%v`,", f)
	}
	if len(field) > 0 {
		field = field[0 : len(field)-1]
	}
	builder.WriteString(field)
	builder.WriteString(" FROM ")
	builder.WriteString("`" + session.statement.TableName + "` ")
	if len(session.statement.QueryStr) > 0 {
		builder.WriteString("WHERE ")
		builder.WriteString(session.statement.QueryStr)
	}
	builder.WriteString(";")
	sqls := builder.String()
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

func (session *Session) Insert(bean ...interface{}) (int64, error) {
	defer session.resetStatement()

	if session.statement.LastError != nil {
		return 0, session.statement.LastError
	}
	return session.insert(bean...)
}

func (session *Session) insert(bean ...interface{}) (int64, error) {
	sqls := []string{}
	args := []interface{}{}
	for _, b := range bean {
		session.statement.StructToField(b)
		builder := strings.Builder{}
		builder.WriteString("INSERT INTO ")
		builder.WriteString("`" + session.statement.TableName + "` ")

		field := "("
		for _, f := range session.statement.Fields {
			field += fmt.Sprintf("`%v`,", f)
		}
		if len(field) > 0 {
			field = field[0 : len(field)-1]
		}
		builder.WriteString(field + ")")

		builder.WriteString("VALUES ( ")
		vs := ""
		for _, f := range session.statement.Field {
			vs += "?,"
			args = append(args, session.statement.GetValueByFieldWithOutID(b, fmt.Sprintf("%v", f)))
		}
		if len(vs) > 0 {
			vs = vs[0 : len(vs)-1]
		}
		builder.WriteString(vs)
		builder.WriteString(") ")
		builder.WriteString(";")
		sqls = append(sqls, builder.String())
	}
	join := strings.Join(sqls, "\n")
	fmt.Println("insert sql  ", join)
	fmt.Println("insert args  ", args)
	res, err := session.engine.DB().Exec(join, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (session *Session) Update(bean interface{}) (int64, error) {
	defer session.resetStatement()

	if session.statement.LastError != nil {
		return 0, session.statement.LastError
	}
	return session.update(bean)
}

func (session *Session) update(bean interface{}) (int64, error) {
	args := []interface{}{}
	session.statement.StructToField(bean)
	builder := strings.Builder{}
	builder.WriteString("UPDATE ")
	builder.WriteString("`" + session.statement.TableName + "` ")
	builder.WriteString("SET ")
	field := ""
	for i, f := range session.statement.Fields {
		args = append(args, session.statement.GetValueByField(bean, fmt.Sprintf("%v", session.statement.Field[i])))
		field += fmt.Sprintf("`%v` = ?,", f)
	}
	if len(field) > 0 {
		field = field[0 : len(field)-1]
	}
	builder.WriteString(field)
	if len(session.statement.QueryStr) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(session.statement.QueryStr)
	}
	builder.WriteString(";")
	sqls := builder.String()
	args = append(args, session.statement.Args...)
	fmt.Println("update sql  ", sqls)
	fmt.Println("update args  ", args)
	res, err := session.engine.DB().Exec(sqls, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (session *Session) Delete(bean interface{}) (int64, error) {
	defer session.resetStatement()

	if session.statement.LastError != nil {
		return 0, session.statement.LastError
	}
	return session.delete(bean)
}
func (session *Session) delete(bean interface{}) (int64, error) {
	args := []interface{}{}
	session.statement.StructToField(bean)
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM ")
	builder.WriteString("`" + session.statement.TableName + "` ")
	builder.WriteString("WHERE ")
	field := ""
	for i, f := range session.statement.Fields {
		args = append(args, session.statement.GetValueByField(bean, fmt.Sprintf("%v", session.statement.Field[i])))
		field += fmt.Sprintf("`%v` = ? and", f)
	}
	if len(session.statement.QueryStr) > 0 {
		field += " " + session.statement.QueryStr
	} else {
		if len(field) > 0 {
			field = field[0 : len(field)-3]
		}
	}
	builder.WriteString(field)
	builder.WriteString(";")
	sqls := builder.String()
	args = append(args, session.statement.Args...)
	fmt.Println("Delete sql  ", sqls)
	fmt.Println("Delete args  ", args)
	res, err := session.engine.DB().Exec(sqls, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()

}

//func (session *Session) Select(fields ...interface{}) *Session {
//	session.statement.Fields = fields
//	session.statement.IsSelect = true
//	return session
//}
