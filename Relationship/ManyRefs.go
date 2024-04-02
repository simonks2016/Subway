package Relationship

import (
	errors2 "errors"
	"fmt"
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/DataAdapter"
	"github.com/simonks2016/Subway/errors"
	"reflect"
)

type ManyReferencesInterface[dataModel any] interface {
	Query() ([]*dataModel, error)
	Remove(...string) error
	Add(...string) error
	DeleteAll() error
	Has(string) (bool, error)
	GetCount() (int, error)
	Output() string
	Rebuild(string, *Core.OperationLib) *ManyReferencesInterface[dataModel]
	New(*Core.OperationLib)
}

type ManyRefs[dataModel any] struct {
	dataModel    []dataModel
	keyName      string
	operationLib *Core.OperationLib
	//ManyReferencesInterface[dataModel]
}

func (this *ManyRefs[dataModel]) New(lib *Core.OperationLib) {
	this.operationLib = lib
}

func (this *ManyRefs[dataModel]) Query() ([]*dataModel, error) {

	if len(this.keyName) <= 0 {
		return nil, errors.ErrMissingTheKeyName
	}

	if this.operationLib == nil || this.operationLib.Fuel == nil {
		return nil, errors.ErrNotSetSubway
	}
	var response []*dataModel

	members, err := this.operationLib.SMembers(this.keyName)
	if err != nil {
		return nil, err
	}

	for _, member := range members {

		var i3 dataModel
		var pTypeName = reflect.TypeOf(i3).Name()
		//if the input parameter is ptr
		if reflect.TypeOf(i3).Kind() == reflect.Ptr {
			//
			pTypeName = reflect.TypeOf(i3).Elem().Name()
		}
		//new document id
		docId := Core.NewDocumentId(pTypeName, member)
		//get data
		err, i2 := this.operationLib.GetByte(docId)
		if err != nil {
			if !errors2.Is(err, errors.ErrNil) {
				return nil, err
			} else {
				continue
			}
		}
		//new dsa
		dsa := DataAdapter.NewDataAdapter[dataModel]("", &i3)
		//un marshal the agreement
		if err = dsa.UnMarshal(i2); err != nil {
			return nil, err
		}
		//append
		response = append(response, dsa.Data)
	}
	if len(response) <= 0 {
		return nil, errors.ErrNil
	}

	return response, nil
}

func (this *ManyRefs[dataModel]) Remove(dataId ...interface{}) error {

	if len(this.keyName) <= 0 {
		return errors.ErrMissingTheKeyName
	}

	if this.operationLib == nil || this.operationLib.Fuel == nil {
		return errors.ErrNotSetSubway
	}
	return this.operationLib.SRemove(this.keyName, dataId...)
}

func (this *ManyRefs[dataModel]) GetCount() (int, error) {

	if len(this.keyName) <= 0 {
		return 0, errors.ErrMissingTheKeyName
	}

	if this.operationLib == nil || this.operationLib.Fuel == nil {
		return 0, errors.ErrNotSetSubway
	}

	return this.operationLib.SCard(this.keyName)
}

func (this *ManyRefs[dataModel]) Has(dataId string) (bool, error) {
	if len(this.keyName) <= 0 {
		return false, errors.ErrMissingTheKeyName
	}

	if this.operationLib == nil || this.operationLib.Fuel == nil {
		return false, errors.ErrNotSetSubway
	}
	return this.operationLib.SIsMember(this.keyName, dataId)
}

func (this *ManyRefs[dataModel]) Add(dataIds ...string) error {
	if len(this.keyName) <= 0 {
		return errors.ErrMissingTheKeyName
	}

	if this.operationLib == nil || this.operationLib.Fuel == nil {
		return errors.ErrNotSetSubway
	}

	var args []interface{}

	for _, id := range dataIds {
		args = append(args, id)
	}
	return this.operationLib.SAdd(this.keyName, args...)
}

func (this *ManyRefs[dataModel]) Output() string {

	if len(this.keyName) <= 0 {
		panic("You didn't set many refs")
	}
	return this.keyName
}
func (this *ManyRefs[dataModel]) Rebuild(keyName string, ol *Core.OperationLib) *ManyRefs[dataModel] {
	return &ManyRefs[dataModel]{
		keyName:      keyName,
		operationLib: ol,
	}
}

func (this *ManyRefs[dataModel]) HasAndCall(d1 []string, call func(string2 string) error) error {
	if len(this.keyName) <= 0 {
		return errors.ErrMissingTheKeyName
	}

	if this.operationLib == nil || this.operationLib.Fuel == nil {
		return errors.ErrNotSetSubway
	}
	members, err := this.operationLib.SMembers(this.keyName)
	if err != nil {
		return err
	}

	for _, member := range members {
		for _, s := range d1 {
			if s == member {
				if err = call(member); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func NewManyRefs[currentDataModel, targetDataModel any](CurrentDataId, fieldName string) *ManyRefs[targetDataModel] {
	var current currentDataModel
	return &ManyRefs[targetDataModel]{
		keyName: NewKeyId(Core.GetViewModelName(current), CurrentDataId, fieldName),
	}
}

func NewKeyId(ViewModelName, ViewModelId, FieldName string) string {
	return fmt.Sprintf("Many-Refs-%s-%s-%s", ViewModelName, ViewModelId, FieldName)

}
