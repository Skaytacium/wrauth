package main

import (
	"fmt"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// only depend on url and cookie, since each user has a different cookie
// and Authelia shouldn't care about the network if its asking for
// further authentication.
// this is not a security risk since if somebody knows your Authlia
// cookie, it's over anyway
// someday these will be SIMD... (assembly i'm looking at you)
// func AuthHash(url []byte, cookie []byte, hash *uint256.Int) {
// 	if len(cookie) < 25 {
// 		if !CompareSlice(cookie[0:17], []byte("authelia_session")) {
// 			Log.Fatalf("cookie not from Authelia: %v", string(cookie))
// 		}
// 		Log.Fatalf("cookie not long enough: %v", string(cookie))
// 	}

// 	// start after authelia_session= (17) (256 bit cookie)         start after https:// (8)
// 	hash.Xor(uint256.NewInt(0).SetBytes32(cookie[17:49]), uint256.NewInt(0).SetBytes(url[8:]))
// }

// generate one hash per url
// 64 bits is enough for this (i don't expect
// this to be more than 50-60 values for the
// average homelab)
// https://en.wikipedia.org/wiki/Birthday_problem#Probability_table
// func IDHash(id Identity, url string) uint64 {
// 	var hash uint64

// 	// XOR is commutative and associative
// 	for i := 0; i < len(id.User) && i < 8; i++ {
// 		hash ^= uint64(id.User[i]) << (8 * i)
// 	}
// 	for _, g := range id.Groups {
// 		for i := 0; i < len(g) && i < 8; i++ {
// 			hash ^= uint64(g[i]) << (8 * i)
// 		}
// 	}
// 	for i := 0; i < len(url) && i < 8; i++ {
// 		hash ^= uint64(url[8+i]) << (8 * i)
// 	}

// 	return hash
// }

func addMatch(ip IP, name string) error {
	if CFind(&Matches, func(a Match) bool {
		return CompareUIP(&ip, &a.Ip)
	}) == nil {
		_, ok := Db.Users[name]
		if !ok {
			return fmt.Errorf("user %v not in Authelia database", name)
		}

		Matches = append(Matches, Match{
			Ip: ip,
			Id: name,
		})
	} else {
		return fmt.Errorf("rule for IP %v has duplicates", ip)
	}
	return nil
}

func convHeaders(h map[string]string) []byte {
	var HTTP []byte

	for k, v := range h {
		HTTP = append(HTTP, []byte(k)...)
		HTTP = append(HTTP, []byte(": ")...)
		HTTP = append(HTTP, []byte(v)...)
		HTTP = append(HTTP, []byte("\r\n")...)
	}

	return HTTP
}

// no clue why these are so inefficient O(n(n + n^4)), but they're only
// called on file update. and yeah it's inefficient but it's not too slow,
// since n is usually quite small
func AddMatches() error {
	for _, v := range Db.Rules {
		for _, k := range v.Pubkeys {
			for _, wg := range WGs {
				for _, ip := range CFind(&wg.Peers, func(a wgtypes.Peer) bool {
					return a.PublicKey.String() == k
				}).AllowedIPs {
					if err := addMatch(IP{Addr: ToUint32([4]byte(ip.IP)), Mask: ToUint32([4]byte(ip.Mask))}, v.User); err != nil {
						return err
					}
				}
			}
		}
		for _, i := range v.Ips {
			if err := addMatch(i, v.User); err != nil {
				return err
			}
		}
	}
	return nil
}

func AddHeaders() error {
	for _, d := range Db.Headers {
		for _, u := range d.Urls {
			for _, s := range d.Subjects {
				headerMap := HeaderCache[u]
				if s.User != "" {
					headerMap.User = map[string][]byte{}
					headerMap.User[s.User] = append(headerMap.User[s.User], convHeaders(d.Headers)...)
				}
				if s.Group != "" {
					headerMap.Group = map[string][]byte{}
					headerMap.Group[s.Group] = append(headerMap.Group[s.Group], convHeaders(d.Headers)...)
				}
				HeaderCache[u] = headerMap
			}
		}
	}

	Log.Debugf("%+v", HeaderCache)

	return nil
}
