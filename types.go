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

type HTMethod int

const (
	HTGet HTMethod = iota
	HTPost
	HTPut
)

func (m HTMethod) String() string {
	switch m {
	case HTGet:
		return "GET"
	case HTPost:
		return "POST"
	case HTPut:
		return "PUT"
	}
	return ""
}

type HTStat int

const (
	HT200 HTStat = 0
	HT401        = 1
	HT403        = 3
	HT404        = 4
)

func (s HTStat) String() string {
	switch s {
	case HT200:
		return "200 OK"
	case HT401:
		return "401 Unauthorized"
	case HT403:
		return "403 Forbidden"
	case HT404:
		return "404 Not Found"
	}
	return ""
}

type HTReq struct {
	Method HTMethod
	Path   []byte
	// only these 4 headers matter
	XRemote IP
	// not set to HTMethod because it'll just be passed
	XMethod []byte
	XURL    []byte
	Cookie  []byte
}

func (h HTReq) String() string {
	return fmt.Sprintf(
		"%v %v from %v with\nMethod: %v\nURL: %v\nCookie: %v",
		h.Method,
		string(h.Path),
		h.XRemote,
		string(h.XMethod),
		string(h.XURL),
		string(h.Cookie),
	)
}

type HTRes struct {
	Stat    HTStat
	CLength int
	CType   []byte
	// only headers required from Authelia response
	Location []byte
	SCookie  []byte
	// Date string
	Content []byte
}

type SubReq struct {
	data  []byte
	notif chan int
}
