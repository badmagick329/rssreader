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
	subRouter, created := r.getSubRouter(basePath)
	subRouter.Router.Method(method, path, handler)
	if created {
		r.Router.Mount(subRouter.BasePath, subRouter.Router)
	}
	newSlice := append(*r.SubRouters, *subRouter)
	r.SubRouters = &newSlice
}

func (r *Router) getSubRouter(basePath string) (*SubRouter, bool) {
	log.Printf("Searching for basePath %s. SubRouters are: %v\n", basePath, r.SubRouters)
	for _, subRouter := range *r.SubRouters {
		if subRouter.BasePath == basePath {
			log.Printf("Base path %s already exists\n", subRouter.BasePath)
			return &subRouter, false
		}
	}
	return &SubRouter{
		chi.NewRouter(),
		basePath,
	}, true
}

func (r *Router) Run() {
	srv := &http.Server{
		Addr:    ":" + r.Port,
		Handler: r.Router,
	}
	log.Fatal(srv.ListenAndServe())
}
