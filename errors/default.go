package errors

import "errors"

var (
	ErrUnable2ConnectRedis = errors.New("unable to connect redis")
	ErrNil                 = errors.New("return data is nil")
	ErrNotSetSubway        = errors.New("you have not set up Subway")
	ErrAddFailed           = errors.New("add cache failed")
	ErrNotJson             = errors.New("the data structure is not JSON")
	ErrMissingTheKeyName = errors.New("missing the key name")
	ErrUnableGenerateMap = errors.New("unable to generate Map data type")
)
