package Relationship

import (
	"fmt"
	"github.com/simonks2016/Subway/Basic"
)

type BasicLineController struct {
	LineControllers
	OL      Basic.OperationLib
	SetName string
}

func (this BasicLineController) GetCount() (int64, error) {

	err, i := this.OL.ZCard(this.SetName)
	if err != nil {
		return 0, err
	}
	return i, err
}
func (this BasicLineController) Add(dataId string, score float64) error {
	return this.OL.ZAdd(this.SetName, score, dataId)
}
func (this BasicLineController) Remove(keys ...string) error {

	var args []interface{}

	for _, key := range keys {
		args = append(args, key)
	}

	err, _ := this.OL.ZRemove(this.SetName, args...)
	if err != nil {
		return err
	}
	return nil
}
func (this BasicLineController) Get(start, end int64, desc bool) ([]string, error) {

	fmt.Println(this.SetName)
	err, i := this.OL.ZRange(this.SetName, start, end, desc)
	if err != nil {
		return nil, err
	}
	return i, err
}

func NewBasicLineControllers(SetName string, OL Basic.OperationLib) *BasicLineController {

	return &BasicLineController{SetName: SetName, OL: OL}
}
