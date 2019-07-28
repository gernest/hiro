package qrcode

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gernest/alien"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/testutil"
	uuid "github.com/satori/go.uuid"
)

func TestQR2(t *testing.T) {
	ctx, err := testutil.New("qrcode", "someSecret")
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()
	req := ctx.Req
	t.Run("create", func(ts *testing.T) {
		ts.Run("bad body", func(ts *testing.T) {
			r := req("POST", "/v1/qr/", nil)
			w := httptest.NewRecorder()
			Create(w, r)
			if w.Code != http.StatusBadRequest {
				ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
			}
			a := apiError(ts, w)
			if a.Message != keys.BadJSON {
				t.Errorf("expected %s got %s", keys.BadJSON, a.Message)
			}
		})
		ts.Run("bad token string", func(ts *testing.T) {
			r := req("POST", "/v1/qr/", testutil.ReqData(&models.QRReq{}))
			orig := ctx.Claims.Id
			ctx.Claims.Id = ""
			w := httptest.NewRecorder()
			Create(w, r)
			if w.Code != http.StatusBadRequest {
				ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
			}
			ctx.Claims.Id = orig
			a := apiError(ts, w)
			if a.Message != keys.BadToken {
				t.Errorf("expected %s got %s", keys.BadToken, a.Message)
			}
		})
		ts.Run("unknown token", func(ts *testing.T) {
			r := req("POST", "/v1/qr/", testutil.ReqData(&models.QRReq{}))

			orig := ctx.Claims.Id
			ctx.Claims.Id = uuid.NewV4().String()
			w := httptest.NewRecorder()
			Create(w, r)
			if w.Code != http.StatusBadRequest {
				ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
			}
			ctx.Claims.Id = orig
			a := apiError(ts, w)
			if a.Message != keys.BadToken {
				t.Errorf("expected %s got %s", keys.BadToken, a.Message)
			}
		})

		r := req("POST", "/v1/qr/", testutil.ReqData(&models.QRReq{}))

		w := httptest.NewRecorder()
		Create(w, r)
		if w.Code != http.StatusOK {
			ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
		}
		e := &models.QR{}
		err := json.Unmarshal(w.Body.Bytes(), e)
		if err != nil {
			ts.Fatal(err)
		}
	})
	t.Run("create with context", func(ts *testing.T) {
		r := req("POST", "/v1/qr/", testutil.ReqData(&models.QRReq{
			ShouldRedirect: true,
			RedirectURL:    keys.RootURL,
		}))

		w := httptest.NewRecorder()
		Create(w, r)
		if w.Code != http.StatusOK {
			ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
		}
		e := &models.QR{}
		err := json.Unmarshal(w.Body.Bytes(), e)
		if err != nil {
			ts.Fatal(err)
		}
		ts.Run("view", func(ts *testing.T) {
			// view created info
			link := "/v1/qr/" + e.UUID.String()
			r = req("GET", link, nil)

			param, err := alien.ParseParams(link, "/v1/qr/:uuid")
			if err != nil {
				ts.Fatal(err)
			}
			r.Header.Add("_alien", param)
			w = httptest.NewRecorder()
			View(w, r)
			if w.Code != http.StatusOK {
				ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
			}
			ve := &models.QR{}
			err = json.Unmarshal(w.Body.Bytes(), ve)
			if err != nil {
				ts.Fatal(err)
			}
			if ve.UUID != e.UUID {
				ts.Errorf("expected %s got %s", e.UUID, ve.UUID)
			}

			//view with invalid uuid
			link = "/v1/qr/bad"
			r = req("GET", link, nil)

			param, err = alien.ParseParams(link, "/v1/qr/:uuid")
			if err != nil {
				ts.Fatal(err)
			}
			r.Header.Add("_alien", param)
			w = httptest.NewRecorder()
			View(w, r)
			if w.Code != http.StatusBadRequest {
				ts.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
			}

			// view with non existing uuid
			link = "/v1/qr/" + uuid.NewV4().String()
			r = req("GET", link, nil)

			param, err = alien.ParseParams(link, "/v1/qr/:uuid")
			if err != nil {
				ts.Fatal(err)
			}
			r.Header.Add("_alien", param)
			w = httptest.NewRecorder()
			View(w, r)
			if w.Code != http.StatusNotFound {
				ts.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
			}

			ts.Run("update", func(ts *testing.T) {

				// To avoid extra code we are reusing the qrcode generated in the parent
				// context for updates.
				//
				//TODO: make this self contained test case.
				link = "/v1/qr/" + ve.UUID.String()
				r = req("POST", link, nil)

				param, err = alien.ParseParams(link, "/v1/qr/:uuid")
				if err != nil {
					ts.Fatal(err)
				}
				r.Header.Add("_alien", param)
				w = httptest.NewRecorder()
				Update(w, r)
				if w.Code != http.StatusBadRequest {
					ts.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
				}

				link = "/v1/qr/bad"
				r = req("POST", link, nil)

				param, err = alien.ParseParams(link, "/v1/qr/:uuid")
				if err != nil {
					ts.Fatal(err)
				}
				r.Header.Add("_alien", param)
				w = httptest.NewRecorder()
				Update(w, r)
				if w.Code != http.StatusBadRequest {
					ts.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
				}

				link = "/v1/qr/" + uuid.NewV4().String()
				r = req("POST", link, testutil.ReqData(ve))

				param, err = alien.ParseParams(link, "/v1/qr/:uuid")
				if err != nil {
					ts.Fatal(err)
				}
				r.Header.Add("_alien", param)
				w = httptest.NewRecorder()
				Update(w, r)
				if w.Code != http.StatusNotFound {
					ts.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
				}

				ve.ShouldRedirect = true
				link = "/v1/qr/" + ve.UUID.String()
				r = req("POST", link, testutil.ReqData(ve))

				param, err = alien.ParseParams(link, "/v1/qr/:uuid")
				if err != nil {
					ts.Fatal(err)
				}
				r.Header.Add("_alien", param)
				w = httptest.NewRecorder()
				Update(w, r)
				if w.Code != http.StatusOK {
					ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
				}
				ue := &models.QR{}
				err = json.Unmarshal(w.Body.Bytes(), ue)
				if err != nil {
					ts.Fatal(err)
				}
				if !e.ShouldRedirect {
					ts.Error("expected to be true")
				}
			})
		})
	})
	t.Run("delete", func(ts *testing.T) {
		r := req("POST", "/v1/qr/", testutil.ReqData(&models.QRReq{}))

		w := httptest.NewRecorder()
		Create(w, r)
		if w.Code != http.StatusOK {
			ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
		}
		e := &models.QR{}
		err := json.Unmarshal(w.Body.Bytes(), e)
		if err != nil {
			ts.Fatal(err)
		}

		link := "/v1/qr/bad"
		r = req("DELETE", link, nil)

		param, err := alien.ParseParams(link, "/v1/qr/:uuid")
		if err != nil {
			ts.Fatal(err)
		}
		r.Header.Add("_alien", param)
		w = httptest.NewRecorder()
		Delete(w, r)
		if w.Code != http.StatusBadRequest {
			ts.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
		}
		link = "/v1/qr/" + e.UUID.String()
		r = req("DELETE", link, nil)

		param, err = alien.ParseParams(link, "/v1/qr/:uuid")
		if err != nil {
			ts.Fatal(err)
		}
		r.Header.Add("_alien", param)
		w = httptest.NewRecorder()
		Delete(w, r)
		if w.Code != http.StatusOK {
			ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("list", func(ts *testing.T) {
		link := "/v1/qr/"
		r := req("GET", link, nil)

		w := httptest.NewRecorder()
		List(w, r)
		if w.Code != http.StatusOK {
			ts.Fatalf("expected %d got %d", http.StatusOK, w.Code)
		}

		o := models.QRList{}
		err = json.Unmarshal(w.Body.Bytes(), &o)
		if err != nil {
			ts.Fatal(err)
		}
		if len(o.QRCodes) == 0 {
			t.Error("expected more qrcode")
		}
		if o.Total < 1 {
			t.Error("expected total to be set")
		}
	})
}

func apiError(t *testing.T, w *httptest.ResponseRecorder) *models.APIError {
	m := &models.APIError{}
	err := json.Unmarshal(w.Body.Bytes(), m)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
