package scan

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gernest/hiro/bus"
	"github.com/gernest/hiro/headers"
	"github.com/gernest/hiro/resource"
	"github.com/gernest/hiro/util"
	"github.com/labstack/echo/v4"
	"github.com/ory/ladon"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func Scan(rctx echo.Context) error {
	ev := bus.ScanEvent{
		Timestamp: time.Now().Format(time.RFC3339),
	}
	r := rctx.Request()
	ctx := util.RequestContext(rctx)
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
		zap.Namespace("scan"),
	)
	pid := rctx.Param("uuid")
	ev.ID = pid
	id, err := uuid.FromString(pid)
	if err != nil {
		log.Debug("scan cant parse uuid",
			zap.Error(err),
			zap.String("uuid", pid),
		)
		ev.Status = http.StatusBadRequest
		return util.BadRequest(rctx)
	}
	o, err := ctx.DB.GetQR(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return util.NotFound(rctx)
		}
		log.Debug("scan fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		ev.Status = http.StatusInternalServerError
		return util.Internal(rctx)
	}
	if ctx.Claims != nil {
		// the client/user who is making this call is authenticated. There are wto
		// scenarios at play.
		//
		// 	- The qrcode is private and the client has permission to access it
		// - The qrcode is public in which case we will use the client metadata for
		// tracking.
		tkID, err := uuid.FromString(ctx.Claims.Id)
		if err != nil {
			log.Debug("checking token",
				zap.Error(err),
			)
			ev.Status = http.StatusBadRequest
			return util.BadRequest(rctx)
		}
		token, err := ctx.DB.GetToken(r.Context(), tkID)
		if err != nil {
			log.Debug("scan checking token",
				zap.Error(err),
			)
			ev.Status = http.StatusBadRequest
			return util.BadRequest(rctx)
		}
		usr := token.Subject.String()
		ev.User = usr
		if o.Private {
			err = ctx.Warden.IsAllowed(&ladon.Request{
				Resource: resource.QR,
				Action:   "view",
				Subject:  usr,
				Context: ladon.Context{
					"user": o.UserID.String(),
				},
			})
			if err != nil {
				log.Debug("scan checking access permissions",
					zap.Error(err),
				)
				ev.Status = http.StatusForbidden
				return util.Forbid(rctx)
			}
		}
	}
	if !headers.IsJSONContent(r.Header) &&
		o.ShouldRedirect && o.RedirectURL != "" {
		ev.Status = http.StatusFound
		return rctx.Redirect(http.StatusFound, o.RedirectURL)
	}
	ev.Status = http.StatusOK
	return rctx.JSON(http.StatusOK, o)
}
