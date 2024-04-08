package Relationship

import (
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/DataAdapter"
	errors2 "github.com/simonks2016/Subway/errors"
)

type ReferencesInterface[dataModel any] interface {
	Query() (*dataModel, error)
	Edit(string) error
	Delete() error
	Output() string
	Rebuild(string, *Core.OperationLib)
	New(lib *Core.OperationLib)
}

type Ref[dataModel any] struct {
	keyName string
	//dataModel    *dataModel
	operationLib *Core.OperationLib
	dataId       string
	//ReferencesInterface[dataModel]
}

func (m *Ref[dataModel]) New(lib *Core.OperationLib) {
	m.operationLib = lib
}

func (m *Ref[dataModel]) Query() (*dataModel, error) {

	if m == nil {
		return nil, nil
	}
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
	var d dataModel
	m.keyName = Core.NewDocumentId(Core.GetViewModelName(d), dataId)
}
func (m *Ref[dataModel]) Delete() {
	m.keyName = ""
}
func (m *Ref[dataModel]) Output() string {
	return m.dataId
}

func (m *Ref[dataModel]) Rebuild(dataId string, ol *Core.OperationLib) *Ref[dataModel] {

	var d dataModel
	return &Ref[dataModel]{
		keyName:      Core.NewDocumentId(Core.GetViewModelName(d), dataId),
		operationLib: ol,
		dataId:       dataId,
	}
}

func NewRef[dataModel any](dataId string) *Ref[dataModel] {
	var d dataModel
	return &Ref[dataModel]{keyName: Core.NewDocumentId(Core.GetViewModelName(d), dataId), dataId: dataId}
}
