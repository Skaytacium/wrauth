package main

import "fmt"

type IP struct {
	Addr uint32
	Mask uint32
}

func (ip IP) String() string {
	var bytes = To4Byte(ip.Addr)
	return fmt.Sprintf("%v.%v.%v.%v/%v", bytes[0], bytes[1], bytes[2], bytes[3], Bits(ip.Mask))
}

func (ip *IP) UnmarshalYAML(data []byte) error {
	err := FastUCIDR(sanitize(data), &ip.Addr, &ip.Mask)
	if err != nil {
		return err
	}
	return nil
}

type Match struct {
	User
	Ip   IP
	Name string
}

func (m Match) String() string {
	return fmt.Sprintf("Match %v on %v", m.Name, m.Ip)
}

func sanitize(data []byte) []byte {
	if data[0] == []byte("\"")[0] || data[0] == []byte("'")[0] {
		return data[1 : len(data)-1]
	}
	return data
}
