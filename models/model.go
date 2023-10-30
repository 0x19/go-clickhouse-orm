package models

type Model interface {
	TableName() string
	ToMap() map[string]interface{}
}
