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
	Log(LogDebug, "%+v", header.RequestURI())

	if Compare(header.RequestURI(), []byte("/auth")) {
		if len(header.Peek("X-Forwarded-For")) > 0 {
			Log(LogDebug, "hi")
		}

		return fasthttp.RequestConfig{
			// 0 sets it to default, so 1 it is
			MaxRequestBodySize: 1,
		}
	}

	return fasthttp.RequestConfig{}
}

func handle(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(200)
}
