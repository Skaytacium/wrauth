package main

import (
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
var Matches []IP

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
		Log(LogFatal, "couldn't start WireGuard client: %v", err)
	}
	for _, inf := range Conf.Interfaces {
		dev, err := wgclient.Device(inf.Name)
		if err != nil {
			Log(LogFatal, "couldn't get device %v: %v", inf.Name, err)
		}

		WGs = append(WGs, dev)
	}
	defer wgclient.Close()

	CachePubkeys()

	Log(LogDebug, "%+v", Matches)

	Log(LogInfo, "listening on: %v", Conf.Address)
	Listen()
}
