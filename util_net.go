package main

import (
	"fmt"
	"net"
)

// most of these are highly optimized prematurely, not for
// actual speed gains, but because i was having fun

// ~1.8x faster than net.ParseIP
// >2x faster without safety
func ParseUIP(data []byte, addr *uint32) error {
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
func ParseUCIDR(data []byte, addr *uint32, mask *uint32) error {
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

// ~15x faster than net.IP.Equal
// uses the 2nd IP's mask
func CompareUIP(a, b *IP) bool {
	return (a.Addr^b.Addr)&b.Mask == 0
}

func ConvIP(ip net.IPNet) IP {
	return IP{
		Addr: ToUint32([4]byte(ip.IP)),
		Mask: ToUint32([4]byte(ip.Mask)),
	}
}

func GetHost(url []byte) []byte {
	// skip https:// (8)   ASCII /
	if i := LFind(url[8:], 0x2f); i != 0xffffffff {
		return url[8 : i+8]
	} else {
		return url[8:]
	}
}

func MatchGlobURL(glob, url string) bool {
	// ASCII *
	if glob[0] == 0x2A {
		if glob[2:] == url[LFind([]byte(url), 0x2E)+1:] {
			return true
		}
	} else if glob == url {
		return true
	}
	return false
}

// HTTP parsers, faster than O(n)
// you'll need this...
// https://www.utf8-chartable.de/unicode-utf8-table.pl?names=2
func HTAuthReqParse(data []byte, h *HTAuthReq) {
	// current index, previous index, headers received
	n, p, c := 1, 1, 0

	switch data[n] {
	case 0x45:
		// default method
		// h.Method = HTGet
		n = 4
	case 0x4f:
		h.Method = HTPost
		n = 5
	case 0x55:
		h.Method = HTPut
		n = 5
	}
	p = n
	n += LFind(data[n:], 0x20)
	h.Path = data[p:n]

	// skip _HTTP/1.1\r\n (11)
	n += 11
	p = n

	// atrocious
	for data[n] != 0x0d {
		switch data[n] {
		case 0x58:
			// skip till farthest diff character, then till starting
			// X-Forwarded-For: -> : +2
			// X-Original-Method: -> o +4
			// X-Original-URL: -> _ +1
			n += 15
			switch data[n] {
			case 0x3a:
				n += 2
				p = n
				n += LFind(data[n:], 0x0d)
				ParseUIP(data[p:n], &h.XRemote.Addr)
				// all received IPs are /32 by default
				h.XRemote.Mask = 0xffffffff
				c++
			case 0x6f:
				n += 4
				p = n
				n += LFind(data[n:], 0x0d)
				h.XMethod = data[p:n]
				c++
			case 0x20:
				n++
				p = n
				n += LFind(data[n:], 0x0d)
				h.XURL = data[p:n]
				c++
			default:
				n += LFind(data[n:], 0x0d)
			}
		case 0x43:
			// skip till starting
			// ookie:_ (8)
			n += 8
			p = n
			n += LFind(data[n:], 0x0d)
			h.Cookie = data[p:n]
			c++
		default:
			n += LFind(data[n:], 0x0d)
		}
		// skip \r\n
		n += 2
		// no need to parse further
		if c == 4 {
			break
		}
	}
}

func HTAuthResParse(data []byte, h *HTAuthRes) {
	h.Stat = HTStat(data[11] - 0x30)
	if h.Stat != HT200 {
		return
	}

	n, p := 17, 17

	for data[n] != 0x0d {
		switch data[n] {
		case 0x52:
			// skip Remote- (7)
			n += 7
			switch data[n] {
			case 0x55:
				// skip User:_ (6)
				n += 6
				p = n
				n += LFind(data[n:], 0x0d)
				h.Id = data[p:n]
				break
			default:
				n += LFind(data[n:], 0x0d)
			}
		default:
			n += LFind(data[n:], 0x0d)
		}
		n += 2
	}
}

func HTAuthResGen(res []byte, id string, user *User, h HTStat) int {
	n := copy(res, "HTTP/1.1 ")
	n += copy(res[n:], HTStatName[h])
	n += copy(res[n:], "\r\n")

	if h == HT200 {
		var i int

		n += copy(res[n:], "Remote-User: ")
		n += copy(res[n:], id)
		n += copy(res[n:], "\r\n")

		if len(user.Groups) > 0 {
			n += copy(res[n:], "Remote-Groups: ")
			for i = 0; i < len(user.Groups)-1; i++ {
				n += copy(res[n:], user.Groups[i])
				n += copy(res[n:], ",")
			}
			n += copy(res[n:], user.Groups[i])
			n += copy(res[n:], "\r\n")
		}

		n += copy(res[n:], "Remote-Name: ")
		n += copy(res[n:], user.DisplayName)
		n += copy(res[n:], "\r\n")

		n += copy(res[n:], "Remote-Email: ")
		n += copy(res[n:], user.Email)
		n += copy(res[n:], "\r\n")
	}

	n += copy(res[n:], "Content-Length: 0\r\n")

	return n
}

func AddHeaders(h map[string]string) []byte {
	t, n := make([]byte, 2048), 0
	for k, v := range h {
		n += copy(t[n:], k)
		n += copy(t[n:], []byte(": "))
		n += copy(t[n:], v)
		n += copy(t[n:], []byte("\r\n"))
	}
	return t[:n]
}
