package main

import (
	"net"
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

func BenchmarkNetParseCIDR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		net.ParseCIDR("129.168.1.23/32")
	}
}

func BenchmarkFastParseCIDR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ip IP
		FastParseCIDR([]byte("129.168.1.23/32"), &ip)
	}
}

func BenchmarkNetParseIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		net.ParseIP("129.168.1.23")
	}
}

func BenchmarkFastParseIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ip [4]byte
		FastParseIP([]byte("129.168.1.23"), &ip)
	}
}
