package cmd

import (
	"template-renderer/test"
	"testing"
)

func TestParseDataYaml(t *testing.T) {
	yaml := `a: 1`
	data, err := parseData(yaml)

	test.AssertEqual(t, nil, err)
	test.AssertNotEqual(t, nil, data["a"])
	test.AssertEqual(t, 1, data["a"])
}

func TestParseDataJson(t *testing.T) {
	yaml := `{"a": 1}`
	data, err := parseData(yaml)

	test.AssertEqual(t, nil, err)
	test.AssertNotEqual(t, nil, data["a"])
	test.AssertEqual(t, 1, data["a"])
}

func TestMerge(t *testing.T) {
	data1 := Data{"a": 1, "b": Data{"c": 2}}
	data2 := Data{"a": 2, "b": Data{"d": 3}}

	data3 := data1.merge(data2)

	test.AssertEqual(t, 1, data3["a"])
	test.AssertEqual(t, 2, data3["b"].(Data)["c"])
	test.AssertEqual(t, 3, data3["b"].(Data)["d"])

}
