package Sort

type SortFieldsInterface interface {
	Set() error
	GetValue() float64
	Sort(bool, float64) []string
	Remove() error
}

type SortFields struct {
	keyName string
	value   float64
	_       SortFieldsInterface
}

func (this *SortFields) Set(value float64) error                                        { return nil }
func (this *SortFields) Sort(IsAscendingOrder bool, filterGreaterThan float64) []string { return nil }
func (this *SortFields) GetValue() float64                                              { return this.value }
