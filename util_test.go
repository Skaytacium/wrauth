package main

import (
	"net"
	"testing"
)

func BenchmarkNetParseCIDR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		net.ParseCIDR("129.168.255.235/32")
	}
}

func BenchmarkFastUCIDR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ip, mask uint32
		FastUCIDR([]byte("129.168.255.235/32"), &ip, &mask)
	}
}

func BenchmarkNetParseIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		net.ParseIP("129.168.255.235")
	}
}

func BenchmarkFastUIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ip uint32
		FastUIP([]byte("129.168.255.235"), &ip)
	}
}

func BenchmarkNetCompIP(b *testing.B) {
	a, B := net.IPv4(129, 168, 255, 235), net.IPv4(255, 255, 255, 255)
	for i := 0; i < b.N; i++ {
		if net.IP.Equal(a, B) {

		}
	}
}

func BenchmarkUCompIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if CompareUIP(&IP{
			Addr: 0xf1f2f3f4,
			Mask: 0xffffffff,
		}, &IP{
			Addr: 0xf1f2f3f4,
			Mask: 0xffffff00,
		}) {

		}
	}
}

func BenchmarkHTParse(b *testing.B) {
	req := []byte("GET /auth HTTP/1.1\r\nHost: 127.0.0.1:9092\r\nUser-Agent: curl/8.10.1\r\nAccept: */*\r\nX-Forwarded-For: 10.0.0.32\r\nX-Original-Method: GET\r\nX-Original-URL: https://home.skaytacium.com\r\n\r\n")
	parse := HTReq{}

	for i := 0; i < b.N; i++ {
		FastHTParse(req, &parse)
	}
}
