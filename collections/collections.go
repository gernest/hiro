package collections

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/resource"
	"github.com/gernest/hiro/util"
	"github.com/labstack/echo/v4"
	"github.com/ory/ladon"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// Create creates a new collection.
func Create(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("collections-create"),
	)
	if ctx.Claims == nil {
		return util.Forbid(rctx)
	}
	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		log.Debug("checking token",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		log.Debug(" checking token",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	usr := token.Subject.String()
	err = ctx.Warden.IsAllowed(&ladon.Request{
		Resource: resource.Collections,
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
		return util.Forbid(rctx)
	}

	m := &models.CollectionReq{}
	err = json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		log.Debug("fail to unmarshal request body",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	now := time.Now()
	c, err := ctx.DB.CreateCollection(r.Context(), &models.Collection{
		UUID:      uuid.NewV4(),
		Name:      m.Name,
		Color:     m.Color,
		AccountID: token.Subject,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		log.Debug("saving new collection",
			zap.Error(err),
		)
		return util.Internal(rctx)
	}
	return rctx.JSON(http.StatusOK, c)
}

func View(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("collections-view"),
	)
	if ctx.Claims == nil {
		return util.Forbid(rctx)
	}
	params := r.URL.Query()
	pid := params.Get("uuid")
	id, err := uuid.FromString(pid)
	if err != nil {
		log.Debug("cant parse uuid",
			zap.Error(err),
			zap.String("uuid", pid),
		)
		return util.BadRequest(rctx)
	}
	o, err := ctx.DB.GetCollection(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return util.NotFound(rctx)
		}
		log.Debug("fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		return util.Internal(rctx)
	}
	return rctx.JSON(http.StatusOK, o)
}

func List(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("collections-list"),
	)
	if ctx.Claims == nil {
		return util.Forbid(rctx)
	}
	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		log.Debug("checking token",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		log.Debug("checking token",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	usr := token.Subject.String()
	err = ctx.Warden.IsAllowed(&ladon.Request{
		Resource: resource.Collections,
		Action:   "list",
		Subject:  usr,
		Context: ladon.Context{
			"user": usr,
		},
	})
	if err != nil {
		log.Debug("checking access permissions",
			zap.Error(err),
		)
		return util.Forbid(rctx)
	}
	opts, err := util.ListOptions(r)
	if err != nil {
		log.Debug("checking query params",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	o, err := ctx.DB.ListCollections(r.Context(), token.Subject, opts)
	if err != nil {
		log.Debug("fail to retrieve collcetions list",
			zap.Error(err),
		)
		return util.Internal(rctx)
	}
	return rctx.JSON(http.StatusOK, o)
}

func Delete(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("collections-delete"),
	)
	if ctx.Claims == nil {
		return util.Forbid(rctx)
	}
	params := r.URL.Query()
	pid := params.Get("uuid")
	id, err := uuid.FromString(pid)
	if err != nil {
		log.Debug("qr.delete cant parse uuid",
			zap.Error(err),
			zap.String("uuid", pid),
		)
		return util.BadRequest(rctx)
	}
	err = ctx.DB.DeleteCollection(r.Context(), id)
	if err != nil {
		log.Debug("fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		return util.Internal(rctx)
	}
	return util.Ok(rctx)
}

// Assign assigns a collection to the qrcode.
func Assign(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("collections-assign"),
	)
	if ctx.Claims == nil {
		return util.Forbid(rctx)
	}
	m := &models.CollectionAssignReq{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		log.Debug("collections.create fail to unmarshal request body",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	err = ctx.DB.AssignCollection(r.Context(), m.ID, m.QRID)
	if err != nil {
		log.Debug(" saving new collection",
			zap.Error(err),
		)
		return util.Internal(rctx)
	}
	return util.Ok(rctx)
}

func DeAssign(rctx echo.Context) error {
	ctx := util.RequestContext(rctx)
	r := rctx.Request()
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("collections-deasign"),
	)
	if ctx.Claims == nil {
		return util.Forbid(rctx)
	}

	m := &models.CollectionAssignReq{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		log.Debug("fail to unmarshal request body",
			zap.Error(err),
		)
		return util.BadRequest(rctx)
	}
	err = ctx.DB.DeAssignCollection(r.Context(), m.ID, m.QRID)
	if err != nil {
		log.Debug("collection.assign saving new collection",
			zap.Error(err),
		)
		return util.Internal(rctx)
	}
	return util.Ok(rctx)
}
