package server

import "net/http"

type Server struct {
	Mux *http.ServeMux
}

func NewServer() *Server {
	return &Server{
		Mux: http.NewServeMux(),
	}
}

func (s *Server) Start(port string) error {
	return http.ListenAndServe(port, s.Mux)
}

func (s *Server) ServeStaticDir(url string, path string) {
	s.Mux.Handle(url, http.FileServer(http.Dir(path)))
}

func (s *Server) handleMethod(url string, callback func(Ctx), method string) {
	s.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ctx := Ctx{Req: r, Res: &w}
		callback(ctx)
	})
}

func (s *Server) Get(url string, callback func(Ctx)) {
	s.handleMethod(url, callback, "GET")
}

func (s *Server) Post(url string, callback func(Ctx)) {
	s.handleMethod(url, callback, "POST")
}

func (s *Server) Patch(url string, callback func(Ctx)) {
	s.handleMethod(url, callback, "PATCH")
}

func (s *Server) Put(url string, callback func(Ctx)) {
	s.handleMethod(url, callback, "PUT")
}

func (s *Server) Delete(url string, callback func(Ctx)) {
	s.handleMethod(url, callback, "DELETE")
}



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

