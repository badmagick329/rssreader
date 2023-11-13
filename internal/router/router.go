package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type SubRouter struct {
	Router   *chi.Mux
	BasePath string
}

type Router struct {
	Router     *chi.Mux
	Port       string
	SubRouters *[]SubRouter
}

func New(port string) Router {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	return Router{
		Router:     r,
		Port:       port,
		SubRouters: &[]SubRouter{},
	}
}

func (r *Router) AddRoute(method, basePath, path string, handler http.HandlerFunc) {
	subRouter := r.getSubRouter(basePath)
	subRouter.Router.Method(method, path, handler)
	r.Router.Mount(subRouter.BasePath, subRouter.Router)
}

func (r *Router) getSubRouter(basePath string) *SubRouter {
	for _, subRouter := range *r.SubRouters {
		if subRouter.BasePath == basePath {
			return &subRouter
		}
	}
	return &SubRouter{
		chi.NewRouter(),
		basePath,
	}
}

func (r *Router) Run() {
	srv := &http.Server{
		Addr:    ":" + r.Port,
		Handler: r.Router,
	}
	log.Fatal(srv.ListenAndServe())
}
