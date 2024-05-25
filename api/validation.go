package api

import (
	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/bucket/model/apperror"
)

func checkBucketCreation(name string, schema []model.Field) error {
	if len(name) == 0 {
		return apperror.New(model.InvalidBucketName)
	}

	if len(schema) == 0 {
		return apperror.New(model.SchemaMissing)
	}

	for _, field := range schema {
		if len(field.Name) == 0 {
			return apperror.New(model.InvalidSchema, field.Name)
		}

		if len(field.Type) == 0 {
			return apperror.New(model.InvalidSchema, field.Name)
		}

		if field.Type != model.StringDataType &&
			field.Type != model.NumberDataType &&
			field.Type != model.BoolDataType {
			return apperror.New(model.InvalidSchema, field.Name)
		}
	}

	return nil
}
