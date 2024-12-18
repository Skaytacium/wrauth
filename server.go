package main

import (
	"path/filepath"

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
	req, res, id := HTAuthReq{}, make([]byte, 2048), ""
	n := HTAuthResGen(res, "", nil, HT403)
	var m *Match

	// reqs will be max 1kB, TCP buffer should be able to handle that
	data, err := c.Next(-1)
	if err != nil {
		Log.Errorln("server: reading request:", err)
	}

	if err := HTAuthReqParse(data, &req); err != nil {
		Log.Warnln("server: parsing request:", err)
	}
	// used a lot, convenience
	reqdom, valc := UFStr(GetHost(req.XURL)), len(req.Cookie) >= 49 && UFStr(req.Cookie[:17]) == "authelia_session="
	if len(req.XURL) < 12 {
		Log.Errorln("server: requested domain isn't valid")
		Log.Debugln(UFStr(data))
		goto response
	}

	Log.Debugln("IP:", req.XRemote)
	Log.Debugln("domain:", reqdom)
	Log.Debugln("resource:", UFStr(GetResource(req.XURL)))
	Log.Debugln("cookie:", UFStr(req.Cookie))

	m = CFind(&Matches, func(m Match) bool {
		return CompareUIP(&req.XRemote, &m.Ip)
	})
	if m != nil {
		Log.Debugln("IP matched user:", m.Id)
		_, allowed := Cache[reqdom][m.Id]
		r, regex := Regexps[reqdom][m.Id]
		if !allowed {
			Log.Debugln("no direct matches, checking bypasses")
			_, allowed = Cache[reqdom]["*"]
			r, regex = Regexps[reqdom]["*"]
		}
		if !allowed {
			Log.Debugln("no bypasses, checking globs")
			for u, sub := range Cache {
				if g, err := filepath.Match(u, reqdom); g && err == nil {
					_, allowed = sub[m.Id]
					r, regex = Regexps[u][m.Id]
					if allowed {
						break
					}
				}
			}
		}
		if allowed && regex {
			Log.Debugln("regex match required:", r.String())
			allowed = r.Match(GetResource(req.XURL))
		}
		Log.Debugln("allowed in domain:", allowed)
		if allowed {
			user := Db.Users[m.Id]
			n = HTAuthResGen(res, m.Id, &user, HT200)
			id = m.Id
			// skip cache check
			goto headers
		}
		goto response
	}
	Log.Debugln("IP didn't match any rules, checking cache")
	if Conf.Caching && valc {
		AuthCache.RLock()
		if cid, ok := AuthCache.cache[reqdom][UFStr(req.Cookie[17:17+32])]; ok {
			Log.Debugln("cached as user:", cid)
			if cid != "" {
				user := Db.Users[cid]
				n = HTAuthResGen(res, cid, &user, HT200)
				id = cid
			}
			AuthCache.RUnlock()
			goto headers
		}
		// i hate writing this again
		AuthCache.RUnlock()
		Log.Debugln("cache missed")
	}

	// it takes ~300/330us for the entire request,
	// minimum time it could take is 200us, if somehow
	// there was no overhead in proxying.
	// the issue is, 320us is >1.5x the actual time,
	// so the speed is also <1/1.5x.
	// except for shaving off maybe 40-50us by implementing
	// an alternative to go channels (batshit crazy), there
	// seems to be no way to make this any faster.
	// the best way to optimize something is to not do it in
	// the first place, so what we WILL use is a cache.
	//
	// if true is used for allowing goto statements
	if true {
		Log.Debugln("subrequesting Authelia")
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
			Log.Errorln("server: writing Authelia subrequest:", err)
		}

		// ~225us
		// not including the TCP connection initiation time.
		n = <-notif
		// ~30-50us
		Conns <- cc

		subres := HTAuthRes{}
		HTAuthResParse(res, &subres)
		Log.Debugln("subrequest status:", subres.Stat)
		// this does mean that a new entry will be created for each
		// Authelia request, and yeah that's an edge case that would
		// take some programming to account for, but if that edge
		// case is happening, there are bigger problems to fix in the setup.
		//
		// also, don't cache on 401s, these are meant to be asked to
		// the Authelia server (so that it can actually perform authentication)
		// ASCII 1
		if Conf.Caching && subres.Stat != HT401 && valc {
			// convert using string here, so as to copy the byte slice
			// instead of reusing, since it's dependent on `res`, which
			// doesn't get released only reused
			id = string(subres.Id)
			Log.Debugln("subrequest user:", id)

			AuthCache.RLock()
			sub, ok := AuthCache.cache[reqdom]
			AuthCache.RUnlock()

			Log.Debugln("caching subrequest response")
			if !ok {
				sub = make(map[string]string)
			}
			sub[UFStr(req.Cookie[17:32+17])] = id
			AuthCache.Lock()
			AuthCache.cache[reqdom] = sub
			AuthCache.Unlock()

			user := Db.Users[UFStr(subres.Id)]
			n = HTAuthResGen(res, UFStr(subres.Id), &user, subres.Stat)
			goto headers
		}
		Log.Debugln("not caching subrequest")
		goto response
	}
headers:
	if h := Cache[reqdom][id]; len(h) > 0 {
		Log.Debugln("custom headers found")
		n += copy(res[n:], h)
	}

response:
	n += copy(res[n:], "\r\n")
	if _, err = c.Write(res[:n]); err != nil {
		Log.Errorln("server: writing TCP buffer:", err)
	}
	return gnet.None
}
