package subway

import (
	"github.com/simonks2016/Subway/Core"
	errors2 "github.com/simonks2016/Subway/errors"
)

func Expire[VM any](DataId string, ExpirationSeconds int64) error {

	if Subway == nil {
		return errors2.ErrNotSetSubway
	}
	var lib = Subway.GetLib()
	var vm VM
	//get View model Name
	vm_name := Core.GetViewModelName(vm)
	//set the expiry
	return lib.Expire(Core.NewDocumentId(vm_name, DataId), ExpirationSeconds)
}
