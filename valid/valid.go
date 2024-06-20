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
		return apperror.InvalidBucketName.NewError(name)
	}

	matched := bucketNameRegExp.MatchString(name)

	if !matched {
		return apperror.InvalidBucketName.NewError(name)
	}

	return nil
}

func FieldName(name string) error {
	if len(name) == 0 || len(name) > 30 {
		return apperror.InvalidFieldName.NewError(name)
	}

	matched := fieldNameRegExp.MatchString(name)

	if !matched {
		return apperror.InvalidFieldName.NewError(name)
	}

	return nil
}

func DataType(dataType model.DataType) error {
	if len(dataType) == 0 {
		return apperror.InvalidFieldType.NewError(dataType)
	}

	if dataType != model.StringDataType &&
		dataType != model.NumberDataType &&
		dataType != model.BoolDataType {
		return apperror.InvalidFieldType.NewError(dataType)
	}

	return nil
}

func Schema(schema []model.Field) error {
	if len(schema) == 0 {
		return apperror.SchemaMissing.NewError()
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
		return apperror.InvalidKey.NewError(value)
	}

	matched := keyRegExp.MatchString(value)

	if !matched {
		return apperror.InvalidKey.NewError(value)
	}

	return nil
}

func Object(obj model.Object, schema []model.Field) error {
	fieldMap := toFieldMap(schema)

	for name, value := range obj {
		field, found := fieldMap[name]
		if !found {
			return apperror.UnknownField.NewError(name)
		}

		if !field.Type.ValidValue(value) {
			return apperror.InvalidField.NewError(name)
		}
	}

	for _, field := range schema {
		if !field.Required {
			continue
		}

		if _, found := obj[field.Name]; !found {
			return apperror.MissingField.NewError(field.Name)
		}
	}

	return nil
}

func Criteria(criteria url.Values, schema []model.Field) error {
	fieldMap := toFieldMap(schema)

	for name := range criteria {
		if _, found := fieldMap[name]; !found {
			return apperror.UnknownField.NewError(name)
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
