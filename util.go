package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func CompareSlice[T comparable](a []T, b []T) bool {
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

func Find[T any](a *[]T, c func(a T) bool) *T {
	for _, x := range *a {
		if c(x) {
			return &x
		}
	}
	return nil
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

func ToUint(data [4]byte) uint32 {
	var tmp uint32

	tmp |= uint32(data[0]) << 24
	tmp |= uint32(data[1]) << 16
	tmp |= uint32(data[2]) << 8
	tmp |= uint32(data[3])

	return tmp
}

func To4Byte(data uint32) [4]byte {
	var tmp [4]byte

	tmp[0] = byte(data >> 24)
	tmp[1] = byte(data >> 16)
	tmp[2] = byte(data >> 8)
	tmp[3] = byte(data)

	return tmp
}

func Bits(data uint32) byte {
	var n byte

	for data != 0 {
		n += byte(data & 1)
		data >>= 1
	}

	return n
}

func FFind(data []byte, query byte) int {
	for i, v := range data {
		if query == v {
			return i
		}
	}
	return 0
}
