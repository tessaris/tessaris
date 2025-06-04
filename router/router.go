package router

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	inertiaGo "github.com/petaki/inertia-go"
	"github.com/tessaris/tessaris/config"
	"github.com/tessaris/tessaris/inertia"
)

type Router struct {
	prod              bool
	mux               *http.ServeMux
	in                *inertiaGo.Inertia
	vite              *inertia.Vite
	port              int
	trimTrailingSlash bool
}

type Route interface {
	getPath() string
}
type Routes []Route

type Middleware func(next http.HandlerFunc) http.HandlerFunc

type RouteGroup struct {
	Prefix     string
	Middleware Middleware
	Routes     Routes
}

func (r RouteGroup) getPath() string {
	return r.Prefix
}

type HttpHandlerRoute struct {
	Path    string
	Handler http.HandlerFunc
}

func (r HttpHandlerRoute) getPath() string {
	return r.Path
}

type InertiaProps map[string]interface{}

type InertiaRoute struct {
	Path  string
	Page  string
	Props InertiaProps
}

func (r InertiaRoute) getPath() string {
	return r.Path
}

const PUBLIC_DIR = "public"

func New(prod bool, cfg *config.Config) *Router {
	r := &Router{
		prod:              prod,
		mux:               http.NewServeMux(),
		in:                inertia.CreateInertiaClient(cfg),
		vite:              inertia.NewVite(prod),
		port:              cfg.Port,
		trimTrailingSlash: cfg.TrimTrailingSlash,
	}

	if prod {
		r.in.EnableSsrWithDefault()
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

func (r *Router) getRoutePath(prefix string, path string) string {
	p := path

	if r.trimTrailingSlash && prefix != "" && len(path) == 1 && path[0] == '/' {
		p = ""
	}

	return prefix + p
}

func (r *Router) serveFile(w http.ResponseWriter, req *http.Request) {
	filePath := filepath.Join(PUBLIC_DIR, filepath.Clean("/"+req.URL.Path))

	// Check if the file exists
	if fileExists(filePath) {
		http.ServeFile(w, req, filePath)
		return
	}

	// If no file is found, return 404
	http.NotFound(w, req)
}

func (r *Router) registerHandler(path string, handler http.Handler) {
	if path == "/" {
		r.mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/" {
				handler.ServeHTTP(w, req)

				return
			}

			r.serveFile(w, req)
		}))

		return
	}

	r.mux.Handle(path, handler)
}

func (r *Router) registerInertiaRoute(route InertiaRoute, prefix string) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		req = r.requestAddViewData(req, route)

		err := r.in.Render(w, req, route.Page, route.Props)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// r.mux.Handle(r.getRoutePath(prefix, route.Path), r.in.Middleware(http.HandlerFunc(handler)))

	r.registerHandler(r.getRoutePath(prefix, route.Path), r.in.Middleware(http.HandlerFunc(handler)))
}

func (r *Router) registerRoute(route Route, prefix string) {
	switch typedRoute := route.(type) {
	case InertiaRoute:
		r.registerInertiaRoute(typedRoute, prefix)
	case HttpHandlerRoute:
		r.registerHandler(r.getRoutePath(prefix, typedRoute.Path), typedRoute.Handler)
	}
}

func (r *Router) registerRouteGroup(routeGroup RouteGroup, prefix string) {
	r.registerRoutes(routeGroup.Routes, prefix+routeGroup.Prefix)
}

func (r *Router) registerRoutes(routes Routes, prefix string) {
	for _, route := range routes {
		switch typedRoute := route.(type) {
		case RouteGroup:
			r.registerRouteGroup(typedRoute, prefix)
		default:
			r.registerRoute(route, prefix)
		}
	}
}

func (r *Router) baseUrlRegistered(routes Routes) bool {
	for _, route := range routes {
		switch typedRoute := route.(type) {
		case InertiaRoute:
			if typedRoute.Path == "/" {
				return true
			}
		case HttpHandlerRoute:
			if typedRoute.Path == "/" {
				return true
			}
		}
	}

	return false
}

func (r *Router) Serve(routes Routes) {
	go inertia.RunInertiaServer(r.prod)

	r.registerRoutes(routes, "")

	if !r.baseUrlRegistered(routes) {
		r.mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/" {
				http.NotFound(w, req)

				return
			}

			r.serveFile(w, req)
		}))
	}

	fmt.Printf("Server started on http://localhost:%d\n", r.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.mux)

	if err != nil {
		panic(err)
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
