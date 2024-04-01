package subway

import (
	"errors"
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/DataAdapter"
	"github.com/simonks2016/Subway/Filter"
)

func Del[VM any](ids ...string) error {

	if Subway == nil {
		return errors.New("you have not set up Subway")
	}
	var v VM
	var docIds []interface{}
	var deleteKeys []interface{}
	var deleteHash = make(map[string]string)
	var EditHash = make(map[string][]any)

	var ViewModelName = GetViewModelName(v)

	ol := Core.OperationLib{
		Fuel: Subway.pool,
	}
	for _, s := range ids {
		//delete
		docIds = append(docIds, Core.NewDocumentId(ViewModelName, s))
	}

	err, i := ol.BatchGetStrings(docIds...)
	if err != nil {
		return err
	}

	for _, s := range i {

		var da = DataAdapter.DataAdapter[any]{}
		//simple un marshal
		err, d := da.SimpleUnMarshal([]byte(s))
		if err != nil {
			return err
		}

		for _, s2 := range d.ManyRefs {
			deleteKeys = append(deleteKeys, s2)
		}

		for _, setting := range d.SortFields {
			//add to delete hash / key:key name value: docId
			deleteHash[setting.KeyName] = d.DocId
		}

		for _, setting := range d.FilterField {
			//
			if val, exist := EditHash[setting.KeyName]; !exist {
				EditHash[setting.KeyName] = []any{
					setting.Value,
				}
			} else {
				//append
				val = append(val, setting.Value)
				//
				EditHash[setting.KeyName] = val
			}
		}
	}
	//
	deleteKeys = append(deleteKeys, docIds...)

	//upgrade the filter field
	for keyName, field := range EditHash {

		hashMap, err := ol.MGetHashMap(keyName, field...)
		if err != nil {
			return err
		}

		var data = make(map[any]any)

		for index, a := range hashMap {

			v1 := a.([]uint8)
			s1 := Filter.SplitString(string(v1))
			if len(s1) <= 0 {
				continue
			}
			//newSlice
			newSLice := removeSpecifiedElem[string](s1, ids)
			//new ids
			newIdsString := Filter.Merge2String(newSLice)
			//copy to data
			data[field[index]] = newIdsString
		}

		err = ol.MSetHashMap(keyName, data)
		if err != nil {
			return err
		}
	}

	return ol.BatchDelete(deleteKeys...)
}

func removeSpecifiedElem[T string | int | float64 | float32 | int64 | int8 | int16 | int32](s []T, d []T) []T {

	var newS []T
	for _, s2 := range s {
		var exist = false
	loop2:
		for _, s3 := range d {
			if s2 == s3 {
				exist = true
				break loop2
			}
		}
		if !exist {
			newS = append(newS, s2)
		}
	}
	return newS
}
