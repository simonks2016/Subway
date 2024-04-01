package Relationship

import (
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/DataAdapter"
	errors2 "github.com/simonks2016/Subway/errors"
	"reflect"
)

type ReferencesInterface[dataModel any] interface {
	Query() (*dataModel, error)
	Edit(string) error
	Delete() error
	Output() string
	Rebuild(string)
}

type Ref[dataModel any] struct {
	keyName   string
	dataModel    *dataModel
	operationLib *Core.OperationLib
	//ReferencesInterface[dataModel]
}

func (m *Ref[dataModel]) Query() (*dataModel, error) {

	if len(m.keyName) <= 0 {
		return nil, errors2.ErrMissingTheKeyName
	}
	if m.operationLib == nil || m.operationLib.Fuel == nil {
		return nil, errors2.ErrNotSetSubway
	}
	err, s := m.operationLib.GetByte(m.keyName)
	if err != nil {
		return nil, err
	}
	var c dataModel
	//new dsa
	dsa := DataAdapter.NewDataAdapter[dataModel]("", &c)
	//unmarshal the agreement
	if err = dsa.UnMarshal(s); err != nil {
		return nil, err
	}
	return dsa.Data, nil
}
func (m *Ref[dataModel]) Edit(dataId string) {
	m.keyName = Core.NewDocumentId(reflect.TypeOf(m.dataModel).Name(), dataId)
}
func (m *Ref[dataModel]) Delete() {
	m.keyName = ""
}
func (m *Ref[dataModel]) Output() string {
	return m.keyName
}

func (m *Ref[dataModel]) Rebuild(keyName string, ol *Core.OperationLib) *Ref[dataModel] {

	m = &Ref[dataModel]{
		keyName:      keyName,
		operationLib: ol,
	}

	return m
}

func NewRef[dataModel any](dataId string) *Ref[dataModel] {

	return &Ref[dataModel]{keyName: dataId}
}
