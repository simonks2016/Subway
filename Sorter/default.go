package Sorter

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/simonks2016/Subway/Core"
	errors2 "github.com/simonks2016/Subway/errors"
	"sort"
	"strconv"
)

type SortFieldsInterface interface {
	Set() error
	GetValue() float64
	Sort(bool, float64) []string
	Remove() error
	SortFilteredData(bool, []string, int, int)
	SetRedisConn(*redis.Conn)
	Upgrade() error
}

type SortField struct {
	keyName      string
	fieldName    string
	docId        string
	value        float64
	_            SortFieldsInterface
	operationLib *Core.OperationLib
}

func (this *SortField) Remove() error {

	err := this.operationLib.DelHashMap(this.keyName, this.docId)
	if err != nil {
		return err
	}
	return nil
}

func (this *SortField) Set(value float64) error {

	if this.operationLib == nil {
		return errors2.ErrNotSetSubway
	}
	err := this.operationLib.SetHashMap(this.keyName, this.docId, value)
	if err != nil {
		return err
	}
	return nil
}
func (this *SortField) getAllData() (map[string]float64, error) {

	data, err := this.operationLib.GetALLHashMap(this.keyName)
	if err != nil {
		return nil, err
	}

	var response = make(map[string]float64)
	var num = len(data) / 2
	if len(data)%2 != 0 {
		return nil, errors2.ErrUnableGenerateMap
	}

	for i := 0; i < num; i++ {

		var o1 = i + (i * 1)
		var o2 = o1 + 1
		//get the value
		var key, value = data[o1].([]uint8), data[o2].([]uint8)
		var k1, v1 = string(key), string(value)
		//parse float
		float, err := strconv.ParseFloat(v1, 64)
		if err != nil {
			return nil, err
		}
		response[k1] = float
	}

	return response, nil
}

func (this *SortField) Sort(IsAscendingOrder bool, filterGreaterThan float64) []string {

	d, err := this.getAllData()
	if err != nil {
		return nil
	}

	var m AscendingAlgorithm
	//make the data
	for dataId, Score := range d {
		//
		if filterGreaterThan >= -1 {
			if Score <= filterGreaterThan {
				continue
			}
		}
		m = append(m, insideSortStruct{
			DataId: dataId,
			Score:  Score,
		})
	}

	if IsAscendingOrder {
		sort.Sort(m)
	} else {
		sort.Sort(sort.Reverse(m))
	}

	var response []string
	//make response
	for _, sortStruct := range m {
		response = append(response, sortStruct.DataId)
	}
	return response

}
func (this *SortField) GetValue() float64 {

	return this.value
}
func (this *SortField) SortFilteredData(IsAscendingOrder bool, FilterDataIds []string, offset, limit int) []string {
	d, err := this.getAllData()
	if err != nil {
		return nil
	}

	var m AscendingAlgorithm
	//make the data
	for _, id := range FilterDataIds {
		if val, exist := d[id]; exist {
			m = append(m, insideSortStruct{
				DataId: id,
				Score:  val,
			})
		}
	}

	if offset > m.Len() {
		return nil
	}

	if IsAscendingOrder {
		sort.Sort(m)
	} else {
		sort.Sort(sort.Reverse(m))
	}

	var end = offset + limit

	if end > m.Len() {
		end = m.Len()
	}

	var newData = m[offset:end]
	var response []string

	for _, datum := range newData {
		//response
		response = append(response, datum.DataId)
	}
	return response
}
func (this *SortField) Output() (string, float64) {
	return this.keyName, this.value
}

func (this *SortField) Rebuild(keyName, fieldName, docId string, value float64, ol *Core.OperationLib) *SortField {

	return &SortField{
		keyName:      keyName,
		value:        value,
		fieldName:    fieldName,
		docId:        docId,
		operationLib: ol,
	}
}

func CreateSorter(DataModelName, fieldName string, ol *Core.OperationLib) *SortField {

	return &SortField{
		keyName:      NewKeyId(DataModelName, fieldName),
		value:        0,
		operationLib: ol,
	}
}

func NewSorter(dataModelName, fieldName string, fieldValue float64) *SortField {

	s := &SortField{
		keyName:   NewKeyId(dataModelName, fieldName),
		fieldName: fieldName,
		value:     fieldValue,
	}
	return s
}

func NewKeyId(dataModelName, fieldName string) string {
	return fmt.Sprintf("Sorter-%s-%s", dataModelName, fieldName)
}
