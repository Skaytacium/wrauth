package main

import (
	"regexp"

	"go.uber.org/zap"
)

// YAML
// # db
// ## rules
type Rule struct {
	Ips     []IP
	Pubkeys []string
	User    string
}

// ## indentification
type Identity struct {
	Users  []string
	Groups [][]string
}

// ## access control rules
type Access struct {
	Identity `yaml:",inline"`
	Domains  []string
	Resource regexp.Regexp
	Headers  map[string]string
}

// ## Authelia users
type User struct {
	Disabled    bool
	DisplayName string
	Email       string
	Groups      []string
}

type DB struct {
	Rules  []Rule
	Admins Identity
	Access []Access
	Users  map[string]User
}

// # wrauth
// ## Authelia config
type Authelia struct {
	Address     string
	Db          string
	Connections int
	Cache       int
	Ping        int
}

// ## WG interfaces
type Interface struct {
	Name  string
	Addr  IP
	Conf  string
	Shake int
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
