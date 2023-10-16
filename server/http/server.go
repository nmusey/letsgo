package http

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

func (s *Server) RegisterHandler(handler *Handler) {
	s.Mux.HandleFunc(handler.path, func(w http.ResponseWriter, r *http.Request) {
		ctx := Ctx{Req: r, Res: &w}

        for _, middleware := range handler.Middleware {
            if err := middleware (ctx); err != nil {
                ctx.SetStatus(http.StatusBadRequest)
                ctx.Send(err.Error())
                return
            }
        }

        if err := handler.callback(ctx); err != nil {
            errorMessage := os.Getenv("ENVIRONMENT") == "dev" ? err.Error() : "Internal Server Error"
            ctx.SetStatus(http.StatusInternalServerError)
            ctx.Send(errorMessage)
        }
	})
    .Methods(handler.method)
}
