package main

// no clue why this is so inefficient, but it's only called on startup
func CachePubkeys() {
	keys := []string{}

	for _, v := range Db.Rules {
		for _, k := range v.Pubkeys {
			keys = append(keys, k)
		}
	}

	for _, wg := range WGs {
		for _, p := range wg.Peers {
			for _, k := range keys {
				if k == p.PublicKey.String() {
					for _, ip := range p.AllowedIPs {
						ip = ip
						// Matches = append(Matches, ip)
					}
				}
			}
		}
	}
}
