package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route represents all the routes of the API
type Route struct {
	URI                   string
	Method                string
	Function              func(http.ResponseWriter, *http.Request)
	RequireAuthentication bool
}

func Configure(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, loginRoute)

	for _, route := range routes {
		r.HandleFunc(route.URI, route.Function).Methods(route.Method)
	}

	//r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return r
}
