package controllers

import (
	"fmt"
	"net/http"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/middleware"
	"github.com/golly-go/plugins/gql"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
)

// Routes is the entry point to the systems routes
func Initializer(a golly.Application) error {
	a.Routes().
		Use(middleware.Recoverer, middleware.RequestLogger).
		Use(accounts.JWTMiddleware).
		Use(middleware.Cors(middleware.CorsOptions{
			AllowedHeaders:   a.Config.GetStringSlice("cors.headers"),
			AllowedOrigins:   a.Config.GetStringSlice("cors.origins"),
			AllowedMethods:   a.Config.GetStringSlice("cors.methods"),
			AllowCredentials: true,
		})).
		Mount("graphql", gql.NewGraphQL()).
		Namespace("/", func(r *golly.Route) {
			r.Get("/status", func(wctx golly.WebContext) { wctx.RenderStatus(http.StatusOK) })
			r.Post("/test", func(wctx golly.WebContext) {
				fmt.Printf("Request Body: %s", string(wctx.RequestBody()))
			})

		})

	return nil
}
