package http

import (
    "net/http"
)

type Ctx struct {
	Req *http.Request
	Res *http.ResponseWriter
}

func (c *Ctx) GetParam(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Ctx) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}

func (c *Ctx) SetHeader(key string, value string) {
	(*c.Res).Header().Set(key, value)
}

func (c *Ctx) SetStatus(status int) string {
    (*c.Res).WriteHeader(key)
}

func (c *Ctx) GetBody() string {
	buf := make([]byte, 1024)
	n, _ := c.Req.Body.Read(buf)
	return string(buf[0:n])
}

func (c *Ctx) SendStatus(status int) {
	(*c.Res).WriteHeader(status)
}

func (c *Ctx) Send(body string) {
	(*c.Res).Write([]byte(body))
}

func (c *Ctx) SendJson(body string) {
	c.SetHeader("Content-Type", "application/json")
	c.Send(body)
}

func (c *Ctx) SendJson(body interface{}) {
    c.SetHeader("Content-Type", "application/json")
    c.Send(body)
}

