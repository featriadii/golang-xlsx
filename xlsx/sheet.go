package xlsx

import (
	"errors"
	"golang-xlsx/reflects"
	"log"
	"reflect"
)

type XLSXSheetMapping struct {
	ModelName  string
	SheetName  string
	TargetPtr  any
	StructType reflect.Type
}

func MappingXLSXSheetMappingFromStruct[T any]() []XLSXSheetMapping {
	sheetModels, err := reflects.MappingStructFieldFromStruct[T]()

	if err != nil {
		panic(err)
	}

	xlsxMappings := make([]XLSXSheetMapping, len(sheetModels))

	for i, sheetModel := range sheetModels {
		sheetModelTyp := sheetModel.Type
		sliceSheetModel := reflect.SliceOf(sheetModelTyp)

		xlsxMappings[i] = XLSXSheetMapping{
			ModelName:  sheetModel.Name,
			SheetName:  reflects.GetJSONTagFromStructField(sheetModel),
			TargetPtr:  reflect.New(sliceSheetModel).Interface(),
			StructType: sheetModelTyp,
		}
	}

	return xlsxMappings
}

func MappingToStruct(xlsxMappings []XLSXSheetMapping, xlsxSheets map[string][][]string) error {
	// loop each sheet
	for _, xlsxMapping := range xlsxMappings {
		sheetName := xlsxMapping.SheetName
		structType := xlsxMapping.StructType
		targetPtr := xlsxMapping.TargetPtr

		// get sheet
		sheet := xlsxSheets[sheetName]

		if sheet == nil {
			return errors.New("sheet not found")
		}

		// get header and records
		header, records, err := GetHeaderRecords(sheet)
		if err != nil {
			return err
		}

		// check if targetPtr is a pointer to a slice of structs
		targetVal := reflect.ValueOf(targetPtr)
		if targetVal.Kind() != reflect.Ptr || targetVal.Elem().Kind() != reflect.Slice || targetVal.Elem().Type().Elem().Kind() != reflect.Struct {
			return errors.New("targetPtr must be a pointer to a slice of structs")
		}

		// create slice of struct
		sliceVal := targetVal.Elem()
		log.Println(sheetName)

		// loop each record
		for _, record := range records {
			newInstance := reflect.New(structType).Elem()
			headerRecords := MappingHeaderRecordToMap(header, record)

			// map record to struct fields
			MappingRecordToNewStruct(headerRecords, newInstance)

			sliceVal = reflect.Append(sliceVal, newInstance)
		}

		targetVal.Elem().Set(sliceVal)
	}

	return nil
}

func MappingRecordToNewStruct(headerRecords map[string]string, val reflect.Value) {
	val, err := reflects.IsValueValidStruct(val)

	if err != nil {
		log.Println(err)
		return
	}

	typ := val.Type()
	structFields, err := reflects.MappingStructFieldFromStructType(typ)

	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < len(structFields); i++ {
		field := structFields[i]
		fieldName := field.Name
		fieldValue := val.Field(i)
		jsonTagName := reflects.GetJSONTagFromStructField(field)

		setValue, err := setValue(headerRecords, jsonTagName, fieldName)

		if err != nil {
			log.Println(err)
			return
		}

		if err := reflects.SetFieldValue(fieldValue, setValue); err != nil {
			log.Println(err)
			return
		}
	}
}

func setValue(headerRecords map[string]string, jsonTagName string, fieldName string) (string, error) {
	setValue := headerRecords[jsonTagName]

	return setValue, nil
}
