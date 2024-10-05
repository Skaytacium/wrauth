package main

import "fmt"

type IP struct {
	Addr uint32
	Mask uint32
}

func (ip IP) String() string {
	var bytes = To4Byte(ip.Addr)
	return fmt.Sprintf("%v.%v.%v.%v/%v", bytes[0], bytes[1], bytes[2], bytes[3], Bits(ip.Mask))
}

type Match struct {
	User
	Ip   IP
	Name string
}

// YAML
// # db
// ## rules
type Rule struct {
	Ips     []IP
	Pubkeys []string
	User    string
}

// ## admins
type Admin struct {
	Ip IP
	// no clue why embedding doesn't work for yaml Unmarshal, but it's just more lines
	Pubkey string
	User   string
	Group  string
}

// ## custom headers
type Header struct {
	Domain  string
	Subject []string
	Headers []map[string]string
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
	Rules  []Rule
	Admins []Admin
	Data   []Header
	Users  map[string]User
}

// # wrauth
// ## WG interfaces
type Interface struct {
	Name   string
	Addr   IP
	Conf   string
	Watch  int
	Subnet IP
	Shake  int
}

// ## final struct
type Config struct {
	Address  string
	External string
	Log      LogLevel
	Theme    string
	Authelia struct {
		Config string
	}
	Interfaces []Interface
}

// # Authelia
// ## networks
type Network struct {
	Name     string
	Networks []IP
}

// ## ugly, but "necessary" final struct
type AutheliaConfiguration struct {
	Server struct {
		Address string
	}
	Authentication_backend struct {
		File struct {
			Path string
		}
	}
	Access_control struct {
		Networks []Network
	}
	Session struct {
		Cookies []struct {
			Domain string
		}
	}
}
