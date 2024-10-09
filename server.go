package main

import (
	"github.com/panjf2000/gnet/v2"
)

type SHandler struct {
	gnet.BuiltinEventEngine
}

func (ev *SHandler) OnOpen(_ gnet.Conn) ([]byte, gnet.Action) {
	Log.Debugln("wrauth connection opened")
	return nil, gnet.None
}

func (ev *SHandler) OnClose(_ gnet.Conn, _ error) gnet.Action {
	Log.Debugln("wrauth connection closed")
	return gnet.Close
}

func (ev *SHandler) OnTraffic(c gnet.Conn) gnet.Action {
	req, res, ask := HTReq{}, make([]byte, 2048), true
	n := copy(res, "HTTP/1.1 403 Forbidden\r\n")

	// reqs will be max 1kB, TCP buffer should be able to handle that
	data, err := c.Next(-1)
	if err != nil {
		Log.Errorf("error while reading request: %v", err)
	}

	FastHTReqParse(data, &req)

	for _, m := range Matches {
		if CompareUIP(&req.XRemote, &m.Ip) {
			n = copy(res, "HTTP/1.1 200 OK\r\n")
			ask = false
			break
		}
	}

	// overall it takes ~300/330us for the entire request,
	// minimum time it could take is 200us, if somehow
	// there was no overhead in proxying.
	// the issue is, 320us is >1.5x the actual time,
	// so the speed is also <1/1.5x.
	// except for shaving off maybe 40-50us by implementing
	// an alternative to go channels (batshit crazy), there
	// seems to be no way to make this any faster.
	// the best way to optimize something is to not do it in
	// the first place, so what we WILL use is a cache.
	if ask {
		// ~40-50 us
		cc := <-Conns
		// ~10-15us
		notif := make(chan int)
		// ~12.5us
		cc.SetContext(SubReq{
			data:  res,
			notif: notif,
		})

		// ~15us
		_, err = cc.Write(data)
		if err != nil {
			Log.Errorf("error while subrequesting Authelia: %v", err)
		}

		// ~225us
		n = <-notif
		// ~30-50us
		Conns <- cc
	} else {
		n += copy(res[n:], "Content-Length: 0\r\n\r\n")
	}

	if _, err = c.Write(res[:n]); err != nil {
		Log.Errorf("error while writing to TCP buffer: %v", err)
	}
	return gnet.None
}
