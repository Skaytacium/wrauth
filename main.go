package main

import (
	"flag"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var Args struct {
	Config string
	DB     string
}

var Conf = Config{
	Address: "127.0.0.1:9092",
	Level:   zap.NewAtomicLevel(),
	Caching: true,
	Theme:   "gruvbox-dark",
	Authelia: Authelia{
		Connections: 64,
		Cache:       300,
		Ping:        25,
	},
}
var Db DB

var Matches []Match
var WGs []struct {
	*wgtypes.Device
	data Interface
}

// on some revelations, maps are the fastest way to do
// anything out here and theyre equally safe

// host -> cookie -> id
var AuthCache struct {
	sync.RWMutex
	cache map[string]map[string]string
}

// host -> user -> headers
var Cache map[string]map[string][]byte

var Log *zap.SugaredLogger
var C *gnet.Client
var Conns chan gnet.Conn

func main() {
	Log = zap.Must(zap.Config{
		Level:    Conf.Level,
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:       "msg",
			LevelKey:         "lvl",
			TimeKey:          "time",
			ConsoleSeparator: ": ",
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeLevel:      zapcore.LowercaseColorLevelEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()).Sugar()

	defer Log.Sync()

	Args.Config = *flag.String("config", "./config.yaml", "location of the configuration file")
	Args.DB = *flag.String("db", "./db.yaml", "location of the database file")
	flag.Parse()

	Log.Debugln("parsing files")
	if err := ParseFiles(); err != nil {
		Log.Fatalln("parsing: ", err)
	}
	Log.Debugln("checking configuration")
	if err := CheckConf(); err != nil {
		Log.Fatalln("configuration: ", err)
	}
	Log.Debugln("checking database")
	if err := CheckDB(); err != nil {
		Log.Fatalln("database: ", err)
	}
	Log.Debugln("caching access control and headers")
	Cache = make(map[string]map[string][]byte)
	if err := AddCache(); err != nil {
		Log.Fatalln("caching: ", err)
	}

	if !Conf.Caching {
		Log.Warnln("caching is disabled for all Authelia requests")
	}

	fswatch, err := fsnotify.NewWatcher()
	if err != nil {
		Log.Fatalln("filesystem watcher creation: ", err)
	}
	defer fswatch.Close()

	go WatchFS(fswatch)

	Log.Debugln("watching configuration directory")
	if err = fswatch.Add(filepath.Dir(Args.Config)); err != nil {
		Log.Fatalln("filesytem watch: ", err)
	}

	wgclient, err := wgctrl.New()
	if err != nil {
		Log.Fatalln("WireGuard client creation: ", err)
	}
	defer wgclient.Close()

	Log.Debugln("adding WireGuard interfaces")
	for _, inf := range Conf.Interfaces {
		dev, err := wgclient.Device(inf.Name)
		if err != nil {
			Log.Fatalln("WireGuard device %v: ", inf.Name, err)
		}
		if dev.Type != wgtypes.LinuxKernel {
			Log.Warnf("wrauth is using userspace WireGuard device %v", inf.Name)
		}

		WGs = append(WGs, struct {
			*wgtypes.Device
			data Interface
		}{dev, inf})
	}

	Log.Debugln("matching IPs to users")
	// needs WireGuard setup
	if err := AddMatches(); err != nil {
		Log.Fatalln("matching: ", err)
	}

	C, err = gnet.NewClient(
		&CHandler{},
		gnet.WithEdgeTriggeredIO(true),
		gnet.WithMulticore(true),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
		// Authelia doesn't care
		// gnet.WithTCPKeepAlive(time.Second *15),
		gnet.WithLogger(Log),
		gnet.WithLogLevel(Conf.Level.Level()),
	)
	if err != nil {
		Log.Fatalln("TCP client creation: ", err)
	}
	if err := C.Start(); err != nil {
		Log.Fatalln("TCP client starting: ", err)
	}
	Conns = make(chan gnet.Conn, Conf.Authelia.Connections)
	Log.Debugln("creating Authelia connections")
	if err = CreateConnections(C); err != nil {
		Log.Fatalln(err)
	}

	go func() {
		tick := time.NewTicker(time.Duration(Conf.Authelia.Ping) * time.Second)
		for {
			<-tick.C
			Log.Debugln("pinging Authelia connections")

			connections := make([]gnet.Conn, Conf.Authelia.Connections)
			for i := 0; i < Conf.Authelia.Connections; i++ {
				connections[i] = <-Conns
			}
			for _, c := range connections {
				if err = PingConnection(c); err != nil {
					Log.Errorln(err)
				}
				<-c.Context().(SubReq).notif
				Conns <- c
			}
		}
	}()

	AuthCache.cache = make(map[string]map[string]string)
	go func() {
		tick := time.NewTicker(time.Duration(Conf.Authelia.Cache) * time.Second)
		for {
			<-tick.C
			Log.Debugln("clearing Authelia cache")
			AuthCache.Lock()
			AuthCache.cache = make(map[string]map[string]string)
			AuthCache.Unlock()
		}
	}()

	Log.Debugln("running server")
	if err = gnet.Run(
		&SHandler{},
		"tcp4://"+Conf.Address,
		gnet.WithEdgeTriggeredIO(true),
		gnet.WithMulticore(true),
		gnet.WithReusePort(true),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
		gnet.WithLogger(Log),
		gnet.WithLogLevel(Conf.Level.Level()),
	); err != nil {
		Log.Fatalln("server creation: ", err)
	}
}
