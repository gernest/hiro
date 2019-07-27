package collections

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/resource"
	"github.com/gernest/hiro/util"
	"github.com/ory/ladon"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// Create creates a new collection.
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
		log.Error("collection.create checking token",
			zap.Error(err),
		)
		return
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("collections.create checking token",
			zap.Error(err),
		)
		return
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
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		log.Error("collections.create checking access permissions",
			zap.Error(err),
		)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadBody}, http.StatusBadRequest)
		log.Error("collections.create can't read body",
			zap.Error(err),
		)
		return
	}
	m := &models.CollectionReq{}
	err = json.Unmarshal(b, m)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadJSON}, http.StatusBadRequest)
		log.Error("collections.create fail to unmarshal request body",
			zap.Error(err),
		)
		return
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
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusInternalServerError)
		log.Error("collection.create saving new collection",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w, c, http.StatusOK)
}

func View(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	if ctx.Claims == nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		return
	}
	params := r.URL.Query()
	pid := params.Get("uuid")
	id, err := uuid.FromString(pid)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadRequest},
			http.StatusBadRequest)
		log.Error("collection.view cant parse uuid",
			zap.Error(err),
			zap.String("uuid", pid),
		)
		return
	}
	o, err := ctx.DB.GetCollection(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			util.WriteJSON(w, &models.APIError{Message: keys.IsNotExist},
				http.StatusNotFound)
		} else {
			util.WriteJSON(w, &models.APIError{Message: keys.InternalError},
				http.StatusInternalServerError)
		}
		log.Error("collection.view fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w, o, http.StatusOK)
}

// List renders json response of all collections owned by the subject who
// issued the request.
//
//This handler is protected, only valid authenticated tokens are processed.
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
		log.Error("collection.list checking token",
			zap.Error(err),
		)
		return
	}
	token, err := ctx.DB.GetToken(r.Context(), tkID)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("collections.list checking token",
			zap.Error(err),
		)
		return
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
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		log.Error("collections.list checking access permissions",
			zap.Error(err),
		)
		return
	}
	opts, err := util.ListOptions(r)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadToken}, http.StatusBadRequest)
		log.Error("collection.list checking query params",
			zap.Error(err),
		)
		return
	}
	o, err := ctx.DB.ListCollections(r.Context(), token.Subject, opts)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError},
			http.StatusInternalServerError)
		log.Error("collection.list fail to retrieve collcetions list",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w, models.CollectionList{
		Collections: o,
	}, http.StatusOK)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	if ctx.Claims == nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		return
	}
	params := r.URL.Query()
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
	err = ctx.DB.DeleteCollection(r.Context(), id)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError},
			http.StatusInternalServerError)
		log.Error("collection.view fail to retrieve stored qrcode info",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w, &models.Status{Status: keys.Success}, http.StatusOK)
}

// Assign assigns a collection to the qrcode.
func Assign(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	if ctx.Claims == nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadBody}, http.StatusBadRequest)
		log.Error("collections.create can't read body",
			zap.Error(err),
		)
		return
	}
	m := &models.CollectionAssignReq{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadJSON}, http.StatusBadRequest)
		log.Error("collections.create fail to unmarshal request body",
			zap.Error(err),
		)
		return
	}
	err = ctx.DB.AssignCollection(r.Context(), m.ID, m.QRID)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusInternalServerError)
		log.Error("collection.assign saving new collection",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w, &models.Status{Status: keys.Success}, http.StatusOK)
}

func DeAssign(w http.ResponseWriter, r *http.Request) {
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	if ctx.Claims == nil {
		util.WriteJSON(w, &models.APIError{Message: keys.Forbidden}, http.StatusForbidden)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadBody}, http.StatusBadRequest)
		log.Error("collections.create can't read body",
			zap.Error(err),
		)
		return
	}
	m := &models.CollectionAssignReq{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.BadJSON}, http.StatusBadRequest)
		log.Error("collections.create fail to unmarshal request body",
			zap.Error(err),
		)
		return
	}
	err = ctx.DB.DeAssignCollection(r.Context(), m.ID, m.QRID)
	if err != nil {
		util.WriteJSON(w, &models.APIError{Message: keys.InternalError}, http.StatusInternalServerError)
		log.Error("collection.assign saving new collection",
			zap.Error(err),
		)
		return
	}
	util.WriteJSON(w, &models.Status{Status: keys.Success}, http.StatusOK)
}
