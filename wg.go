package main

import (
	"fmt"

	"golang.zx2c4.com/wireguard/wgctrl"
)

func LoadWGs(client *wgctrl.Client) error {
	for _, inf := range Conf.Interfaces {
		dev, err := client.Device(inf.Name)
		if err != nil {
			return fmt.Errorf("error while getting WireGuard device: %w", err)
		}
		WGs = append(WGs, dev)
	}

	return nil
}
