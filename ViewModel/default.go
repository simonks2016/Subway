package ViewModel

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/simonks2016/Subway/Basic"
	"github.com/simonks2016/Subway/Relationship"
	"reflect"
	"strings"
)

type ModelOperation[ViewModel any] interface {
	CreateRelationship(string) Relationship.Controllers
	QueryRelationship(string) (Relationship.Controllers, error)
	CreateLine(sortField string,specialKey ...string) Relationship.LineControllers

	LoadRelationship(string) Relationship.Controllers
	LoadLine(string,...string) Relationship.LineControllers

	GetDataId() string

	Update() error
	Remove() error
	Expire(int64) error
	Read(id string) (*ViewModel, error)
	BatchRead(...string) (error, []*ViewModel)
	Exist(string) (bool, error)

	ViewModelName() string
	Set(string, any)

	IsValidViewModel(any) bool
}

type BasicModelOperation[ViewModel any] struct {
	BasicOperationLib Basic.OperationLib
	dsa               *Basic.DSA[ViewModel]
	vm                *ViewModel
}

func (this BasicModelOperation[ViewModel]) IsValidFieldName(fieldName string, targetKind ...reflect.Kind) error {

	exist, FieldType := this.hasField(fieldName)
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

func (this BasicModelOperation[ViewModel]) createCollectionKey(fieldName string) string {

	return this.ViewModelName() + ":" + this.GetDataId() + "-" + fieldName
}

func (this BasicModelOperation[ViewModel]) createLineKey(fieldName string, specialKey ...string) string {

	if len(specialKey)<=0{
		return this.ViewModelName() + "-line-by-" + fieldName
	}
	return this.ViewModelName() + ":"+strings.Join(specialKey,"&&") + "-line-by-" + fieldName
}

func (this BasicModelOperation[VideoModel]) CreateRelationship(fieldName string) Relationship.Controllers {

	err := this.IsValidFieldName(fieldName, reflect.String)
	if err != nil {
		panic(err.Error())
	}
	var key = this.createCollectionKey(fieldName)
	//set key in fieldName
	this.Set(fieldName, key)
	//return
	return Relationship.NewBasicRelationshipControllers(key, this.BasicOperationLib)
}

func (this BasicModelOperation[ViewModel]) QueryRelationship(fieldName string) (Relationship.Controllers, error) {

	err := this.IsValidFieldName(fieldName, reflect.String)
	if err != nil {
		return nil, err
	}
	var key = this.createCollectionKey(fieldName)
	//return
	return Relationship.NewBasicRelationshipControllers(key, this.BasicOperationLib), nil
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
	var docId = Basic.NewDocumentId(this.ViewModelName(), this.GetDataId())
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

	var docId = Basic.NewDocumentId(this.ViewModelName(), id)
	return this.BasicOperationLib.Exist(docId)
}

/*
read data from redis
we need the data id ,then returning the view model and error
*/
func (this BasicModelOperation[ViewModel]) Read(id string) (*ViewModel, error) {

	var docId = Basic.NewDocumentId(this.ViewModelName(), id)

	err, s := this.BasicOperationLib.GetString(docId)
	if err != nil {
		return nil, err
	}

	err = this.dsa.UnMarshal([]byte(s))
	if err != nil {
		return nil, err
	}

	return this.vm, nil
}

/*
this function is create the line from document
the line base on the redis sorted set.

*/

func (this BasicModelOperation[ViewModel]) CreateLine(fieldName string,specialKey ...string) Relationship.LineControllers {

	if err := this.IsValidFieldName(fieldName, reflect.Float64, reflect.Int64, reflect.Int, reflect.Float32, reflect.Int8, reflect.Int32, reflect.Int16); err != nil {
		panic(err.Error())
	}

	var key = this.createLineKey(fieldName,specialKey...)
	this.dsa.AddLine(fieldName)

	return Relationship.NewBasicLineControllers(key, this.BasicOperationLib)
}

func (this BasicModelOperation[ViewModel]) LoadRelationship(fieldName string) Relationship.Controllers {

	if err := this.IsValidFieldName(fieldName, reflect.String); err != nil {
		panic(err.Error())
	}
	if this.dsa.HasRelationship(fieldName) == false {
		return this.CreateRelationship(fieldName)
	}

	var key = this.GetValue(fieldName)
	var key1 = this.createCollectionKey(fieldName)
	if len(key) <= 0 {
		key = key1
		//set value in field
		this.SetValue(fieldName, key1)
	} else {
		if strings.Compare(key1, key) != 0 {
			panic("Collection keys are not as expected")
		}
	}
	return Relationship.NewBasicRelationshipControllers(key, this.BasicOperationLib)
}

func (this BasicModelOperation[ViewModel]) LoadLine(fieldName string,specialKey ...string) Relationship.LineControllers {

	if err := this.IsValidFieldName(fieldName, reflect.Float64, reflect.Int64, reflect.Int, reflect.Float32, reflect.Int8, reflect.Int32, reflect.Int16); err != nil {
		panic(err.Error())
	}

	/*
		if this.dsa.HasLine(fieldName) == false {
			return this.CreateLine(fieldName)
		}*/
	//make key
	var key = this.createLineKey(fieldName,specialKey...)
	//return
	return Relationship.NewBasicLineControllers(key, this.BasicOperationLib)
}

func (this BasicModelOperation[ViewModel]) ViewModelName() string {

	var v ViewModel
	return reflect.TypeOf(v).Name()
}
func (this BasicModelOperation[ViewModel]) BatchRead(ids ...string) (error, []*ViewModel) {

	var args []interface{}
	var response []*ViewModel

	for _, id := range ids {
		args = append(args, Basic.NewDocumentId(this.ViewModelName(), id))
	}
	//batch get
	if err, result := this.BasicOperationLib.BatchGetStrings(args...); err != nil {
		return err, nil
	} else {
		for _, s := range result {
			var v Basic.DSA[ViewModel]
			//UnMarshal
			err = json.Unmarshal([]byte(s), &v)
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

func (this BasicModelOperation[ViewModel]) hasField(fieldName string) (bool, reflect.Kind) {

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

func NewBasicModelOperation[ViewModel any](redis *redis.Pool, vm *ViewModel) BasicModelOperation[ViewModel] {
	return BasicModelOperation[ViewModel]{
		BasicOperationLib: Basic.OperationLib{
			Fuel: redis,
		},
		vm:  vm,
		dsa: Basic.NewDSA[ViewModel]("", vm, nil),
	}
}
