package ViewModel

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/DataAdapter"
	errors2 "github.com/simonks2016/Subway/errors"
	"reflect"
	"strings"
)

type ModelOperation[ViewModel comparable] interface {
	GetDataId() string
	Update() error
	Remove() error
	Expire(int64) error
	Read(id string) (*ViewModel, error)
	BatchRead(...string) (error, []*ViewModel)
	Exist(string) (bool, error)
	Import([]byte) (*ViewModel, error)
	Export() ([]byte, error)
	ViewModelName() string
	Set(string, any)
	IsValidViewModel(any) bool
	HasField(string) bool
}

type BasicModelOperation[ViewModel any] struct {
	BasicOperationLib Core.OperationLib
	dsa               *DataAdapter.DataAdapter[ViewModel]
	vm                *ViewModel
}

func (this BasicModelOperation[ViewModel]) IsValidFieldName(fieldName string, targetKind ...reflect.Kind) error {

	exist, FieldType := this.HasField(fieldName)
	if exist == false {
		return errors.New("the related field does not exist")
	}
	var Types []string
	//if related field is not target type
	for _, kind := range targetKind {
		if FieldType == kind {
			return nil
		}
		Types = append(Types, kind.String())
	}
	return errors.New("the related field type is not (" + strings.Join(Types, ",") + ")")

}

func (this BasicModelOperation[ViewModel]) IsValidViewModel(vm any) bool {

	return reflect.TypeOf(vm).Kind() == reflect.Struct
}

func (this BasicModelOperation[ViewModel]) GetValue(fieldName string) string {

	var fieldValue = reflect.ValueOf(this.vm)
	var fieldType = reflect.TypeOf(*this.vm)

	for i := 0; i < fieldType.NumField(); i++ {
		if strings.Compare(strings.ToLower(fieldType.Field(i).Name), strings.ToLower(fieldName)) == 0 {
			return fieldValue.Elem().Field(i).String()
		}
	}
	return ""
}

func (this BasicModelOperation[ViewModel]) GetDataId() string {

	//fmt.Println(this.vm)
	id := this.GetValue("Id")
	//panic
	if len(id) <= 0 {
		panic("The structure is missing the ID primary key")
	}
	return id
}

func (this BasicModelOperation[ViewModel]) Update() error {
	//new doc ID
	var docId = Core.NewDocumentId(this.ViewModelName(), this.GetDataId())
	//copy to dsa
	this.dsa.DocId = docId
	//marshal dsa
	marshal, err := this.dsa.Marshal()
	if err != nil {
		return err
	}
	return this.BasicOperationLib.SetString(docId, string(marshal))
}

func (this BasicModelOperation[ViewModel]) Remove() error {
	if len(this.dsa.DocId) <= 0 {
		return errors.New("missing document Id")
	}

	return this.BasicOperationLib.Delete(this.dsa.DocId)
}

func (this BasicModelOperation[ViewModel]) Expire(expire int64) error {

	if len(this.dsa.DocId) <= 0 {
		return errors.New("missing document Id")
	}

	return this.BasicOperationLib.Expire(this.dsa.DocId, expire)
}

func (this BasicModelOperation[ViewModel]) Exist(id string) (bool, error) {

	var docId = Core.NewDocumentId(this.ViewModelName(), id)
	return this.BasicOperationLib.Exist(docId)
}

/*
read data from redis
we need the data id ,then returning the view model and error
*/
func (this BasicModelOperation[ViewModel]) Read(id string) (*ViewModel, error) {

	var docId = Core.NewDocumentId(this.ViewModelName(), id)

	err, s := this.BasicOperationLib.GetByte(docId)
	if err != nil {
		return nil, err
	}
	if !json.Valid(s) || len(s) <= 0 {
		return nil, errors2.ErrNotJson
	}

	err = this.dsa.UnMarshal(s)
	if err != nil {
		return nil, err
	}
	return this.dsa.Data, nil
}
func (this BasicModelOperation[ViewModel]) ViewModelName() string {

	var v ViewModel
	return reflect.TypeOf(v).Name()
}
func (this BasicModelOperation[ViewModel]) BatchRead(ids ...string) (error, []*ViewModel) {

	var args []interface{}
	var response []*ViewModel

	for _, id := range ids {
		args = append(args, Core.NewDocumentId(this.ViewModelName(), id))
	}
	//batch get
	if err, result := this.BasicOperationLib.BatchGetStrings(args...); err != nil {
		return err, nil
	} else {
		for _, s := range result {
			if len(s) <= 0 || !json.Valid([]byte(s)) {
				continue
			}
			var v1 ViewModel
			var v = DataAdapter.NewDataAdapter[ViewModel]("", &v1)
			//UnMarshal
			err = v.UnMarshal([]byte(s))
			if err != nil {
				return err, nil
			}
			response = append(response, v.Data)
		}
		return nil, response
	}
}

func (b *BasicModelOperation[ViewModel]) SetValue(fieldName string, value any) {
	var vm = b.vm

	values := reflect.ValueOf(vm)
	T := reflect.TypeOf(*vm)

	for i := 0; i < T.NumField(); i++ {
		if strings.Compare(strings.ToLower(T.Field(i).Name), strings.ToLower(fieldName)) == 0 {
			values.Elem().Field(i).Set(reflect.ValueOf(value))
		}
	}
}

func (this BasicModelOperation[ViewModel]) HasField(fieldName string) (bool, reflect.Kind) {

	T := reflect.TypeOf(*this.vm)

	for i := 0; i < T.NumField(); i++ {
		if strings.Compare(T.Field(i).Name, fieldName) == 0 {
			return true, T.Field(i).Type.Kind()
		}
	}
	return false, reflect.String
}

func (this BasicModelOperation[ViewModel]) Set(fieldName string, value any) {
	this.SetValue(fieldName, value)
}

func (this BasicModelOperation[ViewModel]) Import(data []byte) (*ViewModel, error) {

	err := this.dsa.UnMarshal(data)
	if err != nil {
		return nil, err
	}
	return this.dsa.Data, nil
}

func (this BasicModelOperation[ViewModel]) Export() ([]byte, error) {
	return this.dsa.Marshal()
}

func NewBasicModelOperation[ViewModel any](redis *redis.Pool, vm *ViewModel) BasicModelOperation[ViewModel] {
	return BasicModelOperation[ViewModel]{
		BasicOperationLib: Core.OperationLib{
			Fuel: redis,
		},
		vm:  vm,
		dsa: DataAdapter.NewDataAdapter[ViewModel]("", vm),
	}
}
