package middleware

import (
	"golang-xlsx/xlsx"
	"net/http"

	responsebuilder "github.com/featriadi/golang-libs/response_builder"
	"github.com/xuri/excelize/v2"
)

func RequestXLSXSheetRows(maxSizeMB int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			maxMemory := int64(maxSizeMB) << 20
			r.ParseMultipartForm(maxMemory)

			file, header, err := r.FormFile("file")
			if err != nil {
				responsebuilder.NewResponseBuilder[any](&w).
					Status(http.StatusBadRequest).
					Message(err.Error()).
					Build()
				return
			}
			defer file.Close()

			f, err := excelize.OpenFile(header.Filename)
			if err != nil {
				responsebuilder.NewResponseBuilder[any](&w).
					Status(http.StatusInternalServerError).
					Message(err.Error()).
					Build()
				return
			}
			defer f.Close()

			sheets := f.GetSheetList()
			requestXLSXSheets := make(map[string][][]string)

			for _, sheet := range sheets {
				rows, err := f.GetRows(sheet)

				if err != nil {
					responsebuilder.NewResponseBuilder[any](&w).
						Status(http.StatusInternalServerError).
						Message(err.Error()).
						Build()
					return
				}

				requestXLSXSheets[sheet] = rows
			}

			ctx := r.Context()
			ctx = SetContextKey(ctx, requestXlsxKey, &requestXLSXSheets)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequestXLSXConvertDataStruct[T any]() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get xlsx mapping
			xlsxMappings := xlsx.MappingXLSXSheetMappingFromStruct[T]()

			// get sheets record
			xlsxSheets, err := GetRequestXLSX(r)
			if err != nil {
				responsebuilder.NewResponseBuilder[any](&w).
					Status(err.GetStatus()).
					Message(err.GetMessage()).
					Build()
				return
			}

			if err := xlsx.MappingToStruct(xlsxMappings, *xlsxSheets); err != nil {
				responsebuilder.NewResponseBuilder[any](&w).
					Status(http.StatusInternalServerError).
					Message(err.Error()).
					Build()
				return
			}

			ctx := r.Context()
			ctx = SetContextKey(ctx, requestXlsxMappingKey, &xlsxMappings)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
