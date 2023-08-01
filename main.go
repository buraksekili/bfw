package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func register(dest Destination, targetURL *url.URL) {
	http.HandleFunc(dest.From, func(w http.ResponseWriter, r *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = dest.To
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		}

		proxy.ServeHTTP(w, r)
	})
}

func main() {
	targetURL, err := url.Parse("https://httpbingo.org")
	if err != nil {
		log.Fatal("Error parsing target URL:", err)
	}

	myConf := BFWConf{
		Destinations: []Destination{
			{From: "/hello", To: "/ip"},
			{From: "/tmp", To: "/get"},
		},
	}

	for _, destination := range myConf.Destinations {
		register(destination, targetURL)
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
