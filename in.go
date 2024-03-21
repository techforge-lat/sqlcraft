package sqlcraft

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
)

func buildIn(value any, index int) (string, []any) {
	if value == nil {
		return "", nil
	}

	builder := bytes.Buffer{}
	builder.WriteString("(")

	var args []interface{}

	// Use reflection to handle different types
	valSlice := reflect.ValueOf(value)
	if valSlice.Kind() == reflect.Slice {
		for i := 0; i < valSlice.Len(); i++ {
			builder.WriteString("$")
			builder.WriteString(strconv.Itoa(index + i + 1))
			builder.WriteString(", ")

			args = append(args, valSlice.Index(i).Interface())
		}

		if valSlice.Len() > 0 {
			builder.Truncate(builder.Len() - 2)
		}

		builder.WriteString(")")

		return builder.String(), args
	}

	str, ok := value.(string)
	if !ok {
		return "", nil
	}

	stringValues := strings.Split(str, ",")
	for i, v := range stringValues {
		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(index + i + 1))
		builder.WriteString(", ")

		args = append(args, v)
	}

	builder.Truncate(builder.Len() - 2)
	builder.WriteString(")")

	return builder.String(), args
}
