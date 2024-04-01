package subway

import (
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/DataAdapter"
	"github.com/simonks2016/Subway/Filter"
	errors2 "github.com/simonks2016/Subway/errors"
)

func BulkInsert[vm any](data map[string]*vm) error {

	if Subway == nil {
		return errors2.ErrNotSetSubway
	}
	var lib = Core.OperationLib{
		Fuel: Subway.pool,
	}
	var insertData = make(map[string][]byte)
	//
	//key: key Name
	//value: field:value
	var insertFilterField = make(map[string]map[any][]string)
	var insertSortField = make(map[string]map[any]any)

	for dataId, datum := range data {

		var da = DataAdapter.NewDataAdapter[vm]("", datum)
		//
		da.CallbackDiscoveryFilterField = func(key string, field any) error {

			var val map[any][]string
			var exist bool

			val, exist = insertFilterField[key]
			if !exist {
				//set the new field in map
				insertFilterField[key] = map[any][]string{
					field: {dataId},
				}
				return nil
			}

			v1, exist := val[field]
			if !exist {
				//
				val[field] = []string{dataId}
			} else {
				//
				v1 = append(v1, dataId)
				//
				val[field] = v1
			}

			return nil
		}
		//
		da.CallbackDiscoverySortField = func(key, fieldName string, value float64) error {

			val, exist := insertSortField[key]
			if !exist {
				insertSortField[key] = map[any]any{
					dataId: value,
				}
				return nil
			}

			if _, exist = val[dataId]; !exist {
				//if not exist
				val[dataId] = value
				//
				insertSortField[key] = val
			}
			return nil
		}
		//marshal data
		marshal, err := da.Marshal()
		if err != nil {
			return err
		}
		if _, exist := insertData[da.DocId]; !exist {
			insertData[da.DocId] = marshal
		}
	}

	for keyName, m := range insertFilterField {
		var f1 = make(map[any]any)
		//
		for field, values := range m {
			//Check if the data exists
			existsHashMap, err := lib.ExistsHashMap(keyName, field)
			if err != nil {
				return err
			}
			//if exists...
			if existsHashMap {
				//get the old data
				docIds, err := lib.GetHashMap(keyName, field)
				if err != nil {
					return err
				}
				if docIds != nil {
					//split string
					d2 := Filter.SplitString(Filter.Uint82String(docIds))
					// Copy existing id to new list sales
					values = append(values, d2...)
				}
			}
			f1[field] = Filter.Merge2String(values)
		}
		//set the hash map
		err := lib.MSetHashMap(keyName, f1)
		if err != nil {
			return err
		}
	}

	for keyName, sortData := range insertSortField {
		err := lib.MSetHashMap(keyName, sortData)
		if err != nil {
			return err
		}
	}

	for keyName, bytes := range insertData {
		//set the document
		err := lib.Set(keyName, bytes)
		if err != nil {
			return err
		}
	}

	return nil
}
