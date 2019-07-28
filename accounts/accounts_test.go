package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/resource"
	"github.com/gernest/hiro/testutil"
	"github.com/gernest/hiro/util"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func TestJWT(t *testing.T) {
	secret := "some secret stuff"
	j := &models.JWT{
		Secret: []byte(secret),
	}
	now := time.Now()

	usrID := uuid.NewV4()
	tk := &models.Token{
		UUID:      usrID,
		Issuer:    usrID,
		Subject:   usrID,
		CreatedAt: now,
		ExpiresOn: now.Add(time.Hour),
	}

	token, err := util.GenerateJWTToken(j, tk)
	if err != nil {
		t.Fatal(err)
	}
	c, err := util.ParseJWTToken(j, token)
	if err != nil {
		t.Fatal(err)
	}
	if c.Issuer != tk.Issuer.String() {
		t.Errorf("expected %s got %s", tk.Issuer, c.Issuer)
	}
}

func TestValidation(t *testing.T) {
	g := &models.CreateAccount{
		Name:            "gee",
		Email:           "gernest@examples.com",
		Password:        "pass",
		ConfirmPassword: "pass",
	}
	sample := []struct {
		form *models.CreateAccount
		err  *models.APIError
		desc string
	}{
		{
			form: g,
			desc: "all correct fields",
		},
		{
			form: &models.CreateAccount{
				Name:            g.Name,
				Password:        g.Password,
				ConfirmPassword: g.ConfirmPassword,
			},
			err: &models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: resource.Account,
						Field:    "email",
						Desc:     keys.MissingEmail,
					},
				},
			},
			desc: "missing email",
		},
		{
			form: &models.CreateAccount{
				Name:            g.Name,
				Email:           g.Email,
				ConfirmPassword: g.ConfirmPassword,
			},
			err: &models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: resource.Account,
						Field:    "password",
						Desc:     keys.MissingPassword,
					},
				},
			},
			desc: "missing password",
		},
		{
			form: &models.CreateAccount{
				Name:     g.Name,
				Email:    g.Email,
				Password: g.Password,
			},
			err: &models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: resource.Account,
						Field:    "confirm_password",
						Desc:     keys.ConfirmPasswordMismatch,
					},
				},
			},
			desc: "missing confirm password",
		},
		{
			form: &models.CreateAccount{},
			err: &models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: resource.Account,
						Field:    "name",
						Desc:     keys.MissingUsername,
					},
					{
						Resource: resource.Account,
						Field:    "email",
						Desc:     keys.MissingEmail,
					},
					{
						Resource: resource.Account,
						Field:    "password",
						Desc:     keys.MissingPassword,
					},
				},
			},
			desc: "missing all fields",
		},
		{
			form: &models.CreateAccount{
				Email: "bad email",
			},
			err: &models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: resource.Account,
						Field:    "name",
						Desc:     keys.MissingUsername,
					},
					{
						Resource: resource.Account,
						Field:    "email",
						Desc:     keys.InvalidEmail,
					},
					{
						Resource: resource.Account,
						Field:    "password",
						Desc:     keys.MissingPassword,
					},
				},
			},
			desc: "missing all fields with invalid email",
		},
		{
			form: &models.CreateAccount{
				Name: "some.name",
			},
			err: &models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: resource.Account,
						Field:    "name",
						Desc:     keys.InvalidUsername,
					},
					{
						Resource: resource.Account,
						Field:    "email",
						Desc:     keys.MissingEmail,
					},
					{
						Resource: resource.Account,
						Field:    "password",
						Desc:     keys.MissingPassword,
					},
				},
			},
			desc: "missing all fields with invalid name",
		},
		{
			form: &models.CreateAccount{
				Name: "dashboard",
			},
			err: &models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: resource.Account,
						Field:    "name",
						Desc:     keys.UsernameAlreadyExists,
					},
					{
						Resource: resource.Account,
						Field:    "email",
						Desc:     keys.MissingEmail,
					},
					{
						Resource: resource.Account,
						Field:    "password",
						Desc:     keys.MissingPassword,
					},
				},
			},
			desc: "missing all fields with reserved username",
		},
	}

	for _, v := range sample {
		t.Run(v.desc, func(ts *testing.T) {
			err := v.form.Validate()
			if v.err == nil && err != nil {
				t.Error(err)
			}
			if v.err != nil && err != nil {
				if v.err.Error() != err.Error() {
					t.Errorf("expected %v got %v", v.err, err)
				}
			}
		})
	}

}

const name = "gernestaccounts"
const secret = "someSecret"

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

func TestAccounts(t *testing.T) {
	db, err := testutil.NewDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	l, _ := zap.NewProduction()
	jwt := &models.JWT{Secret: []byte(secret)}
	req := request(
		models.Item{Key: keys.DB, Value: db},
		models.Item{Key: keys.LoggerKey, Value: l},
		models.Item{Key: keys.JwtKey, Value: jwt},
	)
	t.Run("register", func(ts *testing.T) {
		ts.Run("no body", func(ts *testing.T) {
			r := req("POST", "/register", nil)
			w := httptest.NewRecorder()
			Create(w, r)
			if w.Code != http.StatusBadRequest {
				ts.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
			}
		})

		r := req("POST", "/register", reqData(&models.CreateAccount{
			Name:     name,
			Email:    name + "@bqtest.com",
			Password: "pass",
		}))
		w := httptest.NewRecorder()
		Create(w, r)
		if w.Code != http.StatusUnprocessableEntity {
			ts.Fatalf("expected %d got %d", http.StatusUnprocessableEntity, w.Code)
		}
		a := &models.APIError{}
		err := json.Unmarshal(w.Body.Bytes(), a)
		if err != nil {
			ts.Fatal(err)
		}
		if a.Message != keys.FailedValidation {
			ts.Errorf("expected %s got %s", keys.FailedValidation, a.Message)
		}

		r = req("POST", "/register", reqData(&models.CreateAccount{
			Name:            name,
			Email:           name + "@bqtest.com",
			Password:        "pass",
			ConfirmPassword: "pass",
		}))
		w = httptest.NewRecorder()
		Create(w, r)
		if w.Code != http.StatusOK {
			ts.Fatalf("expected %d got %d \n%s", http.StatusOK, w.Code, w.Body.String())
		}
		s := &models.Status{}
		err = json.Unmarshal(w.Body.Bytes(), s)
		if err != nil {
			ts.Fatal(err)
		}
		if s.Status != keys.Success {
			ts.Errorf("expected %s got %s", keys.Success, s.Status)
		}
	})

	t.Run("login", func(ts *testing.T) {
		ts.Run("no body", func(ts *testing.T) {
			r := req("POST", "/login", nil)
			w := httptest.NewRecorder()
			Login(w, r)
			if w.Code != http.StatusBadRequest {
				ts.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
			}
			a := apiError(ts, w)
			if a.Message != keys.BadJSON {
				t.Errorf("expected %s got %s", keys.BadJSON, a.Message)
			}
		})
		ts.Run("missing username ", func(ts *testing.T) {
			r := req("POST", "/login", reqData(&models.Login{
				Password: "pass",
			}))
			w := httptest.NewRecorder()
			Login(w, r)
			if w.Code != http.StatusBadRequest {
				ts.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
			}
			a := apiError(ts, w)
			if a.Message != keys.FailedValidation {
				t.Errorf("expected %s got %s", keys.FailedValidation, a.Message)
			}
		})
		ts.Run("missing password ", func(ts *testing.T) {
			r := req("POST", "/login", reqData(&models.Login{
				Name: name,
			}))
			w := httptest.NewRecorder()
			Login(w, r)
			if w.Code != http.StatusBadRequest {
				ts.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
			}
			a := apiError(ts, w)
			if a.Message != keys.FailedValidation {
				t.Errorf("expected %s got %s", keys.FailedValidation, a.Message)
			}
		})
		ts.Run("unknown user ", func(ts *testing.T) {
			r := req("POST", "/login", reqData(&models.Login{
				Name:     "some dude",
				Password: "pass",
			}))
			w := httptest.NewRecorder()
			Login(w, r)
			if w.Code != http.StatusNotFound {
				ts.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
			}
			a := apiError(ts, w)
			for _, v := range a.Errors {
				if v.Field == "name" {
					if v.Desc != keys.WrongCredentials {
						t.Errorf("expected %s got %v", keys.WrongCredentials, v.Desc)
					}
				}
			}
		})
		r := req("POST", "/login", reqData(&models.Login{
			Name:     name,
			Password: "pass",
		}))
		w := httptest.NewRecorder()
		Login(w, r)

		if w.Code != http.StatusOK {
			ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
		}
		res := &models.LoginRes{}
		err := json.Unmarshal(w.Body.Bytes(), res)
		if err != nil {
			ts.Fatal(err)
		}
		if res.Token == "" {
			ts.Errorf("expected token %s", w.Body.String())
		}
	})
}

func reqData(v interface{}) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

func apiError(t *testing.T, w *httptest.ResponseRecorder) *models.APIError {
	m := &models.APIError{}
	err := json.Unmarshal(w.Body.Bytes(), m)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
