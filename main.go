package main

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var Conf = Config{
	Address: "127.0.0.1:9091",
	Log:     LogInfo,
	Theme:   "gruvbox-dark",
}
var Db DB
var Authelia AutheliaConfiguration
var WGCache []*wgtypes.Device

func main() {
	Parse(&Conf, "./config.yaml")
	Parse(&Db, "./db.yaml")
	Parse(&Authelia, Conf.Authelia.Config)
	Parse(&Db, Authelia.Authentication_backend.File.Path)
	AddDefaults()

	wgclient, err := wgctrl.New()
	defer wgclient.Close()
	if err != nil {
		Log(LogFatal, "couldn't obtain WireGuard devices")
	}

	CacheWG(wgclient)

	Log(LogDebug, WGCache)
	Log(LogDebug, Conf)
	Log(LogDebug, Authelia)
	Log(LogDebug, Db)
}
