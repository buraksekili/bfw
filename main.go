package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func register(listenPath string, target Target) {
	http.HandleFunc(listenPath, func(w http.ResponseWriter, r *http.Request) {
		targetURL := &url.URL{
			Scheme: target.Scheme,
			Host:   target.Host,
			Path:   target.Path,
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = target.Path
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		}

		proxy.ServeHTTP(w, r)
	})
}

func main() {
	myConf := BFWConf{
		Proxy: []Proxy{
			{
				ListenPath: "/",
				Target: Target{
					Scheme: "http", Host: "httpbingo.org", Path: "/get",
				},
			},
			{
				ListenPath: "/tmp",
				Target: Target{
					Scheme: "http", Host: "httpbingo.org", Path: "/ip",
				},
			},
		},
	}

	for _, proxy := range myConf.Proxy {
		register(proxy.ListenPath, proxy.Target)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		fmt.Println("Started listening on 8080")
		errs <- http.ListenAndServe(":8080", nil)
	}()

	fmt.Println("ERROR: ", <-errs)
}
