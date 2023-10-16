package middleware

func LogRequest(ctx Ctx) err {
    log.Println(ctx.Req.Method, ctx.Req.URL.Path)
    return nil
}
