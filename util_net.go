package main

import "fmt"

// ~1.8x faster than net.ParseIP
// >2x faster without safety
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

// ~15x faster than net.IP.Equal
func CompareUIP(a, b *IP) bool {
	return (a.Addr^b.Addr)&b.Mask == 0
}

// HTTP parsers, faster than O(n)
// you'll need this...
// https://www.utf8-chartable.de/unicode-utf8-table.pl?names=2
func FastHTReqParse(data []byte, h *HTReq) {
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
	n += FFind(data[n:], 0x20)
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
				n += FFind(data[n:], 0x0d)
				FastUIP(data[p:n], &h.XRemote.Addr)
				// all received IPs are /32 by default
				h.XRemote.Mask = 0xffffffff
				c++
			case 0x6f:
				n += 4
				p = n
				n += FFind(data[n:], 0x0d)
				h.XMethod = data[p:n]
				c++
			case 0x20:
				n++
				p = n
				n += FFind(data[n:], 0x0d)
				h.XURL = data[p:n]
				c++
			default:
				n += FFind(data[n:], 0x0d)
			}
		case 0x43:
			// skip till starting
			// ookie:_ (8)
			n += 8
			p = n
			n += FFind(data[n:], 0x0d)
			h.Cookie = data[p:n]
			c++
		default:
			n += FFind(data[n:], 0x0d)
		}
		// skip \r\n
		n += 2
		// no need to parse further
		if c == 4 {
			break
		}
	}
}
