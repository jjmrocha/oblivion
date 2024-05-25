package model

type Bucket struct {
	Name   string  `json:"name"`
	Schema []Field `json:"schema"`
}

type DataType string

const (
	StringDataType DataType = "string"
	NumberDataType DataType = "number"
	BoolDataType   DataType = "bool"
)

type Field struct {
	Name     string   `json:"field"`
	Type     DataType `json:"type"`
	Required bool     `json:"not-null"`
}
