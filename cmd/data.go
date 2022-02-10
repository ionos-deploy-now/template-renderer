package cmd

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

type IntermediateValue struct {
	value            string
	updateUsedValues func(newValue string)
}

func (v IntermediateValue) String() string {
	v.updateUsedValues(v.value)
	return v.value
}

func (v IntermediateValue) MarshalJSON() ([]byte, error) {
	return []byte("\"" + v.String() + "\""), nil
}

type Data map[string]interface{}

func ParseInputData(input []string, intermediateInput []string, usedValues *[]string) (Data, error) {
	data := Data{}
	if len(input) > 0 {
		for i := 0; i < len(input); i++ {
			data2, err := parseData(input[i])
			if err != nil {
				return nil, err
			}
			data = data.merge(data2)
		}
	}
	if len(intermediateInput) > 0 {
		for i := 0; i < len(intermediateInput); i++ {
			data2, err := parseData(intermediateInput[i])
			data2 = data2.convertToIntermediateValues(usedValues)
			if err != nil {
				return nil, err
			}
			data = data.merge(data2)
		}
	}
	return data, nil
}

func parseData(input string) (Data, error) {
	var data Data
	err := yaml.Unmarshal([]byte(input), &data)
	return data, err
}

func (d *Data) convertToIntermediateValues(usedValues *[]string) Data {
	for key, value := range *d {
		if subData, ok := value.(Data); ok {
			subData.convertToIntermediateValues(usedValues)
		} else {
			(*d)[key] = IntermediateValue{
				value: fmt.Sprintf("%v", value),
				updateUsedValues: func(newValue string) {
					for _, value := range *usedValues {
						if value == newValue {
							return
						}
					}
					*usedValues = append(*usedValues, newValue)
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
