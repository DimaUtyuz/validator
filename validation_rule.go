package validator

import (
	"reflect"
	"strconv"
	"strings"
)

type validationRule interface {
	validate(v reflect.Value) ValidationError
}

func parseValidationRules(s string) ([]validationRule, error) {
	var validationRules []validationRule

	params := strings.Split(s, ";")
	for _, param := range params {
		parts := strings.Split(param, ":")

		switch parts[0] {
		case "len":
			ln, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, ErrInvalidValidatorSyntax
			}
			validationRules = append(validationRules, &validateLen{ln})

		case "min":
			min, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, ErrInvalidValidatorSyntax
			}
			validationRules = append(validationRules, &validateMin{min})

		case "max":
			max, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, ErrInvalidValidatorSyntax
			}
			validationRules = append(validationRules, &validateMax{max})

		case "in":
			mappedValues := make(map[string]struct{})

			values := strings.TrimSpace(parts[1])
			if values != "" {
				for _, value := range strings.Split(values, ",") {
					mappedValues[value] = struct{}{}
				}
			}

			validationRules = append(validationRules, &validateIn{mappedValues})
		}

	}

	return validationRules, nil
}
