package repo

import "encoding/json"

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
	Indexed  bool     `json:"indexed"`
}

func unmarshalSchema(data []byte) ([]Field, error) {
	schema := make([]Field, 0)
	err := json.Unmarshal(data, &schema)
	return schema, err
}

func marshalSchema(schema []Field) ([]byte, error) {
	data, err := json.Marshal(schema)
	return data, err
}
