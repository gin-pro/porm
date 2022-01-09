package statements

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type Statement struct {
	TableName string

	IsSelect  bool
	Fields    []interface{}
	Field     []interface{}
	QueryStr  string
	UpdateStr string
	DeleteStr string

	Args      []interface{}
	LastError error
}

func (s *Statement) Where(query interface{}, args ...interface{}) *Statement {
	return s.And(query, args...)
}

func (s *Statement) And(query interface{}, args ...interface{}) *Statement {
	switch query.(type) {
	case string:
		s.QueryStr = fmt.Sprintf("%v ", query)
	default:
		s.LastError = errors.New("query is not string")
	}
	s.Args = args
	return s
}

func (s *Statement) Reset() {
	s.TableName = ""
	s.QueryStr = ""
	s.UpdateStr = ""
	s.DeleteStr = ""
	s.Args = nil
	s.LastError = nil
}

func (s *Statement) QueryStructToField(bean interface{}) {
	if s.IsSelect {
		return
	}
	s.StructToField(bean)

}
func (s *Statement) StructToField(bean interface{}) {
	vv := reflect.ValueOf(bean)
	tyo := reflect.TypeOf(bean)
	if vv.Kind() == reflect.Ptr {
		if vv.Elem().Kind() == reflect.Slice {
			tyo = tyo.Elem().Elem()
			if tyo.Kind() == reflect.Ptr {
				vv = reflect.New(tyo.Elem())
			} else {
				vv = reflect.New(tyo)
			}
		}
		vv = vv.Elem()
	}
	ls := make([]interface{}, 0)
	fs := make([]interface{}, 0)
	for i := 0; i < vv.NumField(); i++ {
		f := vv.Type().Field(i)
		tag := f.Tag
		t := tag.Get("porm")
		if t == "" || t == "-" {
			continue
		}
		ls = append(ls, t)
		fs = append(fs, f.Name)
	}
	s.Fields = ls
	s.Field = fs
}

func (s *Statement) GetValueByFieldWithOutID(bean interface{}, fieldName string) interface{} {
	vv := reflect.ValueOf(bean)
	tyo := reflect.TypeOf(bean)
	if vv.Kind() == reflect.Ptr {
		if vv.Elem().Kind() == reflect.Slice {
			tyo = tyo.Elem().Elem()
			if tyo.Kind() == reflect.Ptr {
				vv = reflect.New(tyo.Elem())
			} else {
				vv = reflect.New(tyo)
			}
		}
		vv = vv.Elem()
	}
	f := vv.FieldByName(fieldName)
	if fieldName == "Id" || fieldName == "iD" || fieldName == "ID" {
		if f.Int() <= 0 {
			return sql.NullString{}
		}
	}
	switch f.Kind() {
	case reflect.Int:
		return f.Int()
	case reflect.String:
		return f.String()
	}
	return f.String()
}

func (s *Statement) GetValueByField(bean interface{}, fieldName string) interface{} {
	vv := reflect.ValueOf(bean)
	tyo := reflect.TypeOf(bean)
	if vv.Kind() == reflect.Ptr {
		if vv.Elem().Kind() == reflect.Slice {
			tyo = tyo.Elem().Elem()
			if tyo.Kind() == reflect.Ptr {
				vv = reflect.New(tyo.Elem())
			} else {
				vv = reflect.New(tyo)
			}
		}
		vv = vv.Elem()
	}
	f := vv.FieldByName(fieldName)
	switch f.Kind() {
	case reflect.Int:
		return f.Int()
	case reflect.String:
		return f.String()
	}
	return f.String()
}
