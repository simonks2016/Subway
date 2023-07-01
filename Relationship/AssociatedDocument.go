package Relationship

import (
	"encoding/json"
	"github.com/simonks2016/Subway/Basic"
)

type BasicAssociatedDocument[ViewModel any] struct {
	TargetDataId        string
	TargetViewModelName string
	TargetDocumentId    string
	OL                  Basic.OperationLib
}

func (this BasicAssociatedDocument[ViewModel]) Query() (*ViewModel, error) {

	var vm ViewModel
	err, s := this.OL.GetString(this.TargetDocumentId)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(s), &vm); err != nil {
		return nil, err
	}

	return &vm, nil
}

func NewAssociatedDocument[tvm any](TargetDataId string, TargetViewModelName string) BasicAssociatedDocument[tvm] {

	return BasicAssociatedDocument[tvm]{
		TargetDataId:        TargetDataId,
		TargetViewModelName: TargetViewModelName,
		TargetDocumentId:    Basic.NewDocumentId(TargetViewModelName, TargetDataId),
	}
}
