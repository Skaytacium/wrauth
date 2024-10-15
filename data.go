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
// 			Log.Fatalln("cookie not from Authelia:", string(cookie))
// 		}
// 		Log.Fatalln("cookie not long enough:", string(cookie))
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

// no clue why these are so inefficient O(n(n + n^4)), but they're only
// called on file update. and yeah it's inefficient but it's not too slow,
// since n is usually quite small
func AddMatches() error {
	for _, v := range Db.Rules {
		for _, k := range v.Pubkeys {
			for _, wg := range WGInfs {
				for _, ip := range CFind(&wg.Peers, func(a wgtypes.Peer) bool {
					return a.PublicKey.String() == k
				}).AllowedIPs {
					if err := addMatch(ConvIP(ip), v.User); err != nil {
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

func AddCache() error {
	for _, a := range Db.Access {
		for _, d := range a.Domains {
			if _, ok := Cache[d]; !ok {
				Cache[d] = make(map[string][]byte)
			}
			if len(a.Identity.Users) == 1 && a.Identity.Users[0] == "*" {
				Cache[d]["*"] = append(Cache[d]["*"], AddHeaders(a.Headers)...)
				continue
			}
			for u := range Db.Users {
				if UserIn(u, a.Identity) {
					Cache[d][u] = append(Cache[d][u], AddHeaders(a.Headers)...)
				}
			}
		}
	}
	return nil
}

/*
func AddPeers(wg *wgtypes.Device) {
	for _, p := range wg.Peers {
		for _, ip := range p.AllowedIPs {
			pip := IP{Addr: ToUint32([4]byte(ip.IP)), Mask: ToUint32([4]byte(ip.Mask))}
			if CFind(&WGPeers, func(wip IP) bool {
				return CompareUIP(&pip, &wip)
			}) == nil {
				WGPeers = append(WGPeers, pip)
			}
		}
	}
}
*/
