package main

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func ParseFiles() error {
	if err := ParseYaml(&Conf, Args.Config); err != nil {
		return fmt.Errorf("configuration: %w", err)
	}
	if err := ParseYaml(&Db, Args.DB); err != nil {
		return fmt.Errorf("database: %w", err)
	}
	if err := ParseYaml(&Db, Conf.Authelia.Db); err != nil {
		return fmt.Errorf("Authelia user database: %w", err)
	}

	return nil
}

func CheckConf() error {
	if Conf.External == "" {
		return fmt.Errorf("external address not configured")
	}
	if Conf.Authelia.Address == "" {
		return fmt.Errorf("Authelia address not configured")
	}
	if Conf.Authelia.Db == "" {
		return fmt.Errorf("Authelia user database path not configured")
	}
	if Conf.Authelia.Ping >= 30 {
		return fmt.Errorf("configured ping interval is too large")
	}
	if len(Conf.Interfaces) == 0 {
		return fmt.Errorf("interfaces not configured")
	}

	for i, inf := range Conf.Interfaces {
		if inf.Name == "" {
			return fmt.Errorf("interface name not configured")
		}
		if inf.Addr.Mask == 0 {
			return fmt.Errorf("address not configured for interface %v", inf.Name)
		}
		if inf.Conf == "" {
			Conf.Interfaces[i].Conf = "/etc/wireguard/" + inf.Name + ".conf"
		}
		if inf.Shake == 0 {
			Conf.Interfaces[i].Shake = 150
		}
	}

	return nil
}

func CheckDB() error {
	for _, r := range Db.Rules {
		if r.User == "" {
			return fmt.Errorf("rules: user not configured")
		}
		if len(r.Ips) == 0 && len(r.Pubkeys) == 0 {
			return fmt.Errorf("rules: IPs or public keys not configured for %v", r.User)
		}
	}
	for _, a := range Db.Access {
		if len(a.Domains) == 0 {
			return fmt.Errorf("access: domains not configured")
		}
		if len(a.Users) == 0 && len(a.Groups) == 0 {
			return fmt.Errorf("access: neither users nor groups configured")
		}
	}
	if len(Db.Admins.Users) == 0 && len(Db.Admins.Groups) == 0 {
		return fmt.Errorf("admins not configured")
	}

	return nil
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
			Log.Fatalln("file watch: ", err)
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			if ev.Op == fsnotify.Write && (strings.Contains(ev.Name, Args.Config) || strings.Contains(ev.Name, Args.DB)) {
				if send {
					Log.Infoln("file updated: ", ev.Name)
					if err := ParseFiles(); err != nil {
						Log.Errorf("parsing: %v", err)
					}
					if err := CheckConf(); err != nil {
						Log.Errorf("configuration: %v", err)
					}
					if err := CheckDB(); err != nil {
						Log.Errorf("database: %v", err)
					}
					Matches = nil
					if err := AddMatches(); err != nil {
						Log.Errorf("matching: %v", err)
					}
				}
				send = !send
			}
		}
	}
}
