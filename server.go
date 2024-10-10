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
	req, res, ask := HTAuthReq{}, make([]byte, 2048), true
	n := copy(res, "HTTP/1.1 403 Forbidden\r\n")

	// reqs will be max 1kB, TCP buffer should be able to handle that
	data, err := c.Next(-1)
	if err != nil {
		Log.Errorf("server: reading request: %v", err)
	}

	FastHTAuthReqParse(data, &req)

	for _, m := range Matches {
		if CompareUIP(&req.XRemote, &m.Ip) {
			n = FastHTAuthResGen(res, &m, HT200, req.XURL)
			ask = false
			break
		}
	}
	if _, ok := Cache[CacheHash(req.XURL, req.Cookie)]; ok {
		n = FastHTAuthResGen(res, nil, HT403, req.XURL)
		ask = false
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
		// ~10-15us
		notif := make(chan int)
		// ~40-50 us
		cc := <-Conns
		// ~12.5us
		cc.SetContext(SubReq{
			data:  res,
			notif: notif,
		})

		// ~15us
		_, err = cc.Write(data)
		if err != nil {
			Log.Errorf("server: writing Authelia subrequest: %v", err)
		}

		// ~225us
		// not including the TCP connection initiation time.
		n = <-notif
		// ~30-50us
		Conns <- cc

		// this does mean that a new entry will be created for each
		// Authelia request, and yeah that's an edge case that would
		// take some programming to account for, but if that edge
		// case is happening, there are bigger problems to fix in the setup.
		//
		// also, don't cache on 401s, these are meant to be asked to
		// the Authelia server (so that it can actually perform authentication)
		// ASCII 1
		if res[11] != 0x31 {
			// Cache[CacheHash(req.XURL, req.Cookie)] = HTStat(res[11] - 0x30)
		}
	}

	if _, err = c.Write(res[:n]); err != nil {
		Log.Errorf("server: writing TCP buffer: %v", err)
	}
	return gnet.None
}
