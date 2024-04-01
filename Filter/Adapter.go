package Filter

import (
	"fmt"
	"reflect"
	"strconv"
)

func Convert[FType FieldType](input any) any {

	var a1 FType

	if input == nil {
		return 0
	}
	var t1, t2 = reflect.TypeOf(a1), reflect.TypeOf(input)

	if t1.Kind() != t2.Kind() {

		switch t2.Kind() {
		case reflect.Float64, reflect.Float32:
			return float2Any(t1, input.(float64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int2Any(t1, input.(int))
		case reflect.String:
			return string2Any(t1, input.(string))
		}
		return nil

	} else {
		return input
	}
}

func float2Any(target reflect.Type, current float64) any {

	switch target.Kind() {
	case reflect.Int:
		return int(current)
	case reflect.Int64:
		return int64(current)
	case reflect.Int8:
		return int8(current)
	case reflect.Float32:
		return float32(current)
	case reflect.Float64:
		return current
	case reflect.String:
		return fmt.Sprintf("%.2f", current)
	}
	return nil
}

func int2Any(target reflect.Type, current int) any {
	switch target.Kind() {
	case reflect.Int:
		return current
	case reflect.Int64:
		return int64(current)
	case reflect.Int8:
		return int8(current)
	case reflect.Float32:
		return float32(current)
	case reflect.Float64:
		return float64(current)

	case reflect.String:
		return fmt.Sprintf("%d", current)
	}
	return nil
}

func string2Any(target reflect.Type, current string) any {
	switch target.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int8, reflect.Int64:

		parseInt, err := strconv.ParseInt(current, 0, 64)
		if err != nil {
			return nil
		}

		switch target.Kind() {
		case reflect.Int:
			return int(parseInt)
		case reflect.Int8:
			return int8(parseInt)
		case reflect.Int64:
			return parseInt
		}
		return parseInt
	case reflect.Float32, reflect.Float64:
		parseFloat, err := strconv.ParseFloat(current, 64)
		if err != nil {
			return nil
		}

		switch target.Kind() {
		case reflect.Float32:
			return float32(parseFloat)
		default:
			return parseFloat
		}
	case reflect.String:
		return current
	}
	return nil
}

func Uint82String(v1 any) string {

	switch v1.(type) {
	case []uint8:
		return string(v1.([]uint8))
	case string:
		return v1.(string)
	default:
		return fmt.Sprintf("%v", v1)
	}

}
