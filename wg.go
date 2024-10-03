package main

import (
	"fmt"

	"golang.zx2c4.com/wireguard/wgctrl"
)

func CacheWG(client *wgctrl.Client) {
	for _, inf := range Conf.Interfaces {
		dev, err := client.Device(inf.Name)
		if err != nil {
			Log(LogFatal, fmt.Sprintf("WireGuard device %v does not exist", inf.Name))
		}
		WGCache = append(WGCache, dev)
	}
}
