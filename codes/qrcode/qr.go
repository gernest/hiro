package qrcode

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/resource"
	"github.com/gernest/hiro/util"
	"github.com/labstack/echo/v4"
	"github.com/ory/ladon"
	uuid "github.com/satori/go.uuid"
)

func Create(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("qr-create"),
	)
	if ctx.Claims == nil {
		return rctx.JSON(http.StatusForbidden, models.APIError{Message: keys.Forbidden})
	}
	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		log.Debug("checking token",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadToken})
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		log.Debug("checking token",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadToken})
	}
	usr := token.Subject.String()
	err = ctx.Warden.IsAllowed(&ladon.Request{
		Resource: resource.QR,
		Action:   "create",
		Subject:  usr,
		Context: ladon.Context{
			"user": usr,
		},
	})
	if err != nil {
		log.Debug("checking access permissions",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusForbidden, models.APIError{Message: keys.Forbidden})
	}
	m := &models.QRReq{}
	err = json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		log.Debug("fail to unmarshal request body",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadJSON})
	}
	id := uuid.NewV4()
	now := time.Now().UTC()
	c := &models.QR{
		UUID:           id,
		Name:           m.Name,
		UserID:         token.Issuer,
		ShouldRedirect: m.ShouldRedirect,
		RedirectURL:    m.RedirectURL,
		CreatedAt:      now,
		UpdatedAt:      now,
		Size:           m.Size,
	}
	if m.Context != nil {
		c.Context = m.Context
	}
	err = ctx.DB.CreateQR(r.Context(), c, m.Groups...)
	if err != nil {
		log.Debug("fail to save",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusInternalServerError, models.APIError{Message: keys.InternalError})
	}
	return rctx.JSON(http.StatusOK, c)
}

func View(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("qr-view"),
	)
	if ctx.Claims == nil {
		return rctx.JSON(http.StatusForbidden, models.APIError{Message: keys.Forbidden})
	}
	pid := rctx.Param("uuid")
	id, err := uuid.FromString(pid)
	if err != nil {
		return util.BadRequest(rctx)
	}
	o, err := ctx.DB.GetQR(r.Context(), id)
	if err != nil {
		log.Debug("qr.view fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		if err == sql.ErrNoRows {
			return rctx.JSON(http.StatusNotFound, models.APIError{Message: keys.IsNotExist})
		}
		return rctx.JSON(http.StatusInternalServerError, models.APIError{Message: keys.InternalError})

	}

	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		log.Debug("checking token",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadToken})
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		log.Debug(" checking token",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadToken})
	}
	usr := token.Issuer.String()
	err = ctx.Warden.IsAllowed(&ladon.Request{
		Resource: resource.QR,
		Action:   "view",
		Subject:  usr,
		Context: ladon.Context{
			"user": o.UserID.String(),
		},
	})
	if err != nil {
		log.Debug("checking access permissions",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusForbidden, models.APIError{Message: keys.Forbidden})
	}
	return rctx.JSON(http.StatusOK, o)
}

// List list generated qrcodes.
func List(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("qr-list"),
	)
	if ctx.Claims == nil {
		return rctx.JSON(http.StatusForbidden, models.APIError{Message: keys.Forbidden})
	}
	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		log.Debug("checking token",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadToken})
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		log.Debug("checking token",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadToken})
	}
	opts, err := util.ListOptions(r)
	if err != nil {
		log.Debug("checking query params",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusBadRequest, models.APIError{Message: keys.BadToken})
	}
	o, err := ctx.DB.ListQR(r.Context(), token.Subject, opts)
	if err != nil {
		log.Debug("qr.view fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		return rctx.JSON(http.StatusInternalServerError, models.APIError{Message: keys.InternalError})
	}
	return rctx.JSON(http.StatusOK,
		models.QRList{QRCodes: o, Options: opts,
			Total: ctx.DB.TotalCodes(r.Context(), token.Subject),
		},
	)

}

// Delete deletes generated qrcodes.
func Delete(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("qr-delete"),
	)
	if ctx.Claims == nil {
		return util.Forbid(rctx)
	}
	pid := rctx.Param("uuid")
	id, err := uuid.FromString(pid)
	if err != nil {
		log.Debug("qr.delete cant parse uuid",
			zap.Error(err),
			zap.String("uuid", pid),
		)
		return util.BadRequest(rctx)
	}
	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		log.Debug("checking token",
			zap.Error(err),
		)
		return util.BadToken(rctx)
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		log.Debug("qr.delete checking token",
			zap.Error(err),
		)
		return util.BadToken(rctx)
	}
	code, err := ctx.DB.GetQR(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("missing qrcode",
				zap.String("uuid", id.String()),
				zap.String("subject", token.Subject.String()),
			)
			return util.NotFound(rctx)
		}
		log.Debug("qr.delete getting qrcode",
			zap.Error(err),
		)
		return util.Internal(rctx)
	}
	usr := token.Issuer.String()
	err = ctx.Warden.IsAllowed(&ladon.Request{
		Resource: resource.QR,
		Action:   "view",
		Subject:  usr,
		Context: ladon.Context{
			"user": code.UserID.String(),
		},
	})
	if err != nil {
		log.Debug("checking access permissions",
			zap.Error(err),
		)
		return util.Forbid(rctx)
	}
	err = ctx.DB.DeleteQR(r.Context(), id)
	if err != nil {
		log.Debug("fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		return util.Internal(rctx)
	}
	return rctx.JSON(http.StatusOK, models.Status{Status: keys.Success})
}

// Update updates information stored with the qrcode.
func Update(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("qr-update"),
	)
	if ctx.Claims == nil {
		return util.Forbid(rctx)
	}
	id, err := uuid.FromString(rctx.Param("uuid"))
	if err != nil {
		log.Debug("can't parse uuid",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	info := &models.QR{}

	err = json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		log.Debug("fail to unmarshal request body",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	o, err := ctx.DB.GetQR(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return util.NotFound(rctx)
		}
		log.Debug("qr.update fail to get info data",
			zap.Error(err),
		)
		return util.Internal(rctx)
	}
	o.Name = info.Name
	o.ShouldRedirect = info.ShouldRedirect
	o.RedirectURL = info.RedirectURL
	if info.Context != nil {
		if o.Context == nil {
			o.Context = info.Context
		} else {
			for k, v := range info.Context {
				o.Context[k] = v
			}
		}
	}
	err = ctx.DB.UpdateQR(r.Context(), o)
	if err != nil {
		log.Debug("failed to update",
			zap.Error(err),
		)
		return util.Internal(rctx)

	}
	return rctx.JSON(http.StatusOK, o)
}
