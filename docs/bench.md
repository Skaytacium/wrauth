# benchmarks

## What?

benchmarks, atleast on my machines for wrauth. you can always bench on your machine using `go test -bench=` and one of:
- `fhttpAuthOK`

these were done during development as well, to make sure everything's up to speed.

## And?

a frequency after the CPU is the frequency that it's locked to.
one could also use [unix domain sockets](https://en.wikipedia.org/wiki/Unix_domain_socket) for internal communications (nginx <-> wrauth <-> Authelia). they're much faster[^1] to communicate internally with, since they have no TCP overhead.

## Results?

written in the format: `Requests/sec, Avg Latency/ms, Stdev Latency/ms`  
from:
```
Running 10s test @ <wrauth/authelia loopback>
  6 threads and 30 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.26ms   10.15ms 156.09ms   98.28%
    Req/Sec     5.46k     1.46k    8.07k    68.01%
  323726 requests in 10.02s, 37.66MB read
Requests/sec:  32306.15
Transfer/sec:      3.76MB
```

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

|System|Authelia|Authelia (UDS)|wrauth|wrauth (UDS)|
|---|---|---|---|---|
|i7-4870HQ 1.4GHz|32306.15, 2.26, 10.15<br>5127.97, 12.29, 39.77|
|i7-4870HQ 2GHz|45807.85, 1.25, 4.30<br>9535.91, 3.62, 3.47|
|i7-4870HQ 2.5GHz|54432.05, 0.845, 1.22<br>10334.32, 3.40, 3.31|
|Broadcom BCM2712 (Raspberry Pi 5)|