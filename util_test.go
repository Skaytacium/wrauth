package main

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
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

func BenchmarkFastHTAuthReqParse(b *testing.B) {
	req := []byte("GET /auth HTTP/1.1\r\nHost: 127.0.0.1:9092\r\nUser-Agent: curl/8.10.1\r\nAccept: */*\r\nX-Forwarded-For: 10.0.0.32\r\nX-Original-Method: GET\r\nX-Original-URL: https://home.skaytacium.com\r\n\r\n")
	parse := HTAuthReq{}

	for i := 0; i < b.N; i++ {
		FastHTAuthReqParse(req, &parse)
	}
}

func BenchmarkFastHTAuthResGen(b *testing.B) {
	req, m := make([]byte, 2048), Match{
		Ip: IP{
			Addr: 0xf0f0f0f0,
			Mask: 0xffffffff,
		},
		Id: "test",
		User: User{
			Disabled:    false,
			DisplayName: "Test",
			Email:       "test@mail",
			Groups:      []string{"gtest", "gtest2"},
		},
	}

	for i := 0; i < b.N; i++ {
		FastHTAuthResGen(req, &m, HT200, []byte("hello"))
	}
}

func BenchmarkConcat(b *testing.B) {
	var str string
	for n := 0; n < b.N; n++ {
		str += "x"
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); str != s {
		b.Errorf("unexpected result; got=%s, want=%s", str, s)
	}
}

func BenchmarkBuffer(b *testing.B) {
	var buffer bytes.Buffer
	for n := 0; n < b.N; n++ {
		buffer.WriteString("x")
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); buffer.String() != s {
		b.Errorf("unexpected result; got=%s, want=%s", buffer.String(), s)
	}
}

func BenchmarkCopy(b *testing.B) {
	bs := make([]byte, b.N)
	bl := 0

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bl += copy(bs[bl:], "x")
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); string(bs) != s {
		b.Errorf("unexpected result; got=%s, want=%s", string(bs), s)
	}
}

// Go 1.10
func BenchmarkStringBuilder(b *testing.B) {
	var strBuilder strings.Builder

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		strBuilder.WriteString("x")
	}
	b.StopTimer()

	if s := strings.Repeat("x", b.N); strBuilder.String() != s {
		b.Errorf("unexpected result; got=%s, want=%s", strBuilder.String(), s)
	}
}

func BenchmarkTimeDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now() //.Format(time.RFC1123)
	}
}

func BenchmarkMap(b *testing.B) {
	var cache = map[uint64]bool{
		1: true,
		2: true,
		3: true,
		4: true,
		5: true,
		6: true,
		7: true,
		8: true,
	}

	for i := 0; i < b.N; i++ {
		_ = cache[1]
	}
}

func BenchmarkFFind(b *testing.B) {
	var cache = []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
	}

	for i := 0; i < b.N; i++ {
		FFind(cache, 8)
	}
}
