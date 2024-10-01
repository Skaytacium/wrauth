# docs

## Index?

- [benchmarks](bench.md)
- [configuration reference](config.md)
- [firewall](firewall.md)
- [reverse proxy](reverse.md)

## And?

#### Authentication servers?
currently, wrauth only supports Authelia (and only with the file backend) and I don't plan to support any additional auth servers like [Authentik](https://goauthentik.io/) or LDAP, since this was made for my requirements.

#### Reverse proxies?
due to the aforementioned reasons, only examples for nginx have been given. keep in mind, nginx requires the `ngx_http_js_module` to be installed on your system and loaded in your `nginx.conf`. [Traefik](https://traefik.io/traefik/) *should* be able to work with this with some middleware magic, but I'm not sure since I've never used it.

#### Future?
any PRs that add another auth server, modularize functionality or add more examples new are welcome, but there is no guarantee for their maintenance.