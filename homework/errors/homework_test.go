package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	msg := fmt.Sprintf("%d errors occured:\n", len(e.Errors))
	for _, err := range e.Errors {
		msg += fmt.Sprintf("\t* %s", err.Error())
	}
	return msg + "\n"
}

func Append(err error, errs ...error) *MultiError {
	var me *MultiError
	if err != nil {
		if existingME, ok := err.(*MultiError); ok {
			me = existingME
		} else {
			me = &MultiError{Errors: []error{err}}
		}
	} else {
		me = &MultiError{}
	}
	for _, e := range errs {
		if e != nil {
			me.Errors = append(me.Errors, e)
		}
	}
	return me
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
