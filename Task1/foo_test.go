package foo

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type failTestCase struct {
	args []interface{}
	err  error
}

type successTestCase struct {
	args []interface{}
	res  FooResult
}

type fooSuite struct {
	failTestCases    []failTestCase
	successTestCases []successTestCase
}

var _ = Suite(&fooSuite{})

func (s *fooSuite) SetUpTest(c *C) {
	s.failTestCases = []failTestCase{
		failTestCase{
			args: []interface{}{},
			err:  notEnoughParametersError,
		},
		failTestCase{
			args: []interface{}{0, 0, time.Now(), new(int), 0},
			err:  tooManyParametersError,
		},
		failTestCase{
			args: []interface{}{""},
			err:  firstParameterTypeError,
		},
		failTestCase{
			args: []interface{}{0, ""},
			err:  secondParameterTypeError,
		},
		failTestCase{
			args: []interface{}{0, 0, true},
			err:  thirdParameterTypeError,
		},
		failTestCase{
			args: []interface{}{0, 0, defaultTime, ""},
			err:  fourthParameterTypeError,
		},
	}
	res := FooResult{
		Required{1},
		Optional{2, time.Now(), new(int)},
	}
	s.successTestCases = []successTestCase{
		successTestCase{
			args: []interface{}{res.Required.int},
			res:  FooResult{Required{res.Required.int}, Optional{defaultInt, defaultTime, defaultPointer}},
		},
		successTestCase{
			args: []interface{}{res.Required.int, res.Optional.int},
			res:  FooResult{Required{res.Required.int}, Optional{res.Optional.int, defaultTime, defaultPointer}},
		},
		successTestCase{
			args: []interface{}{res.Required.int, defaultInt, res.Optional.Time},
			res:  FooResult{Required{res.Required.int}, Optional{defaultInt, res.Optional.Time, defaultPointer}},
		},
		successTestCase{
			args: []interface{}{res.Required.int, defaultInt, defaultTime, res.Optional.Pointer},
			res:  FooResult{Required{res.Required.int}, Optional{defaultInt, defaultTime, res.Optional.Pointer}},
		},
		successTestCase{
			args: []interface{}{res.Required.int, res.Optional.int, res.Optional.Time, res.Optional.Pointer},
			res:  res,
		},
	}
}

func (s *fooSuite) TestFoo(c *C) {
	for _, t := range s.failTestCases {
		c.Assert(NewFoo(time.Now).RunWithError(t.args...), Equals, t.err)
	}
	for _, t := range s.successTestCases {
		c.Assert(NewFoo(nil).RunWithPanic(t.args...), Equals, t.res)
	}
}
