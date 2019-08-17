package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/gernest/hiro/access"

	"golang.org/x/crypto/bcrypt"

	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/resource"
	"github.com/gernest/hiro/util"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func Create(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("accounts-create"),
	)
	c := &models.CreateAccount{}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		log.Debug(" can't unmarshal json", zap.Error(err))
		return util.BadRequest(rctx)
	}
	v := c.Validate()
	if v != nil {
		return rctx.JSON(http.StatusUnprocessableEntity, v)
	}
	a := &models.Account{
		UUID:  uuid.NewV4(),
		Email: c.Email,
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Debug("hash password", zap.Error(err))
		return util.Internal(rctx)
	}
	a.Password = string(pass)
	err = ctx.DB.CreateAccount(r.Context(), a)
	if err != nil {
		if NameExists(err) {
			return rctx.JSON(http.StatusBadRequest,
				models.APIError{
					Message: keys.FailedValidation,
					Errors: []models.Message{
						{
							Resource: resource.Account,
							Field:    "username",
							Desc:     keys.UsernameAlreadyExists,
						},
					},
				},
			)
		}
		if EmailExists(err) {
			return rctx.JSON(http.StatusBadRequest,
				models.APIError{
					Message: keys.FailedValidation,
					Errors: []models.Message{
						{
							Resource: resource.Account,
							Field:    "email",
							Desc:     keys.UsernameAlreadyExists,
						},
					},
				},
			)
		}
		log.Debug("save to db", zap.Error(err))
		return util.Internal(rctx)
	}
	if ctx.Warden != nil {
		// we are creating access policy for the user to manage his/her profile.
		for _, policy := range access.NewUserPolicies(a.UUID) {
			err = ctx.Warden.Manager.Create(policy)
			if err != nil {
				log.Debug("creating access policies",
					zap.Error(err),
				)
				return util.Internal(rctx)
			}
		}
	}
	return util.Ok(rctx)
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

func Login(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("accounts-login"),
	)
	c := &models.Login{}
	err := json.NewDecoder(r.Body).Decode(c)
	if err != nil {
		log.Debug("can't unmarshal json", zap.Error(err))
		return util.BadRequest(rctx)
	}
	if c.Name == "" {
		return rctx.JSON(http.StatusBadRequest, models.APIError{
			Message: keys.FailedValidation,
			Errors: []models.Message{
				{
					Resource: "Login",
					Field:    "name",
					Desc:     keys.MissingEmail,
				},
			},
		})
	}

	if c.Password == "" {
		return rctx.JSON(http.StatusBadRequest, models.APIError{
			Message: keys.FailedValidation,
			Errors: []models.Message{
				{
					Resource: "Login",
					Field:    "password",
					Desc:     keys.MissingPassword,
				},
			},
		})
	}
	a, err := ctx.DB.GetAccount(r.Context(), c.Name)
	if err != nil {
		log.Debug("can't find a user", zap.Error(err))
		return util.BadRequest(rctx)
	}

	err = bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(c.Password))
	if err != nil {
		log.Debug("matching password", zap.Error(err))
		return util.Invalid(rctx)
	}
	token := models.DefaultLoginToken(a.UUID)
	err = ctx.DB.CreateToken(r.Context(), token)
	if err != nil {
		log.Debug(" creating token", zap.Error(err))
		return util.Internal(rctx)
	}
	s, err := util.GenerateJWTToken(ctx.JWT, token)
	if err != nil {
		log.Debug("creating jwt token", zap.Error(err))
		return util.Internal(rctx)
	}
	return rctx.JSON(http.StatusOK, models.LoginRes{Token: s})
}
