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
	err := FastUCIDR(Sanitize(data), &ip.Addr, &ip.Mask)
	if err != nil {
		return err
	}
	return nil
}

type Match struct {
	User
	Ip IP
	Id string
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
	HT200 HTStat = iota
	HT401
	HT302
	HT403
	HT404
)

var HTStatName = [5]string{"200 OK", "401 Unauthorized", "302 Found", "403 Forbidden", "404 Not Found"}

type HTAuthReq struct {
	Method HTMethod
	Path   []byte
	// only these 4 headers matter
	XRemote IP
	// not set to HTMethod because it'll just be passed
	XMethod []byte
	XURL    []byte
	Cookie  []byte
}

func (h HTAuthReq) String() string {
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

// type HTAuthRes struct {
// 	Stat    HTStat
// 	Id		string
// 	User 	User
// }

type SubReq struct {
	data  []byte
	notif chan int
}
