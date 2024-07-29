package DataAdapter

import (
	"encoding/json"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/simonks2016/Subway/Core"
	"github.com/simonks2016/Subway/Filter"
	"github.com/simonks2016/Subway/Sorter"
	"github.com/simonks2016/Subway/define"
	errors2 "github.com/simonks2016/Subway/errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const callOutputFunc = "Output"
const callRebuildFunc = "Rebuild"

type DiscoveryFilterField func(key string, field any) error
type DiscoverySortField func(key, fieldName string, value float64) error

type DataAdapter[ViewModel any] struct {
	DocId                        string     `json:"doc_id"`
	Data                         *ViewModel `json:"-"`
	CreateTime                   int64      `json:"create_time"`
	ref                          map[string]string
	manyRefs                     map[string]string
	fields                       map[string]any
	filterFields                 map[string]*define.FilterFieldSetting
	sortFields                   map[string]*define.SortFieldSetting
	CallbackDiscoveryFilterField DiscoveryFilterField
	CallbackDiscoverySortField   DiscoverySortField
	OperationLib                 *Core.OperationLib
}

func NewDataAdapter[ViewModel any](docId string, data *ViewModel) *DataAdapter[ViewModel] {

	return &DataAdapter[ViewModel]{
		DocId:      docId,
		Data:       data,
		CreateTime: time.Now().Unix(),
	}
}

func (this *DataAdapter[ViewModel]) analyze(viewModel *ViewModel) {

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
			this.analyzeDynamicFields(t.Field(i).Name, v.Field(i))
		} else {
			if t.Field(i).Name == "Id" {
				if len(this.DocId) <= 0 {
					this.DocId = Core.NewDocumentId(
						Core.GetViewModelName(this.Data),
						v.Field(i).String())
				}
			}

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

func (this *DataAdapter[ViewModel]) analyzeDynamicFields(FieldName string, FieldValue reflect.Value) {

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

	if FieldValue.IsZero() == true || FieldValue.IsNil() == true {
		//panic(fmt.Sprintf("You did not initialize the field(%s)", FieldName))
		return
	}
	if FieldValue.Type().Kind() == reflect.Ptr {
		fieldValTypeName = FieldValue.Elem().Type().Name()
	} else {
		fieldValTypeName = FieldValue.Type().Name()
	}

	switch true {

	case checkIsRef(fieldValTypeName):
		if FieldValue.IsZero() == false && FieldValue.IsNil() == false {
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
		if FieldValue.IsZero() == false && FieldValue.IsNil() == false {
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
		if FieldValue.IsZero() == false && FieldValue.IsNil() == false {
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
					if this.CallbackDiscoveryFilterField != nil {
						err := this.CallbackDiscoveryFilterField(v1.String(), v2.Interface())
						if err != nil {
							panic(err)
						}
					}
				}
			}

		}
	case checkIsSortField(fieldValTypeName):
		if FieldValue.IsZero() == false && FieldValue.IsNil() == false {
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
					if this.CallbackDiscoverySortField != nil {
						err := this.CallbackDiscoverySortField(
							v1.String(), FieldName, v2.Float())
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

	this.analyze(this.Data)

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

func (this *DataAdapter[ViewModel]) SimpleUnMarshal(dataByte []byte) (error, *define.DataAgreement) {

	var a define.DataAgreement

	if !json.Valid(dataByte) {
		return errors2.ErrNotJson, nil
	}
	err := easyjson.Unmarshal(dataByte, &a)
	if err != nil {
		return err, nil
	}

	return nil, &a
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

		switch true {
		case checkIsDynamicFields(typeName):
			//rebuild the dynamic field
			this.rebuildDynamicFields(typeName, fieldName, currentField, &a)
		case field.Type.Kind() == reflect.Ptr, field.Type.Kind() == reflect.Struct:
			//rebuild struct
			if fieldValue, exist := a.Fields[fieldName]; exist {
				this.rebuildStructField(field.Type, currentField, fieldValue.(map[string]interface{}))
			}

		default:
			if fieldValue, exist := a.Fields[fieldName]; exist {
				if currentField.CanSet() {
					//Are the types of both parties consistent?
					if reflect.ValueOf(fieldValue).Type() == currentField.Type() {
						//set the value
						currentField.Set(reflect.ValueOf(fieldValue))
					} else {
						fv := reflect.ValueOf(fieldValue)

						switch currentField.Type().Kind() {
						case reflect.Slice:
							handleSlice(fv, currentField)
						default:
							if reflect.ValueOf(fieldValue).CanConvert(currentField.Type()) == false {
								panic("Unable to convert target type")
							}
							currentField.Set(fv.Convert(currentField.Type()))
						}

					}
				}
			}

		}
	}
	return nil
}

func (this *DataAdapter[ViewModel]) rebuildDynamicFields(typeName string, fieldName string, currentField reflect.Value, agree *define.DataAgreement) {

	var keyName, value, fn, docId reflect.Value
	var ol = reflect.ValueOf(this.OperationLib)
	var ViewModelName = Core.GetViewModelName(this.Data)

	switch true {
	case checkIsRef(typeName):

		if currentField.MethodByName(callRebuildFunc).IsNil() == false {
			//get the key name from the agreement
			if val, exist := agree.Refs[fieldName]; exist {
				//gen the value list
				var values = []reflect.Value{
					reflect.ValueOf(val),
					ol,
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
					ol,
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
				keyName = reflect.ValueOf(val.KeyName)
				fn = reflect.ValueOf(val.FieldName)
				docId = reflect.ValueOf(agree.DocId)
				value = reflect.ValueOf(val.Value)
			} else {
				keyName = reflect.ValueOf(Sorter.NewKeyId(ViewModelName, fieldName))
				fn = reflect.ValueOf(fieldName)
				docId = reflect.ValueOf("")
				value = reflect.ValueOf(0.0)
			}
			//call the function
			resultValues := currentField.MethodByName(callRebuildFunc).Call([]reflect.Value{
				keyName,
				fn,
				docId,
				value,
				ol,
			})
			for _, value := range resultValues {
				//set the result to dsa
				currentField.Set(value)
			}
		}

	case checkIsFilterField(typeName):

		if currentField.MethodByName(callRebuildFunc).IsNil() == false {
			//get the key name from the agreement
			if val, exist := agree.FilterField[fieldName]; exist {
				keyName = reflect.ValueOf(val.KeyName)
				value = reflect.ValueOf(val.Value)
			} else {
				keyName = reflect.ValueOf(Filter.NewKeyId(ViewModelName, fieldName))
				value = reflect.ValueOf("")
			}

			//call the function
			resultValues := currentField.MethodByName(callRebuildFunc).Call([]reflect.Value{
				keyName,
				value,
				ol,
			})
			for _, newField := range resultValues {
				//set the result to dsa
				currentField.Set(newField)
			}
		}

	default:
		return
	}
}

func (this DataAdapter[ViewModel]) rebuildStructField(t reflect.Type, currentField reflect.Value, data map[string]interface{}) {

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.New(t)
	//init
	initializeStruct(t, v.Elem())
	//
	v1 := v.Elem()
	//copy data to
	for i := 0; i < v1.NumField(); i++ {

		//get the name
		var name = t.Field(i).Name
		var tagName = t.Field(i).Tag.Get("json")
		var fieldType = t.Field(i).Type
		//if exist
		val, exist := data[name]
		if !exist {
			//if tag name is not in data
			if val, exist = data[tagName]; !exist {
				continue
			}
		}
		//if the field is struct or ptr
		if fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Struct {
			this.rebuildStructField(t, v1.Field(i), val.(map[string]interface{}))
		} else {
			v1.Field(i).Set(reflect.ValueOf(val))
		}
	}
	//set the field
	currentField.Set(v)
}

func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)

		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			//new prt
			fv := reflect.New(ft.Type.Elem())
			//init
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
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

func handleSlice(value reflect.Value, targetField reflect.Value) {
	Len := value.Len()
	cf := reflect.MakeSlice(targetField.Type(), Len, (Len-1)*2)
	targetType := targetField.Type()

	for i1 := 0; i1 < Len; i1++ {
		vf := value.Index(i1).Interface()
		svf := fmt.Sprintf("%v", vf)

		switch targetType.Elem().Kind() {
		case reflect.String:
			if reflect.ValueOf(svf).CanConvert(targetType.Elem()) {
				v5 := reflect.ValueOf(svf).Convert(targetType.Elem())
				cf.Index(i1).Set(v5)
			} else {
				panic("cannot be converted to string")
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

			parseInt, err := strconv.ParseInt(svf, 0, 64)
			if err != nil {
				panic(err)
			}
			if reflect.ValueOf(parseInt).CanConvert(targetType.Elem()) {
				v5 := reflect.ValueOf(parseInt).Convert(targetType.Elem())
				cf.Index(i1).Set(v5)
			} else {
				panic("cannot be converted to int")
			}
		case reflect.Float32, reflect.Float64:
			parseFloat, err := strconv.ParseFloat(svf, 64)
			if err != nil {
				panic(err)
			}
			if reflect.ValueOf(parseFloat).CanConvert(targetType.Elem()) {
				v5 := reflect.ValueOf(parseFloat).Convert(targetType.Elem())
				cf.Index(i1).Set(v5)
			} else {
				panic("cannot be converted to float")
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			parseUint, err := strconv.ParseUint(svf, 0, 64)
			if err != nil {
				panic(err)
			}
			if reflect.ValueOf(parseUint).CanConvert(targetType.Elem()) {
				v5 := reflect.ValueOf(parseUint).Convert(targetType.Elem())
				cf.Index(i1).Set(v5)
			} else {
				panic("cannot be converted to int")
			}
		case reflect.Bool:
			parseBool, err := strconv.ParseBool(svf)
			if err != nil {
				panic(err)
			}
			if reflect.ValueOf(parseBool).CanConvert(targetType.Elem()) {
				v5 := reflect.ValueOf(parseBool).Convert(targetType.Elem())
				cf.Index(i1).Set(v5)
			} else {
				panic("cannot be converted to int")
			}
		//[]byte
		default:
			panic("The target type cannot be" + targetType.Name())
		}

	}
	targetField.Set(cf)
}
