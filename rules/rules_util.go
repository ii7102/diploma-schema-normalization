package rules

import "fmt"

func validateEnumValues(baseType BaseType, enumValues []any) (err error) {
	for _, v := range enumValues {
		if v == nil {
			return fmt.Errorf("enum value cannot be nil")
		}
		if !matchesBaseType(baseType, v) {
			return fmt.Errorf("enum value %v is not of %s type", v, baseType.String())
		}
	}
	return nil
}

func matchesBaseType(baseType BaseType, v any) bool {
	switch baseType {
	case Boolean:
		_, ok := v.(bool)
		return ok
	case Integer:
		switch v.(type) {
		case int, int8, int16, int32, int64:
			return true
		}
	case String:
		_, ok := v.(string)
		return ok
	case Float:
		switch v.(type) {
		case float32, float64:
			return true
		}
	}
	return false
}
