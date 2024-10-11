package main

import (
	"fmt"
	"os"
	"unsafe"

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

func CFind[T any](a *[]T, c func(a T) bool) *T {
	for _, x := range *a {
		if c(x) {
			return &x
		}
	}
	return nil
}

func LFind[T byte](data []T, query T) int {
	for i, v := range data {
		if query == v {
			return i
		}
	}
	return 0xffffffff
}

// func BFindU64(data []uint64, query uint64) int {
// 	l, h, m := 0, len(data), 0
// 	for l <= h {
// 		m = (l + h) >> 1
// 		if data[m] == query {
// 			return m
// 		}
// 		if data[m] > query {
// 			h = m - 1
// 		} else {
// 			l = m + 1
// 		}
// 	}
// 	return 0xffffffff
// }

// func BFindU256(data []uint256.Int, query *uint256.Int) int {
// 	l, h, m := 0, len(data), 0
// 	for l <= h {
// 		m = (l + h) >> 1
// 		if data[m].Eq(query) {
// 			return m
// 		}
// 		if data[m].Gt(query) {
// 			h = m - 1
// 		} else {
// 			l = m + 1
// 		}
// 	}
// 	return 0xffffffff
// }

// no clue why generics are needed here, but its a rare operation
func ParseYaml[T any](file *T, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error while opening file %v: %w", path, err)
	}

	if err := yaml.Unmarshal(data, file); err != nil {
		return fmt.Errorf("error while parsing file %v: %v", path, yaml.FormatError(err, true, true))
	}

	return nil
}

func Sanitize(data []byte) []byte {
	if data[0] == []byte("\"")[0] || data[0] == []byte("'")[0] {
		return data[1 : len(data)-1]
	}
	return data
}

func ToUint32(data [4]byte) uint32 {
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

func UFStr(data []byte) string {
	if i := len(data); i != 0 {
		return unsafe.String(&data[0], len(data))
	}
	return ""
}
