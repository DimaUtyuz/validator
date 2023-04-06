package validator

import (
	"fmt"
	"reflect"
)

type validateMax struct {
	max int
}

func (o *validateMax) validate(v reflect.Value) ValidationError {
	switch v.Type().Kind() {
	case reflect.Int:
		intValue := int(v.Int())
		if intValue > o.max {
			return ValidationError{fmt.Errorf("shouldn't be great than %d", o.max)}
		}
	case reflect.String:
		stringValue := v.String()
		if len([]rune(stringValue)) > o.max {
			return ValidationError{fmt.Errorf("shouldn't be longer than %d", o.max),}
		}
	default:
		return ValidationError{fmt.Errorf("validation parameter 'max'' isn't defined for %s", v.Type().Kind())}
	}

	return ValidationError{}
}
