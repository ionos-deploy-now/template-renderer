package cmd

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

type RuntimeValue struct {
	value            string
	updateUsedValues func()
}

func (v RuntimeValue) String() string {
	v.updateUsedValues()
	return v.value
}

func (v RuntimeValue) MarshalJSON() ([]byte, error) {
	return []byte("\"" + v.String() + "\""), nil
}

type Values map[string]interface{}

type Data struct {
	Values            Values
	RuntimeValuesUsed bool
}

func (d Data) ResetUsedRuntimeValues() {
	d.RuntimeValuesUsed = false
}

func ParseInputData(secrets string, runtimeData string, additionalData []string) (*Data, error) {
	data := Data{}
	values := Values{}
	if secrets != "" {
		values2, err := parseData(secrets)
		if err != nil {
			return nil, err
		}
		values = values.merge(Values{"secrets": values2})
	}
	if runtimeData != "" {
		values2, err := parseData(runtimeData)
		values2 = values2.convertToRuntimeValues(&data)
		if err != nil {
			return nil, err
		}
		values = values.merge(Values{"runtime": values2})
	}
	for _, datum := range additionalData {
		values2, err := parseData(datum)
		if err != nil {
			return nil, err
		}
		values = values.merge(values2)
	}
	data.Values = values
	return &data, nil
}

func parseData(input string) (Values, error) {
	var values Values
	err := yaml.Unmarshal([]byte(input), &values)
	return values, err
}

func (v *Values) convertToRuntimeValues(data *Data) Values {
	for key, value := range *v {
		if subData, ok := value.(Values); ok {
			subData.convertToRuntimeValues(data)
		} else {
			stringValue := fmt.Sprintf("%v", value)
			(*v)[key] = RuntimeValue{
				value: stringValue,
				updateUsedValues: func() {
					(*data).RuntimeValuesUsed = true
				},
			}
		}
	}
	return *v
}

func (v *Values) merge(v2 Values) Values {
	for key, value := range v2 {
		if (*v)[key] == nil {
			(*v)[key] = value
		} else {
			v1, ok1 := (*v)[key].(Values)
			v2, ok2 := value.(Values)
			if ok1 && ok2 {
				(*v)[key] = v1.merge(v2)
			}
		}
	}
	return *v
}

func (v Values) String() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
