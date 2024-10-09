package main

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// i hate this function so much
func Store() error {
	if err := ParseYaml(&Conf, Args.Config); err != nil {
		return fmt.Errorf("error while parsing configuration: %w", err)
	}
	if err := ParseYaml(&Db, Args.DB); err != nil {
		return fmt.Errorf("error while parsing database: %w", err)
	}
	if err := ParseYaml(&Db, Conf.Authelia.Db); err != nil {
		return fmt.Errorf("error while parsing Authelia users: %w", err)
	}

	if len(Conf.Interfaces) == 0 {
		return fmt.Errorf("no interfaces configured")
	}
	if len(Db.Admins) == 0 {
		return fmt.Errorf("no admins configured")
	}

	return nil
}

func SetDefaults() {
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

func WatchFS(w *fsnotify.Watcher) {
	// no clue why write requests happen twice, but simple fix
	send := true

	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			Log.Fatalf("error watching files %v", err)
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			if ev.Op == fsnotify.Write && (strings.Contains(ev.Name, Args.Config) || strings.Contains(ev.Name, Args.DB)) {
				if send {
					Log.Infof("file updated: %v", ev.Name)
					if err := Store(); err != nil {
						Log.Errorln(err)
					}
					Matches = nil
					if err := AddMatches(); err != nil {
						Log.Errorf("error while caching rules: %v", err)
					}
				}
				send = !send
			}
		}
	}
}
