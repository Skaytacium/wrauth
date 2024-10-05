# reverse proxy

## nginx

since nginx and Authelia support is first-party, there is no additional configuration required except changing your `auth_request` address and adding site specific `auth_request_set` directives.

### base
```nginx
http {
	# not documented
	include ssl.conf;
	include proxy.conf;

	server {
		server_name extauth.example.com;
		listen *:443 ssl;

		set $uri_authelia <wherever_authelia_is>

		location / {
			proxy_pass $uri_authelia;
		}
	}

	server {
		server_name wrauth.example.com;
		listen *:443 ssl;

		set $uri_authelia <wherever_wrauth_is>

		location / {
			proxy_pass $uri_wrauth;
		}
	}

	server {
		server_name private.examplee.com;
		listen *:443 ssl;

		include auth.conf

		# ...
	}
}
```

### proxy

also check out [Authelia's recommended proxy configuration](https://www.authelia.com/integration/proxies/nginx/#proxyconf)

```nginx
# proxy.conf

# version
proxy_http_version 1.1;

# headers
proxy_set_header Host $host;
proxy_set_header X-Forwarded-Proto $scheme;
proxy_set_header X-Forwarded-Host $http_host;
proxy_set_header X-Forwarded-URI $request_uri;
proxy_set_header X-Forwarded-Ssl on;
proxy_set_header X-Forwarded-Server $host;
proxy_set_header X-Forwarded-For $remote_addr;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Real-Port $remote_port;

# upgrade
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection $http_connection;

# misc
proxy_next_upstream error timeout invalid_header http_500 http_502 http_503;
proxy_cache_bypass $cookie_session;
proxy_no_cache $cookie_session;
```

### auth

```nginx
# auth.conf

set $uri_wrauth http://127.0.0.1:9092/

location /int/auth {
	internal;

	proxy_pass $uri_wrauth/auth;

	proxy_set_header X-Original-Method $request_method;
	proxy_set_header X-Original-URL $scheme://$http_host$request_uri;
	proxy_set_header X-Forwarded-For $remote_addr;

	proxy_set_header Content-Length "";
	proxy_set_header Connection "";

	proxy_pass_request_body off;
}

auth_request /int/auth;

auth_request_set $user $upstream_http_remote_user;
auth_request_set $groups $upstream_http_remote_groups;
auth_request_set $name $upstream_http_remote_name;
auth_request_set $email $upstream_http_remote_email;
# this is whatever site specific header
auth_request_set $data $upstream_http_remote_data;

auth_request_set $redirection_url $upstream_http_location;
error_page 401 =302 $redirection_url;
```