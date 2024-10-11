# reverse proxy

## nginx

since nginx and Authelia support is first-party, there is not much additional configuration required.
- **remove `X-Forwarded-URI` and `X-Forwarded-Ssl` on the auth request**. this is required due to the way wrauth's parser works. ideally, the only headers should be `X-Original-Method`, `X-Original-URL` and `X-Forwarded-For`. any headers that are as long as the aforementioned 3 could cause wrauth to bug out silently.
- change your `auth_request` address.
- add site specific `auth_request_set` directives.

keep in mind, wrauth supports **only HTTP/1.1**, any other version will **not work**. set your reverse proxy configuration accordingly.

### base
```nginx
http {
	# not documented
	include ssl.conf;
	include proxy.conf;

	server {
		server_name extauth.example.com;
		listen *:443 ssl;

		location / {
			proxy_pass <wherever Authelia is>;
		}

		location /api/authz/auth-request {
			return 403;
		}
	}

	server {
		server_name wrauth.example.com;
		listen *:443 ssl;

		location / {
			proxy_pass <wherever wrauth is>;
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

# change this to wrauth's url
set $auth_server http://127.0.0.1:9092/

location /int/auth {
	internal;

	proxy_pass $auth_server/api/authz/auth-request;

	# add these 2 or any more as mentioned in the docs
	proxy_set_header X-Forwarded-URI "";
	proxy_set_header X-Forwarded-Ssl "";

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
# add this accordingly
auth_request_set $<site_specific_header> $upstream_http_remote_<site_specific_header>;

auth_request_set $redirection_url $upstream_http_location;
error_page 401 =302 $redirection_url;
```