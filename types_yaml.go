package main

import "go.uber.org/zap"

// YAML
// # db
// ## rules
type Rule struct {
	Ips     []IP
	Pubkeys []string
	User    string
}

// ## admins
type Identity struct {
	User   string
	Groups []string
}

// ## custom headers
type Headers struct {
	Domains  []string
	Subjects []Identity
	Headers  map[string]string
}

// ## Authelia users
type User struct {
	Disabled    bool
	DisplayName string
	Email       string
	Groups      []string
}

// ## final struct
type DB struct {
	Rules   []Rule
	Admins  []Identity
	Headers []Headers
	Users   map[string]User
}

// # wrauth
// ## WG interfaces
type Interface struct {
	Name string
	Addr IP
	// Conf   string
	// Watch  int
	// Subnet IP
	// Shake  int
}

// ## Authelia config
type Authelia struct {
	Address     string
	Db          string
	Connections int
	Cache       int
	Ping        int
}

// ## final struct
type Config struct {
	Address    string
	External   string
	Caching    bool
	Level      zap.AtomicLevel
	Theme      string
	Authelia   Authelia
	Interfaces []Interface
}
