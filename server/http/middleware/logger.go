package middleware

import (
    "log"
    "github.com/nmusey/letsgo/server/http"
)

func LogRequest(ctx http.Ctx) error {
    log.Println(ctx.Req.Method, ctx.Req.URL.Path)
    return nil
}
