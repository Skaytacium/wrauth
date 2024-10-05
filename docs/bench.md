# benchmarks

## What?

notes on performance for wrauth.  

- configure Authelia with `log.level` as `warn` since it outputs every failed authorization on `stdout`. this is expensive, so expensive that I got a 6x speed improvement from doing this (~5k req/s to ~30k req/s).
- using [unix domain sockets](https://en.wikipedia.org/wiki/Unix_domain_socket) for internal communications (nginx <-> wrauth <-> Authelia) could be faster than the following benchmarks, since they have no TCP overhead.
- use the `GOMAXPROCS` environment variable to limit the no. OS threads that wrauth uses. and also check out [the fasthttp performance optimization tricks](https://pkg.go.dev/github.com/valyala/fasthttp#readme-performance-optimization-tips-for-multi-core-systems).

## Benchmarks?
benchmarks, atleast on my machines for wrauth. made using [wrk](https://github.com/wg/wrk). there are also internal benchmarks that you can check in `_test.go` files and run to see on your machine.

a frequency after the CPU is the frequency that it's locked to.

written in the format: `Requests/sec/thread, Avg Latency/ms, Stdev Latency/ms`  

Authelia (to establish the base speed) format: 
```
network bypass
unauthorized
```

wrauth format:
```
authorized
unauthorized
```

|System|Authelia|wrauth|
|---|---|---|
|i7-4870HQ 1.4GHz|32306.15, 2.26, 10.15<br>5127.97, 12.29, 39.77|
|i7-4870HQ 2GHz|45807.85, 1.25, 4.30<br>9535.91, 3.62, 3.47|
|i7-4870HQ 2.5GHz|54432.05, 0.845, 1.22<br>10334.32, 3.40, 3.31|
|Broadcom BCM2712 (Raspberry Pi 5)|
