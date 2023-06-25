package Basic

type Wheel interface {
	SAdd(string, ...interface{}) error
	SMember(string) []string
	SIsMember(string, string) bool
	SPop(string) (error, string)

	ZAdd(string, float64, ...interface{}) error
	ZRemove(string) error

	SetString(string, string) error
	GetString(string) (error, string)
	BatchGetString(...interface{}) (error, []string)
	Delete(string) error
}
