package Relationship

import (
	"errors"
	"github.com/simonks2016/Subway/Core"
)

type BasicRelationshipControllers struct {
	SetName string
	Controllers
	OperationLib Core.OperationLib
}

func (this BasicRelationshipControllers) GetCount() (int, error) {

	if len(this.SetName) <= 0 {
		return 0, errors.New("collection key cannot be empty")
	}
	return this.OperationLib.SCard(this.SetName)
}

func (this BasicRelationshipControllers) IsMember(docId string) (bool, error) {

	if len(this.SetName) <= 0 {
		return false, errors.New("collection key cannot be empty")
	}
	return this.OperationLib.SIsMember(this.SetName, docId)
}

func (this BasicRelationshipControllers) Add(docIds ...string) error {

	var arg []interface{}

	for _, id := range docIds {
		arg = append(arg, id)
	}
	return this.OperationLib.SAdd(this.SetName, arg...)
}

func (this BasicRelationshipControllers) Pop(count int) ([]string, error) {
	return this.OperationLib.SPop(this.SetName, count)
}

func (this BasicRelationshipControllers) RandMembers(Count int) ([]string, error) {
	return this.OperationLib.SRandMember(this.SetName, Count)
}

func (this BasicRelationshipControllers) Remove(key ...string) error {

	var args []interface{}
	//make args slice
	for _, s := range key {
		args = append(args, s)
	}
	return this.OperationLib.SRemove(this.SetName, args...)
}

func (this BasicRelationshipControllers) Members() ([]string, error) {
	return this.OperationLib.SMembers(this.SetName)
}

func NewBasicRelationshipControllers(SetName string, OL Core.OperationLib) BasicRelationshipControllers {

	return BasicRelationshipControllers{
		OperationLib: OL,
		SetName:      SetName,
	}
}
