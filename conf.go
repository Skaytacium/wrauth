package main

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func StoreFiles() error {
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
	// if Conf.External == "" {
	// 	return fmt.Errorf("external address not configured")
	// }
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

	for _, inf := range Conf.Interfaces {
		if inf.Name == "" {
			return fmt.Errorf("interface name not configured")
		}
		if inf.Addr.Mask == 0 {
			return fmt.Errorf("address not configured for interface %v", inf.Name)
		}
	}

	return nil
}

func CheckDB() error {
	for _, r := range Db.Rules {
		if r.User == "" {
			return fmt.Errorf("user not configured")
		}
		if len(r.Ips) == 0 && len(r.Pubkeys) == 0 {
			return fmt.Errorf("IPs or public keys not configured for %v", r.User)
		}
	}
	for _, d := range Db.Headers {
		if len(d.Urls) == 0 {
			return fmt.Errorf("URLs not configured")
		}
		if len(d.Subjects) == 0 {
			return fmt.Errorf("subjects not configured")
		}
		for _, s := range d.Subjects {
			if s.User == "" && s.Group == "" {
				return fmt.Errorf("neither subject user nor group configured")
			}
		}
		if len(d.Headers) == 0 {
			return fmt.Errorf("headers not configured")
		}
	}
	if len(Db.Admins) == 0 {
		return fmt.Errorf("admins not configured")
	}
	for _, a := range Db.Admins {
		if a.User == "" && a.Group == "" {
			return fmt.Errorf("neither admin user nor group configured")
		}
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
			Log.Fatalf("file watch: %v", err)
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			if ev.Op == fsnotify.Write && (strings.Contains(ev.Name, Args.Config) || strings.Contains(ev.Name, Args.DB)) {
				if send {
					Log.Infof("file updated: %v", ev.Name)
					if err := StoreFiles(); err != nil {
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
					HeaderCache = make(map[string]Header)
					if err := AddHeaders(); err != nil {
						Log.Fatalf("headers: %v", err)
					}
				}
				send = !send
			}
		}
	}
}
