package main

import (
	"net"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/fsnotify/fsnotify"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var Args struct {
	Config string `arg:"-c,--config" help:"location of wrauth configuration" default:"./config.yaml"`
	DB     string `arg:"-d,--databse" help:"location of wrauth database" default:"./db.yaml"`
}

var Conf = Config{
	Address: "127.0.0.1:9092",
	Log:     LogInfo,
	Theme:   "gruvbox-dark",
}
var Db DB
var Authelia AutheliaConfiguration
var WGs []*wgtypes.Device

// list of (subnet mask, ip address)
var Peers []net.IPNet

func main() {
	arg.MustParse(&Args)
	if err := Store(); err != nil {
		Log(LogFatal, "error while parsing: %v", err)
	}
	AddDefaults()

	fswatch, err := fsnotify.NewWatcher()
	if err != nil {
		Log(LogFatal, "couldn't create new watcher: %v", err)
	}
	defer fswatch.Close()

	go WatchConfigs(fswatch)

	// Authelia configs aren't monitored for a reason: Authelia itself doesn't monitor them...
	if err = fswatch.Add(filepath.Dir(Args.Config)); err != nil {
		Log(LogFatal, "couldn't watch wrauth's directory: %v", err)
	}

	wgclient, err := wgctrl.New()
	if err != nil {
		Log(LogFatal, "couldn't obtain WireGuard devices: %v", err)
	}
	defer wgclient.Close()

	Log(LogInfo, "listening on: %v", Conf.Address)
	Listen()
}
