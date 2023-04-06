package validator

import (
	"fmt"
	"reflect"
	"strings"
)

type validateIn struct {
	values map[string]struct{}
}

func (o *validateIn) showValues() string {
	keys := make([]string, len(o.values))
	i := 0
	for k := range o.values {
		keys[i] = k
		i++
	}
	return fmt.Sprintf("{%s}", strings.Join(keys, ", "))
}

func (o *validateIn) validate(v reflect.Value) ValidationError {
	kind := v.Type().Kind()

	if kind == reflect.Int || kind == reflect.String {
		stringValue := fmt.Sprintf("%v", v.Interface())
		if _, ok := o.values[stringValue]; !ok {
			return ValidationError{fmt.Errorf("should be one of %v", o.showValues())}
		}
		return ValidationError{}
	}

	return ValidationError{fmt.Errorf("validation parameter 'in' isn't defined for %s", v.Type().Kind())}
}
