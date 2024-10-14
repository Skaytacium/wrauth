package main

import (
	"flag"
	"path/filepath"
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
var WGs []*wgtypes.Device

// on some revelations, maps are the fastest way to do
// anything out here and theyre equally safe

// host -> cookie -> id
var AuthCache map[string]map[string]string

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

	Log.Infoln("parsing files")
	if err := ParseFiles(); err != nil {
		Log.Fatalf("parsing: %v", err)
	}
	Log.Infoln("checking configuration")
	if err := CheckConf(); err != nil {
		Log.Fatalf("configuration: %v", err)
	}
	Log.Infoln("checking database")
	if err := CheckDB(); err != nil {
		Log.Fatalf("database: %v", err)
	}
	Log.Infoln("caching access control and headers")
	Cache = make(map[string]map[string][]byte)
	if err := AddCache(); err != nil {
		Log.Fatalf("caching: %v", err)
	}

	if !Conf.Caching {
		Log.Warnln("caching is disabled for all Authelia requests")
	}

	fswatch, err := fsnotify.NewWatcher()
	if err != nil {
		Log.Fatalf("filesystem watcher creation: %v", err)
	}
	defer fswatch.Close()

	go WatchFS(fswatch)

	Log.Infoln("watching configuration directory")
	if err = fswatch.Add(filepath.Dir(Args.Config)); err != nil {
		Log.Fatalf("filesytem watch: %v", err)
	}

	wgclient, err := wgctrl.New()
	if err != nil {
		Log.Fatalf("WireGuard client creation: %v", err)
	}
	defer wgclient.Close()

	Log.Infoln("adding WireGuard interfaces")
	for _, inf := range Conf.Interfaces {
		dev, err := wgclient.Device(inf.Name)
		if err != nil {
			Log.Fatalf("WireGuard device %v: %v", inf.Name, err)
		}
		if dev.Type != wgtypes.LinuxKernel {
			Log.Warnf("wrauth is using userspace WireGuard device %v", inf.Name)
		}

		WGs = append(WGs, dev)
	}

	Log.Infoln("matching IPs to users")
	// needs WireGuard setup
	if err := AddMatches(); err != nil {
		Log.Fatalf("matching: %v", err)
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
		Log.Fatalf("TCP client creation: %v", err)
	}
	if err := C.Start(); err != nil {
		Log.Fatalf("TCP client starting: %w", err)
	}
	Conns = make(chan gnet.Conn, Conf.Authelia.Connections)
	Log.Infoln("creating Authelia connections")
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

	AuthCache = make(map[string]map[string]string)
	go func() {
		tick := time.NewTicker(time.Duration(Conf.Authelia.Cache) * time.Second)
		for {
			<-tick.C
			Log.Debugln("clearing Authelia cache")
			AuthCache = make(map[string]map[string]string)
		}
	}()

	Log.Infoln("running server")
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
		Log.Fatalf("server creation: %v", err)
	}
}
