package subway

import (
	"github.com/simonks2016/Subway/Core"
	errors2 "github.com/simonks2016/Subway/errors"
)

func Exists[ViewModel any](Id string) (bool, error) {

	if Subway == nil {
		return false, errors2.ErrNotSetSubway
	}
	var vm ViewModel
	var VMName = GetViewModelName(vm)
	var docId = Core.NewDocumentId(VMName, Id)
	var lib = Subway.GetLib()
	//check
	return lib.Exist(docId)
}
