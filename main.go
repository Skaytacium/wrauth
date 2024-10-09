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
	Theme:   "gruvbox-dark",
	Authelia: struct {
		Address     string
		Db          string
		Connections int
	}{Connections: 64},
}
var Db DB
var WGs []*wgtypes.Device
var Matches []Match
var Log *zap.SugaredLogger
var Conns chan gnet.Conn
var C *gnet.Client

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

	if err := Store(); err != nil {
		Log.Fatalln(err)
	}
	SetDefaults()

	fswatch, err := fsnotify.NewWatcher()
	if err != nil {
		Log.Fatalf("couldn't create new watcher: %v", err)
	}
	defer fswatch.Close()

	go WatchFS(fswatch)

	// Authelia configs aren't monitored for a reason: Authelia itself doesn't monitor them...
	if err = fswatch.Add(filepath.Dir(Args.Config)); err != nil {
		Log.Fatalf("couldn't watch wrauth's directory: %v", err)
	}

	wgclient, err := wgctrl.New()
	if err != nil {
		Log.Fatalf("couldn't start WireGuard client: %v", err)
	}
	for _, inf := range Conf.Interfaces {
		dev, err := wgclient.Device(inf.Name)
		if err != nil {
			Log.Fatalf("couldn't get device %v: %v", inf.Name, err)
		}

		WGs = append(WGs, dev)
	}
	defer wgclient.Close()

	if err = UpdateCache(); err != nil {
		Log.Fatalf("error while caching rules: %v", err)
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
		Log.Fatalf("error while creating client: %v", err)
	}
	if err := C.Start(); err != nil {
		Log.Fatalf("error while starting client: %w", err)
	}
	Conns = make(chan gnet.Conn, Conf.Authelia.Connections)
	if err = CreateConnections(C); err != nil {
		Log.Fatalln(err)
	}

	go func() {
		tick := time.NewTicker(20 * time.Second)

		for {
			<-tick.C
			Log.Debugln("pinging Authelia connections")

			connections := make([]gnet.Conn, Conf.Authelia.Connections)
			for i := 0; i < Conf.Authelia.Connections; i++ {
				connections[i] = <-Conns
			}
			for _, c := range connections {
				if c != nil {
					if err = PingConnection(c); err != nil {
						Log.Errorln(err)
					}
					<-c.Context().(SubReq).notif
					Conns <- c
				}
			}
		}
	}()

	Log.Debugln(BFind([]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 10))

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
		Log.Fatalf("error while creating server: %v", err)
	}
}
