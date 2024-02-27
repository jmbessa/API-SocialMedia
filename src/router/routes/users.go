package routes

import (
	"api/src/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		URI:                   "/users",
		Method:                http.MethodPost,
		Function:              controllers.CreateUser,
		RequireAuthentication: false,
	},
	{
		URI:                   "/users",
		Method:                http.MethodGet,
		Function:              controllers.GetUsers,
		RequireAuthentication: false,
	},
	{
		URI:                   "/users/{userID}",
		Method:                http.MethodGet,
		Function:              controllers.GetUser,
		RequireAuthentication: false,
	},
	{
		URI:                   "/users/{userID}",
		Method:                http.MethodPut,
		Function:              controllers.UpdateUser,
		RequireAuthentication: false,
	},
	{
		URI:                   "/users/{userID}",
		Method:                http.MethodDelete,
		Function:              controllers.DeleteUser,
		RequireAuthentication: false,
	},
}
