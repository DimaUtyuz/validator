package validator

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v {
		messages = append(messages, err.Err.Error())
	}
	return strings.Join(messages, "\n")
}

func Validate(v any) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	validationErrors := validateStruct(value)

	if validationErrors != nil {
		return validationErrors
	}
	return nil
}

func validateStruct(v reflect.Value) ValidationErrors {
	var validationErrors ValidationErrors

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		validationParams, ok := typeField.Tag.Lookup("validate")

		if !ok {
			continue
		}

		if !typeField.IsExported() {
			validationErrors = addValidationError(validationErrors, ErrValidateForUnexportedFields)
			continue
		}

		valueField := v.Field(i)
		kind := valueField.Type().Kind()
		if kind == reflect.Struct {
			validationErrs := validateStruct(valueField)

			for _, validationError := range validationErrs {
				validationErrors = addFieldValidationError(validationErrors, typeField, valueField, validationError)
			}
		}

		validationRules, err := parseValidationRules(validationParams)
		if err != nil {
			validationErrors = addValidationError(validationErrors, err)
			continue
		}

		if kind == reflect.Slice {
			for j := 0; j < valueField.Len(); j++ {
				element := valueField.Index(j)
				validationErrors = validateField(validationRules, typeField, element, validationErrors)
			}
		} else {
			validationErrors = validateField(validationRules, typeField, valueField, validationErrors)
		}

	}

	return validationErrors
}

func validateField(validationRules []validationRule, typeField reflect.StructField, valueField reflect.Value, validationErrors ValidationErrors) ValidationErrors {
	for _, rule := range validationRules {
		validationError := rule.validate(valueField)
		validationErrors = addFieldValidationError(validationErrors, typeField, valueField, validationError)
	}
	return validationErrors
}

func addValidationError(validationErrors ValidationErrors, err error) ValidationErrors {
	return append(validationErrors, ValidationError{err})
}

func addFieldValidationError(validationErrors ValidationErrors, t reflect.StructField, v reflect.Value, validationError ValidationError) ValidationErrors {
	if validationError.Err != nil {
		validationErrors = addValidationError(
			validationErrors,
			fmt.Errorf("field: %s, value: %v, error: %w", t.Name, v.Interface(), validationError.Err),
		)
	}
	return validationErrors
}
