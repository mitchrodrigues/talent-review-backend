package accounts

import (
	"github.com/golly-go/golly"
)

func Initializer(app golly.Application) error {
	InitGraphQL()

	app.Routes().Namespace("/access", func(r *golly.Route) {
		r.Mount("/webhooks", WebhookController{})

		// For now till i can fix insomnia
		if !golly.Env().IsProduction() {
			r.Get("/token", func(wctx golly.WebContext) {
				token := DecodeAuthorizationHeader(wctx.Request().Header.Get("Authorization"))
				wctx.RenderJSON(map[string]string{
					"token": token,
				})

			})
		}
	})

	return initializeJWKMiddleware(app)
}
