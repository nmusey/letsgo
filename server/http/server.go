package http

import (
	"net/http"
	"os"
)

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
	s.Mux.HandleFunc(handler.Path, func(w http.ResponseWriter, r *http.Request) {
        ctx := Ctx{Req: r, Res: &w}

        if handler.Method != r.Method {
            ctx.SetStatus(http.StatusMethodNotAllowed)
            return
        }

        for _, middleware := range handler.Middleware {
            if err := middleware (ctx); err != nil {
                ctx.SetStatus(http.StatusBadRequest)
                ctx.Send(err.Error())
                return
            }
        }

        if err := handler.Callback(ctx); err != nil {
            errorMessage := "Internal Server Error"
            if os.Getenv("ENVIRONMENT") == "dev" {
                err.Error() 
            }

            ctx.SetStatus(http.StatusInternalServerError)
            ctx.Send(errorMessage)
        }
    })
}
