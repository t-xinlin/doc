package main

import (
	"fmt"
	"net/http"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		// next handler
		next.ServeHTTP(wr, r)
		timeElapsed := time.Since(timeStart)
		fmt.Println(timeElapsed)
	})
}

func main() {
	r := NewRouter()
	//r.Use(logger)
	//r.Use(timeout)
	//r.Use(ratelimit)
	//r.Add("/", http.HandlerFunc(hello))
	http.Handle("/", timeMiddleware(http.HandlerFunc(hello)))
	for k, v := range r.mux {
		http.Handle(k, v)
	}
	if err := http.ListenAndServe(":8080", nil); err != nil {

	}
}

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) Add(route string, h http.Handler) {
	var mergedHandler = h

	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}

	r.mux[route] = mergedHandler
}
