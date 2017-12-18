package foo

import (
	"errors"
	"reflect"
	"time"
)

type Required struct {
	int
}

type Pointer interface{}

type Optional struct {
	int
	time.Time
	Pointer
}

type FooResult struct {
	Required
	Optional
}

var (
	notEnoughParametersError = errors.New("Not enough parameters")
	tooManyParametersError   = errors.New("Too many parameters")
	firstParameterTypeError  = errors.New("Parameter type is not int")
	secondParameterTypeError = firstParameterTypeError
	thirdParameterTypeError  = errors.New("Parameter type is not time.Time")
	fourthParameterTypeError = errors.New("Parameter type does not point to smth")
	inappropriatePanicError  = errors.New("Inappropriate panic")

	defaultInt     int         = 10
	defaultPointer interface{} = nil
	defaultTime    time.Time   = time.Time{}
)

type Foo struct {
	getDefaultTime func() time.Time
}

func NewFoo(fn func() time.Time) *Foo {
	if fn == nil {
		return &Foo{getDefaultTime: func() time.Time { return defaultTime }}
	}

	return &Foo{getDefaultTime: fn}
}

func (f *Foo) RunWithError(args ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = inappropriatePanicError
			}
		}
	}()
	f.RunWithPanic(args...)
	return
}

func (f *Foo) RunWithPanic(args ...interface{}) (res FooResult) {
	res.Optional.int = defaultInt
	res.Optional.Time = f.getDefaultTime()
	res.Optional.Pointer = defaultPointer
	if len(args) < 1 {
		panic(notEnoughParametersError)
	}
	var ok bool
	for i, arg := range args {
		switch i {
		case 0:
			res.Required.int, ok = arg.(int)
			if !ok {
				panic(firstParameterTypeError)
			}

		case 1:
			res.Optional.int, ok = arg.(int)
			if !ok {
				panic(secondParameterTypeError)
			}

		case 2:
			res.Optional.Time, ok = arg.(time.Time)
			if !ok {
				panic(thirdParameterTypeError)
			}

		case 3:
			val := reflect.ValueOf(arg)
			if !val.IsValid() {
				arg = nil
			}
			if val.IsValid() && val.Kind() != reflect.Ptr {
				panic(fourthParameterTypeError)
			}
			res.Optional.Pointer = arg

		default:
			panic(tooManyParametersError)
		}
	}
	return
}
