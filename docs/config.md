# config

## What?

the configuration is split into 2 parts:
- wrauth db, which is a `db.yaml` file.
- wrauth configuration, which is a `config.yaml` file.

both files **must be in the same directory**. they are parsed on program start and on change and their position can be specified using command line arguments (check `wrauth --help`).

use `$<variable_name>` to add dynamic variables to some options (listed with `# DYNAMIC`)  
currently supported variables:
```
user:   the user making the request.
group:  the group that the aforementioned user is in. if multiple groups are present, they're seperated by a comma and space.
ip:     the ip that has been requested with
```

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
    # OR: REQUIRED: the publickey (from WireGuard) to match
  - pubkey: 'MJ6JoquFLTf419V5dzkcV1z8TY8SIuPyaSH/1SBBP1o='
    user: 'bob'

# OPTIONAL: the site specific headers to add. also sequential.
data:
    # REQUIRED: the domain to match
    # DYNAMIC
  - domain: '/^(db|test)\.example\.com$/'
    # OPTIONAL: a specific set of users/groups to match. same as Authelia subject
    # DEFAULT: match all
    subject:
      - 'group:devs'
    # REQUIRED: the headers to add
    headers:
      # MINIMUM: 1
        # REQUIRED: header name: header value
        # DYNAMIC
        - X-AuthDB-User: "$user"
        - X-AuthDB-Roles: "devdb"

# REQUIRED: site admins, who can control all peers
admins:
  # MINIMUM: 1
    # EITHER: REQUIRED: the addresses to allow admin access to
  - ips: '172.16.0.10'
    # OR: REQUIRED: the public key(s) to allow admin access to
  - pubkey: 'MJ6JoquFLTf419V5dzkcV1z8TY8SIuPyaSH/1SBBP1o='
    # OR: REQUIRED: the user to match
  - user: 'admin'
    # OR: REQUIRED: the group to match
  - group: 'admins'
```

### wrauth configuration

```yaml
# OPTIONAL: the address to listen on. use the 'unix:' prefix to specify a unix domain path
# DEFAULT: 127.0.0.1:9092
address: '127.0.0.1:9093'
# REQUIRED: the full external address
external: 'https://wrauth.example.com'
# OPTIONAL: the log level
# DEFAULT: info
log: 'debug'
# OPTIONAL: the theme (currently only gruvbox-dark)
# DEFAULT: gruvbox-dark
theme: 'gruvbox-dark'

# REQUIRED: Authelia related configuration
authelia:
  # REQUIRED: the Authelia configuration
  config: '/opt/authelia/configuration.yml'

# REQUIRED: the wireguard interfaces to manage, and their respective addresses
interfaces:
  # MINIMUM: 1
    # REQUIRED: name of the interface
  - name: 'wg0'
    # REQUIRED: listening address (subnet mask defaults to 32)
    addr: '10.0.0.1'
    # OPTIONAL: the configuration file
    # DEFAULT: /etc/wireguard/<name>.conf
    conf: '/etc/wireguard/wg.conf'
    # OPTIONAL: the duration in seconds after which the peer list cache is updated (happens on a request that misses cache as well)
    # DEFAULT: 15
    watch: 5
  - name: 'wg1'
    addr: '172.16.0.1'
    # OPTIONAL: to internally mark that only addresses from this IP range will be allowed
    # DEFAULT: <listening_address>/24
    subnet: '172.16.0.0/16'
    # OPTIONAL: the duration in seconds from the last handshake after which the connection is considered "closed"
    # DEFAULT: 150
    shake: 300
```