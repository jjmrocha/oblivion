package valid

import (
	"regexp"

	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/repo"
)

const (
	_BucketNameRegExp = "^[a-zA-Z][a-zA-Z0-9_]*[a-zA-Z0-9]$"
	_FieldNameRegExp  = "^[a-zA-Z][a-zA-Z0-9_]*[a-zA-Z0-9]$"
	_KeyRegExp        = "^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$"
)

var (
	bucketNameRegExp = regexp.MustCompile(_BucketNameRegExp)
	fieldNameRegExp  = regexp.MustCompile(_FieldNameRegExp)
	keyRegExp        = regexp.MustCompile(_KeyRegExp)
)

func BucketName(name string) error {
	if len(name) == 0 || len(name) > 30 {
		return model.Error(model.InvalidBucketName, name)
	}

	matched := bucketNameRegExp.MatchString(name)

	if !matched {
		return model.Error(model.InvalidBucketName, name)
	}

	return nil
}

func FieldName(name string) error {
	if len(name) == 0 || len(name) > 30 {
		return model.Error(model.InvalidFieldName, name)
	}

	matched := fieldNameRegExp.MatchString(name)

	if !matched {
		return model.Error(model.InvalidFieldName, name)
	}

	return nil
}

func FieldDataType(dataType repo.DataType) error {
	if len(dataType) == 0 {
		return model.Error(model.InvalidFieldType, dataType)
	}

	if dataType != repo.StringDataType &&
		dataType != repo.NumberDataType &&
		dataType != repo.BoolDataType {
		return model.Error(model.InvalidFieldType, dataType)
	}

	return nil
}

func Schema(schema []repo.Field) error {
	if len(schema) == 0 {
		return model.Error(model.SchemaMissing)
	}

	for _, field := range schema {
		if err := FieldName(field.Name); err != nil {
			return err
		}

		if err := FieldDataType(field.Type); err != nil {
			return err
		}
	}

	return nil
}

func Key(value string) error {
	if len(value) == 0 || len(value) > 50 {
		return model.Error(model.InvalidKey, value)
	}

	matched := keyRegExp.MatchString(value)

	if !matched {
		return model.Error(model.InvalidKey, value)
	}

	return nil
}
