// simple https reverse proxy
// Copyright (C) 2017-2019  geosoft1  geosoft1@gmail.com
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// +build https

package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var configFile = flag.String("conf", "conf.json", "configuration file")
var httpsAddress = flag.String("https", ":8090", "https address")
var verbose = flag.Bool("verbose", false, "explain what is being done")

var config map[string]interface{}

func NewReverseProxy(scheme, host string) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: scheme,
		Host:   host,
	})
}

func Register(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if *verbose {
			log.Printf("request %s%s", r.RemoteAddr, r.RequestURI)
		}
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		p.ServeHTTP(w, r)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: %s [options]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(filepath.Join(folder, *configFile))
	if err != nil {
		log.Fatalln(err)
	}

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalln(err)
	}

	for path, host := range config["routes"].(map[string]interface{}) {
		log.Printf("%s -> %s", path, host)
		if strings.HasPrefix(path, "#") {
			// skip comments
			continue
		}
		u, err := url.Parse(host.(string))
		if err != nil {
			// skip invalid hosts
			log.Println(err)
			continue
		}
		http.HandleFunc(path, Register(NewReverseProxy(u.Scheme, u.Host)))
	}

	// allow you to use self signed certificates
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt
	log.Printf("start https server on %s", *httpsAddress)
	if err := http.ListenAndServeTLS(*httpsAddress, filepath.Join(folder, "server.crt"), filepath.Join(folder, "server.key"), nil); err != nil {
		log.Fatalln(err)
	}
}
