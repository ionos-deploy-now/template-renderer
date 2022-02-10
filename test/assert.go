package test

import "testing"

func AssertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Fatalf("\nExpected: %v (%T)\nActual:   %v (%T)", expected, expected, actual, actual)
	}
}

func AssertNotEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected == actual {
		t.Fatalf("\nExpected not: %v (%T)\nActual: %v (%T)", expected, expected, actual, actual)
	}
}
