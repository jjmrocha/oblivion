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
