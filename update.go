package subway

import (
	"errors"
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/DataAdapter"
	"github.com/simonks2016/Subway/Filter"
	"reflect"
	"strings"
)

func Update[VM any](data VM, docId string) error {

	if Subway == nil {
		return errors.New("you have not set up Subway")
	}

	var ol = Core.OperationLib{
		Fuel: Subway.pool,
	}
	var ViewModelName = reflect.TypeOf(data).Name()
	//if is the ptr
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		ViewModelName = reflect.TypeOf(data).Elem().Name()
	}

	var da = DataAdapter.NewDataAdapter[VM](docId, &data)

	da.CallbackSort = func(key, fieldName string, value float64) error {
		//set hash map
		err := ol.SetHashMap(key, docId, value)
		if err != nil {
			return err
		}
		return nil
	}
	da.CallbackFilter = func(key string, field any) error {

		var d1 = []string{docId}

		exist, err := ol.Exist(key)
		if err != nil {
			return err
		}

		if exist {
			docIds, err := ol.GetHashMap(key, field)
			if err != nil {
				return err
			}

			if docIds != nil {
				//split string
				d2 := Filter.SplitString(Filter.Uint82String(docIds))
				for _, s := range d2 {
					if strings.Compare(s, docId) == 0 {
						return nil
					}
				}
				d1 = append(d1, d2...)
			}
		}
		//set hash map
		err = ol.SetHashMap(key, field, Filter.Merge2String(d1))
		if err != nil {
			return err
		}
		return nil
	}

	//marshal the data
	marshal, err := da.Marshal()
	if err != nil {
		return err
	}

	err = ol.SetString(Core.NewDocumentId(ViewModelName, docId), string(marshal))
	if err != nil {
		return err
	}
	return nil
}
