package main

import (
	"fmt"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// only depend on url and cookie, since each user has a different cookie
// and Authelia shouldn't care about the network if its asking for
// further authentication.
func CacheHash(url []byte, cookie []byte) uint64 {
	var hash uint64

	// start after https:// (8)
	hash |= uint64(ToUint([4]byte(url[7:11]))) << 32
	hash |= uint64(ToUint([4]byte(url[11:15])))
	// start after authelia_session= (17)
	hash ^= uint64(ToUint([4]byte(cookie[16:20]))) << 32
	hash ^= uint64(ToUint([4]byte(cookie[20:24])))

	return hash
}

func DataHash(sub []Identity, id []byte, url []byte) uint64 {
	var hash uint64

	for _, s := range sub {
		// order is important
		
	}

}

func addMatch(ip IP, name string) error {
	if Find(&Matches, func(a Match) bool {
		return CompareUIP(&ip, &a.Ip)
	}) == nil {
		user, ok := Db.Users[name]
		if !ok {
			return fmt.Errorf("user %v not in Authelia database", name)
		}

		Matches = append(Matches, Match{
			User: user,
			Ip:   ip,
			Id:   name,
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
			for _, wg := range WGs {
				for _, ip := range Find(&wg.Peers, func(a wgtypes.Peer) bool {
					return a.PublicKey.String() == k
				}).AllowedIPs {
					if err := addMatch(IP{Addr: ToUint([4]byte(ip.IP)), Mask: ToUint([4]byte(ip.Mask))}, v.User); err != nil {
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

func AddData() error {
	for _, d := range Db.Headers {
		for _, u := range d.Urls {
			for _, ss := range d.Subjects {
				Log.Debugln()
				for _, 
			}
		}
	}
}
