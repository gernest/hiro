package accounts

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gernest/hiro/access"

	"golang.org/x/crypto/bcrypt"

	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/resource"
	"github.com/gernest/hiro/util"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// Create creates a new account. This only accepts json as payload. If minio was
// configured as the default blob store a new bucket will be created with user's
// name which will be used to store the qrcodes the user will generate.
//
// Password is hashed with bcrypt before storing to the database.
func Create(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	c := &models.CreateAccount{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadBody}, http.StatusBadRequest)
		log.Error("account.create can't read body", zap.Error(err))
		return
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadJSON}, http.StatusBadRequest)
		log.Error("account.create can't marshal json", zap.Error(err))
		return
	}
	v := c.Validate()
	if v != nil {
		util.WriteJSON(w, v, http.StatusUnprocessableEntity)
		return
	}
	a := &models.Account{
		UUID:  uuid.NewV4(),
		Name:  c.Name,
		Email: c.Email,
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusBadRequest)
		log.Error("account.create hash password", zap.Error(err))
		return
	}
	a.Password = string(pass)
	err = ctx.DB.CreateAccount(r.Context(), a)
	if err != nil {
		if NameExists(err) {
			util.WriteJSON(w,
				&models.APIError{
					Message: keys.FailedValidation,
					Errors: []models.Message{
						{
							Resource: resource.Account,
							Field:    "username",
							Desc:     keys.UsernameAlreadyExists,
						},
					},
				},
				http.StatusBadRequest)
			return
		}
		if EmailExists(err) {
			util.WriteJSON(w,
				&models.APIError{
					Message: keys.FailedValidation,
					Errors: []models.Message{
						{
							Resource: resource.Account,
							Field:    "email",
							Desc:     keys.UsernameAlreadyExists,
						},
					},
				},
				http.StatusBadRequest)
			return
		}
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusBadRequest)
		log.Error("account.create save to db", zap.Error(err))
		return
	}
	if ctx.Minio != nil {
		if err = ctx.Minio.MakeBucket(a.UUID.String(), "us-east-1"); err != nil {
			util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusBadRequest)
			log.Error("creating bucket",
				zap.String("bucket", a.UUID.String()), zap.Error(err))
			return
		}
	}
	if ctx.Warden != nil {
		// we are creating access policy for the user to manage his/her profile.
		for _, policy := range access.NewUserPolicies(a.UUID) {
			err = ctx.Warden.Manager.Create(policy)
			if err != nil {
				log.Error("accounts.create creating access policies",
					zap.Error(err),
				)
				util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusInternalServerError)
				return
			}
		}
	}
	util.WriteJSON(w, &models.Status{Status: keys.Success}, http.StatusOK)
}

// NameExists returns true if the error is for a account name already exists.
func NameExists(err error) bool {
	return err.Error() ==
		`pq: duplicate key value violates unique constraint "accounts_name_key"`
}

// EmailExists returns true if the error is for a account email already exists.
func EmailExists(err error) bool {
	return err.Error() ==
		`pq: duplicate key value violates unique constraint "accounts_email_key"`
}

// Login authenticates a user and issue a jwt token.
func Login(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	c := &models.Login{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadBody}, http.StatusBadRequest)
		log.Error("account.Login can't read body", zap.Error(err))
		return
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadJSON}, http.StatusBadRequest)
		log.Error("account.Login can't marshal json", zap.Error(err))
		return
	}
	if c.Name == "" {
		util.WriteJSON(w,
			&models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: "Login",
						Field:    "name",
						Desc:     keys.MissingEmail,
					},
				},
			},
			http.StatusBadRequest,
		)
		return
	}

	if c.Password == "" {
		util.WriteJSON(w,
			&models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: "Login",
						Field:    "password",
						Desc:     keys.MissingPassword,
					},
				},
			},
			http.StatusBadRequest,
		)
		return
	}
	a, err := ctx.DB.GetAccount(r.Context(), c.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			util.WriteJSON(w, &models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: "Login",
						Field:    "name",
						Desc:     keys.WrongCredentials,
					},
				},
			}, http.StatusNotFound)
			return
		}
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusInternalServerError)
		log.Error("account.Login can't find a user", zap.Error(err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(c.Password))
	if err != nil {
		util.WriteJSON(w,
			&models.APIError{
				Message: keys.FailedValidation,
				Errors: []models.Message{
					{
						Resource: "Login_req",
						Field:    "password",
						Desc:     "wrong password",
					},
				},
			},
			http.StatusBadRequest,
		)
		return
	}
	token := models.DefaultLoginToken(a.UUID)
	err = ctx.DB.CreateToken(r.Context(), token)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusInternalServerError)
		log.Error("account.Login creating token", zap.Error(err))
		return
	}
	s, err := util.GenerateJWTToken(ctx.JWT, token)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusInternalServerError)
		log.Error("account.Login creating jwt token", zap.Error(err))
		return
	}
	util.WriteJSON(w, &models.LoginRes{Token: s}, http.StatusOK)
}
