package qrcode

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/gernest/alien"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/resource"
	"github.com/gernest/hiro/util"
	"github.com/ory/ladon"

	uuid "github.com/satori/go.uuid"
)

// Create generates a new qrcode which encodes a unique url to associate the
// generated image with metadata that is stored in the database.
//
// The request must be authorized through JWT token. The authentication logic is
// done inside this handler. This expects the payload to be json data in the
// request body.
//
// Sample request body
// {
//   "width": 100,
//   "height": 100,
//   "should_redirect": true,
//   "redirect_url": "/home"
// }
func Create(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	if ctx.Claims == nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		return
	}
	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("qr.create checking token",
			zap.Error(err),
		)
		return
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("qr.create checking token",
			zap.Error(err),
		)
		return
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
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		log.Error("qr.create checking access permissions",
			zap.Error(err),
		)
		return
	}
	m := &models.QRReq{}
	err = json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadJSON}, http.StatusBadRequest)
		log.Error("qr.create fail to unmarshal request body",
			zap.Error(err),
		)
		return
	}
	id := uuid.NewV4()
	content := fmt.Sprintf("%s/scan/%s", ctx.Host, id)
	now := time.Now().UTC()
	queryImageURL := make(url.Values)
	queryImageURL.Set("uuid", id.String())
	queryImageURL.Set("usr", token.Issuer.String())
	queryImageURL.Set("width", "100")
	queryImageURL.Set("height", "100")
	c := &models.QR{
		UUID:           id,
		Name:           m.Name,
		UserID:         token.Issuer,
		URL:            content,
		ImageURL:       fmt.Sprintf("%s/img?%s", ctx.ImageHost, queryImageURL.Encode()),
		ShouldRedirect: m.ShouldRedirect,
		RedirectURL:    m.RedirectURL,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if m.Context != nil {
		c.Context = m.Context
	}
	err = ctx.DB.CreateQR(r.Context(), c, m.Groups...)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError},
			http.StatusInternalServerError)
		log.Error("qr.create fail to save",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w, c, http.StatusOK)
}

func getDimension(q url.Values) (int, int, error) {
	w, h := q.Get("width"), q.Get("height")
	width, err := strconv.Atoi(w)
	if err != nil {
		return 0, 0, err
	}
	height, err := strconv.Atoi(h)
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

// View displays information related to the QRCode.
func View(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	if ctx.Claims == nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		return
	}
	params := alien.GetParams(r)
	pid := params.Get("uuid")
	id, err := uuid.FromString(pid)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadRequest},
			http.StatusBadRequest)
		log.Error("qr.view cant parse uuid",
			zap.Error(err),
			zap.String("uuid", pid),
		)
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
		log.Error("qr.view fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		return
	}

	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("qr.list checking token",
			zap.Error(err),
		)
		return
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("qr.list checking token",
			zap.Error(err),
		)
		return
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
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		log.Error("qr.view checking access permissions",
			zap.Error(err),
		)
		return
	}

	util.WriteJSON(w, o, http.StatusOK)
}

// List list generated qrcodes.
func List(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	if ctx.Claims == nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		return
	}
	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("qr.list checking token",
			zap.Error(err),
		)
		return
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("qr.list checking token",
			zap.Error(err),
		)
		return
	}
	opts, err := util.ListOptions(r)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("qr.list checking query params",
			zap.Error(err),
		)
		return
	}
	o, err := ctx.DB.ListQR(r.Context(), token.Subject, opts)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError},
			http.StatusInternalServerError)
		log.Error("qr.view fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w,
		models.QRList{QRCodes: o, Options: opts,
			Total: ctx.DB.TotalCodes(r.Context(), token.Subject),
		},
		http.StatusOK)
}

// Delete deletes generated qrcodes.
func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	if ctx.Claims == nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		return
	}
	params := alien.GetParams(r)
	pid := params.Get("uuid")
	id, err := uuid.FromString(pid)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadRequest},
			http.StatusBadRequest)
		log.Error("qr.delete cant parse uuid",
			zap.Error(err),
			zap.String("uuid", pid),
		)
		return
	}
	tkID, err := uuid.FromString(ctx.Claims.Id)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("qr.list checking token",
			zap.Error(err),
		)
		return
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("qr.delete checking token",
			zap.Error(err),
		)
		return
	}
	code, err := ctx.DB.GetQR(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			util.WriteJSON(w, &models.APIError{Message: http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
			log.Info("qr.delete missing qrcode",
				zap.String("uuid", id.String()),
				zap.String("subject", token.Subject.String()),
			)
			return
		}
		log.Error("qr.delete getting qrcode",
			zap.Error(err),
		)
		return
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
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		log.Error("qr.view checking access permissions",
			zap.Error(err),
		)
		return
	}
	err = ctx.DB.DeleteQR(r.Context(), id)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError},
			http.StatusInternalServerError)
		log.Error("qr.view fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w, &models.Status{Status: keys.Success}, http.StatusOK)
}

// Update updates information stored with the qrcode.
func Update(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	if ctx.Claims == nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		return
	}
	params := alien.GetParams(r)
	id, err := uuid.FromString(params.Get("uuid"))
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Invalid}, http.StatusBadRequest)
		log.Error("qr.update can't parse uuid",
			zap.Error(err),
		)
		return
	}
	info := &models.QR{}

	err = json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadJSON}, http.StatusBadRequest)
		log.Error("qr.update fail to unmarshal request body",
			zap.Error(err),
		)
		return
	}
	o, err := ctx.DB.GetQR(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			util.WriteJSON(w, models.APIError{Message: keys.IsNotExist},
				http.StatusNotFound)
		} else {
			util.WriteJSON(w, models.APIError{Message: keys.BadRequest},
				http.StatusInternalServerError)
		}
		log.Error("qr.update fail to get info data",
			zap.Error(err),
		)
		return
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
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusBadRequest)
		log.Error("qr.update failed to update",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w, o, http.StatusOK)
}
