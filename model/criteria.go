package model

import (
	"net/url"
)

type Criteria map[string][]any

func Convert(criteria url.Values, schema []Field) (Criteria, error) {
	output := make(Criteria)

	for _, field := range schema {
		values, found := criteria[field.Name]

		if !found {
			continue
		}

		options := make([]any, 0)

		for _, value := range values {
			converted, err := field.Type.Convert(value)
			if err != nil {
				return nil, err
			}

			options = append(options, converted)
		}

		output[field.Name] = options
	}

	return output, nil
}
