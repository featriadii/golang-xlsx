package reflects

import (
	"fmt"
	"golang-xlsx/parser"
	"log"
	"reflect"
	"time"
)

var (
	stringType  = reflect.TypeOf("")
	intType     = reflect.TypeOf(int(0))
	int8Type    = reflect.TypeOf(int8(0))
	int16Type   = reflect.TypeOf(int16(0))
	int32Type   = reflect.TypeOf(int32(0))
	int64Type   = reflect.TypeOf(int64(0))
	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))
	boolType    = reflect.TypeOf(false)
	timeType    = reflect.TypeOf(time.Time{})
)

func SetFieldValue(field reflect.Value, setValue string) error {
	if field.Kind() == reflect.Ptr {
		if setValue == "" {
			field.Set(reflect.Zero(field.Type()))
			return nil
		}

		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		field = field.Elem()
	}

	fieldType := field.Type()
	fieldName := fieldType.Name()

	switch fieldType {
	case stringType:
		field.SetString(setValue)
	case intType, int8Type, int16Type, int32Type, int64Type:
		value, err := parser.ParseStringToInt[int64](setValue)
		if err != nil {
			return fmt.Errorf("error parsing field %s: %v", fieldName, err)
		}

		field.SetInt(value)
	case float32Type, float64Type:
		value, err := parser.ParseStringToFloat[float64](setValue)
		if err != nil {
			return fmt.Errorf("error parsing field %s: %v", fieldName, err)
		}

		field.SetFloat(value)
	case boolType:
		value, err := parser.ParseStringToBool(setValue)
		if err != nil {
			log.Println("setValue", setValue)
			return fmt.Errorf("error parsing field %s: %v", fieldName, err)
		}

		field.SetBool(value)
	case timeType:
		parsedDate, err := parser.ParseStringToDate(setValue)
		if err != nil {
			log.Println("setValue", setValue)
			return fmt.Errorf("error parsing field %s: %v", fieldName, err)
		}

		field.Set(reflect.ValueOf(parsedDate))
	default:
		return fmt.Errorf("unsupported field type for field %s", fieldName)
	}

	return nil
}
