package routes

import (
	"api/src/controllers"
	"net/http"
)

var postsRoutes = []Route{
	{
		URI:                   "/posts",
		Method:                http.MethodPost,
		Function:              controllers.CreatePost,
		RequireAuthentication: true,
	},
	{
		URI:                   "/posts",
		Method:                http.MethodGet,
		Function:              controllers.GetPosts,
		RequireAuthentication: true,
	},
	{
		URI:                   "/posts/{postId}",
		Method:                http.MethodGet,
		Function:              controllers.GetPost,
		RequireAuthentication: true,
	},
	{
		URI:                   "/posts/{postId}",
		Method:                http.MethodPut,
		Function:              controllers.UpdatePost,
		RequireAuthentication: true,
	},
	{
		URI:                   "/posts/{postId}",
		Method:                http.MethodDelete,
		Function:              controllers.DeletePost,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}/posts",
		Method:                http.MethodGet,
		Function:              controllers.GetPostsPerUser,
		RequireAuthentication: true,
	},
	{
		URI:                   "/posts/{postId}/like",
		Method:                http.MethodPost,
		Function:              controllers.LikePost,
		RequireAuthentication: true,
	},
}
