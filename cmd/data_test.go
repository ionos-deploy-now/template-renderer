package cmd

import (
	"fmt"
	"template-renderer/test"
	"testing"
)

func TestParseDataYaml(t *testing.T) {
	yaml := `a: 1`
	values, err := parseData(yaml)

	test.AssertEqual(t, nil, err)
	test.AssertNotEqual(t, nil, values["a"])
	test.AssertEqual(t, 1, values["a"])
}

func TestParseDataJson(t *testing.T) {
	yaml := `{"a": 1}`
	values, err := parseData(yaml)

	test.AssertEqual(t, nil, err)
	test.AssertNotEqual(t, nil, values["a"])
	test.AssertEqual(t, 1, values["a"])
}

func TestMerge(t *testing.T) {
	values1 := Values{"a": 1, "b": Values{"c": 2}}
	values2 := Values{"a": 2, "b": Values{"d": 3}}

	values3 := values1.merge(values2)

	test.AssertEqual(t, 1, values3["a"])
	test.AssertEqual(t, 2, values3["b"].(Values)["c"])
	test.AssertEqual(t, 3, values3["b"].(Values)["d"])
}

func TestConvert(t *testing.T) {
	data := Data{}
	values1 := Values{"a": 1, "b": Values{"c": 2}}
	values2 := values1.convertToRuntimeValues(&data)
	data.Values = values2

	test.AssertEqual(t, "1", fmt.Sprintf("%v", data.Values["a"]))
	test.AssertEqual(t, "2", fmt.Sprintf("%v", data.Values["b"].(Values)["c"]))

	test.AssertEqual(t, true, data.RuntimeValuesUsed)
}
