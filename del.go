package subway

import (
	"errors"
	"github.com/simonks2016/Subway/Core"
	"reflect"
)

func Del[VM any](id ...string) error {

	if Subway == nil {
		return errors.New("you have not set up Subway")
	}
	var v VM
	var dataIds []interface{}
	var ty = reflect.TypeOf(v)

	ol := Core.OperationLib{
		Fuel: Subway.pool,
	}
	for _, s := range id {
		dataIds = append(dataIds, Core.NewDocumentId(ty.Name(), s))
	}
	return ol.BatchDelete(dataIds...)
}
