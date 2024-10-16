package main

import (
	"flag"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.zx2c4.com/wireguard/wgctrl"
)

var Args struct {
	Config string
	DB     string
}

var Conf Config
var Db DB
var Matches []Match
var WGClient *wgctrl.Client

// on some revelations, maps are the fastest way to do
// anything out here and theyre equally safe

// host -> cookie -> id
var AuthCache struct {
	sync.RWMutex
	cache map[string]map[string]string
}

// host -> user -> headers
var Cache map[string]map[string][]byte

// host -> user -> regexp
var Regexps map[string]map[string]*regexp.Regexp

var Log *zap.SugaredLogger
var C *gnet.Client
var Conns chan gnet.Conn

func main() {
	Clear()

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

	if err := LoadFiles(); err != nil {
		Log.Fatalln("loading files:", err)
	}

	fswatch, err := fsnotify.NewWatcher()
	if err != nil {
		Log.Fatalln("filesystem watcher creation:", err)
	}
	defer fswatch.Close()

	go WatchFS(fswatch)

	Log.Debugln("watching directory:", filepath.Dir(Args.Config))
	if err = fswatch.Add(filepath.Dir(Args.Config)); err != nil {
		Log.Fatalln("filesytem watch:", err)
	}

	WGClient, err = wgctrl.New()
	if err != nil {
		Log.Fatalln("WireGuard client creation:", err)
	}
	defer WGClient.Close()

	if err := LoadData(); err != nil {
		Log.Fatalln("processing data:", err)
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
		Log.Fatalln("TCP client creation:", err)
	}
	if err := C.Start(); err != nil {
		Log.Fatalln("TCP client starting:", err)
	}

	Conns = make(chan gnet.Conn, Conf.Authelia.Connections)
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
		Log.Fatalln("server creation:", err)
	}
}
