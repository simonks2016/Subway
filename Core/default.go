package Core

import "fmt"

type Wheel interface {
	SAdd(string, ...interface{}) error
	SMembers(string) ([]string, error) //
	SIsMember(string, string) (bool, error)
	SPop(string, int) ([]string, error)
	SCard(string) (int, error)
	SRandMember(string, int) ([]string, error)
	SRemove(SetName string, key ...interface{}) error
	SIncr(...interface{}) ([]string, error)
	SUnion(...string) ([]string, error)

	ZAdd(string, float64, ...interface{}) error
	ZRemove(string, ...interface{}) (error, int)
	ZRange(string, int64, int64, bool) (error, []string)
	ZCard(string) (error, int64)
	ZRangeBySore(string, int64, int64, int64, int64) (error, []string)
	ZIsMember(string, string) (error, bool)

	SetString(string, string) error
	GetString(string) (error, string)
	GetByte(string) (error, []byte)
	BatchGetStrings(...interface{}) (error, []string)
	Delete(string) error
	BatchDelete(...interface{}) error

	SetHashMap(interface{}, interface{}, ...interface{}) error
	GetHashMap(interface{}, interface{}) (any, error)
	MSetHashMap(interface{}, map[any]any) error
	DelHashMap(interface{}, interface{}) error
	GetFieldsHashMap(interface{}) ([]string, error)
	GetALLHashMap(interface{}) (map[any]any, error)

	Persist(interface{}) error
	Keys() ([]string, error)

	NewDocumentId(string, string) string
}

type DocumentIds []string
type DocumentId string

func NewDocumentId(ViewModelName string, dataId string) string {
	return fmt.Sprintf("%s-%s", ViewModelName, dataId)
}

func (this DocumentId) ToString() string {
	return string(this)
}
