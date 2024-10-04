package main

import (
	"testing"

	"github.com/valyala/fasthttp"
)

// use GOMAXPROCS for threads
// no. connections is at the end of each function
func BenchmarkFhttpAuth30(b *testing.B) {
	c := &fasthttp.Client{
		Name:                          "wrauth/0.1.0",
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}

	b.SetParallelism(30)

	b.RunParallel(func(p *testing.PB) {
		req := fasthttp.AcquireRequest()
		req.SetRequestURI("http://127.0.0.1:9091/api/authz/auth-request")
		req.Header.Add("X-Original-Method", "GET")
		req.Header.Add("X-Original-URL", "https://home.skaytacium.com")
		req.Header.Add("X-Forwarded-For", "10.0.0.31")

		res := fasthttp.AcquireResponse()

		for p.Next() {
			err := c.Do(req, res)
			if err != nil {
				b.Fatal(err)
			}
		}

		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)

		return
	})
}
