package Filter

import (
	"fmt"
	"testing"
)

func TestSetFilter(t *testing.T) {

	type a1 struct {
		Name  string `json:"name"`
		State int    `json:"state"`
	}
	var a = a1{
		Name:  "simonks",
		State: 1,
	}

	filter, err := CreateFilter[int](&a, "State")
	if err != nil {
		return
	}

	filter.Set(1, "a", "c", "d", "a1")
	filter.Set(0, "b")
	filter.Set(2, "a2", "a3", "a4")

	d1 := filter.GetSameConditions(2, Less[int])
	fmt.Println(d1)
	fmt.Println(filter.Get())

	t.Deadline()
}
