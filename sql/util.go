package sql

import (
	"reflect"

	"database/sql"
)

// ParseSQLRows will parse the row
func ParseSQLRows(rows *sql.Rows, schema interface{}) ([]interface{}, error) {
	var parsedRows []interface{}

	// Fetch rows
	for rows.Next() {
		newSchema := reflect.New(reflect.ValueOf(schema).Elem().Type()).Interface()

		s := reflect.ValueOf(newSchema).Elem()

		var fields []interface{}
		for i := 0; i < s.NumField(); i++ {
			fields = append(fields, s.Field(i).Addr().Interface())
		}

		err := rows.Scan(fields...)
		if err != nil {
			return nil, err
		}
		parsedRows = append(parsedRows, newSchema)
	}

	return parsedRows, nil
}
