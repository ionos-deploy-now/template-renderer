package types

import (
	"fmt"
	"strings"
)

type Path string

func (p Path) Append(paths ...interface{}) Path {
	values := []string{string(p)}
	for _, value := range paths {
		if value != "" {
			values = append(values, fmt.Sprintf("%s", value))
		}
	}
	return Path(strings.Join(values, "/"))
}

func (p Path) String() string {
	return string(p)
}
