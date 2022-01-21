package cmd

import (
	"gopkg.in/yaml.v3"
)

type Data map[string]interface{}

func ParseInputData(input []string) Data {
	if len(input) == 0 {
		return Data{}
	}
	data := parseData(input[0])
	for i := 1; i < len(input); i++ {
		data = data.merge(parseData(input[i]))
	}
	return data
}

func parseData(input string) Data {
	var data Data
	handleError(yaml.Unmarshal([]byte(input), &data))
	return data
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
