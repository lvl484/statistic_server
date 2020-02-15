package main

import "net/http"

//Middleware - standard middleware for server
func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		handler.ServeHTTP(w, r)
	})
}