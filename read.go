package subway

import (
	"errors"
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/DataAdapter"
	"github.com/simonks2016/Subway/Filter"
	"github.com/simonks2016/Subway/Sorter"
	errors2 "github.com/simonks2016/Subway/errors"
	"reflect"
)

func Read[vm any](dataId string) (*vm, error) {

	if Subway == nil {
		return nil, errors2.ErrNotSetSubway
	}
	var ol = &Core.OperationLib{
		Fuel: Subway.pool,
	}
	var v vm
	var ViewModelName = reflect.TypeOf(v).Name()
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		ViewModelName = reflect.TypeOf(v).Elem().Name()
	}
	//new document id
	var DocId = Core.NewDocumentId(ViewModelName, dataId)
	//get byte from redis
	err, s := ol.GetByte(DocId)
	if err != nil {
		return nil, err
	}
	//new data adapter
	var da = DataAdapter.NewDataAdapter[vm](DocId, &v)
	da.OperationLib = ol
	//unmarshal json
	err = da.UnMarshal(s)
	if err != nil {
		return nil, err
	}

	return da.Data, nil
}

type SortRequest struct {
	SortFieldName    string `json:"sort_field_name"`
	IsAscendingOrder bool   `json:"is_ascending_order"`
}

type QueryRequest[ConditionType Filter.FieldType] struct {
	FieldName   string        `json:"field_name"`
	Condition   ConditionType `json:"condition_type"`
	FilterFunc  Filter.CompareFunc[ConditionType]
	SortRequest *SortRequest `json:"sort_request"`
}

func List[vm any, CType Filter.FieldType](request *QueryRequest[CType], offset, limit int) ([]*vm, error) {

	if Subway == nil {
		return nil, errors.New("you have not set up Subway")
	}
	var v vm
	var ol = &Core.OperationLib{
		Fuel: Subway.pool,
	}
	var ViewModelName = reflect.TypeOf(v).Name()
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		ViewModelName = reflect.TypeOf(v).Elem().Name()
	}
	//create filter
	filter, err := Filter.CreateFilter[CType](ViewModelName, request.FieldName, request.Condition, ol)
	if err != nil {
		return nil, err
	}
	//Get data with the same conditions
	var FF = request.FilterFunc
	if FF == nil {
		FF = Filter.Equal[CType]
	}
	dataIds := filter.GetSameConditions(request.Condition, FF)
	if request.SortRequest != nil {
		//create Sorter
		sorter := Sorter.CreateSorter(ViewModelName, request.SortRequest.SortFieldName, ol)
		//Sort filtered data
		dataIds = sorter.SortFilteredData(request.SortRequest.IsAscendingOrder, dataIds, offset, limit)
	}
	//arg
	var args []interface{}
	//append to args
	for _, id := range dataIds {
		args = append(args, Core.NewDocumentId(ViewModelName, id))
	}

	//Get documents in batches
	err, i := ol.BatchGetStrings(args...)
	if err != nil {
		return nil, err
	}
	//response
	var response []*vm
	//unmarshal
	for _, s := range i {
		var v1 vm
		var da = DataAdapter.NewDataAdapter[vm]("", &v1)
		//
		da.OperationLib = ol
		//unmarshal
		err = da.UnMarshal([]byte(s))
		if err != nil {

			return nil, err
		}
		response = append(response, da.Data)
	}

	return response, nil
}
