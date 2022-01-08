package statements

import (
	"errors"
	"fmt"
)

type Statement struct {
	TableName string

	QueryStr  string
	InsertStr string
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
	s.InsertStr = ""
	s.UpdateStr = ""
	s.DeleteStr = ""
	s.Args = nil
	s.LastError = nil
}
