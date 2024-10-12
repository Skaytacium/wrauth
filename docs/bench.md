# benchmarks

## What?

notes on performance for wrauth.  

- configure Authelia with `log.level` as `warn` since it outputs every failed authorization on `stdout`. this is expensive, so expensive that I got a 6x speed improvement from doing this (~5k req/s to ~30k req/s).
- using [unix domain sockets](https://en.wikipedia.org/wiki/Unix_domain_socket) for internal communications (nginx <-> wrauth <-> Authelia) is much faster than the following benchmarks, since those have no TCP overhead.
- use the `GOMAXPROCS` environment variable to limit the no. OS threads that wrauth uses. and also check out [the fasthttp performance optimization tricks](https://pkg.go.dev/github.com/valyala/fasthttp#readme-performance-optimization-tips-for-multi-core-systems).

## Benchmarks?
benchmarks, atleast on my machines for wrauth. made using [wrk](https://github.com/wg/wrk). there are also internal benchmarks that you can check in `_test.go` files and run to see on your machine. no benchmarks provided for unix domain sockets because [wrk doesn't support it](https://github.com/wg/wrk/issues/400).

a frequency after the CPU is the frequency that it's locked to.

written in the format: `Requests/sec/thread, Avg Latency/ms, Stdev Latency/ms`  

### TODO

**while these are being done,** you can assume wrauth to go ~5.5x faster than Authelia on 1 connection (4k->22k Req/sec), and ~2.5x faster on 64 connections (28k -> 85k Req/sec)

Authelia (to establish the base speed) format: 
```
network bypass
unauthorized
```

wrauth format:
```
IP authorized
Cache authorized
Cache unauthorized
uncached
```

|System|Authelia|wrauth|
|---|---|---|
|i7-4870HQ 1.4GHz|
|i7-4870HQ 2GHz|
|i7-4870HQ 2.5GHz|
|Broadcom BCM2712 (Raspberry Pi 5)|
