package main

import "github.com/panjf2000/gnet/v2"

type EvHandler struct {
	gnet.BuiltinEventEngine
}

func (ev *EvHandler) OnTraffic(c gnet.Conn) gnet.Action {
	req := HTReq{}
	// reqs will be max 1kB, TCP buffer should be able to handle that
	data, err := c.Next(-1)
	if err != nil {
		Log.Errorf("error while reading request: %v", err)
	}
	Log.Debugln(string(data))
	FastHTParse(data, &req)
	Log.Debugln(req)
	return gnet.Close
}
