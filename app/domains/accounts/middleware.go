package accounts

import (
	"context"
	"regexp"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/orm"
	"github.com/golly-go/plugins/passport"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/workos"
)

var (
	usersHeaderMatcher = regexp.MustCompile(`[bB]earer\s(.+)`)
	jwkSet             jwk.Set
)

func JWTMiddleware(next golly.HandlerFunc) golly.HandlerFunc {
	return func(c golly.WebContext) {
		var user User

		ident := identity.Identity{}
		token := DecodeAuthorizationHeader(c.Request().Header.Get("Authorization"))

		tok, err := jwt.ParseString(token, jwt.WithKeySet(jwkSet))

		if err != nil {
			c.Logger().Debugf("cannot parse token: %v", err)
			goto next
		}

		user, err = FindUserByIDPId(c.Context, tok.Subject())
		if err != nil {
			c.Logger().Debugf("cannot find idp user: %v", err)
			goto next
		}

		ident = identity.Identity{
			UID:            user.ID,
			OrganizationID: user.OrganizationID,
		}

	next:
		c.Context = passport.ToContext(c.Context, ident)
		next(c)
	}
}

func ScopeDBMiddleware(next golly.HandlerFunc) golly.HandlerFunc {
	return func(c golly.WebContext) {
		ident := identity.FromContext(c.Context)

		db := orm.DB(c.Context).Scopes(common.OrganizationIDScope(ident.OrganizationID))

		c.Context = orm.SetDBOnContext(c.Context, db)

		next(c)
	}
}

// DecodeAuthorizationHeader removes the "Bearer"
func DecodeAuthorizationHeader(header string) string {
	if token := usersHeaderMatcher.FindStringSubmatch(header); len(token) > 1 {
		return token[1]
	}

	return ""
}

func NewKeySet(ctx golly.Context) (interface{}, error) {
	return nil, nil
}

func initializeJWKMiddleware(app golly.Application) error {
	ctx, cancel := context.WithCancel(context.Background())
	jwkCache := jwk.NewCache(ctx)

	golly.Events().On(golly.EventAppShutdown, func(gctx golly.Context, e golly.Event) error {
		cancel()
		return nil
	})

	url, err := workos.Client(app.NewContext(ctx)).JWKSURL()
	if err != nil {
		return errors.WrapFatal(err)
	}

	err = jwkCache.Register(
		url.String(),
		jwk.WithMinRefreshInterval(10*time.Minute),
	)

	if err != nil {
		return err
	}

	set, err := jwkCache.Refresh(ctx, url.String())
	if err != nil {
		return err
	}

	jwkSet = set
	return nil
}
