package scan

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gernest/alien"
	"github.com/gernest/hiro/bus"
	"github.com/gernest/hiro/headers"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/resource"
	"github.com/gernest/hiro/util"
	"github.com/ory/ladon"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// Scan is the endpoint which matches the url embedded in the qrcode.This tracks
// useful metrics which are used for analytics purpose.
//
// The messages/collected metrics are snt on message bus if the bus was
// configured when the server was started. Successful scan can either redirect
// to the url specified in the qrcode information, or it will serve default
// qrcode viewing page.
//
// If Content-Type header is application/json then the qrcode info will be
// served as json, this will however not redirect.
func Scan(w http.ResponseWriter, r *http.Request) {
	ev := bus.ScanEvent{
		User:      new(uuid.UUID).String(),
		Timestamp: time.Now().Format(time.RFC3339),
	}
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	params := alien.GetParams(r)
	pid := params.Get("uuid")
	ev.ID = pid
	id, err := uuid.FromString(pid)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadRequest},
			http.StatusBadRequest)
		log.Debug("scan cant parse uuid",
			zap.Error(err),
			zap.String("uuid", pid),
		)
		ev.Status = http.StatusBadRequest
		return
	}
	o, err := ctx.DB.GetQR(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			util.WriteJSON(w, &models.APIError{Message: keys.IsNotExist},
				http.StatusNotFound)
		} else {
			util.WriteJSON(w, &models.APIError{Message: keys.InternalError},
				http.StatusInternalServerError)
		}
		log.Debug("scan fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		ev.Status = http.StatusInternalServerError
		return
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
			util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
			log.Debug("qr.list checking token",
				zap.Error(err),
			)
			ev.Status = http.StatusBadRequest
			return
		}
		token, err := ctx.DB.GetToken(r.Context(), tkID)
		if err != nil {
			util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
			log.Debug("scan checking token",
				zap.Error(err),
			)
			ev.Status = http.StatusBadRequest
			return
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
				util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
				log.Debug("scan checking access permissions",
					zap.Error(err),
				)
				ev.Status = http.StatusForbidden
				return
			}
		}
	}
	if !headers.IsJSONContent(r.Header) &&
		o.ShouldRedirect && o.RedirectURL != "" {
		ev.Status = http.StatusFound
		http.Redirect(w, r, o.RedirectURL, http.StatusFound)
		return
	}
	ev.Status = http.StatusOK
	util.WriteJSON(w, o, http.StatusOK)
}
