package cmd

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

type Data map[string]interface{}

func ParseInputData(input []string) (Data, error) {
	if len(input) == 0 {
		return Data{}, nil
	}
	data, err := parseData(input[0])
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(input); i++ {
		data2, err := parseData(input[i])
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
