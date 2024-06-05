package valid

import (
	"net/url"
	"regexp"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/model"
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
		return apperror.New(apperror.InvalidBucketName, name)
	}

	matched := bucketNameRegExp.MatchString(name)

	if !matched {
		return apperror.New(apperror.InvalidBucketName, name)
	}

	return nil
}

func FieldName(name string) error {
	if len(name) == 0 || len(name) > 30 {
		return apperror.New(apperror.InvalidFieldName, name)
	}

	matched := fieldNameRegExp.MatchString(name)

	if !matched {
		return apperror.New(apperror.InvalidFieldName, name)
	}

	return nil
}

func DataType(dataType model.DataType) error {
	if len(dataType) == 0 {
		return apperror.New(apperror.InvalidFieldType, dataType)
	}

	if dataType != model.StringDataType &&
		dataType != model.NumberDataType &&
		dataType != model.BoolDataType {
		return apperror.New(apperror.InvalidFieldType, dataType)
	}

	return nil
}

func Schema(schema []model.Field) error {
	if len(schema) == 0 {
		return apperror.New(apperror.SchemaMissing)
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
		return apperror.New(apperror.InvalidKey, value)
	}

	matched := keyRegExp.MatchString(value)

	if !matched {
		return apperror.New(apperror.InvalidKey, value)
	}

	return nil
}

func Object(obj model.Object, schema []model.Field) error {
	fieldMap := toFieldMap(schema)

	for name, value := range obj {
		field, found := fieldMap[name]
		if !found {
			return apperror.New(apperror.UnknownField, name)
		}

		if !MatchesDataType(value, field.Type) {
			return apperror.New(apperror.InvalidField, name)
		}
	}

	for _, field := range schema {
		if !field.Required {
			continue
		}

		if _, found := obj[field.Name]; !found {
			return apperror.New(apperror.MissingField, field.Name)
		}
	}

	return nil
}

func MatchesDataType(value any, dataType model.DataType) bool {
	switch dataType {
	case model.StringDataType:
		_, ok := value.(string)
		return ok
	case model.NumberDataType:
		_, ok := value.(float64)
		return ok
	case model.BoolDataType:
		_, ok := value.(bool)
		return ok
	}

	return false
}

func Criteria(criteria url.Values, schema []model.Field) error {
	fieldMap := toFieldMap(schema)

	for name := range criteria {
		if _, found := fieldMap[name]; !found {
			return apperror.New(apperror.UnknownField, name)
		}
	}

	return nil
}

func toFieldMap(schema []model.Field) map[string]model.Field {
	fieldMap := make(map[string]model.Field)
	for _, field := range schema {
		fieldMap[field.Name] = field
	}

	return fieldMap
}
