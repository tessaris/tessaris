package router

import (
	"fmt"
	"net/http"

	inertiaGo "github.com/petaki/inertia-go"
	"github.com/tesseris-go/tesseris/inertia"
)

type Router struct {
	prod bool
	mux  *http.ServeMux
	in   *inertiaGo.Inertia
	vite *inertia.Vite
}

type InertiaProps map[string]interface{}

type InertiaRoute struct {
	Path  string
	Page  string
	Props InertiaProps
}

func New(prod bool) *Router {
	r := &Router{
		prod: prod,
		mux:  http.NewServeMux(),
		in:   inertia.CreateInertiaClient(),
		vite: inertia.NewVite(prod),
	}

	return r
}

func (r *Router) requestAddViewData(req *http.Request, route InertiaRoute) *http.Request {
	if !r.prod {
		req = req.WithContext(
			r.in.WithViewData(
				r.in.WithViewData(
					req.Context(),
					"reactRefresh",
					r.vite.ReactRefresh(),
				),
				"viteTags",
				r.vite.ViteTags([]string{
					"@vite/client",
					"resources/css/app.css",
					"resources/js/app.tsx",
					fmt.Sprintf("resources/js/pages/%s.tsx", route.Page),
				}),
			),
		)
	} else {
		req = req.WithContext(
			r.in.WithViewData(
				req.Context(),
				"viteTags",
				r.vite.ViteTags([]string{
					"resources/css/app.css",
					"resources/js/app.tsx",
					fmt.Sprintf("resources/js/pages/%s.tsx", route.Page),
				}),
			),
		)
	}

	return req
}

func (r *Router) registerRoute(route InertiaRoute) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req = r.requestAddViewData(req, route)

		err := r.in.Render(w, req, route.Page, route.Props)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (r *Router) registerRoutes(routes []InertiaRoute) {
	for _, route := range routes {
		r.mux.HandleFunc(route.Path, r.registerRoute(route))
	}
}

func (r *Router) Serve(routes []InertiaRoute) {
	go inertia.RunInertiaServer()

	r.registerRoutes(routes)

	// Serve static files from the resources folder
	r.mux.Handle("/__assets/", http.StripPrefix("/__assets/", http.FileServer(http.Dir(".."))))

	fmt.Println("Server started on http://localhost:3000")
	http.ListenAndServe(":3000", r.mux)
}
