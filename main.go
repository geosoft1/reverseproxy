// simple http reverse proxy
// Copyright (C) 2017  geosoft1  geosoft1@gmail.com
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
)

func NewReverseProxy(target string) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   target,
	})
}

func Handle(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("request:", r.RemoteAddr, "want", r.RequestURI)
		//Many webservers are configured to not serve pages if a request doesn’t appear from the same host.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		p.ServeHTTP(w, r)
	}
}

var Config map[string]interface{}

func main() {
	log.Print("init logger")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		os.Exit(1)
	}

	log.Print("load configuration")
	f, err := os.Open(filepath.ToSlash(pwd + "/conf.json"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	if err := json.NewDecoder(f).Decode(&Config); err != nil {
		log.Println(err.Error())
		return
	}

	log.Print("init routes")
	for Path, Target := range Config["routes"].(map[string]interface{}) {
		// avoid add comments as route
		if Path != "#" {
			http.HandleFunc(Path, Handle(NewReverseProxy(Target.(string))))
			log.Printf("%s > %s", Path, Target)
		}
	}

	Address := fmt.Sprintf("%s:%s", Config["ip"], Config["port"])
	log.Print("start listening on " + Address)
	http.ListenAndServe(Address, nil)
}
