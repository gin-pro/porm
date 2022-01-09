package porm

type Interface interface {
	Insert(...interface{}) (int64, error)
	Delete(...interface{}) (int64, error)
	Update(bean interface{}, condiBeans ...interface{}) (int64, error)
	Where(interface{}, ...interface{}) *Session
}

type Table interface {
	TableName() string
}
