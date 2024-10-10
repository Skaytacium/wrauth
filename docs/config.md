# config

## What?

the configuration is split into 2 parts:
- wrauth db, which is a `db.yaml` file.
- wrauth configuration, which is a `config.yaml` file.

both files **must be in the same directory**. they are parsed on program start and on change and their position can be specified using command line arguments (check `wrauth --help`).  
all IP addresses **must** be in CIDR notation (have a subnet mask at the end, use `/32` as an equivalent to no subnet).

### wrauth db

```yaml
# OPTIONAL: the rules to match. only IPs that match one of these will be authorized. applied sequentially.
rules:
    # EITHER: REQUIRED: the addresses to match
  - ips:
    - '10.0.0.0/30'
    - '10.0.0.10/32'
    # REQUIRED: the user to allot
    user: 'alice'
    # OR: REQUIRED: the public keys (from WireGuard) to match
  - pubkeys: 
      - 'MJ6JoquFLTf419V5dzkcV1z8TY8SIuPyaSH/1SBBP1o='
    user: 'bob'

# OPTIONAL: the site specific headers to add. also sequential.
headers:
    # REQUIRED: the X-Forwarded-URLs to match
  - urls: 
    - 'https://test.example.com'
    - 'https://devdb.example.com'
    # REQUIRED: a specific set of identities to match
    subjects:
      # MINIMUM: 1
        # EITHER: REQUIRED: the user to allow admin access to
      - - user: 'databaseguy'
        # OR: REQUIRED: the group to allow admin access to
      - - group: 'devs'
        - group: 'sys'
    # REQUIRED: the headers to add
    headers:
      # MINIMUM: 1
        # REQUIRED: header name: header value
        - X-AuthDB-Roles: "devdb"

# REQUIRED: site admins, who can control all peers, same as data.subjects
admins:
  # MINIMUM: 1
    # EITHER: REQUIRED: the user to allow admin access to
  - - user: 'admin'
    # OR: REQUIRED: the group to allow admin access to
    - group: 'admins'
```

### wrauth configuration

```yaml
# OPTIONAL: the address to listen on. use the 'unix:' prefix to specify a unix domain path
# NOTE: wrauth doesn't support ipv6 (and doesn't plan to)
# DEFAULT: 127.0.0.1:9092
# NORELOAD
address: '127.0.0.1:9093'
# REQUIRED: the full external address
external: 'https://wrauth.example.com'
# OPTIONAL: the log level 
# DEFAULT: info
# NORELOAD
level: 'debug'
# OPTIONAL: the theme (currently only gruvbox-dark)
# DEFAULT: gruvbox-dark
theme: 'gruvbox-dark'

# REQUIRED: the Authelia configuration
authelia:
  # REQUIRED: Authelia's listening address
  address: '127.0.0.1:9091'
  # REQUIRED: Authelia's user database
  db: '/opt/authelia/users.yaml'
  # OPTIONAL: how many connections to keep open with Authelia
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

# REQUIRED: the wireguard interfaces to manage, and their respective addresses
interfaces:
  # MINIMUM: 1
    # REQUIRED: name of the interface
  - name: 'wg0'
    # REQUIRED: listening address (subnet mask defaults to 32)
    addr: '10.0.0.1'
  - name: 'wg1'
    addr: '172.16.0.1'
```