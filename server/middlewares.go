package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gernest/hiro/headers"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// MustHeader makesure the headers are present and matches.
func MustHeader(head ...http.Header) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()
			for _, hd := range head {
				for key := range hd {
					a := strings.ToLower(r.Header.Get(key))
					b := strings.ToLower(hd.Get(key))
					if a != b {
						return util.Forbid(ctx)
					}
				}
			}
			return h(ctx)
		}
	}
}

func jsonHeader() http.Header {
	h := make(http.Header)
	h.Set(headers.ContentType, headers.ApplicationJSON)
	return h
}

// User decodes jwt token if present and injects it into request context.
func User(log *zap.Logger, jwt *models.JWT) echo.MiddlewareFunc {
	ulog := log.With(zap.Namespace("user-middleware"))
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()
			tk, err := extract(r)
			if err != nil {
				ulog.Debug("extracting token", zap.Error(err))
			} else {
				s, err := util.ParseJWTToken(jwt, tk)
				if err != nil {
					ulog.Debug("parsing token", zap.Error(err))
				} else {
					if err := s.Valid(); err != nil {
						ulog.Debug("validating token", zap.Error(err))
					} else {
						ctx.Set(keys.Session, s)
					}
				}
			}
			return h(ctx)
		}
	}
}

// CtxMiddleware injects the values into the request context.
func CtxMiddleware(items ...models.Item) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			for _, v := range items {
				ctx.Set(v.Key, v.Value)
			}
			return h(ctx)
		}
	}
}

// silence the linter which keeps complaining about using strings as keys to
// context.WithValue.
func helper(k string) interface{} {
	return k
}

func extract(r *http.Request) (string, error) {
	bearer := "Bearer"
	authHeader := r.Header.Get(headers.Authorization)
	if authHeader != "" {
		components := strings.SplitN(authHeader, " ", 2)
		if len(components) != 2 || components[0] != bearer {
			return "", errors.New("can't get bearer token")
		}
		return components[1], nil
	}
	return "", errors.New("bearer token can't be empty")
}

func Logger() middleware.LoggerConfig {
	c := middleware.DefaultLoggerConfig
	c.Format = "${time_rfc3339_nano} [${id}] [${method}] ${uri} ${status}\n"
	return c
}
