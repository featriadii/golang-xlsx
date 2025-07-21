package xlsx

import "errors"

func GetHeaderRecords(rows [][]string) ([]string, [][]string, error) {
	if len(rows) < 1 {
		return nil, nil, errors.New("invalid CSV File: Empty file")
	}

	if len((rows)[0]) < 1 {
		return nil, nil, errors.New("invalid CSV File: Header not found")
	}

	if len(rows) < 2 {
		return nil, nil, errors.New("invalid CSV File: Records not found")
	}

	return (rows)[0], (rows)[1:], nil
}

func MappingHeaderRecordToMap(header []string, records []string) map[string]string {
	result := make(map[string]string, len(header))

	for i := 0; i < len(header); i++ {
		if i < len(records) {
			result[header[i]] = records[i]
			continue
		}

		result[header[i]] = ""
	}

	return result
}
