package repo

import (
	"encoding/json"

	"github.com/jjmrocha/oblivion/model"
)

func unmarshalSchema(data []byte) ([]model.Field, error) {
	schema := make([]model.Field, 0)
	err := json.Unmarshal(data, &schema)
	return schema, err
}

func marshalSchema(schema []model.Field) ([]byte, error) {
	data, err := json.Marshal(schema)
	return data, err
}
