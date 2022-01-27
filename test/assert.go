package test

import "testing"

func AssertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Fatalf("Expected: %v, Actual: %v", expected, actual)
	}
}

func AssertNotEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected == actual {
		t.Fatalf("Expected: %v, Actual: %v", expected, actual)
	}
}
