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
  # REQUIRED(0): the addresses to match
  - ips:
    - '10.0.0.0/30'
    - '10.0.0.10'
    # REQUIRED: the user to assign to the IPs
    user: 'alice'
  - ips: '172.16.0.3'
    user: 'bob'

# REQUIRED: the site specific headers to add. also sequential.
site:
  # REQUIRED(0): the domain to match
  - domain: '/^(db|test)\.example\.com$/'
    # OPTIONAL: a specific set of users/groups to match. same as Authelia subject
    # DEFAULT: nil (match all)
    subject: 'group:devs'
    # REQUIRED: the headers to add
    data:
      # REQUIRED(1): header_name: value
      X-Auth-DB-Roles: 'devdb'
  - domain: 'metrics.example.com'
    data:
      Authorization: "Basic ZGV2OnBhc3N3b3Jk"
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

# REQUIRED: the wireguard interfaces to manage, and their respective addresses
interfaces:
  # REQUIRED(1): name and listening address
  - wg0: '10.0.0.1/32'
    # OPTIONAL: the configuration file
    # DEFAULT: /etc/wireguard/<name>.conf
    conf: '/etc/wireguard/wg.conf'
  - wg1: '172.16.0.1/32'
    # OPTIONAL: to internally mark that only addresses from this IP range will be allowed
    # DEFAULT: <listening_address>/24
    subnet: '172.16.0.0/16'
```