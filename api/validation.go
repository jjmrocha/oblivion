package api

import (
	"regexp"

	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/repo"
)

func checkBucketCreation(name string, schema []repo.Field) error {
	if len(name) == 0 {
		return model.Error(model.InvalidBucketName)
	}

	if len(schema) == 0 {
		return model.Error(model.SchemaMissing)
	}

	for _, field := range schema {
		matched, err := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]*[a-zA-Z0-9]$", name)
		if err != nil {
			return model.ErrorWithReason(model.UnexpectedError, err)
		}

		if !matched || len(field.Name) == 0 || len(field.Name) > 30 {
			return model.Error(model.InvalidSchema, field.Name)
		}

		if len(field.Type) == 0 {
			return model.Error(model.InvalidSchema, field.Name)
		}

		if field.Type != repo.StringDataType &&
			field.Type != repo.NumberDataType &&
			field.Type != repo.BoolDataType {
			return model.Error(model.InvalidSchema, field.Name)
		}
	}

	return nil
}
