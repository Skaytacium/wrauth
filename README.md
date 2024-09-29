# wrauth

## What?

wrauth is a [WireGuard](https://www.wireguard.com/) management interface and IPv4 authentication provider that

- has a web UI to manage (create, show, delete, link) WireGuard peers written in plain HTML, CSS and (sadly) some JS.
- comes with [nginx](https://nginx.org/) [auth_request](https://nginx.org/en/docs/http/ngx_http_auth_request_module.html) capabilities out of the box.
- supports existing [Authelia](https://www.authelia.com/) access control rules (with minimal additions) and user database.
- allows site-specific headers and data to be specified.
- is HTTPS TLSv1.3[^a] with only, which every major browser[^b] & library[^c] supports since ~2019 (except Internet Explorer of course).
- [is multithreaded and goes fast.](#benchmarks)

## Why?

I needed something to authenticate users based on their WireGuard IP addresses better than basic [nginx access](https://nginx.org/en/docs/http/ngx_http_access_module.html) and also easily manage creating new peers, quickly generate QR codes and configuration files.  
the drive for automation is always present but this seemed to be a good project to take on and I'd finally be contributing back to the selfhosted community I've taken so much from.  

## How?

wrauth is written in [Go](https://go.dev/) and it uses

## And?

### nginx?

example proxy config:
```nginx
server {

}
```

### Benchmarks?

### Security?

- since public IPv4 addresses are extremely variable[^1] and IPv6 is not supported, it is **absolutely necessary** for WireGuard subnets to be a subset of any of the private use address ranges[^2]. These are official address ranges that are "meant for" virtual networks and are not publicly assigned[^3]. Could be any one of:
    - `10.0.0.0        -   10.255.255.255  (10/8 prefix)`
    - `172.16.0.0      -   172.31.255.255  (172.16/12 prefix)`
    - `192.168.0.0     -   192.168.255.255 (192.168/16 prefix)`
- while this provides host authentication, there is **no guarantee the user is the same**.
- ingress filtering has been best practice since a while[^4], but firewalling on the server is also **seriously recommended**. some simple `nftables` rules to make sure packets are being routed in and out of the same subnet and *only* on the WireGuard interface could prevent basic spoofing (but not [DoS](https://en.wikipedia.org/wiki/Denial-of-service_attack) or [MiTM](https://en.wikipedia.org/wiki/Man-in-the-middle_attack) attacks if your network is already compromised).

```nftables
flush ruleset

table ip firewall {
	chain ingress {
		type filter hook input 
	}
}
```


[^a]: [RFC8446](https://www.rfc-editor.org/rfc/rfc8446.html) and a [nice article about it](https://blog.cloudflare.com/rfc-8446-aka-tls-1-3/)
[^b]: [TLSv1.3 browser adoption statistics](https://caniuse.com/tls1-3)
[^c]: [OpenSSL](https://wiki.openssl.org/index.php/TLS1.3) & [LibreSSL](https://ftp.openbsd.org/pub/OpenBSD/LibreSSL/libressl-3.1.1-relnotes.txt)
[^1]: [IPv4 exhaustion (rise in dynamic addresses)](https://en.wikipedia.org/wiki/IPv4_address_exhaustion)
[^2]: [BCP5/RFC1918](https://www.rfc-editor.org/rfc/rfc1918.html)
[^3]: [IANA IPv4 Special-Purpose Address Registry](https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml)
[^4]: [BCP38/RFC2827](https://www.rfc-editor.org/rfc/rfc2827.html)