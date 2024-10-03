package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Rule struct {
	Ips    []string
	Pubkey string
	User   string
}

type Admin struct {
	Rule
	Group string
}

type Header struct {
	Domain  string
	Subject []string
	Headers [][2]string
}

type DB struct {
	Rules  []Rule
	Data   []Header
	Admins []Admin
}

type InterfaceConf struct {
	Name   string
	Addr   string
	Conf   string
	Watch  int
	Subnet string
	Shake  int
}

type AutheliaConf struct {
	Config string
	Userdb string
	Login  string
}

type Config struct {
	Address    string
	External   string
	Log        LogLevel
	Theme      string
	Authelia   AutheliaConf
	Interfaces []InterfaceConf
}

func ParseConfig(conf *Config) {
	data, err := os.ReadFile("./config.yaml")

	if err != nil {
		Log(LogFatal, err)
	}

	if err := yaml.Unmarshal(data, conf); err != nil {
		Log(LogFatal, yaml.FormatError(err, true, true))
	}
}

func ParseDB(db *DB) {
	data, err := os.ReadFile("./db.yaml")

	if err != nil {
		Log(LogFatal, err)
	}

	if err := yaml.Unmarshal(data, db); err != nil {
		Log(LogFatal, yaml.FormatError(err, true, true))
	}
}
