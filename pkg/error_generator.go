package pkg

import (
	"github.com/pkg/errors"
)

type MyError struct {
	funcName string
}

func NewMyError(fn string) *MyError{
	myErr := &MyError{funcName: fn}
	return myErr

}

func (er *MyError) Wrap(err error, errorMessage string) error{
	if err == nil {
		return errors.Wrap(errors.New(er.funcName), errorMessage)
	}
	return errors.Wrap(err, errorMessage)
}