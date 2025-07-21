package middleware

import (
	"context"
	"golang-xlsx/xlsx"
	"net/http"
	"reflect"

	"github.com/featriadi/golang-libs/errors"
)

type contextKey string

func SetContextKey(ctx context.Context, key contextKey, value any) context.Context {
	ctx = context.WithValue(ctx, key, value)

	return ctx
}

const requestXlsxKey contextKey = "requestXLSX"

func GetRequestXLSX(r *http.Request) (*map[string][][]string, *errors.Error) {
	data, ok := r.Context().Value(requestXlsxKey).(*map[string][][]string)
	if !ok {
		return nil, errors.NewError(http.StatusInternalServerError, "Invalid or missing type assestion from context", nil)
	}

	return data, nil
}

const requestXlsxMappingKey contextKey = "requestXLSXMapping"

func GetRequestXLSXData[T any](r *http.Request) (*[]T, *errors.Error) {
	xlsxMapping, ok := r.Context().Value(requestXlsxMappingKey).(*[]xlsx.XLSXSheetMapping)
	if !ok {
		return nil, errors.NewError(http.StatusInternalServerError, "Invalid or missing type assestion from context", nil)
	}

	var structTarget T
	structType := reflect.TypeOf(structTarget)
	structName := structType.Name()

	for _, mapping := range *xlsxMapping {
		if mapping.ModelName == structName {
			dataAddress := mapping.TargetPtr
			dataSlice := reflect.ValueOf(dataAddress)
			dataInterface := dataSlice.Interface().(*[]T)

			return dataInterface, nil
		}
	}

	return nil, errors.NewError(http.StatusInternalServerError, "Invalid or missing type assestion from context", nil)
}
