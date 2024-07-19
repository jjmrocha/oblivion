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
		return apperror.InvalidBucketName.New(name)
	}

	matched := bucketNameRegExp.MatchString(name)

	if !matched {
		return apperror.InvalidBucketName.New(name)
	}

	return nil
}

func FieldName(name string) error {
	if len(name) == 0 || len(name) > 30 {
		return apperror.InvalidFieldName.New(name)
	}

	matched := fieldNameRegExp.MatchString(name)

	if !matched {
		return apperror.InvalidFieldName.New(name)
	}

	return nil
}

func DataType(dataType model.DataType) error {
	if len(dataType) == 0 {
		return apperror.InvalidFieldType.New(dataType)
	}

	if dataType != model.StringDataType &&
		dataType != model.NumberDataType &&
		dataType != model.BoolDataType {
		return apperror.InvalidFieldType.New(dataType)
	}

	return nil
}

func Schema(schema []model.Field) error {
	if len(schema) == 0 {
		return apperror.SchemaMissing.New()
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
		return apperror.InvalidKey.New(value)
	}

	matched := keyRegExp.MatchString(value)

	if !matched {
		return apperror.InvalidKey.New(value)
	}

	return nil
}

func Object(obj model.Object, schema []model.Field) error {
	fieldMap := toFieldMap(schema)

	for name, value := range obj {
		field, found := fieldMap[name]
		if !found {
			return apperror.UnknownField.New(name)
		}

		if !field.Type.ValidValue(value) {
			return apperror.InvalidField.New(name)
		}
	}

	for _, field := range schema {
		if !field.Required {
			continue
		}

		if _, found := obj[field.Name]; !found {
			return apperror.MissingField.New(field.Name)
		}
	}

	return nil
}

func Criteria(criteria url.Values, schema []model.Field) error {
	fieldMap := toFieldMap(schema)

	for name := range criteria {
		if _, found := fieldMap[name]; !found {
			return apperror.UnknownField.New(name)
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
