package Filter

import (
	"fmt"
	"github.com/simonks2016/Subway/Core"
	errors2 "github.com/simonks2016/Subway/errors"
	"reflect"
	"strings"
)

type FieldType interface {
	~int | ~int8 | ~int64 | ~float32 | ~float64 | ~string
}
type CompareFunc[fieldType FieldType] func(fieldType, fieldType) bool

type FilterField[fieldType FieldType] struct {
	keyName      string
	keyValue     fieldType
	fieldName    string
	DataIds      []string
	_map         map[fieldType][]string
	operationLib *Core.OperationLib
}

func NewKeyId(ViewModelName string, FieldName string) string {
	return fmt.Sprintf("filter-%s-%s", ViewModelName, FieldName)
}

func (this *FilterField[fieldType]) GetSameConditions(Conditions fieldType, Compare CompareFunc[fieldType]) []string {

	if this._map == nil {
		if err := this.getAllData(); err != nil {
			return nil
		}
	}

	var data []string
	//get the
	for key, value := range this._map {
		if Compare(key, Conditions) == true {
			data = append(data, value...)
		}
	}
	return data
}
func (this *FilterField[fieldType]) Get() fieldType {

	return this.keyValue
}
func (this *FilterField[fieldType]) Set(value fieldType, docIds ...string) error {

	//If the existing map does not exist, request redis
	if this._map == nil {
		if err := this.getAllData(); err != nil {
			return err
		}
	}
	//If the element exists in the map
	if val, exists := this._map[value]; exists {
		docIds = append(docIds, val...)
	}
	//redis hash map set
	err := this.operationLib.SetHashMap(this.keyName, value, Merge2String(docIds))
	if err != nil {
		return err
	}
	return nil
}

func RemoveElement(data []string, id string) []string {

	for i := 0; i < len(data); i++ {
		if data[i] == id {
			data = append(data[:i], data[i+1:]...)
			i--
		}
	}
	return data
}

func (this *FilterField[fieldType]) RemoveDocId(docId string) error {

	if this._map == nil {
		if err := this.getAllData(); err != nil {
			return err
		}
	}

	var modifyField = make(map[any]any)

	//loop to map
	for fieldName, item := range this._map {
		for _, s := range item {
			if strings.Compare(s, docId) == 0 {
				//add the modified field
				modifyField[fieldName] = strings.Join(RemoveElement(item, docId), "&")
				continue
			}
		}
	}
	//If the modified field is empty
	if len(modifyField) <= 0 || modifyField == nil {
		return nil
	} else {
		err := this.operationLib.MSetHashMap(this.keyName, modifyField)
		if err != nil {
			return err
		}
	}
	//save on redis
	return nil
}
func (this *FilterField[fieldType]) getAllData() error {

	hashMap, err := this.operationLib.GetALLHashMap(this.keyName)
	if err != nil {
		return err
	}

	var _newMap = make(map[fieldType][]string)
	var num = len(hashMap) / 2
	if len(hashMap)%2 != 0 {
		return errors2.ErrUnableGenerateMap
	}

	for i := 0; i < num; i++ {
		var key, value = hashMap[i].([]uint8), hashMap[i+1].([]uint8)
		k1 := string(key)
		v1 := string(value)
		//convert to key
		k2 := Convert[fieldType](k1)
		//add to map
		_newMap[k2.(fieldType)] = SplitString(v1)
	}

	this._map = _newMap

	return nil
}

func SplitString(ids string) []string {
	return strings.Split(ids, "&&")
}
func Merge2String(ids []string) string {
	return strings.Join(ids, "&&")
}

func (this *FilterField[fieldType]) Output() (string, fieldType) {
	return this.keyName, this.keyValue
}

func (this *FilterField[fieldType]) Rebuild(keyName string, value any, ol *Core.OperationLib) *FilterField[fieldType] {

	if value != nil && !reflect.ValueOf(value).IsZero() {
		v1 := Convert[fieldType](value)
		return &FilterField[fieldType]{
			keyName:      keyName,
			keyValue:     v1.(fieldType),
			operationLib: ol,
		}
	}

	return &FilterField[fieldType]{
		keyName:      keyName,
		operationLib: ol,
	}
}

func NewFilter[fieldType FieldType](dataModelName string, fieldName string, fieldVale fieldType, DataIds ...string) *FilterField[fieldType] {

	return &FilterField[fieldType]{
		keyName:   NewKeyId(dataModelName, fieldName),
		keyValue:  fieldVale,
		fieldName: fieldName,
		DataIds:   DataIds,
	}

}

func CreateFilter[fieldType FieldType](dataModelName, FieldName string, FieldVale fieldType, ol *Core.OperationLib) (*FilterField[fieldType], error) {

	c := &FilterField[fieldType]{
		keyValue:     FieldVale,
		keyName:      NewKeyId(dataModelName, FieldName),
		fieldName:    FieldName,
		_map:         make(map[fieldType][]string),
		operationLib: ol,
	}
	//get map
	if err := c.getAllData(); err != nil {
		return nil, err
	}
	return c, nil
}
