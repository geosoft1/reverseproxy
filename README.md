reverseproxy
====
[![version](https://img.shields.io/badge/version-1.0.1-blue.svg)](https://github.com/geosoft1/reverseproxy/archive/master.zip)
[![license](https://img.shields.io/badge/license-gpl-blue.svg)](https://github.com/geosoft1/reverseproxy/blob/master/LICENSE)

Simple [reverse proxy](https://en.wikipedia.org/wiki/Reverse_proxy) server. Useful for accessing web applications on various servers (or VMs) through a single domain.

## How it works?

![reverseproxy](https://user-images.githubusercontent.com/6298396/36028867-5e549ea4-0da9-11e8-8ecf-62546e95ca5c.png)

Just complete the `conf.json` file and run the server. Example:

     {
         "ip":"",
         "port":"8080",
         "routes":{
		     "/upload":"192.168.88.160:8080",
     	     "/Downloads/":"192.168.88.164:8000",
     	     "#":"the pattern / matches all paths not matched by other registered patterns",
     	     "/":"192.168.88.161"
         }
     }

## Configuration details

     "ip":"",

No ip mean `localhost` on hosting server. Is no need to change this.

     "port":"8080",

The server listening on this port. Remeber to forward the port `80` to this port if your connection pass through a router. No root right are required if you run on big ports (eg. `8080`).

## Routes

Routes has the folowing structure

     "path":"target"

The path is what you request and the target is what you get (eg. if your domain is `example.com` then `/` mean `example.com/` and `/upload` mean `example.com/upload`).

`#` path mean a comment and is not added to routes. Put the text in target. `#something` don't mean a comment.

The reverse proxy add your path to the target, so be prepared to handle this path. For example the folowing will get an error page.

     "/upload":"google.com"

Use `/` path for main site which have index page on `/`. Use sufixes for other web services which have the sufix as main page.

Remeber that a route like `/name/` mean match any starting with `name` (eg. `/api/` match also `/api/bla` and so on).

Do not repeat the routes because the server will take always tha last route to a target.

