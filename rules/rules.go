package rules

import "fmt"

type Field string

func (f Field) String() string {
	return string(f)
}

type BaseType int

const (
	Boolean BaseType = iota
	Integer
	String
	Float
	Object
)

var baseTypeNames = map[BaseType]string{
	Boolean: "boolean",
	Integer: "integer",
	String:  "string",
	Float:   "float",
	Object:  "object",
}

func (bt BaseType) String() string {
	return baseTypeNames[bt]
}

type FieldType struct {
	baseType     BaseType
	isArray      bool
	enumValues   []any
	objectFields map[Field]FieldType
}

func (ft FieldType) BaseType() BaseType {
	return ft.baseType
}

func (ft FieldType) IsArray() bool {
	return ft.isArray
}

func (ft *FieldType) SetIsArray(isArray bool) {
	ft.isArray = isArray
}

func (ft FieldType) EnumValues() []any {
	return ft.enumValues
}

func (ft *FieldType) SetEnumValues(enumValues ...any) error {
	if ft.baseType == Object {
		return fmt.Errorf("enums of objects are not supported")
	}
	if err := validateEnumValues(ft.baseType, enumValues); err != nil {
		return err
	}
	ft.enumValues = enumValues
	return nil
}

func (ft *FieldType) AddEnumValue(enumValue any) error {
	if ft.enumValues == nil {
		return fmt.Errorf("enum values are not set")
	}
	ft.enumValues = append(ft.enumValues, enumValue)
	return nil
}

func (ft FieldType) ObjectFields() map[Field]FieldType {
	return ft.objectFields
}

func (ft FieldType) String() string {
	str := ft.baseType.String()
	if ft.enumValues != nil {
		str = fmt.Sprintf("enum<%s>", str)
	}
	if ft.isArray {
		str = fmt.Sprintf("array<%s>", str)
	}
	return str
}

func (ft FieldType) Field() Field {
	return Field(ft.String())
}

func BooleanType() FieldType {
	return FieldType{baseType: Boolean}
}

func IntegerType() FieldType {
	return FieldType{baseType: Integer}
}

func StringType() FieldType {
	return FieldType{baseType: String}
}

func FloatType() FieldType {
	return FieldType{baseType: Float}
}

func ObjectType(objectFields map[Field]FieldType) (FieldType, error) {
	if objectFields == nil {
		return FieldType{}, fmt.Errorf("object fields cannot be nil")
	}

	if len(objectFields) == 0 {
		return FieldType{}, fmt.Errorf("object fields cannot be empty")
	}

	for field, fieldType := range objectFields {
		if fieldType.baseType == Object {
			return FieldType{}, fmt.Errorf("object field '%s' cannot be of type object (nested objects are not supported)", field)
		}
	}

	return FieldType{
		baseType:     Object,
		objectFields: objectFields,
	}, nil
}

func ArrayOf(fieldType FieldType) FieldType {
	fieldType.SetIsArray(true)
	return fieldType
}

func EnumOf(fieldType FieldType, enumValues ...any) (FieldType, error) {
	if enumValues == nil {
		return FieldType{}, fmt.Errorf("enum values cannot be nil")
	}

	if len(enumValues) == 0 {
		return FieldType{}, fmt.Errorf("enum values cannot be empty")
	}

	return fieldType, fieldType.SetEnumValues(enumValues...)
}
