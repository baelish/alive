package server

import (
	"reflect"

	"go.uber.org/zap"
)

var logger *zap.Logger

func logStructDetails(v any) []zap.Field {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// If it's a pointer, dereference
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	var fields []zap.Field
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Use JSON tag if present
		if tag := fieldType.Tag.Get("json"); tag != "" {
			name := tag
			if commaIdx := indexComma(tag); commaIdx != -1 {
				name = tag[:commaIdx]
			}
			if name != "" && name != "-" {
				fieldName = name
			}
		}

		fields = append(fields, zap.Any(fieldName, field.Interface()))
	}
	return fields
}

func indexComma(tag string) int {
	for i, l := range tag {
		if l == ',' {
			return i
		}
	}
	return -1
}
