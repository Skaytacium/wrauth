package main

import (
	"github.com/valyala/fasthttp"
)

func Listen() {
	srv := &fasthttp.Server{
		Logger:         Logger,
		Name:           "wrauth/0.1.0",
		Handler:        handle,
		HeaderReceived: check,
		ErrorHandler: func(_ *fasthttp.RequestCtx, err error) {
			Log(LogError, "server error: %v", err)
		},
	}

	if err := srv.ListenAndServe(Conf.Address); err != nil {
		Log(LogFatal, "couldn't start server: %v", err)
	}
}

func check(header *fasthttp.RequestHeader) fasthttp.RequestConfig {
	if CompareSlice(header.RequestURI(), []byte("/auth")) {
		return fasthttp.RequestConfig{
			// 0 sets it to default, so 1 it is
			MaxRequestBodySize: 1,
		}
	}

	return fasthttp.RequestConfig{}
}

func handle(ctx *fasthttp.RequestCtx) {
	if recip := ctx.Request.Header.Peek("X-Forwarded-For"); len(recip) > 0 {
		var addr uint32
		FastUIP(recip, &addr)

		ctx.SetStatusCode(403)

		for _, ip := range Matches {
			if CompareUIP(&IP{Addr: addr, Mask: 0xffffffff}, &ip.Ip) {
				ctx.SetStatusCode(200)
			}
		}

	} else {
		Log(LogError, "no X-Forwarded-For header found in auth request")
		ctx.SetStatusCode(400)
	}
}
