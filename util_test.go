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

func BenchmarkFastHTReqParse(b *testing.B) {
	req := []byte("GET /auth HTTP/1.1\r\nHost: 127.0.0.1:9092\r\nUser-Agent: curl/8.10.1\r\nAccept: */*\r\nX-Forwarded-For: 10.0.0.32\r\nX-Original-Method: GET\r\nX-Original-URL: https://home.skaytacium.com\r\n\r\n")
	parse := HTReq{}

	for i := 0; i < b.N; i++ {
		FastHTReqParse(req, &parse)
	}
}

// func BenchmarkFastHTAuthResParse(b *testing.B) {
// 	req := []byte("HTTP/1.1 401 Unauthorized\r\nDate: Tue, 08 Oct 2024 23:40:48 GMT\r\nContent-Type: text/html; charset=utf-8\r\nContent-Length: 107\r\nX-Content-Type-Options: nosniff\r\nReferrer-Policy: strict-origin-when-cross-origin\r\nPermissions-Policy: accelerometer=(), autoplay=(), camera=(), display-capture=(), geolocation=(), gyroscope=(), keyboard-map=(), magnetometer=(), microphone=(), midi=(), payment=(), picture-in-picture=(), screen-wake-lock=(), sync-xhr=(), xr-spatial-tracking=(), interest-cohort=()\r\nX-Frame-Options: DENY\r\nX-Dns-Prefetch-Control: off\r\nLocation: https://auth.skaytacium.com/?rd=https%3A%2F%2Fhome.skaytacium.com&rm=GET\r\nSet-Cookie: authelia_session=Rs^ePFhV3NpDG^KhTBMBg*HIfI-HGbTZ; expires=Wed, 09 Oct 2024 00:40:48 GMT; domain=skaytacium.com; path=/; HttpOnly; secure; SameSite=Lax\r\n\r\n<a href=\"https://auth.skaytacium.com/?rd=https%3A%2F%2Fhome.skaytacium.com&amp;rm=GET\">401 Unauthorized</a>")
// 	parse := HTRes{}

// 	for i := 0; i < b.N; i++ {
// 		FastHTResParse(req, &parse)
// 	}
// }

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

func BenchmarkBFind(b *testing.B) {
	var cache = []uint64{
		0, 1, 2, 3, 4, 5, 6, 7, 8,
	}

	for i := 0; i < b.N; i++ {
		if BFind(cache, 5) != 5 {
			b.Errorf("um")
		}
	}
}
