# Header Proxy

The header proxy listens on a port, injects the specified headers into every request and then forwards it.

This is not a production ready project. It was developed and is intended for development purposes only. It was 
originally created to mimic [oauth2_proxy](https://github.com/pusher/oauth2_proxy) locally.

## Usage

```
$ go run main.go --listen=0.0.0.0:3000 --forward=0.0.0.0:80 --header=foo=bar --header=fizz=buzz
2019/10/17 20:39:55 Listening on 0.0.0.0:3000
2019/10/17 20:39:55 Forwarding to 0.0.0.0:80
2019/10/17 20:39:55 Injecting headers map[fizz:buzz foo:bar]
```

### Docker

```
$ docker pull jjshanks/header-proxy
$ docker run --rm jjshanks/header-proxy --listen=0.0.0.0:3000 --forward=0.0.0.0:8080 --header=fizz=buzz --header="X-Forwarded-For=client, proxy1, proxy2"
2019/10/20 18:36:18 Listening on 0.0.0.0:3000
2019/10/20 18:36:18 Forwarding to 0.0.0.0:8080
2019/10/20 18:36:18 Injecting headers map[X-Forwarded-For:client, proxy1, proxy2 fizz:buzz]
```
