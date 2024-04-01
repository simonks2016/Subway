package subway

import (
	"errors"
	"fmt"
)

func newBlackKey(viewModelName, dataId string) string {
	return fmt.Sprintf("black_key:%s-%s", viewModelName, dataId)
}

func SetBackKey(ViewModelName, DataId string, expireTime int64) error {

	if Subway == nil {
		return errors.New("you have not set up Subway")
	}
	var key = newBlackKey(ViewModelName, DataId)
	var op = Subway.GetLib()
	return op.SetStringEx(key, DataId, expireTime)
}
func ExistBackKey(ViewModel, DataId string) (bool, error) {

	if Subway == nil {
		return false, errors.New("you have not set up Subway")
	}

	var key = newBlackKey(ViewModel, DataId)
	var op = Subway.GetLib()
	return op.Exist(key)
}
