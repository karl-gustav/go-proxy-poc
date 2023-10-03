package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/hello/*", func(w http.ResponseWriter, r *http.Request) {
		serveReverseProxy("/hello", "http://localhost:8081", w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Serving http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func serveReverseProxy(prefix, target string, w http.ResponseWriter, r *http.Request) {
	// parse the targetURL
	targetURL, _ := url.Parse(target)

	r.RequestURI = strings.TrimPrefix(r.RequestURI, prefix)
	r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Update the headers to allow for SSL redirection
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(w, r)
}
