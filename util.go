package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func Compare(a []byte, b []byte) bool {
	ret := true

	if len(a) != len(b) {
		return false
	} else {
		for i := range a {
			ret = (a[i] == b[i])
		}
	}

	return ret
}

// no clue why generics are needed here, but its a rare operation
func ParseYaml[T any](file *T, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error while opening file %v: %w", path, err)
	}

	if err := yaml.Unmarshal(data, file); err != nil {
		return fmt.Errorf("error while parsing file %v: %w", path, err)
	}

	return nil
}
