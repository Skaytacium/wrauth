# benchmarks

## What?

notes on performance for wrauth.  

- configure Authelia with `log.level` as `warn` since it outputs every failed authorization on `stdout`. this is expensive, so expensive that I got a 6x speed improvement from doing this (~5k req/s to ~30k req/s).
- using [unix domain sockets](https://en.wikipedia.org/wiki/Unix_domain_socket) for internal communications (nginx <-> wrauth <-> Authelia) is much faster than the following benchmarks, since those have no TCP overhead.
- use the `GOMAXPROCS` environment variable to limit the no. OS threads that wrauth uses. and also check out [the fasthttp performance optimization tricks](https://pkg.go.dev/github.com/valyala/fasthttp#readme-performance-optimization-tips-for-multi-core-systems).

## Benchmarks?
benchmarks on loopback, atleast on my machines for wrauth. made using [wrk](https://github.com/wg/wrk). there are also internal benchmarks that you can check in `_test.go` files and run to see on your machine. no benchmarks provided for unix domain sockets because [wrk doesn't support it](https://github.com/wg/wrk/issues/400).

a frequency after the CPU is the frequency that it's locked to.

averaged, roughly rounded and written in the format: `Requests/sec/thread, Avg Latency/us, Stdev Latency/us`  

Authelia (to establish the base speed) format: 
```
authorized
unauthorized
```

wrauth format:
```
authorized
unauthorized
uncached
```

| System | Authelia 1t 1c | Authelia 2t 64c | wrauth 1t 1c | wrauth 2t 64c |
|---|---|---|---|---|
| i7-4870HQ 1.4GHz | 4600, 225, 175 <br> 4600, 225, 160 | 23800, 3100, 3000 <br> 23800, 3100, 2900 | 21750, 43, 15 <br> 22500, 41, 15 <br> 2700, 365, 120| 82500, 430, 180 <br> 85000, 430, 180 <br> 15000, 4300, 1900 |
| i7-4870HQ 2GHz | 7300, 144, 130 <br> 7500, 142, 130 | 34000, 2200, 2200 <br> 34000, 2200, 2150 | 31755, 30, 12 <br> 32750, 28, 10 <br> 4000, 250, 100 | 120000, 300, 120 <br> 125000, 300, 120 <br> 21500, 2900, 1000 |
| i7-4870HQ 2.5GHz | 9400, 115, 115 <br> 9600, 110, 110 | 42500, 1800, 1800 <br> 43000, 1750, 1700 | 40300, 23, 8 <br> 41500, 22, 8 <br> 5500, 185, 70 | 156000, 230, 90 <br> 156000, 230, 90 <br> 27000, 2400, 1200 |
| Broadcom BCM2712 (Raspberry Pi 5) | <br> | <br> | <br> <br> | <br> <br> |
