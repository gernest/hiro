package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/ory/ladon"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/gernest/hiro/headers"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/query"
	"github.com/gernest/hiro/templates"
)

// Context groups important utilities for request processing.
type Context struct {
	JWT       *models.JWT
	Host      string
	ImageHost string
	Logger    *zap.Logger
	DB        *query.SQL
	Claims    *jwt.StandardClaims
	Warden    *ladon.Ladon
}

func RequestContext(ctx echo.Context) *Context {
	c := &Context{}
	if v := ctx.Get(keys.JwtKey); v != nil {
		c.JWT = v.(*models.JWT)
	}

	if v := ctx.Get(keys.LoggerKey); v != nil {
		c.Logger = v.(*zap.Logger)
	}
	if v := ctx.Get(keys.DB); v != nil {
		c.DB = v.(*query.SQL)
	}
	if v := ctx.Get(keys.Session); v != nil {
		c.Claims = v.(*jwt.StandardClaims)
	}
	if v := ctx.Get(keys.Warden); v != nil {
		c.Warden = v.(*ladon.Ladon)
	}
	return c
}

// WriteJSON sets content type header and marshals v to the response body.
func WriteJSON(w http.ResponseWriter, v interface{}, status int, h ...http.Header) {
	d, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	for _, head := range h {
		for key, value := range head {
			for _, i := range value {
				w.Header().Set(key, i)
			}
		}
	}
	w.Header().Set(headers.ContentType, headers.ApplicationJSON)
	w.WriteHeader(status)
	w.Write(d)
}

// AddToken add bearer token to request header.
func AddToken(r *http.Request, token string) {
	r.Header.Set(headers.Authorization, "Bearer "+token)
}

// AddContentType add content-type header..
func AddContentType(r *http.Request, contentType string) {
	r.Header.Set(headers.ContentType, contentType)
}

// Checksum calculates sha256 checksum of r contents.
func Checksum(r io.Reader) (string, error) {
	h := sha256.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// InfoCMD returns a commmand which will display information about the binary.
func InfoCMD(date, commit, version string) cli.Command {
	return cli.Command{
		Name:  "info",
		Usage: "displays information about the binary",
		Action: func(ctx *cli.Context) error {
			return templates.Write(os.Stdout, "info", map[string]interface{}{
				"Version":   version,
				"BuildDate": date,
				"Commit":    commit,
				"BuildInfo": buildInfo(),
			})
		},
	}
}

func buildInfo() map[string]string {
	return map[string]string{
		"GOOS":      runtime.GOOS,
		"GOARCH":    runtime.GOARCH,
		"GOVERSION": runtime.Version(),
	}
}

// OutJSON marshals v to json and writes it to stdout.
func OutJSON(o interface{}) error {
	b, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, string(b))
	return nil
}

// ListOptions gets additional options when querying a list of items. This
// checks for limit and offset query parameters which are treated as integers.
//
// Lack of offset defaults to 0 and lack of limit defaults to the value
// specified by MaxLimit constant.
func ListOptions(r *http.Request) (*models.ListOptions, error) {
	u := r.URL.Query()
	o := &models.ListOptions{
		Limit: models.MaxLimit,
	}
	if limit := u.Get("limit"); limit != "" {
		v, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
		o.Limit = v
	}

	if offset := u.Get("offset"); offset != "" {
		v, err := strconv.Atoi(offset)
		if err != nil {
			return nil, err
		}
		o.Offset = v
	}
	return o, nil
}

func GenerateJWTToken(j *models.JWT, tk *models.Token) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Audience:  tk.Subject.String(),
		ExpiresAt: tk.ExpiresOn.Unix(),
		Id:        tk.UUID.String(),
		IssuedAt:  tk.CreatedAt.Unix(),
		Issuer:    tk.Issuer.String(),
		NotBefore: tk.NotBefore.Unix(),
		Subject:   tk.Subject.String(),
	})
	return token.SignedString(j.Secret)
}

// Parse decodes tokenString .
func ParseJWTToken(j *models.JWT, tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return j.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*jwt.StandardClaims), nil
}

func Forbid(ctx echo.Context) error {
	return ctx.JSON(http.StatusForbidden, models.APIError{Message: keys.Forbidden})
}

func BadRequest(ctx echo.Context) error {
	return ctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadRequest})
}

func BadToken(ctx echo.Context) error {
	return ctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadToken})
}

func NotFound(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotFound, models.APIError{Message: http.StatusText(http.StatusNotFound)})
}

func Internal(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, models.APIError{Message: keys.InternalError})
}

func Invalid(ctx echo.Context) error {
	return ctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.FailedValidation})
}

func Ok(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.Status{Status: keys.Success})
}
