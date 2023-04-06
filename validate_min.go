package validator

import (
	"fmt"
	"reflect"
)

type validateMin struct {
	min int
}

func (o *validateMin) validate(v reflect.Value) ValidationError {
	switch v.Type().Kind() {
	case reflect.Int:
		intValue := int(v.Int())
		if intValue < o.min {
			return ValidationError{fmt.Errorf("shouldn't be less than %d", o.min)}
		}
	case reflect.String:
		stringValue := v.String()
		if len([]rune(stringValue)) < o.min {
			return ValidationError{fmt.Errorf("shouldn't be shorter than %d", o.min)}
		}
	default:
		return ValidationError{fmt.Errorf("validation parameter 'min' isn't defined for %s", v.Type().Kind())}
	}
	return ValidationError{}
}
