package cmd

import (
	"log"
	"strings"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func joinPath(s ...string) string {
	return strings.Join(s, "/")
}
