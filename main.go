package main

// "golang.zx2c4.com/wireguard/wgctrl"
// "github.com/valyala/fasthttp"
// "github.com/valyala/quicktemplate"
// "github.com/goccy/go-yaml"
// "os"
// "fmt"

var Conf Config
var Db DB

func main() {
	Log(LogInfo, "starting wrauth")
	ParseConfig(&Conf)
	ParseDB(&Db)
	Log(LogDebug, Conf)
	Log(LogDebug, Db)
}
