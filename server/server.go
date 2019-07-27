package server

import (
	"context"
	"html/template"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/minio/minio-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gernest/alien"
	"github.com/gernest/hiro/access"
	"github.com/gernest/hiro/accounts"
	"github.com/gernest/hiro/assets"
	"github.com/gernest/hiro/bus"
	"github.com/gernest/hiro/collections"
	"github.com/gernest/hiro/config"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/meta"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/prom"
	"github.com/gernest/hiro/qrcode"
	"github.com/gernest/hiro/query"
	"github.com/gernest/hiro/scan"
	"github.com/gernest/hiro/templates"
	"github.com/gernest/hiro/util"

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
	var nio *minio.Client
	var err error
	wares := []models.Item{
		models.Item{Key: keys.DB, Value: db},
		models.Item{Key: keys.LoggerKey, Value: l},
		models.Item{Key: keys.JwtKey, Value: jwt},
		models.Item{Key: keys.Host, Value: cfg.Host},
		models.Item{Key: keys.ImageHost, Value: cfg.ImageHost},
	}
	if cfg.Minio != nil {
		l.Info("Detected minio, we will use minio as qrcode stowage")
		m := cfg.Minio
		nio, err = minio.New(m.Endpoint, m.AccessKey, m.AccessSecret, false)
		if err != nil {
			l.Fatal("initializing minio client",
				zap.Error(err),
			)
		}
		wares = append(wares, models.Item{Key: keys.Minio, Value: nio})
		l.Info("connected to minio server", zap.String(
			"endpoint", m.Endpoint,
		))
	}
	warden, err := access.New(db.DB(), nil)
	if err != nil {
		l.Fatal("initializing warden",
			zap.Error(err),
		)
	}

	if cfg.NSQ != nil {
		l.Info("Detected nsq, we will use nsq to broadcast messages")
		b, err := bus.NewProducer(cfg.NSQ.NSQD)
		if err != nil {
			if err != nil {
				l.Fatal("initializing nsq client",
					zap.Error(err),
				)
			}
		}
		b.Logger = l
		wares = append(wares, models.Item{Key: keys.NSQ, Value: b})
		l.Info("connected to nsqlookupd", zap.String(
			"endpoint", cfg.NSQ.LookupD,
		))
		warden.AuditLogger = access.NewAuditLogger(b)
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

	opts := []struct {
		methods []string
		route   string
	}{
		{
			methods: []string{"GET", "HEAD"},
			route:   "/static/css/*",
		},
		{
			methods: []string{"GET", "HEAD"},
			route:   "/static/js/*",
		},
		{
			methods: []string{"GET", "HEAD"},
			route:   "/static/img/*",
		},
	}
	for _, v := range opts {
		for _, method := range v.methods {
			mux.AddRoute(method, v.route, func(w http.ResponseWriter, r *http.Request) {
				s.ServeHTTP(w, r)
			})
		}
	}
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
	ctx := util.RequestContext(r.Context())
	log := ctx.Logger.With(
		zap.String("url", r.URL.String()),
	)
	m := &meta.Meta{
		OpenGraph: &meta.OpenGraph{
			URL:         keys.WebsiteURL,
			Title:       keys.WebsiteTitle,
			Description: keys.WebsiteDescription,
		},
		Twitter: &meta.Twitter{
			Card:        "summary",
			Site:        "@bqservice",
			Creator:     "@gernesti",
			Title:       keys.WebsiteTitle,
			Description: keys.WebsiteDescription,
		},
	}
	mg, err := m.Map()
	if err != nil {
		log.Error("home.getting meta properties",
			zap.Error(err),
		)
	}
	err = templates.Write(w, "html/home.html", map[string]interface{}{
		"meta":  mg,
		"Title": "BQ: high performacne qrcode service",
	})
	if err != nil {
		log.Error("rendering home page template",
			zap.Error(err),
		)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	m := &meta.Meta{
		OpenGraph: &meta.OpenGraph{
			URL:         keys.WebsiteURL,
			Title:       keys.WebsiteTitle,
			Description: keys.WebsiteDescription,
		},
		Twitter: &meta.Twitter{
			Card:        "summary",
			Site:        "@bqservice",
			Creator:     "@gernesti",
			Title:       keys.WebsiteTitle,
			Description: keys.WebsiteDescription,
		},
	}
	mg, err := m.Map()
	if err != nil {
		//TODO handle error
	}
	err = templates.Write(w, "html/home.html", map[string]interface{}{
		"meta":        mg,
		"Title":       "BQ: high performacne qrcode service",
		"ActiveRoute": template.JS(r.URL.Path),
	})
	if err != nil {
		//TODO handle error
	}
}

func Privacy(w http.ResponseWriter, r *http.Request) {
	m := &meta.Meta{
		OpenGraph: &meta.OpenGraph{
			URL:         keys.WebsiteURL,
			Title:       keys.WebsiteTitle,
			Description: keys.WebsiteDescription,
		},
		Twitter: &meta.Twitter{
			Card:        "summary",
			Site:        "@bqservice",
			Creator:     "@gernesti",
			Title:       keys.WebsiteTitle,
			Description: keys.WebsiteDescription,
		},
	}
	mg, err := m.Map()
	if err != nil {
		//TODO handle error
	}
	err = templates.Write(w, "html/privacy.html", map[string]interface{}{
		"meta":        mg,
		"Title":       "BQ: high performacne qrcode service",
		"ActiveRoute": template.JS(r.URL.Path),
	})
	if err != nil {
		//TODO handle error
	}
}
