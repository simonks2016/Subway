package Core

import "reflect"

func HandleArray[T any](d1 reflect.Value) []T {

	var r []T
	var d any
	for i := 0; i < d1.NumField(); i++ {
		if d1.Field(i).CanConvert(reflect.TypeOf(d)) {
			f1 := d1.Field(i).Convert(reflect.TypeOf(d))
			reflect.ValueOf(r).Field(i).Set(f1)
		}
	}

	return r
}
