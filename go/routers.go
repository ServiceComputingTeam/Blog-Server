/*
 * Blog for service computing
 *
 *
 *
 * API version: 1.0.0
 * Contact: 895118352@qq.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func NewUserRouter() *mux.Router {
	router := mux.NewRouter().PathPrefix("/user").Subrouter().StrictSlash(true)
	for _, route := range userRoutes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

var userRoutes = Routes{
	Route{
		"CreateUser",
		strings.ToUpper("Post"),
		"/",
		CreateUser,
	},

	Route{
		"GetBlogByUser",
		strings.ToUpper("Get"),
		"/blogs",
		GetBlogByUser,
	},

	Route{
		"PublishBlog",
		strings.ToUpper("Post"),
		"/blogs",
		PublishBlog,
	},

	Route{
		"UpdateUser",
		strings.ToUpper("Patch"),
		"/",
		UpdateUser,
	},

	Route{
		"UserLogin",
		strings.ToUpper("Put"),
		"/login",
		UserLogin,
	},
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"GetALLBlog",
		strings.ToUpper("Get"),
		"/blogs",
		GetALLBlog,
	},

	Route{
		"GetBlogByTitle",
		strings.ToUpper("Get"),
		"/blogs/{username}/{title}",
		GetBlogByTitle,
	},

	Route{
		"GetReviewByBlog",
		strings.ToUpper("Get"),
		"/blogs/{username}/{title}/reviews",
		GetReviewByBlog,
	},

	Route{
		"Review",
		strings.ToUpper("Post"),
		"/blogs/{username}/{title}/reviews",
		AddReview,
	},

	Route{
		"GetBlogByLabel",
		strings.ToUpper("Get"),
		"/labels/{labelname}/blogs",
		GetBlogByLabel,
	},

	Route{
		"GetBlogByUsername",
		strings.ToUpper("Get"),
		"/users/{username}/blogs",
		GetBlogByUsername,
	},
}
