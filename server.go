package main

import "github.com/panjf2000/gnet/v2"

type EvHandler struct {
	gnet.BuiltinEventEngine
}

func (ev *EvHandler) OnTraffic(c gnet.Conn) gnet.Action {
	Log.Debug(c.RemoteAddr())
	return gnet.Close
}
