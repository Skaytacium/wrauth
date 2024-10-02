# config

## What?

the configuration is split into 2 parts:
- wrauth db, which is a neighbouring `db.yaml` file.
- wrauth configuration, which is a neighbouring `config.yaml` file.

both files are parsed on program start and on change.

### wrauth db

```yaml
# REQUIRED: the rules to match. applied sequentially.
rules:
  # MINIMUM: 0
    # EITHER: REQUIRED: the addresses to match (named networks from Authelia can be used here)
  - ips:
    - '10.0.0.0/30'
    - '10.0.0.10'
    # REQUIRED: the user to allot
    user: 'alice'
    # OR: REQUIRED: the publickey to match
  - pubkey: 'MJ6JoquFLTf419V5dzkcV1z8TY8SIuPyaSH/1SBBP1o='
    user: 'bob'

# REQUIRED: the site specific headers to add. also sequential.
site:
  # MINIMUM: 0
    # REQUIRED: the domain to match
  - domain: '/^(db|test)\.example\.com$/'
    # OPTIONAL: a specific set of users/groups to match. same as Authelia subject
    # DEFAULT: nil (match all)
    subject: 'group:devs'
    # REQUIRED: the headers to add
    data:
      # MINIMUM: 1
      # REQUIRED: header_name: value
      X-Auth-DB-Roles: 'devdb'

# REQUIRED: site admins, who can control all peers
admins:
  # MINIMUM: 1
    # EITHER: REQUIRED: same as rules.ips
  - ips: '172.16.0.10'
    # OR: REQUIRED: same as rules.pubkey
  - pubkey: 'MJ6JoquFLTf419V5dzkcV1z8TY8SIuPyaSH/1SBBP1o='
    # OR: REQUIRED: the user to match
  - user: 'admin'
    # OR: REQUIRED: the group to match
  - group: 'admins'
```

### wrauth configuration

```yaml
# OPTIONAL: the address to listen on. use the 'unix:' prefix to specify a unix domain path
# DEFAULT: '127.0.0.1:9092'
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
  # OPTIONAL: the Authelia user database
  # DEFAULT: authelia.config->authetication_backend.file.path
  userdb: '/var/db/users.yaml'
  # OPTIONAL: the login page for Authelia
  # DEFAULT: authelia.config->session.cookies.authelia_url
  login: 'https://extauth.example.com/login'

# REQUIRED: the wireguard interfaces to manage, and their respective addresses
interfaces:
  # MINIMUM: 1
    # REQUIRED: name and listening address
  - wg0: '10.0.0.1/32'
    # OPTIONAL: the configuration file
    # DEFAULT: /etc/wireguard/<name>.conf
    conf: '/etc/wireguard/wg.conf'
    # OPTIONAL: the duration in seconds after which the peer list cache is updated (happens on a request that misses cache as well)
    # DEFAULT: 15
    watch: 5
  - wg1: '172.16.0.1/32'
    # OPTIONAL: to internally mark that only addresses from this IP range will be allowed
    # DEFAULT: <listening_address>/24
    subnet: '172.16.0.0/16'
    # OPTIONAL: the duration in seconds from the last handshake after which the connection is considered "closed"
    # DEFAULT: 150
    handshake: 300
```