package main

import (
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/fsnotify/fsnotify"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var Args struct {
	Config string `arg:"-c,--config" help:"location of wrauth configuration" default:"./config.yaml"`
	DB     string `arg:"-d,--databse" help:"location of wrauth database" default:"./db.yaml"`
}

var Conf Config
var Db DB
var Authelia AutheliaConfiguration
var WGs []*wgtypes.Device

func Store() {
	if err := Parse(&Conf, Args.Config); err != nil {
		Log(LogFatal, "error while parsing main config: %v", err)
	}
	if err := Parse(&Db, Args.DB); err != nil {
		Log(LogFatal, "error while parsing database config: %v", err)
	}
	if err := Parse(&Authelia, Conf.Authelia.Config); err != nil {
		Log(LogFatal, "error while parsing Authelia configuration: %v", err)
	}
	if err := Parse(&Db, Authelia.Authentication_backend.File.Path); err != nil {
		Log(LogFatal, "error while parsing Authelia user database: %v", err)
	}

	Authelia.Server.Address = strings.Replace(Authelia.Server.Address, "tcp4", "http", 1)

	Log(LogDebug, "%+v", Conf)
	Log(LogDebug, "%+v", Authelia)
	Log(LogDebug, "%+v", Db)
}

func main() {
	arg.MustParse(&Args)
	Store()
	AddDefaults()

	wgclient, err := wgctrl.New()
	if err != nil {
		Log(LogFatal, "couldn't obtain WireGuard devices: %v", err)
	}
	defer wgclient.Close()

	LoadWGs(wgclient)

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

	Log(LogInfo, "wrauth ready")

	<-make(chan bool)
}
