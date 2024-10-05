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

// ~1.8x faster than net.ParseIP
// >2x faster without safety
func FastIP(data []byte, addr *[4]byte) error {
	var rarr = [3]byte{1, 10, 100}
	var tmp, rad, n, set byte = 0, 0, 3, 0

	for i := len(data) - 1; i > -1; i-- {
		// ASCII .
		if data[i] == 0x2E {
			addr[n] = tmp
			tmp = 0
			rad = 0
			n--
			set++
			continue
		}
		// ASCII 0
		tmp += (data[i] - 0x30) * rarr[rad]
		rad++
	}

	addr[0] = tmp

	if set != 3 {
		return fmt.Errorf("address not in proper format")
	} else {
		return nil
	}
}

// 2-4ns slower
func FastUIP(data []byte, addr *uint32) error {
	var rarr = [3]byte{1, 10, 100}
	var rad, n, set byte
	var tmp uint32

	*addr = 0
	for i := len(data) - 1; i > -1; i-- {
		// ASCII .
		if data[i] == 0x2E {
			*addr |= tmp << n
			tmp = 0
			rad = 0
			n += 8
			set++
			continue
		}
		// ASCII 0
		tmp += uint32((data[i] - 0x30) * rarr[rad])
		rad++
	}

	*addr |= tmp << 24

	if set != 3 {
		return fmt.Errorf("address not in proper format")
	} else {
		return nil
	}
}

// ~7x faster than net.ParseCIDR
// ~8.5x faster without safety
func FastCIDR(data []byte, addr *[4]byte, mask *uint32) error {
	var rarr = [3]byte{1, 10, 100}
	var tmp, rad, n, set byte = 0, 0, 3, 0

	for i := len(data) - 1; i > -1; i-- {
		// ASCII /
		if data[i] == 0x2F {
			*mask = uint32(0xffffffff << (32 - tmp))
			tmp = 0
			rad = 0
			set++
			continue
			// ASCII .
		} else if data[i] == 0x2E {
			addr[n] = tmp
			tmp = 0
			rad = 0
			n--
			set++
			continue
		}
		// ASCII 0
		tmp += (data[i] - 0x30) * rarr[rad]
		rad++
	}

	addr[0] = tmp

	if set != 4 {
		return fmt.Errorf("address not in CIDR format")
	} else {
		return nil
	}
}

// 2-4ns slower
func FastUCIDR(data []byte, addr *uint32, mask *uint32) error {
	var rarr = [3]byte{1, 10, 100}
	var rad, n, set byte
	var tmp uint32

	*addr = 0
	for i := len(data) - 1; i > -1; i-- {
		// ASCII /
		if data[i] == 0x2F {
			*mask = 0xffffffff << (32 - tmp)
			rad = 0
			tmp = 0
			set++
			continue
			// ASCII .
		} else if data[i] == 0x2E {
			*addr |= tmp << n
			tmp = 0
			rad = 0
			n += 8
			set++
			continue
		}
		// ASCII 0
		tmp += uint32((data[i] - 0x30) * rarr[rad])
		rad++
	}

	*addr |= tmp << 24

	if set != 4 {
		return fmt.Errorf("address not in CIDR format")
	} else {
		return nil
	}
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

func CompareUIP(a, b IP) bool {
	return (a.Addr^b.Addr)&b.Mask == 0
}

func Find[T any](a *[]T, c func(a T) bool) *T {
	for _, x := range *a {
		if c(x) {
			return &x
		}
	}
	return nil
}
