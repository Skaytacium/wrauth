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

// >2x faster than net.ParseIP
func FastParseIP(ip []byte, dst *[4]byte) {
	var rarr = [3]byte{1, 10, 100}
	var tmp byte
	var i, rad byte = 0, 2

	for _, v := range ip {
		// ASCII .
		if v == 0x2E {
			dst[i] = tmp / rarr[rad+1]
			rad = 2
			tmp = 0
			i++
			continue
		}
		tmp += (v - 0x30) * rarr[rad]
		rad--
	}

	dst[3] = tmp / rarr[rad+1]
}

// ~4x faster than net.ParseCIDR
// probably because it has very basic error checking
func FastParseCIDR(ip []byte, dst *IP) error {
	var rarr = [3]byte{1, 10, 100}
	var tmp byte
	var i, rad byte = 0, 2
	var set bool = false

	for _, v := range ip {
		// ASCII /
		if v == 0x2F {
			dst.Addr[3] = tmp / rarr[rad+1]
			if !set {
				rad = 1
				tmp = 0
				set = true
				continue
			}
			// ASCII .
		} else if v == 0x2E {
			dst.Addr[i] = tmp / rarr[rad+1]
			rad = 2
			tmp = 0
			i++
			continue
		}
		// ASCII 0
		tmp += (v - 0x30) * rarr[rad]
		rad--
	}
	dst.Mask = uint32(0xffffffff << (32 - (tmp / rarr[rad+1])))

	if !set {
		return fmt.Errorf("data not in CIDRformat")
	} else {
		return nil
	}
}
