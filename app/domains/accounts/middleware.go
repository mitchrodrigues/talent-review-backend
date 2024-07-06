package accounts

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/plugins/orm"
	"github.com/golly-go/plugins/passport"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sirupsen/logrus"

	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/workos"
)

var (
	usersHeaderMatcher = regexp.MustCompile(`[bB]earer\s(.+)`)
	jwkSet             jwk.Set
	jwkURL             *url.URL
	jwkCache           *jwk.Cache

	lock sync.RWMutex
)

func JWTMiddleware(next golly.HandlerFunc) golly.HandlerFunc {
	return func(c golly.WebContext) {
		var user User

		ident := identity.Identity{}
		token := DecodeAuthorizationHeader(c.Request().Header.Get("Authorization"))

		if token != "" {
			tok, err := ParseTokenString(c.Context, token)
			if err != nil {
				c.Context.Logger().Debugf("cannot parse token: %v", err)
				goto next
			}

			fmt.Printf("%#v\n", tok.Subject())

			user, err = FindUserByIDPId(c.Context, tok.Subject())
			if err != nil {
				c.Logger().Debugf("cannot find idp user: %v", err)
				goto next
			}

			ident = identity.Identity{
				UID:            user.ID,
				OrganizationID: user.OrganizationID,
			}
		} else {
			c.Logger().Debug("empty token")
		}

		goto next

	next:
		c.Context = passport.ToContext(c.Context, ident)
		c.SetLogger(
			c.Logger().
				WithFields(logrus.Fields{"user_id": ident.UID, "organization_id": ident.OrganizationID}),
		)

		next(c)
	}
}

func ParseTokenString(gctx golly.Context, token string) (tok jwt.Token, err error) {
	var retried bool = false

retry:
	tok, err = jwt.ParseString(token, jwt.WithKeySet(jwkSet))
	if err == nil {
		return
	}

	if retried {
		return
	}

	if !strings.Contains(err.Error(), "failed to find key with key ID") {
		err = errors.WrapGeneric(err)
		return
	}

	lock.Lock()
	jwkSet, err = jwkCache.Refresh(gctx.Context(), jwkURL.String())
	lock.Unlock()

	err = errors.WrapGeneric(err)
	if err != nil {
		return
	}

	retried = true
	goto retry
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

func NewKeySet(ctx golly.Context) (jwk.Set, error) {
	return nil, nil
}

func initializeJWKMiddleware(app golly.Application) error {
	ctx, cancel := context.WithCancel(context.Background())

	jwkCache = jwk.NewCache(ctx)

	// Register a callback to cancel the context on app shutdown
	golly.Events().On(golly.EventAppShutdown, func(gctx golly.Context, e golly.Event) error {
		cancel()
		return nil
	})

	// Fetch the JWK Set URL from the WorkOS client
	url, err := workos.Client(app.NewContext(ctx)).JWKSURL()
	if err != nil {
		return errors.WrapGeneric(err)
	}

	jwkURL = url

	app.Logger.Debugf("JWKS URL: %s", url.String())

	// Register the JWK set URL with the cache, with a minimum refresh interval
	err = jwkCache.Register(
		jwkURL.String(),
		jwk.WithMinRefreshInterval(10*time.Minute),
	)

	if err != nil {
		return errors.WrapGeneric(err)
	}

	// Initial refresh of the JWK set
	set, err := jwkCache.Refresh(ctx, jwkURL.String())
	if err != nil {
		return errors.WrapGeneric(err)
	}

	jwkSet = set
	app.Logger.Infof("Successfully initialized JWK middleware with URL: %s", url.String())
	return nil
}
