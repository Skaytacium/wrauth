package main

// "golang.zx2c4.com/wireguard/wgctrl"
// "github.com/valyala/fasthttp"
// "github.com/valyala/quicktemplate"
// "github.com/goccy/go-yaml"
// "os"
// "fmt"

var Conf = Config{
	Address: "127.0.0.1:9091",
	Log:     LogInfo,
	Theme:   "gruvbox-dark",
}
var Db DB
var Authelia AutheliaConfiguration

func main() {
	Parse(&Conf, "./config.yaml")
	Parse(&Db, "./db.yaml")
	Parse(&Authelia, Conf.Authelia.Config)
	Parse(&Db, Authelia.Authentication_backend.File.Path)
	AddDefaults()
	Log(LogDebug, Conf)
	Log(LogDebug, Authelia)
	Log(LogDebug, Db)
}
