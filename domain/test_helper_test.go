package domain_test

import (
	"reflect"
	"testing"
)

type asserter struct {
	t *testing.T
}

func (a *asserter) eq(expected, actual interface{}) {
	a.t.Helper()
	if expected != actual {
		a.t.Errorf("Expected %v to equal %v", actual, expected)
	}
}

func (a *asserter) neq(expected, actual interface{}) {
	a.t.Helper()
	if expected == actual {
		a.t.Errorf("Expected %v to not equal %v", actual, expected)
	}
}

func (a *asserter) isNil(actual interface{}) {
	a.t.Helper()
	if !a._isReallyNil(actual) {
		a.t.Errorf("Expected [type: %T] %v to be nil", actual, actual)
	}
}

func (a *asserter) isNotNil(actual interface{}) {
	a.t.Helper()
	if a._isReallyNil(actual) {
		a.t.Error("Expected actual to be not nil")
	}
}

func (a *asserter) err(expected string, e error) {
	a.t.Helper()
	if nil == e {
		a.t.Errorf("Expected error %q, but got nil", expected)
		return
	}

	if e.Error() != expected {
		a.t.Errorf("Expected error %q, but got %q", expected, e)
	}
}

func (a *asserter) ok(err error) {
	a.t.Helper()
	if nil != err {
		a.t.Errorf("Expected no error, but got %v", err)
	}
}

func (a *asserter) is(b bool) {
	a.t.Helper()
	if !b {
		a.t.Errorf("Expected to be true")
	}
}

func (a *asserter) not(b bool) {
	a.t.Helper()
	if b {
		a.t.Errorf("Expected to be false")
	}
}

func (a *asserter) _isReallyNil(i interface{}) bool {
	a.t.Helper()
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func asserts(t *testing.T) *asserter {
	return &asserter{t: t}
}
