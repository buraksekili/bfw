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

func main() {
	targetURL, err := url.Parse("https://httpbingo.org")
	if err != nil {
		log.Fatal("Error parsing target URL:", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	myConf := BFWConf{
		Destinations: []Destination{
			{From: "/hello", To: "/get"},
		},
	}

	for _, destination := range myConf.Destinations {
		http.HandleFunc(destination.From, func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("handling", destination.From)
			proxy.Director = func(req *http.Request) {
				req.URL.Scheme = targetURL.Scheme
				req.URL.Host = targetURL.Host
				req.URL.Path = destination.To
				req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
			}

			proxy.ServeHTTP(w, r)
		})
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
