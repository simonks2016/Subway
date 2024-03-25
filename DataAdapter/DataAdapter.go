package DataAdapter

import (
	"encoding/json"
	"github.com/mailru/easyjson"
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/define"
	errors2 "github.com/simonks2016/Subway/errors"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const callOutputFunc = "Output"
const callRebuildFunc = "Rebuild"
const callUpgradeFunc = "Upgrade"
const callSetRedis = "SetRedisConn"

type FilterCallback func(key string, field any) error
type SortCallback func(key, fieldName string, value float64) error

type DataAdapter[ViewModel any] struct {
	DocId          string     `json:"doc_id"`
	Data           *ViewModel `json:"-"`
	CreateTime     int64      `json:"create_time"`
	ref            map[string]string
	manyRefs       map[string]string
	fields         map[string]any
	filterFields   map[string]*define.FilterFieldSetting
	sortFields     map[string]*define.SortFieldSetting
	CallbackFilter FilterCallback
	CallbackSort   SortCallback
	OperationLib   *Core.OperationLib
}

func NewDataAdapter[ViewModel any](docId string, data *ViewModel) *DataAdapter[ViewModel] {

	return &DataAdapter[ViewModel]{
		DocId:      docId,
		Data:       data,
		CreateTime: time.Now().Unix(),
	}
}

func (this *DataAdapter[ViewModel]) Analyze(viewModel *ViewModel) {

	var t = reflect.TypeOf(*viewModel)
	var v = reflect.ValueOf(*viewModel)

	if this.fields == nil {
		this.fields = make(map[string]any)
	}

	//Determine whether it is a dynamic field
	var DetermineDynamicFields = func(t1 reflect.Type) bool {

		var name string
		//Determine whether it is a pointer type
		if t1.Kind() == reflect.Ptr {
			//If so, get his subject to judge
			name = t1.Elem().Name()
		} else {
			name = t1.Name()
		}

		return checkIsRef(name) ||
			checkIsManyRef(name) ||
			checkIsFilterField(name) ||
			checkIsSortField(name)
	}
	//Traverse each field in the structure
	for i := 0; i < t.NumField(); i++ {
		//Determine whether it is a dynamic field
		if DetermineDynamicFields(t.Field(i).Type) == true {
			this.AnalyzeDynamicFields(t.Field(i).Name, v.Field(i))
		} else {
			//if exist fields
			if _, exist := this.fields[t.Field(i).Name]; !exist {
				//
				fieldName := t.Field(i).Name
				fieldValue := v.Field(i).Interface()
				this.fields[fieldName] = fieldValue
			}
		}
	}
	return
}

func (this *DataAdapter[ViewModel]) AnalyzeDynamicFields(FieldName string, FieldValue reflect.Value) {

	if this.ref == nil {
		this.ref = make(map[string]string)
	}
	if this.manyRefs == nil {
		this.manyRefs = make(map[string]string)
	}
	if this.filterFields == nil {
		this.filterFields = make(map[string]*define.FilterFieldSetting)
	}

	if this.sortFields == nil {
		this.sortFields = make(map[string]*define.SortFieldSetting)
	}

	var fieldValTypeName string

	if FieldValue.Type().Kind() == reflect.Ptr {
		fieldValTypeName = FieldValue.Elem().Type().Name()
	} else {
		fieldValTypeName = FieldValue.Type().Name()
	}

	switch true {

	case checkIsRef(fieldValTypeName):
		if FieldValue.IsZero() == false {
			if !FieldValue.MethodByName(callOutputFunc).IsZero() &&
				!FieldValue.MethodByName(callOutputFunc).IsNil() {
				//call func
				values := FieldValue.MethodByName(callOutputFunc).Call([]reflect.Value{})
				//
				for _, value := range values {
					//
					if _, exist := this.ref[FieldName]; !exist {
						this.ref[FieldName] = value.String()
					}
				}
			}
		}
	case checkIsManyRef(fieldValTypeName):
		if FieldValue.IsZero() == false {
			if !FieldValue.MethodByName(callOutputFunc).IsZero() &&
				!FieldValue.MethodByName(callOutputFunc).IsNil() {
				values := FieldValue.MethodByName(callOutputFunc).Call([]reflect.Value{})
				//
				for _, value := range values {
					if _, exist := this.ref[FieldName]; !exist {

						this.manyRefs[FieldName] = value.String()

					}
				}
			}
		}
	case checkIsFilterField(fieldValTypeName):
		if FieldValue.IsZero() == false {
			if FieldValue.MethodByName(callOutputFunc).IsZero() == false &&
				FieldValue.MethodByName(callOutputFunc).IsNil() == false {
				//call func
				values := FieldValue.MethodByName(callOutputFunc).Call([]reflect.Value{})
				//
				if len(values) == 2 {
					v1 := values[0]
					v2 := values[1]
					//
					this.filterFields[FieldName] = &define.FilterFieldSetting{
						FieldName: FieldName,
						Value:     v2.Interface(),
						KeyName:   v1.String(),
					}
					//callback
					if this.CallbackFilter != nil {
						err := this.CallbackFilter(v1.String(), v2.Interface())
						if err != nil {
							panic(err)
						}
					}
				}
			}

		}
	case checkIsSortField(fieldValTypeName):
		if FieldValue.IsZero() == false {
			if FieldValue.MethodByName(callOutputFunc).IsZero() == false &&
				FieldValue.MethodByName(callOutputFunc).IsNil() == false {
				values := FieldValue.MethodByName(callOutputFunc).Call([]reflect.Value{})
				//
				if len(values) == 2 {
					v1 := values[0]
					v2 := values[1]
					this.sortFields[FieldName] = &define.SortFieldSetting{
						FieldName: FieldName,
						Value:     v2.Float(),
						KeyName:   v1.String(),
					}
					if this.CallbackSort != nil {
						err := this.CallbackSort(v1.String(), FieldName, v2.Float())
						if err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}
	return
}

func (this *DataAdapter[ViewModel]) Marshal() ([]byte, error) {

	this.Analyze(this.Data)

	var a = define.DataAgreement{
		DocId:       this.DocId,
		CreateTime:  time.Now().Unix(),
		Refs:        this.ref,
		ManyRefs:    this.manyRefs,
		FilterField: this.filterFields,
		SortFields:  this.sortFields,
		Fields:      this.fields,
	}

	return easyjson.Marshal(&a)

}

func (this *DataAdapter[ViewModel]) UnMarshal(dataByte []byte) error {

	var a define.DataAgreement

	if !json.Valid(dataByte) {
		return errors2.ErrNotJson
	}
	err := easyjson.Unmarshal(dataByte, &a)
	if err != nil {
		return err
	}

	var t1 = reflect.TypeOf(*this.Data)
	var v1 = reflect.ValueOf(this.Data)

	for i := 0; i < t1.NumField(); i++ {

		var field = t1.Field(i)
		var fieldName = t1.Field(i).Name

		//current field
		currentField := v1.Elem().Field(i)
		typeName := currentField.Type().Name()

		if field.Type.Kind() == reflect.Ptr {
			typeName = t1.Field(i).Type.Elem().Name()

		}

		//check
		if checkIsDynamicFields(typeName) {
			//rebuild the dynamic field
			this.rebuildDynamicFields(typeName, fieldName, currentField, &a)
		} else {
			if fieldValue, exist := a.Fields[fieldName]; exist {
				if currentField.CanSet() {
					//set the value
					currentField.Set(reflect.ValueOf(fieldValue))
				}
			}
		}
	}
	return nil
}

func (this *DataAdapter[ViewModel]) rebuildDynamicFields(typeName string, fieldName string, currentField reflect.Value, agree *define.DataAgreement) {

	switch true {
	case checkIsRef(typeName):
		if currentField.MethodByName(callRebuildFunc).IsNil() == false {
			//get the key name from the agreement
			if val, exist := agree.Refs[fieldName]; exist {
				//gen the value list
				var values = []reflect.Value{
					reflect.ValueOf(val),
					reflect.ValueOf(this.OperationLib),
				}
				//call the function
				resultValues := currentField.MethodByName(callRebuildFunc).Call(values)
				for _, value := range resultValues {
					//set the result to dsa
					currentField.Set(value)
				}
			}
		}
	case checkIsManyRef(typeName):
		if currentField.MethodByName(callRebuildFunc).IsNil() == false {
			//get the key name from the agreement
			if val, exist := agree.ManyRefs[fieldName]; exist {
				//gen the value list
				var values = []reflect.Value{
					reflect.ValueOf(val),
					reflect.ValueOf(this.OperationLib),
				}
				//call the function
				resultValues := currentField.MethodByName(callRebuildFunc).Call(values)
				for _, value := range resultValues {
					//set the result to dsa
					currentField.Set(value)
				}
			}
		}

	case checkIsSortField(typeName):

		if currentField.MethodByName(callRebuildFunc).IsNil() == false {
			//get the key name from the agreement
			if val, exist := agree.SortFields[fieldName]; exist {

				//gen the value list
				var values = []reflect.Value{
					reflect.ValueOf(val.KeyName),
					reflect.ValueOf(val.FieldName),
					reflect.ValueOf(agree.DocId),
					reflect.ValueOf(val.Value),
					reflect.ValueOf(this.OperationLib),
				}
				//call the function
				resultValues := currentField.MethodByName(callRebuildFunc).Call(values)
				for _, value := range resultValues {
					//set the result to dsa
					currentField.Set(value)
				}
			}
		}

	case checkIsFilterField(typeName):

		if currentField.MethodByName(callRebuildFunc).IsNil() == false {
			//get the key name from the agreement
			if val, exist := agree.FilterField[fieldName]; exist {

				//gen the value list
				var values = []reflect.Value{
					reflect.ValueOf(val.KeyName),
					reflect.ValueOf(val.Value),
					reflect.ValueOf(this.OperationLib),
				}
				//call the function
				resultValues := currentField.MethodByName(callRebuildFunc).Call(values)
				for _, value := range resultValues {
					//set the result to dsa
					currentField.Set(value)
				}
			}
		}

	default:
		return
	}
}

func checkIsRef(name string) bool {
	compile, err := regexp.Compile("^Ref\\[[a-zA-Z0-9.//·]+\\]$")
	if err != nil {
		return false
	}
	return compile.Match([]byte(name))
}

func checkIsManyRef(name string) bool {
	compile, err := regexp.Compile("^ManyRefs\\[[a-zA-Z0-9.//·]+\\]$")
	if err != nil {
		return false
	}
	return compile.Match([]byte(name))
}

func checkIsFilterField(name string) bool {
	compile, err := regexp.Compile("^FilterField\\[[a-zA-Z0-9.\\*]+\\]$")
	if err != nil {
		return false
	}
	return compile.Match([]byte(name))
}

func checkIsSortField(name string) bool {
	compile, err := regexp.Compile("^SortField$")
	if err != nil {
		return false
	}
	return compile.Match([]byte(name))
}

func checkIsDynamicFields(name string) bool {
	return checkIsRef(name) || checkIsManyRef(name) || checkIsFilterField(name) || checkIsSortField(name)
}

func sliceHas(d []string, s string) bool {

	for _, s2 := range d {
		if strings.Compare(s, s2) == 0 {
			return true
		}
	}
	return false
}
