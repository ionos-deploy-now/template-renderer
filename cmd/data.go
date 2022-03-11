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

type Data map[string]interface{}

func ParseInputData(secrets string, runtimeData string, additionalData []string, runtimePlaceholderCount *int) (Data, error) {
	data := Data{}
	if secrets != "" {
		data2, err := parseData(secrets)
		if err != nil {
			return nil, err
		}
		data = data.merge(Data{"secrets": data2})
	}
	if runtimeData != "" {
		data2, err := parseData(runtimeData)
		data2 = data2.convertToRuntimeValues(runtimePlaceholderCount)
		if err != nil {
			return nil, err
		}
		data = data.merge(Data{"runtime": data2})
	}
	for _, datum := range additionalData {
		data2, err := parseData(datum)
		if err != nil {
			return nil, err
		}
		data = data.merge(data2)
	}
	return data, nil
}

func parseData(input string) (Data, error) {
	var data Data
	err := yaml.Unmarshal([]byte(input), &data)
	return data, err
}

func (d *Data) convertToRuntimeValues(runtimePlaceholderCount *int) Data {
	for key, value := range *d {
		if subData, ok := value.(Data); ok {
			subData.convertToRuntimeValues(runtimePlaceholderCount)
		} else {
			(*d)[key] = RuntimeValue{
				value: fmt.Sprintf("%v", value),
				updateUsedValues: func() {
					*runtimePlaceholderCount++
				},
			}
		}
	}
	return *d
}

func (d *Data) merge(d2 Data) Data {
	for k, v := range d2 {
		if (*d)[k] == nil {
			(*d)[k] = v
		} else {
			v1, ok1 := (*d)[k].(Data)
			v2, ok2 := v.(Data)
			if ok1 && ok2 {
				(*d)[k] = v1.merge(v2)
			}
		}
	}
	return *d
}

func (d Data) String() string {
	bytes, err := json.Marshal(d)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
