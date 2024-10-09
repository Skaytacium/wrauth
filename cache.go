package main

import (
	"fmt"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func cache(ip IP, name string) error {
	if Find(&Matches, func(a Match) bool {
		return CompareUIP(&ip, &a.Ip)
	}) == nil {
		user, ok := Db.Users[name]
		if !ok {
			return fmt.Errorf("user %v not found in Authelia database", name)
		}

		Matches = append(Matches, Match{
			User: user,
			Ip:   ip,
			Name: name,
		})
	} else {
		return fmt.Errorf("rule for IP %v has duplicates", ip)
	}
	return nil
}

// no clue why this is so inefficient O(n(n + n^4)), but it's only called on file update
// and yeah it's inefficient but it's not slow, n is usually quite small
func UpdateCache() error {
	for _, v := range Db.Rules {
		for _, k := range v.Pubkeys {
			for _, wg := range WGs {
				for _, ip := range Find(&wg.Peers, func(a wgtypes.Peer) bool {
					return a.PublicKey.String() == k
				}).AllowedIPs {
					if err := cache(IP{Addr: ToUint([4]byte(ip.IP)), Mask: ToUint([4]byte(ip.Mask))}, v.User); err != nil {
						return err
					}
				}
			}
		}
		for _, i := range v.Ips {
			if err := cache(i, v.User); err != nil {
				return err
			}
		}
	}
	return nil
}
