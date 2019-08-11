package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gernest/hiro/headers"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/util"
	"go.uber.org/zap"
)

// MustHeader makesure the headers are present and matches.
func MustHeader(head ...http.Header) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, hd := range head {
				for key := range hd {
					a := r.Header.Get(key)
					b := hd.Get(key)
					if a != b {
						util.WriteJSON(w, &models.APIError{Message: "Fobidden"}, http.StatusForbidden)
						return
					}
				}
			}
			h.ServeHTTP(w, r)
		})
	}

}

func jsonHeader() http.Header {
	h := make(http.Header)
	h.Set(headers.ContentType, headers.ApplicationJSON)
	return h
}

// User decodes jwt token if present and injects it into request context.
func User(log *zap.Logger, jwt *models.JWT) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tk, err := extract(r)
			if err != nil {
				// log.Error("jwt",
				// 	zap.String("url", r.URL.Path),
				// 	zap.Error(err))
				//TODO: handle error
			} else {
				s, err := util.ParseJWTToken(jwt, tk)
				if err != nil {
					log.Error("jwt", zap.Error(err))
				} else {
					if err := s.Valid(); err != nil {
						log.Error("jwt", zap.Error(err))
					} else {
						ctx := context.WithValue(r.Context(),
							helper(keys.Session), s)
						h.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
			h.ServeHTTP(w, r)
		})
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

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s \n", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
