package core

import (
	"database/sql"
	"errors"
	"reflect"
)

type Rows struct {
	*sql.Rows
	db *DB
}

func (rs *Rows) ScanMap(dest interface{}) error {
	vv := reflect.ValueOf(dest)
	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Map {
		return errors.New("dest should be a map's pointer")
	}

	columns, err := rs.Columns()
	if err != nil {
		return err
	}

	newDest := make([]interface{}, len(columns))
	vvv := vv.Elem()

	slice := reflect.MakeSlice(reflect.SliceOf(vvv.Type().Elem()), len(columns), len(columns))
	for i := range columns {
		newDest[i] = slice.Index(i).Addr().Interface()
	}

	err = rs.Rows.Scan(newDest...)
	if err != nil {
		return err
	}

	for i, name := range columns {
		vname := reflect.ValueOf(name)
		vvv.SetMapIndex(vname, reflect.ValueOf(newDest[i]).Elem())
	}
	return nil
}

func (rs *Rows) ScanStructByIndex(dest ...interface{}) error {
	if len(dest) == 0 {
		return errors.New("at least one struct")
	}
	vvvs := make([]reflect.Value, len(dest))
	for i, s := range dest {
		vv := reflect.ValueOf(s)
		if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Struct {
			return errors.New("dest should be a struct's pointer")
		}
		vvvs[i] = vv.Elem()
	}

	columns, err := rs.Columns()
	if err != nil {
		return err
	}

	newDest := make([]interface{}, len(columns))
	var j = 0
	for _, vvs := range vvvs {
		for i := 0; i < vvs.NumField(); i++ {
			tag := vvs.Type().Field(i).Tag
			t := tag.Get("porm")
			if t == "" || t == "-" {
				continue
			}
			newDest[j] = vvs.Field(i).Addr().Interface()
			j = j + 1
		}
	}
	return rs.Rows.Scan(newDest...)
}

type Row struct {
	rows *Rows
	err  error
}

func ErrorRow(err error) *Row {
	return &Row{
		err: err,
	}
}

func NewRow(rows *Rows, err error) *Row {
	return &Row{rows, err}
}

func (row *Row) Columns() ([]string, error) {
	if row.err != nil {
		return nil, row.err
	}
	return row.rows.Columns()
}

func (row *Row) Scan(dest ...interface{}) error {
	if row.err != nil {
		return row.err
	}
	defer row.rows.Close()

	for _, v := range dest {
		if _, ok := v.(sql.RawBytes); ok {
			return errors.New("sql: RawBytes isn't allowed on Row.Scan")
		}
	}

	if !row.rows.Next() {
		if err := row.rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	err := row.rows.ScanStructByIndex(dest...)
	if err != nil {
		return err
	}
	return row.rows.Close()
}

func (row *Row) ScanStructByIndex(dest ...interface{}) error {
	if row.err != nil {
		return row.err
	}
	defer row.rows.Close()

	if !row.rows.Next() {
		if err := row.rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}
	err := row.rows.ScanStructByIndex(dest)
	if err != nil {
		return err
	}
	return row.rows.Close()
}

// ScanMap scan data to a map's pointer
func (row *Row) ScanMap(dest interface{}) error {
	if row.err != nil {
		return row.err
	}
	defer row.rows.Close()

	if !row.rows.Next() {
		if err := row.rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}
	err := row.rows.ScanMap(dest)
	if err != nil {
		return err
	}

	return row.rows.Close()
}
