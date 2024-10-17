# docs

## Index?

- [performance](bench.md)
- [configuration reference](config.md)
- [firewall](firewall.md)
- [reverse proxy](reverse.md)

## And?

### authentication servers
currently, wrauth only supports Authelia (and only with the file backend) and I don't plan to support any additional auth servers like [Authentik](https://goauthentik.io/) or LDAP, since this was made for my requirements.

### reverse proxies
due to the aforementioned reasons, only examples for nginx have been given. [Traefik](https://traefik.io/traefik/) *should* be able to work with this with some middleware magic, but I'm not sure since I've never used it.

### contribution
any PRs that add another auth server, modularize functionality or add more examples new are welcome, but there is no guarantee for their maintenance.