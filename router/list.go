package router

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

func (r *Router) ListRoutes(routes Routes) {
	// Order routes alphabetically by path
	sortedRoutes := make([]Route, len(routes))
	copy(sortedRoutes, routes)

	sort.Slice(sortedRoutes, func(i, j int) bool {
		return sortedRoutes[i].getPath() < sortedRoutes[j].getPath()
	})

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	r.listRoutesGroup(sortedRoutes, "")

	w.Flush()
}

func (r *Router) listRoutesGroup(routes Routes, prefix string) {
	for _, route := range routes {
		switch typedRoute := route.(type) {
		case RouteGroup:
			r.listRoutesGroup(typedRoute.Routes, prefix+typedRoute.Prefix)
		case InertiaRoute:
			fmt.Println("Path:", r.getRoutePath(prefix, typedRoute.Path), "\tType: InertiaRoute", "\tPage:", typedRoute.Page)
		case HttpHandlerRoute:
			fmt.Println("Path:", r.getRoutePath(prefix, typedRoute.Path), "\tType: HttpHandlerRoute")
		}
	}
}
