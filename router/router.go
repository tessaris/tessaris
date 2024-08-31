package router

import (
	"fmt"
	"net/http"

	inertiaGo "github.com/petaki/inertia-go"
	"github.com/tesseris-go/tesseris/inertia"
)

type Router struct {
	mux *http.ServeMux
	in  *inertiaGo.Inertia
}

type InertiaProps map[string]interface{}

type InertiaRoute struct {
	Path  string
	Page  string
	Props InertiaProps
}

func New() *Router {
	r := &Router{
		mux: http.NewServeMux(),
		in:  inertia.CreateInertiaClient(),
	}

	return r
}

func (r *Router) registerRoute(route InertiaRoute) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := r.in.Render(w, req, route.Page, route.Props)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (r *Router) RegisterRoutes(routes []InertiaRoute) {
	for _, route := range routes {
		r.mux.HandleFunc(route.Path, r.registerRoute(route))
	}
}

func (r *Router) Serve(routes []InertiaRoute) {
	inertia.CreateInertiaServer()
	r.RegisterRoutes(routes)

	// Serve static files from the resources folder
	r.mux.Handle("/__assets/", http.StripPrefix("/__assets/", http.FileServer(http.Dir(".."))))

	fmt.Println("Server started on http://localhost:3000")
	http.ListenAndServe(":3000", r.mux)
}
