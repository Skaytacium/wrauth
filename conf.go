package main

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func (ip *IP) UnmarshalYAML(data []byte) error {

	// ugh
	if data[0] == []byte("\"")[0] || data[0] == []byte("'")[0] {
		data = data[1 : len(data)-1]
	}

	err := FastUCIDR(data, &ip.Addr, &ip.Mask)

	if err != nil {
		return err
	}

	return nil
}

func Store() error {
	if err := ParseYaml(&Conf, Args.Config); err != nil {
		return fmt.Errorf("error while parsing configuration: %w", err)
	}
	if err := ParseYaml(&Db, Args.DB); err != nil {
		return fmt.Errorf("error while parsing database: %w", err)
	}
	if err := ParseYaml(&Authelia, Conf.Authelia.Config); err != nil {
		return fmt.Errorf("error while parsing Authelia configuration: %w", err)
	}
	if err := ParseYaml(&Db, Authelia.Authentication_backend.File.Path); err != nil {
		return fmt.Errorf("error while parsing Authelia users: %w", err)
	}

	Authelia.Server.Address = strings.Replace(Authelia.Server.Address, "tcp4", "http", 1)

	if len(Conf.Interfaces) == 0 {
		return fmt.Errorf("no interfaces configured")
	}
	if len(Db.Admins) == 0 {
		return fmt.Errorf("no admins configured")
	}

	return nil
}

func AddDefaults() {
	for i, inf := range Conf.Interfaces {
		if inf.Conf == "" {
			Conf.Interfaces[i].Conf = "/etc/wireguard/" + inf.Name
		}
		if inf.Watch == 0 {
			Conf.Interfaces[i].Watch = 15
		}
		if inf.Subnet.Mask == 0 {
			Conf.Interfaces[i].Subnet = IP{
				// default subnet mask of 24
				Addr: inf.Addr.Addr & 0xffffff00,
				Mask: 0xffffff00,
			}
		}
		if inf.Shake == 0 {
			Conf.Interfaces[i].Shake = 150
		}
	}
}

func WatchConfigs(w *fsnotify.Watcher) {
	// no clue why write requests happen twice, but simple fix
	send := true

	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			Log(LogFatal, "error watching files %v", err)
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			if ev.Op == fsnotify.Write && (strings.Contains(ev.Name, Args.Config) || strings.Contains(ev.Name, Args.DB)) {
				if send {
					Log(LogInfo, "file updated: %v", ev.Name)
					if err := Store(); err != nil {
						Log(LogError, "error while parsing: %v", err)
					}
					Matches = nil
					if err := UpdateCache(); err != nil {
						Log(LogError, "error while caching: %v", err)
					}
				}
				send = !send
			}
		}
	}
}
