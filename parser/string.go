package parser

import (
	"fmt"
	"strconv"
	"time"
)

func ParseStringToPointerDate(s string) *time.Time {
	if s == "" {
		return nil
	}

	layouts := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"20026-01-02T15:04:05.000Z07:00",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		if result, err := time.Parse(layout, s); err == nil {
			return &result
		}
	}

	return nil
}

func ParseStringToDate(s string) (time.Time, error) {
	date := ParseStringToPointerDate(s)

	if date == nil {
		return time.Time{}, fmt.Errorf("invalid date format")
	}

	return *date, nil
}

func ParseStringToInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](s string) (T, error) {
	if s == "" {
		return 0, nil
	}

	var result T
	var err error
	switch any(result).(type) {
	case int:
		var v int64
		v, err = strconv.ParseInt(s, 10, 0)
		result = T(v)
	case int8:
		var v int64
		v, err = strconv.ParseInt(s, 10, 8)
		result = T(v)
	case int16:
		var v int64
		v, err = strconv.ParseInt(s, 10, 16)
		result = T(v)
	case int32:
		var v int64
		v, err = strconv.ParseInt(s, 10, 32)
		result = T(v)
	case int64:
		var v int64
		v, err = strconv.ParseInt(s, 10, 64)
		result = T(v)
	default:
		return result, fmt.Errorf("unsupported type")
	}

	return result, err
}

func ParseStringToFloat[T ~float32 | ~float64](s string) (T, error) {
	if s == "" {
		return 0, nil
	}

	var result T
	var err error
	switch any(result).(type) {
	case float32:
		var v float64
		v, err = strconv.ParseFloat(s, 32)
		result = T(v)
	case float64:
		var v float64
		v, err = strconv.ParseFloat(s, 64)
		result = T(v)
	default:
		return result, fmt.Errorf("unsupported type")
	}

	return result, err
}

func ParseStringToBool(s string) (bool, error) {
	switch s {
	case "", "0", "N":
		s = "false"
	case "1", "Y":
		s = "true"
	default:
		return false, fmt.Errorf("invalid bool format")
	}

	return strconv.ParseBool(s)
}
