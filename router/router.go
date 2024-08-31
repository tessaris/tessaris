package router

import (
	"fmt"
	"net/http"

	"github.com/iskandervdh/tesseris/inertia"
	inertiaGo "github.com/petaki/inertia-go"
)

type Router struct {
	mux   *http.ServeMux
	in    *inertiaGo.Inertia
	inSer *inertia.InertiaServer
}

func New() *Router {
	router := &Router{
		mux:   http.NewServeMux(),
		in:    inertia.CreateInertiaClient(),
		inSer: inertia.CreateInertiaServer(),
	}

	return router
}

func (router *Router) homeHandler(w http.ResponseWriter, r *http.Request) {
	err := router.in.Render(w, r, "home", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *Router) Serve() {
	r.mux.HandleFunc("/", r.homeHandler)

	// Serve static files from the resources folder
	r.mux.Handle("/__assets/", http.StripPrefix("/__assets/", http.FileServer(http.Dir(".."))))

	fmt.Println("Server started on http://localhost:3000")
	http.ListenAndServe(":3000", r.mux)
}
