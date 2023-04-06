package validator

import (
	"fmt"
	"reflect"
)

type validateLen struct {
	len int
}

func (o *validateLen) validate(v reflect.Value) ValidationError {
	switch v.Type().Kind() {
	case reflect.String:
		stringValue := v.String()
		if len([]rune(stringValue)) != o.len {
			return ValidationError{fmt.Errorf("should has fixed length %d", o.len)}
		}
	default:
		return ValidationError{fmt.Errorf("validation parameter 'len' isn't defined for %s", v.Type().Kind())}
	}
	return ValidationError{}
}
