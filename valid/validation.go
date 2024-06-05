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

func DataType(dataType repo.DataType) error {
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

		if err := DataType(field.Type); err != nil {
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

func Object(obj repo.Object, schema []repo.Field) error {
	fieldMap := make(map[string]repo.Field)
	for _, field := range schema {
		fieldMap[field.Name] = field
	}

	for name, value := range obj {
		field, found := fieldMap[name]
		if !found {
			return model.Error(model.UnknownField, name)
		}

		if !MatchesDataType(value, field.Type) {
			return model.Error(model.InvalidField, name)
		}
	}

	for _, field := range schema {
		if !field.Required {
			continue
		}

		if _, found := obj[field.Name]; !found {
			return model.Error(model.MissingField, field.Name)
		}
	}

	return nil
}

func MatchesDataType(value any, dataType repo.DataType) bool {
	switch dataType {
	case repo.StringDataType:
		_, ok := value.(string)
		return ok
	case repo.NumberDataType:
		_, ok := value.(float64)
		return ok
	case repo.BoolDataType:
		_, ok := value.(bool)
		return ok
	}

	return false
}
