package model

import "strconv"

type DataType string

const (
	StringDataType DataType = "string"
	NumberDataType DataType = "number"
	BoolDataType   DataType = "bool"
)

func (d DataType) Convert(value string) (any, error) {
	switch d {
	case NumberDataType:
		converted, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}

		return converted, nil
	case BoolDataType:
		converted, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}

		return converted, nil
	}

	return value, nil
}

func (d DataType) ValidValue(value any) bool {
	switch d {
	case StringDataType:
		_, ok := value.(string)
		return ok
	case NumberDataType:
		_, ok := value.(float64)
		return ok
	case BoolDataType:
		_, ok := value.(bool)
		return ok
	}

	return false
}
