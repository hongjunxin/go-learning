package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type location struct {
	path    string
	handler http.HandlerFunc
}

type route struct {
	pattern string
	handler http.HandlerFunc
}

func init() {
	for _, loc := range locations {
		http.HandleFunc(loc.path, loc.handler)
	}
}

func main() {
	runHttpServer()
}

func runHttpServer() error {
	host := fmt.Sprintf("0.0.0.0:%v", 80)

	srv := &http.Server{
		Addr: host,
		//IdleTimeout: time.Duration(1e9 * conf.IdleTimeout),
		Handler: nil,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("server: Server.ListenAndServe() failed, error=%v", err)
		return err
	}

	return nil
}

var locations = []location{
	{"/", routeServer},
}

var routes = []route{
	{"target_ips[0-9]*.jsp$", iplistServer},
}

func routeServer(w http.ResponseWriter, req *http.Request) {
	for _, r := range routes {
		reg, err := regexp.Compile(r.pattern)
		if err != nil {
			fmt.Printf("server: regexp.Compile() failed, request=%v", req.URL.Path)
			continue
		}
		if reg.MatchString(req.URL.Path) {
			r.handler(w, req)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "404 Page Not Found!")
}

func iplistServer(w http.ResponseWriter, req *http.Request) {
	f := func(c rune) bool {
		return c == '/'
	}
	ss := strings.FieldsFunc(req.URL.Path, f)
	path := "iplist.d/" + ss[len(ss)-1]
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("read '%v' failed, err='%v'\n", path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(file))
}
