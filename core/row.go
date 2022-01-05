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
			newDest[j] = vvs.Field(i).Addr().Interface()
			j = j + 1
		}
	}
	return rs.Rows.Scan(newDest...)
}
