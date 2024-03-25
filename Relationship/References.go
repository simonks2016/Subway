package Relationship

type References[dataModel comparable] struct {
	Name      string
	dataModel dataModel
}

func (this *References[dataModel]) Load() (*dataModel, error) {

	return nil, nil
}
