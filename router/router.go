package router

import (
	"fmt"
	"net/http"

	inertiaGo "github.com/petaki/inertia-go"
	"github.com/tesseris-go/tesseris/config"
	"github.com/tesseris-go/tesseris/inertia"
)

type Router struct {
	prod bool
	mux  *http.ServeMux
	in   *inertiaGo.Inertia
	vite *inertia.Vite
	port int
}

type Route interface{}
type Routes []interface{}

type RouteGroup struct {
	Prefix string
	Routes Routes
}

type HttpHandlerRoute struct {
	Path    string
	Handler http.HandlerFunc
}

type InertiaProps map[string]interface{}

type InertiaRoute struct {
	Path  string
	Page  string
	Props InertiaProps
}

func New(prod bool, c *config.Config) *Router {
	r := &Router{
		prod: prod,
		mux:  http.NewServeMux(),
		in:   inertia.CreateInertiaClient(),
		vite: inertia.NewVite(prod),
		port: c.Port,
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

func (r *Router) registerHandler(path string, handler http.HandlerFunc) {
	if path == "/" {
		r.mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				w.WriteHeader(404)
				w.Write([]byte("404 Not Found"))
				return
			}

			handler(w, r)
		}))

		return
	}

	r.mux.HandleFunc(path, handler)
}

func (r *Router) registerInertiaRoute(route InertiaRoute) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		req = r.requestAddViewData(req, route)

		err := r.in.Render(w, req, route.Page, route.Props)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	r.registerHandler(route.Path, handler)
}

func (r *Router) registerRoute(route Route) {
	if inertiaRoute, ok := route.(InertiaRoute); ok {
		r.registerInertiaRoute(inertiaRoute)
	} else if httpHandlerRoute, ok := route.(HttpHandlerRoute); ok {
		r.registerHandler(httpHandlerRoute.Path, httpHandlerRoute.Handler)
	}
}

func (r *Router) registerRouteGroup(routeGroup RouteGroup) {
	for _, route := range routeGroup.Routes {
		if inertiaRoute, ok := route.(InertiaRoute); ok {
			inertiaRoute.Path = routeGroup.Prefix + inertiaRoute.Path

			r.registerInertiaRoute(inertiaRoute)
		}
	}

}

func (r *Router) registerRoutes(routes Routes) {
	for _, route := range routes {
		if routeGroup, ok := route.(RouteGroup); ok {
			r.registerRouteGroup(routeGroup)
		} else {
			r.registerRoute(route)
		}
	}
}

func (r *Router) Serve(routes Routes) {
	go inertia.RunInertiaServer()

	r.registerRoutes(routes)

	// Serve static files from the resources folder
	r.mux.Handle("/__assets/", http.StripPrefix("/__assets/", http.FileServer(http.Dir(".."))))

	fmt.Printf("Server started on http://localhost:%d\n", r.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.mux)

	if err != nil {
		panic(err)
	}
}
