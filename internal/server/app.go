package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gaspoll/internal/config"
	"gaspoll/internal/constant"
	"gaspoll/internal/entity"
	"gaspoll/internal/repository"
	"gaspoll/internal/usecase/auth"
	"gaspoll/internal/usecase/mitra"
	"gaspoll/internal/usecase/order"
	"gaspoll/internal/usecase/product"
	"gaspoll/internal/usecase/reward"
	"gaspoll/internal/usecase/settlement"
	"gaspoll/internal/usecase/voucher"
	"gaspoll/internal/view"
)

type App struct {
	cfg        config.Config
	repo       *repository.Repository
	echo       *echo.Echo
	session    *sessions.CookieStore
	auth       auth.Service
	product    product.Service
	order      order.Service
	voucher    voucher.Service
	mitra      mitra.Service
	settlement settlement.Service
	reward     reward.Service
}

func New(cfg config.Config, repo *repository.Repository) *App {
	e := echo.New()
	// Add middleware to set Cache-Control for static assets
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/assets/") {
				c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			return next(c)
		}
	})
	e.HideBanner = true

	// serve static assets from web/assets at /assets/
	e.Static("/assets", "web/assets")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	app := &App{
		cfg:     cfg,
		repo:    repo,
		echo:    e,
		session: sessions.NewCookieStore([]byte(cfg.SessionKey)),
		reward:  reward.Service{},
	}
	app.auth = auth.Service{Users: repo.User}
	app.product = product.Service{Products: repo.Product}
	app.voucher = voucher.Service{Vouchers: repo.Voucher}
	app.mitra = mitra.Service{Mitras: repo.Mitra}
	app.settlement = settlement.Service{Settlements: repo.Settlement}
	app.order = order.Service{
		Users:    repo.User,
		Products: repo.Product,
		Orders:   repo.Order,
		Vouchers: repo.Voucher,
		Mitras:   repo.Mitra,
		Rewards:  app.reward,
	}

	e.Renderer = view.NewRenderer()
	e.Use(app.injectUser())
	app.routes()

	return app
}

func (a *App) routes() {
	a.echo.GET("/", a.handleHome)
	a.echo.GET("/login", a.showLogin)
	a.echo.POST("/login", a.handleLogin)
	a.echo.POST("/register", a.handleRegister)

	a.echo.GET("/products", a.listProducts)
	a.echo.POST("/orders", a.withAuth(a.handleCheckout))

	a.echo.GET("/dashboard", a.withAuth(a.handleDashboard))
	a.echo.GET("/voucher/:id/qr", a.withAuth(a.renderVoucherQR))
	a.echo.GET("/mitra/:id", a.showMitra)
	a.echo.GET("/mitra/:id/redeem", a.showMitraRedeem)

	a.echo.GET("/ipl", a.showIPL)
	a.echo.POST("/ipl/pay", a.handleIPLPAY)

	a.echo.GET("/mitra/redeem", a.showRedeem)
	a.echo.POST("/mitra/redeem", a.handleRedeem)

	a.echo.GET("/admin/settlement", a.showSettlement)
}

func (a *App) Run() error {
	return a.echo.Start(":" + a.cfg.Port)
}

func (a *App) injectUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if user := a.currentUser(c); user != nil {
				c.Set(constant.ContextUserKey, user)
			}
			return next(c)
		}
	}
}

func (a *App) currentUser(c echo.Context) *entity.User {
	sess, err := a.session.Get(c.Request(), constant.SessionName)
	if err != nil {
		return nil
	}
	val, ok := sess.Values["uid"].(string)
	if !ok || val == "" {
		return nil
	}
	user, err := a.repo.User.FindUserByID(val)
	if err != nil {
		return nil
	}
	return user
}

func (a *App) withAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := a.currentUser(c)
		if user == nil {
			return c.Redirect(http.StatusFound, "/login")
		}
		c.Set(constant.ContextUserKey, user)
		return next(c)
	}
}

func (a *App) render(c echo.Context, name string, data map[string]interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}
	if _, ok := data["user"]; !ok {
		if u := a.currentUser(c); u != nil {
			data["user"] = u
		}
	}
	// provide the content template name derived from the page name
	// e.g. for "home.html" -> "content_home"
	base := strings.TrimSuffix(name, ".html")
	data["contentName"] = "content_" + base

	if _, ok := data["title"]; !ok {
		data["title"] = "Gaspoll — Driver Commerce & Rewards"
	}
	if _, ok := data["metaDescription"]; !ok {
		data["metaDescription"] = "Gaspoll memudahkan driver topup saldo, beli paket data, dan mendapatkan reward bengkel mitra dengan cepat dan aman."
	}
	if _, ok := data["metaImage"]; !ok {
		data["metaImage"] = "/assets/logo-wgaspoll.png?v=1"
	}
	if _, ok := data["metaURL"]; !ok {
		scheme := c.Scheme()
		if scheme == "" {
			scheme = "https"
		}
		host := c.Request().Host
		path := c.Request().URL.RequestURI()
		if host != "" {
			data["metaURL"] = scheme + "://" + host + path
		} else {
			data["metaURL"] = path
		}
	}
	return c.Render(http.StatusOK, name, data)
}
