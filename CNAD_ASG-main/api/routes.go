package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// proxyHandler proxies requests to a given service
func proxyHandler(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)
	}
}

// SetupRoutes configures the routes for the API gateway
func SetupRoutes(router *http.ServeMux) {
	// Proxy requests to individual services
	router.HandleFunc("/user/", proxyHandler("http://localhost:8081"))
	router.HandleFunc("/billing/", proxyHandler("http://localhost:8082"))
	router.HandleFunc("/booking/", proxyHandler("http://localhost:8083"))
	router.HandleFunc("/vehicle/", proxyHandler("http://localhost:8084"))
}
