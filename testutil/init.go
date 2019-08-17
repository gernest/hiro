package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"go.uber.org/zap"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gernest/hiro/access"
	"github.com/gernest/hiro/headers"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/query"
	"github.com/gernest/hiro/util"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type Context struct {
	Account *models.Account
	Token   string
	Claims  *jwt.StandardClaims
	DB      *query.SQL
	items   []models.Item
}

func TestContext(ctx *Context, r *http.Request, w http.ResponseWriter) echo.Context {
	e := echo.New()
	c := e.NewContext(r, w)
	for _, v := range ctx.items {
		c.Set(v.Key, v.Value)
	}
	c.Set(keys.ContextKey, ctx)
	return c
}

const TestEmail = "hiro@examples.com"
const TestSecret = "somesecret"

func New(db *query.SQL) (*Context, error) {

	l, _ := zap.NewProduction()
	jwt := &models.JWT{Secret: []byte(TestSecret)}

	uid := uuid.NewV4()
	a := &models.Account{
		UUID:     uid,
		Email:    TestEmail,
		Password: "pass",
	}
	err := db.CreateAccount(context.Background(), a)
	if err != nil {
		return nil, err
	}
	tk := models.DefaultLoginToken(a.UUID)
	err = db.CreateToken(context.Background(), tk)
	if err != nil {
		return nil, err
	}
	token, err := util.GenerateJWTToken(jwt, tk)
	if err != nil {
		return nil, err
	}
	sd, err := util.ParseJWTToken(jwt, token)
	if err != nil {
		return nil, err
	}
	warden, err := access.New(db.DB(), nil)
	if err != nil {
		return nil, err
	}
	for _, policy := range access.NewUserPolicies(a.UUID) {
		err = warden.Manager.Create(policy)
		if err != nil {
			return nil, err
		}
	}
	return &Context{
		Account: a,
		Token:   token,
		Claims:  sd,
		DB:      db,
		items: []models.Item{
			models.Item{Key: keys.DB, Value: db},
			models.Item{Key: keys.LoggerKey, Value: l},
			models.Item{Key: keys.JwtKey, Value: jwt},
			models.Item{Key: keys.Session, Value: sd},
			models.Item{Key: keys.Warden, Value: warden},
		},
	}, nil
}

func NewDB() (*query.SQL, error) {
	s, err := query.New("postgres", os.Getenv("HIRO_DB_CONN"))
	if err != nil {
		return nil, err
	}
	if err := s.Up(context.Background()); err != nil {
		return nil, err
	}
	return s, nil
}

func request(items ...models.Item) func(string, string, io.Reader) *http.Request {
	return func(method, target string, body io.Reader) *http.Request {
		r := httptest.NewRequest(method, target, body)
		ctx := r.Context()
		for _, v := range items {
			ctx = context.WithValue(ctx, v.Key, v.Value)
		}
		return r.WithContext(ctx)
	}
}

func ReqData(v interface{}) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

func (c *Context) Close() error {
	if err := c.DB.Down(context.Background()); err != nil {
		return err
	}
	return c.DB.Close()
}

func (c *Context) SetHeader(r *http.Request) {
	r.Header.Set(headers.ContentType, headers.ApplicationJSON)
	r.Header.Set(headers.Authorization, "Bearer "+c.Token)
}

func (c *Context) Req(method string, target string, body io.Reader) *http.Request {
	r := httptest.NewRequest(method, target, body)
	c.SetHeader(r)
	return r
}
