package reflects

import (
	"errors"
	"reflect"
)

func IsTypeValidStruct(structTarget reflect.Type) (reflect.Type, error) {
	// check if type is a pointer
	if structTarget.Kind() == reflect.Ptr {
		structTarget = structTarget.Elem()
	}

	// check if type is a struct
	if structTarget.Kind() != reflect.Struct {
		return nil, errors.New("type is not a struct")
	}

	return structTarget, nil
}

func IsValueValidStruct(value reflect.Value) (reflect.Value, error) {
	// check if value is a pointer
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// check if value is a struct
	if value.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("value is not a struct")
	}

	return value, nil
}

func GetStructFieldFromType(structTarget reflect.Type, fieldName string) (reflect.StructField, error) {
	structTarget, err := IsTypeValidStruct(structTarget)

	if err != nil {
		return reflect.StructField{}, err
	}

	field, found := structTarget.FieldByName(fieldName)
	if !found {
		return reflect.StructField{}, errors.New("field not found")
	}

	return field, nil
}

func MappingStructFieldFromStructType(structTarget reflect.Type) ([]reflect.StructField, error) {
	structTarget, err := IsTypeValidStruct(structTarget)

	if err != nil {
		return nil, err
	}

	numFields := structTarget.NumField()
	fields := make([]reflect.StructField, numFields)

	for i := 0; i < numFields; i++ {
		field := structTarget.Field(i)
		fields[i] = field
	}

	return fields, nil
}

func MappingStructFieldFromStructInterface(structInterface interface{}) ([]reflect.StructField, error) {
	structTarget := reflect.TypeOf(structInterface)

	return MappingStructFieldFromStructType(structTarget)
}

func MappingStructFieldFromStruct[T any]() ([]reflect.StructField, error) {
	var structType T
	structTarget := reflect.TypeOf(structType)

	return MappingStructFieldFromStructType(structTarget)
}
