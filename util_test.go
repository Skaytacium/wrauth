package main

import (
	"bytes"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.zx2c4.com/wireguard/wgctrl"
)

func BenchmarkNetParseCIDR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		net.ParseCIDR("129.168.255.235/32")
	}
}

func BenchmarkFastUCIDR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ip, mask uint32
		ParseUCIDR([]byte("129.168.255.235/32"), &ip, &mask)
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
		ParseUIP([]byte("129.168.255.235"), &ip)
	}
}

func BenchmarkNetIPEqual(b *testing.B) {
	a, B := net.IPv4(129, 168, 255, 235), net.IPv4(255, 255, 255, 255)
	for i := 0; i < b.N; i++ {
		if net.IP.Equal(a, B) {

		}
	}
}

func BenchmarkCompareUIP(b *testing.B) {
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

func BenchmarkIPConv(b *testing.B) {
	ip := net.IPNet{
		IP:   net.IPv4(129, 168, 255, 235),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}
	for i := 0; i < b.N; i++ {
		t := IP{
			Addr: ToUint32([4]byte(ip.IP)),
			Mask: ToUint32([4]byte(ip.Mask)),
		}
		if t.Mask != 0xffffff00 {
			b.Error()
		}
	}
}

func BenchmarkFastHTAuthReqParse(b *testing.B) {
	req := []byte("GET /auth HTTP/1.1\r\nHost: 127.0.0.1:9092\r\nUser-Agent: curl/8.10.1\r\nAccept: */*\r\nX-Forwarded-For: 10.0.0.32\r\nX-Original-Method: GET\r\nX-Original-URL: https://home.skaytacium.com\r\n\r\n")
	parse := HTAuthReq{}

	for i := 0; i < b.N; i++ {
		HTAuthReqParse(req, &parse)
	}
}

func BenchmarkFastHTAuthResParse(b *testing.B) {
	res := []byte("HTTP/1.1 200 OK\r\nDate: Fri, 11 Oct 2024 05:20:23 GMT\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Length: 6\r\nRemote-User: sid\r\nRemote-Groups: admins\r\nRemote-Name: Skaytacium\r\nRemote-Email: sidk@tuta.io\r\n\r\n200 OK")
	parse := HTAuthRes{}

	for i := 0; i < b.N; i++ {
		HTAuthResParse(res, &parse)
	}
}

func BenchmarkFastHTAuthResGen(b *testing.B) {
	req, m, user := make([]byte, 2048), Match{
		Ip: IP{
			Addr: 0xf0f0f0f0,
			Mask: 0xffffffff,
		},
		Id: "test",
	}, User{
		Disabled:    false,
		DisplayName: "test",
		Email:       "test@test.com",
		Groups:      []string{"test", "example"},
	}

	for i := 0; i < b.N; i++ {
		HTAuthResGen(req, m.Id, &user, HT200)
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
	cache := map[uint64]bool{
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

func BenchmarkLFind(b *testing.B) {
	var list []byte
	var i byte
	for i = 0; i < 255; i++ {
		list = append(list, i)
	}

	for i := 0; i < b.N; i++ {
		LFind(list, 253)
	}
}

// func BenchmarkBFindU64(b *testing.B) {
// 	var list []uint64
// 	var i uint64
// 	for i = 0; i < 255; i++ {
// 		list = append(list, i)
// 	}

// 	for i := 0; i < b.N; i++ {
// 		BFindU64(list, 253)
// 	}
// }

// func BenchmarkBFindU256(b *testing.B) {
// 	var list []uint256.Int
// 	var i uint64
// 	for i = 0; i < 255; i++ {
// 		list = append(list, *uint256.NewInt(i))
// 	}

// 	for i := 0; i < b.N; i++ {
// 		BFindU256(list, uint256.NewInt(253))
// 	}
// }

// func BenchmarkAuthHash(b *testing.B) {
// 	k := uint256.NewInt(0)
// 	for i := 0; i < b.N; i++ {
// 		AuthHash([]byte("https://test.example.com"), []byte("authelia_sesson: 32bitsofdatashit32bitsofdatashit"), k)
// 	}
// }

// func BenchmarkIDHash(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		IDHash(Identity{
// 			User:   "test",
// 			Groups: []string{"testgroup"},
// 		}, "https://test.example.com")
// 	}
// }

func BenchmarkMapHash(b *testing.B) {
	cache := map[string]map[string]map[string]bool{
		"https://test.example.com": {
			"morehashing": {
				"lasthash": true,
			},
		},
	}

	for i := 0; i < b.N; i++ {
		_ = cache["https://test.example.com"]["morehashing"]["lasthash"]
	}
}

func BenchmarkSyncMapHash(b *testing.B) {
	cache := sync.Map{}
	cache.Store("test", "test")
	for i := 0; i < b.N; i++ {
		cache.Load("test")
	}
}

func BenchmarkRWMutex(b *testing.B) {
	mut := sync.RWMutex{}
	for i := 0; i < b.N; i++ {
		mut.Lock()
		mut.Unlock()
	}
}

func BenchmarkRWMutexR(b *testing.B) {
	mut := sync.RWMutex{}
	for i := 0; i < b.N; i++ {
		mut.RLock()
		mut.RUnlock()
	}
}

func BenchmarkStringToByte(b *testing.B) {
	bytearr := []byte("some medium sized byte slice containing random paraphernalia")
	for i := 0; i < b.N; i++ {
		_ = string(bytearr)
	}
}

func BenchmarkByteToString(b *testing.B) {
	bytearr := "some medium sized byte slice containing random paraphernalia"
	for i := 0; i < b.N; i++ {
		_ = []byte(bytearr)
	}
}

func BenchmarkUnsafeString(b *testing.B) {
	bytearr := []byte("some medium sized byte slice containing random paraphernalia")
	for i := 0; i < b.N; i++ {
		_ = unsafe.String(&bytearr[0], len(bytearr))
	}
}

func BenchmarkGetHostString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UFStr(GetHost([]byte("https://some.bull.shit.a")))
	}
}

func BenchmarkUserIn(b *testing.B) {
	uid, id := "test", Identity{
		Users:  []string{},
		Groups: [][]string{{"1"}, {"2", "3"}},
	}
	for i := 0; i < b.N; i++ {
		UserIn(uid, id)
	}
}

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

func BenchmarkTimeSub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now().Add(time.Duration(50) * time.Second).After(time.Now())
	}
}

func BenchmarkWireGuard(b *testing.B) {
	wg, _ := wgctrl.New()
	dev, _ := wg.Device("wg0")

	for i := 0; i < b.N; i++ {
		for _, p := range dev.Peers {
			p = p
		}
	}
}
