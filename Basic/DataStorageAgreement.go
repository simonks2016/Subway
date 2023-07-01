package Basic

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

type DSA[ViewModel any] struct {
	DocId              string                 `json:"doc_id"`
	Data               *ViewModel             `json:"data"`
	CreateTime         int64                  `json:"create_time"`
	ExtraData          map[string]interface{} `json:"extra_data"`
	RelationshipField  []string               `json:"relationship_field"`
	AssociatedDocument map[string]string      `json:"associated_document"`
	Line               []string               `json:"line"`
}

func NewDSA[ViewModel any](docId string, data *ViewModel, ExtraData map[string]interface{}) *DSA[ViewModel] {

	return &DSA[ViewModel]{
		DocId:              docId,
		Data:               data,
		ExtraData:          ExtraData,
		CreateTime:         time.Now().Unix(),
		AssociatedDocument: make(map[string]string),
		Line:               []string{},
		RelationshipField:  []string{},
	}
}

func (this *DSA[ViewModel]) Marshal() ([]byte, error) {
	return json.Marshal(this)
}

func (this *DSA[ViewModel]) UnMarshal(dataByte []byte) error {

	return json.Unmarshal(dataByte, this)
}

func (this *DSA[ViewModel]) AddRelationship(fieldName string) {
	this.RelationshipField = append(this.RelationshipField, fieldName)
}

func (this *DSA[ViewModel]) AddAssociatedDocument(fieldName string, TargetViewModel any) {

	var ViewModelName = reflect.TypeOf(TargetViewModel).Name()
	this.AssociatedDocument[fieldName] = ViewModelName
}

func (this *DSA[ViewModel]) AddLine(sortFieldName string) {

	this.Line = append(this.Line, sortFieldName)
}

func (this *DSA[ViewModel]) HasRelationship(fieldName string) bool {
	return sliceHas(this.RelationshipField, fieldName)
}

func (this *DSA[ViewModel]) HasLine(f string) bool{
	return sliceHas(this.Line,f)
}

func sliceHas(d []string, s string) bool {

	for _, s2 := range d {
		if strings.Compare(s, s2) == 0 {
			return true
		}
	}
	return false
}
