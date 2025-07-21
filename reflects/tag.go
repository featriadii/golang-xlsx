package reflects

import (
	"reflect"
	"strings"
)

func GetJSONTagFromStructField(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	jsonTagName := strings.Split(jsonTag, ",")[0]

	if jsonTagName == "" {
		return strings.ToLower(field.Name)
	}

	return jsonTagName
}

func GetJsonTagFromStructInterface(structInterface interface{}, fieldName string) string {
	structTarget := reflect.TypeOf(structInterface)

	field, err := GetStructFieldFromType(structTarget, fieldName)
	if err != nil {
		return ""
	}

	return GetJSONTagFromStructField(field)
}

func GetJsonTagFromStruct[T any](fieldName string) string {
	var structType T

	return GetJsonTagFromStructInterface(structType, fieldName)
}

func MappingJsonTagFromStruct[T any]() []string {
	var structType T
	structTarget := reflect.TypeOf(structType)

	numFields := structTarget.NumField()
	jsonTags := make([]string, numFields)

	for i := 0; i < numFields; i++ {
		field := structTarget.Field(i)
		jsonTags[i] = GetJSONTagFromStructField(field)
	}

	return jsonTags
}
