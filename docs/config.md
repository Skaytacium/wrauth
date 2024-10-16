# config

## What?

the configuration is split into 2 parts:
- wrauth db, which is a `db.yaml` file.
- wrauth configuration, which is a `config.yaml` file.

both files **must be in the same directory**. they are parsed on program start and on change and their position can be specified using command line arguments (check `wrauth --help`).  

all IP addresses **must** be in CIDR notation (have a subnet mask at the end, use `/32` as an equivalent to no subnet).

options with `EITHER:` and `OR:` are not mutually exclusive, i.e. `rules` entries can have both `ips` and `pubkeys`, but must have atleast one of them (thats why the prefix).

wrauth doesn't watch the Authelia user database, but it parses that on file change. to reload it, just write something (like a comment) to one of the watched `yaml` files and it will be reloaded.

access rules are NOT matched sequentially, it is specifically in the order:
- bypasses
- direct
- globs

### wrauth db

```yaml
# REQUIRED: IP to user matching rules
rules:
    # EITHER: the addresses to match
  - ips:
    - '10.0.0.0/30'
    - '10.0.0.10/32'
    # REQUIRED: the user to designate
    user: 'alice'
    # OR: the public keys (from WireGuard) to match
  - pubkeys: [ 'MJ6JoquFLTf419V5dzkcV1z8TY8SIuPyaSH/1SBBP1o=' ]
    user: 'bob'

# REQUIRED: access control rules
access:
  # REQUIRED: the domains to match
  - domains: [ 'private.example.com' ]
    # OPTIONAL: regex to match for path
    # NOTE: https://github.com/google/re2/wiki/Syntax
    resource: '/(cpanel|database).*'
    # EITHER: the users to allow
    users: [ 'admin' ]
    # OR: the groups to allow
    groups:
      - [ 'trusted' ]
  # NOTE: https://pkg.go.dev/path/filepath#Match
  # NOTE: you CANNOT bypass a glob, this is by design
  - domains: [ '*.example.com' ]
    users: [ 'superadmin' ]
  - domains: [ 'public.example.com' ]
    # NOTE: allow all users (this is how you bypass)
    users: [ '*' ]
  - domains: 
    - 'test.example.com'
    - 'devdb.example.com'
    users: 
      - 'databaseguy'
      - 'admins'
    groups:
      # NOTE: matches users who are in 'maindbs' or in both 'sys' and 'devs'
      - - 'maindbs'
      - - 'devs' 
        - 'sys'
    # OPTIONAL: any site specific headers to add
    headers:
      # OPTIONAL: header name: header value
      X-AuthDB-Roles: "devdb"

# REQUIRED: site admin rules
admins:
  users: [ 'admin' ]
  groups:
    - - 'admins'
    - - 'sys'
      - 'trusted'
```

### wrauth configuration

```yaml
# OPTIONAL: the port to listen on.
# NOTE: wrauth doesn't support ipv6 (and doesn't plan to)
# NOTE: wrauth only listens on 127.0.0.1. this is by design, it is not meant to be used outside of a reverse proxy.
# DEFAULT: 9092
# NORELOAD
address: '9093'
# REQUIRED: the full external address
external: 'https://wrauth.example.com'
# OPTIONAL: enable or disable caching
# NOTE: this will reduce performance SIGNIFICANTLY if disabled, since proxying Authelia directly reduces performance by ~40%
# DEFAULT: true
# NORELOAD
caching: true
# OPTIONAL: the log level 
# DEFAULT: info
# NORELOAD
level: 'debug'
# OPTIONAL: the theme (currently only gruvbox-dark)
# DEFAULT: gruvbox-dark
theme: 'gruvbox-dark'

# REQUIRED: Authelia configuration
authelia:
  # REQUIRED: Authelia's listening address
  address: '127.0.0.1:9091'
  # REQUIRED: Authelia's user database
  db: '/opt/authelia/users.yaml'
  # OPTIONAL: no. connections to keep open with Authelia
  # DEFAULT: 64
  # NORELOAD
  connections: 32
  # OPTIONAL: how often to clear cache in seconds
  # DEFAULT: 300
  # NORELOAD
  cache: 600
  # OPTIONAL: how often to ping Authelia in seconds (must be below 30 to keep connections alive)
  # DEFAULT: 25
  # NORELOAD
  ping: 28

# REQUIRED: WireGuard interfaces
interfaces:
  # REQUIRED: interface name
  - name: 'wg0'
    # REQUIRED: listening address
    addr: '10.0.0.1/32'
    # OPTIONAL: the configuration file
    # DEFAULT: /etc/wireguard/<name>.conf
    conf: '/etc/wireguard/wg.conf'
  - name: 'wg1'
    addr: '172.16.0.1'
    # OPTIONAL: time from the last handshake to consider a connection closed
    # DEFAULT: 150
    shake: 300
```

### Authelia configuration

this is just the recommended Authelia configuration demonstration how it should be used with wrauth. notice how **there are no conflicting network based rules**, this is important, so that there are no conflicts in authentication, e.g. wrauth disallows or doesn't match `10.0.0.5` (intended) but Authelia has bypassed that entire network (accidental), so wrauth requests Authelia and it immediately responds with 200 OK, causing wrauth to reply with that.

```yaml
  rules:
    # public, allow everybody
    - domain:
      - 'example.com'
      - 'www.example.com'
      - 'public.example.com'
      policy: 'bypass'
    # semi-private, allow everybody who's either on the VPN or authenticated
    - domain: [ 'semi.example.com' ]
      policy: 'bypass'
      networks: 'vpn'
    - domain: [ 'semi.example.com' ]
      policy: 'two_factor'
    # private, allow only specific users
    - domain: '*.example.com'
      policy: 'two_factor'
      subject: 'group:admins' 
    - domain: [ 'db.example.com', 'git.example.com' ]
      policy: 'two_factor'
      subject: 'group:devs'
```