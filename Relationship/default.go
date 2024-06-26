package Relationship

type Controllers interface {
	GetCount() (int, error)
	IsMember(string) (bool, error)
	Add(...string) error
	Pop(count int) ([]string, error)
	RandMembers(int) ([]string, error)
	Remove(...string) error
	Members() ([]string, error)
	DeleteAll() error
}

type LineControllers interface {
	GetCount() (int64, error)
	Add(string, float64) error
	Remove(...string) error
	Get(int64, int64, bool) ([]string, error)
	IsMember(string) (bool, error)
	DeleteAll() error
}
