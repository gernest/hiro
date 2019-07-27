package util

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/gernest/hiro/bus"
	"github.com/ory/ladon"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/minio/minio-go"
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
	Minio     *minio.Client
	NSQ       *bus.Producer
	Warden    *ladon.Ladon
}

// RequestContext returns *Context with the fields filled with values taken from
// ctx. This will panic if the values are not present in the context. This will
// silently ignore fields which have no values set in the context.
func RequestContext(ctx context.Context) *Context {
	c := &Context{}
	if v := ctx.Value(keys.JwtKey); v != nil {
		c.JWT = v.(*models.JWT)
	}
	if v := ctx.Value(keys.Host); v != nil {
		c.Host = v.(string)
	}
	if v := ctx.Value(keys.ImageHost); v != nil {
		c.ImageHost = v.(string)
	}
	if v := ctx.Value(keys.LoggerKey); v != nil {
		c.Logger = v.(*zap.Logger)
	}
	if v := ctx.Value(keys.DB); v != nil {
		c.DB = v.(*query.SQL)
	}
	if v := ctx.Value(keys.Session); v != nil {
		c.Claims = v.(*jwt.StandardClaims)
	}
	if v := ctx.Value(keys.Minio); v != nil {
		c.Minio = v.(*minio.Client)
	}
	if v := ctx.Value(keys.NSQ); v != nil {
		c.NSQ = v.(*bus.Producer)
	}
	if v := ctx.Value(keys.Warden); v != nil {
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

// HostFlag is a flag for specifying the host on which the service is running.
// It defaults to http://localhost:8000.
func HostFlag() cli.Flag {
	return cli.StringFlag{
		Name:   "host",
		Usage:  "connection string to postgres database",
		EnvVar: "BQ_HOST",
		Value:  "http://localhost:8000",
	}
}

func ImageHostFlag() cli.Flag {
	return cli.StringFlag{
		Name:   "image-host",
		Usage:  "dns resolvable url to where images are served",
		EnvVar: "BQ_IMAGE_HOST",
		Value:  "http://localhost:8011",
	}
}

// ConnFlag returns a flag for specifying database connection string.
func ConnFlag() cli.Flag {
	return cli.StringFlag{
		Name:   "db-conn",
		Usage:  "connection string to postgres database",
		EnvVar: "BQ_DB_CONN",
	}
}

// DriverFlag returns a flag for specifying database driver.
func DriverFlag() cli.Flag {
	return cli.StringFlag{
		Name:   "driver",
		Usage:  "database driver to use",
		EnvVar: "BQ_DB_DRIVER",
		Value:  "postgres",
	}
}

func SecretFlag() cli.Flag {
	return cli.StringFlag{
		Name:   "secret",
		Usage:  "hmac jwt secret",
		EnvVar: "BQ_JWT_SECRET",
		Value:  "secret",
	}
}

// HighlightFlag is syntax highlight flag.
func HighlightFlag() cli.Flag {
	return cli.BoolFlag{
		Name:   "highlight",
		Usage:  "Highlight output",
		EnvVar: "BQ_HIGHLIGHT",
	}
}

// MinioFlags flags for minio store.
func MinioFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "minio-endpoint",
			EnvVar: "MINIO_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "minio-access-key",
			EnvVar: "MINIO_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "minio-access-secret",
			EnvVar: "MINIO_SECRET_KEY",
		},
	}
}

// NSQFlags flags for nsq.
func NSQFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "nsq-address",
			EnvVar: "NSQ_ADDRESS",
		},
		cli.StringFlag{
			Name:   "nsq-lookup",
			EnvVar: "NSQLOOKUPD_ADDRESS",
		},
	}
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
