package server

import (
	"context"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gernest/alien"
	"github.com/gernest/hiro/access"
	"github.com/gernest/hiro/accounts"
	"github.com/gernest/hiro/assets"
	"github.com/gernest/hiro/collections"
	"github.com/gernest/hiro/config"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/prom"
	"github.com/gernest/hiro/qrcode"
	"github.com/gernest/hiro/query"
	"github.com/gernest/hiro/scan"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	"go.uber.org/zap"
)

// ServeAPI serves the qrcode generation API.host is used to reference the
// url to the current server a.k.a host.
func ServeAPI(ctx context.Context, db *query.SQL, cfg *config.Config) error {
	h := Handler(ctx, db, cfg)
	return http.ListenAndServe(":8000", h)
}

// Handler  returns *WrapHandler with all engopints registered.
func Handler(ctx context.Context, db *query.SQL, cfg *config.Config) http.Handler {
	mux := alien.New()
	l, _ := zap.NewProduction()
	jwt := &models.JWT{Secret: []byte(cfg.Secret)}
	wares := []models.Item{
		models.Item{Key: keys.DB, Value: db},
		models.Item{Key: keys.LoggerKey, Value: l},
		models.Item{Key: keys.JwtKey, Value: jwt},
	}
	warden, err := access.New(db.DB(), nil)
	if err != nil {
		l.Fatal("initializing warden",
			zap.Error(err),
		)
	}

	wares = append(wares, models.Item{Key: keys.Warden, Value: warden})
	mux.Use(CtxMiddleware(wares...))
	mux.Use(User(l, jwt))
	mux.Get("/", Home)
	mux.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
	api := mux.Group("/v1")

	//qrcode api
	api.Post("/qr", qrcode.Create)
	api.Get("/qr", qrcode.List)
	api.Get("/qr/:uuid", qrcode.View)
	api.Post("/qr/:uuid", qrcode.Update)
	api.Delete("/qr/:uuid", qrcode.Delete)

	// collections
	api.Post("/collections", collections.Create)
	api.Get("/collections", collections.List)
	api.Get("/collections/view", collections.View)
	api.Delete("/collections/delete", collections.Delete)
	api.Post("/collections/assign", collections.Assign)
	api.Post("/collections/deassign", collections.DeAssign)

	//accounts
	mux.Post("/register", accounts.Create)
	mux.Post("/login", accounts.Login)

	// scan
	mux.Get("/scan/:uuid", scan.Scan)

	mux.Get("/privacy", Privacy)

	//static assets
	s := gziphandler.GzipHandler(Static())

	mux.AddRoute(http.MethodGet, "/static/*", func(w http.ResponseWriter, r *http.Request) {
		s.ServeHTTP(w, r)
	})
	mux.NotFoundHandler(http.HandlerFunc(NotFound))
	return prom.Wrap(mux)
}

// CtxMiddleware injects the values into the request context.
func CtxMiddleware(items ...models.Item) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			for _, v := range items {
				ctx = context.WithValue(ctx, v.Key, v.Value)
			}
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Static is a handlerfor serving static assets.
func Static() http.Handler {
	return http.FileServer(
		&assetfs.AssetFS{
			Asset:     assets.Asset,
			AssetDir:  assets.AssetDir,
			AssetInfo: assets.AssetInfo,
		},
	)
}

// Home renders the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	b, err := assets.Asset("static/index.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func Privacy(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
