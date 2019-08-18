package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/gernest/hiro/access"
	"github.com/gernest/hiro/accounts"
	"github.com/gernest/hiro/assets"
	"github.com/gernest/hiro/codes/qrcode"
	"github.com/gernest/hiro/collections"
	"github.com/gernest/hiro/config"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/prom"
	"github.com/gernest/hiro/query"
	"github.com/gernest/hiro/scan"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	"go.uber.org/zap"
)

// ServeAPI serves the qrcode generation API.host is used to reference the
// url to the current server a.k.a host.
func ServeAPI(ctx context.Context, db *query.SQL, cfg *config.Config) error {
	h := Handler(ctx, db, cfg)
	port := fmt.Sprintf(":%d", cfg.Port)
	return http.ListenAndServe(port, h)
}

// Handler  returns *WrapHandler with all engopints registered.
func Handler(ctx context.Context, db *query.SQL, cfg *config.Config) http.Handler {
	mux := echo.New()
	mux.Use(middleware.LoggerWithConfig(Logger()))
	mux.Use(middleware.Recover())
	l, _ := zap.NewDevelopment()
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
	mux.GET("/", Home)
	mux.GET("/user/*", Home)
	mux.GET("/metrics", Metrics)
	codes := mux.Group("/codes")

	//qrcode api
	qr := codes.Group("/qrcode")
	qr.POST("/qr", qrcode.Create)
	qr.GET("/qr", qrcode.List)
	qr.GET("/qr/:uuid", qrcode.View)
	qr.POST("/qr/:uuid", qrcode.Update)
	qr.DELETE("/qr/:uuid", qrcode.Delete)

	// collections
	co := mux.Group("/collections")
	co.POST("/collections", collections.Create)
	co.GET("/collections", collections.List)
	co.GET("/collections/view", collections.View)
	co.DELETE("/collections/delete", collections.Delete)
	co.POST("/collections/assign", collections.Assign)
	co.POST("/collections/deassign", collections.DeAssign)

	//accounts
	mux.POST("/api/register", accounts.Create)
	mux.POST("/api/login/account", accounts.Login)

	// scan
	mux.GET("/scan/:uuid", scan.Scan)
	//static assets
	s := gziphandler.GzipHandler(Static())
	mux.GET("/static/*", func(ctx echo.Context) error {
		s.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	})
	return prom.Wrap(mux)
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
func Home(ctx echo.Context) error {
	b, err := assets.Asset("static/index.html")
	if err != nil {
		// w.Write([]byte(err.Error()))
		return err
	}
	return ctx.HTML(http.StatusOK, string(b))
}

func Metrics(ctx echo.Context) error {
	promhttp.Handler().ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}
