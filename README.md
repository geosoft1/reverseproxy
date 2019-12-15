reverseproxy
====
[![version](https://img.shields.io/badge/version-2.0.1-blue.svg)](https://github.com/geosoft1/reverseproxy/archive/master.zip)
[![license](https://img.shields.io/badge/license-gpl-blue.svg)](https://github.com/geosoft1/reverseproxy/blob/master/LICENSE)

Simple [reverse proxy](https://en.wikipedia.org/wiki/Reverse_proxy) server. Useful for accessing web applications on various servers (or VMs) through a single domain.

## How it works?

![reverseproxy](https://user-images.githubusercontent.com/6298396/36028867-5e549ea4-0da9-11e8-8ecf-62546e95ca5c.png)

Just complete the `conf.json` file and run the server. Example:

	{
	    "routes": {
	        "#": "the pattern / matches all paths not matched by other registered patterns",
	        "/": "http://192.168.88.250",
	        "/wrong": "192.168.88.250:8080",
	        "/upload": "http://192.168.88.250:8080",
	        "/hello": "https://192.168.88.250:8090",
	        "/static/": "http://192.168.88.250:8080",
	        "#/disabled": "192.168.88.250:8080"
	    }
	}

## Getting started

To compile the reverse proxy server use

	go build

If you still want just an HTTP reverse proxy, compile with

	go build http.go

or for HTTPS

	go build https.go

Note that `Register` function (see [main.go](https://github.com/geosoft1/reverseproxy/blob/d5dce6d78fb666405cead40cf0c14fb7278f620a/main.go#L48), [http.go](https://github.com/geosoft1/reverseproxy/blob/d5dce6d78fb666405cead40cf0c14fb7278f620a/http.go#L47) and [https.go](https://github.com/geosoft1/reverseproxy/blob/d5dce6d78fb666405cead40cf0c14fb7278f620a/https.go#L48)) have some headers commented. Change as you wish for dealing with applications which need [CORS](https://en.wikipedia.org/wiki/Cross-origin_resource_sharing).

## Parameters

Name|Description
---|---
`-conf`|Cache file name, default value `cache.json`.
`-http`|Listening address and port for HTTP server, default value `8080`.
`-https`|Listening address and port for HTTPS server, default value `8090`.
`-https-enabled`|Enable HTTPS server. Default `false`.
`-verbose`|Enable verbose mode for middleware.

## Routes

Routes has the folowing structure

     "path":"host"

The path is what you request and the host is what you get. The reverse proxy always add the path to the host (eg. if your host address is `example.com` then the path `/` mean `example.com/` and `/upload` mean `example.com/upload`). 

Paths starting with `#` are comments and are not added to routes.

A path like `/name/` match any request starting with `name` (eg. `/api/` match also `/api/bla` and so on).

Hosts must be a complete url address and port.

Do not repeat the routes because the server will take always the last route to a host.

## Testing the server

	curl --verbose http://localhost:8080/hello

For HTTPS use

	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt
	curl --insecure --verbose https://localhost:8090/hello

## Faq

### Why the HTTPS is not enabled by default?

HTTPS server need some valid certificates which you may not have. If you need only a HTTP server is no reason to generate cerificates just to run the program.

### Should I use http or https in the host address?

Yes, prefixes are mandatory to tell the server in which chain to put the route. Omitting that will skip the route.