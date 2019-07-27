package collections

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/testutil"
	uuid "github.com/satori/go.uuid"
)

func TestCollections(t *testing.T) {
	ctx, err := testutil.New("collections", "http://localhost:8000", "someSecret")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("create", func(ts *testing.T) {
		r := ctx.Req("POST", "/v1/collections", testutil.ReqData(
			&models.CollectionReq{
				Name:  "test",
				Color: "blue",
			},
		))
		w := httptest.NewRecorder()
		Create(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		c := &models.Collection{}
		err = json.Unmarshal(w.Body.Bytes(), c)
		if err != nil {
			t.Fatal(err)
		}

		ts.Run("view", func(ts *testing.T) {
			query := make(url.Values)
			query.Set("uuid", c.ID())
			r := ctx.Req("GET", "/v1/collections/view?"+query.Encode(), nil)
			w := httptest.NewRecorder()
			View(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("expected %d got %d", http.StatusOK, w.Code)
			}
			view := &models.Collection{}
			err := json.Unmarshal(w.Body.Bytes(), view)
			if err != nil {
				t.Fatal(err)
			}
		})
		ts.Run("list", func(ts *testing.T) {
			r := ctx.Req("GET", "/v1/collections", nil)
			w := httptest.NewRecorder()
			List(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("expected %d got %d", http.StatusOK, w.Code)
			}
			list := &models.CollectionList{}
			err := json.Unmarshal(w.Body.Bytes(), list)
			if err != nil {
				t.Fatal(err)
			}
		})
		ts.Run("asssign", func(ts *testing.T) {
			code := &models.QR{
				UUID:   uuid.NewV4(),
				UserID: ctx.Account.UUID,
			}
			err := ctx.DB.CreateQR(context.Background(), code)
			if err != nil {
				ts.Fatal(err)
			}
			r := ctx.Req("POST", "/v1/collections", testutil.ReqData(
				&models.CollectionAssignReq{
					ID:   c.UUID,
					QRID: code.UUID,
				},
			))
			w := httptest.NewRecorder()
			Assign(w, r)
			if w.Code != http.StatusOK {
				ts.Errorf("expected %d got %d", http.StatusOK, w.Code)
			}
			assignStatus := &models.Status{}
			err = json.Unmarshal(w.Body.Bytes(), assignStatus)
			if err != nil {
				ts.Fatal(err)
			}
		})
		ts.Run("delete", func(ts *testing.T) {
			query := make(url.Values)
			query.Set("uuid", c.ID())
			r := ctx.Req("POST", "/v1/collections/delete?"+query.Encode(), nil)
			w := httptest.NewRecorder()
			Delete(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("expected %d got %d", http.StatusOK, w.Code)
			}
			deletedStatus := &models.Status{}
			err := json.Unmarshal(w.Body.Bytes(), deletedStatus)
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}
