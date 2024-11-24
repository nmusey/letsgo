package middleware

import "net/http"

type CorsMiddleware struct {}

func (cm CorsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
}
