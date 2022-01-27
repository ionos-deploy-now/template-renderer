package test

import "testing"

func AssertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Fatalf("\nExpected: %v\nActual: %v", expected, actual)
	}
}

func AssertNotEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected == actual {
		t.Fatalf("\nExpected not: %v\nActual: %v", expected, actual)
	}
}
