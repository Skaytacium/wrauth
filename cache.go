package main

import (
	"fmt"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	// "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// no clue why this is so inefficient O(n(n + n^4)), but it's only called on file update
func UpdateCache() error {
	for _, v := range Db.Rules {
		for _, k := range v.Pubkeys {
			for _, wg := range WGs {
				for _, ip := range Find(&wg.Peers, func(a wgtypes.Peer) bool {
					return a.PublicKey.String() == k
				}).AllowedIPs {
					store := IP{Addr: ToUint([4]byte(ip.IP)), Mask: ToUint([4]byte(ip.Mask))}

					if Find(&Matches, func(a IP) bool {
						return CompareUIP(store, a)
					}) == nil {
						Matches = append(Matches, store)
					} else {
						return fmt.Errorf("rule for public key %v (%v) has duplicates", k, store)
					}
				}
			}
		}
		for _, i := range v.Ips {
			if Find(&Matches, func(a IP) bool {
				return CompareUIP(i, a)
			}) == nil {
				Matches = append(Matches, i)
			} else {
				return fmt.Errorf("rule for IP %v has duplicates", i)
			}
		}
	}
	return nil
}
