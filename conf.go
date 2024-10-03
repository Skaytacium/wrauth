package main

import (
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type Rule struct {
	Ips    []string
	Pubkey string
	User   string
}

type Admin struct {
	Ips []string
	// embedding rule doesn't work for yaml
	Pubkey string
	User   string
	Group  string
}

type Header struct {
	Domain  string
	Subject []string
	Headers []map[string]string
}

type User struct {
	Disabled    bool
	DisplayName string
	Email       string
	Groups      []string
}

type DB struct {
	Rules  []Rule
	Admins []Admin
	Data   []Header
	Users  map[string]User
}

type Interface struct {
	Name   string
	Addr   string
	Conf   string
	Watch  int
	Subnet string
	Shake  int
}

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

type Network struct {
	Name     string
	Networks []string
}

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

func Parse[T any](file *T, path string) {
	data, err := os.ReadFile(path)

	if err != nil {
		Log(LogFatal, err)
	}

	if err := yaml.Unmarshal(data, file); err != nil {
		Log(LogFatal, yaml.FormatError(err, true, true))
	}
}

func AddDefaults() {
	for i, inf := range Conf.Interfaces {
		if inf.Conf == "" {
			Conf.Interfaces[i].Conf = "/etc/wireguard/" + inf.Name
		}
		if inf.Watch == 0 {
			Conf.Interfaces[i].Watch = 15
		}
		if inf.Subnet == "" {
			Conf.Interfaces[i].Subnet = strings.Join(strings.Split(inf.Addr, ".")[0:3], ".") + ".0" + "/24"
		}
		if inf.Shake == 0 {
			Conf.Interfaces[i].Shake = 150
		}
	}
}
