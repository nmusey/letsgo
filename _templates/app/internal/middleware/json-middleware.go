package middleware

import "net/http"

type JsonMiddleware struct {}

func (j JsonMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "application/json")
}
