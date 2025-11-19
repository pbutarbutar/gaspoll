package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"gaspoll/internal/auth"
	"gaspoll/internal/mitra"
	"gaspoll/internal/order"
	"gaspoll/internal/payment"
	"gaspoll/internal/product"
	"gaspoll/internal/reward"
	"gaspoll/internal/settlement"
	"gaspoll/internal/user"
	"gaspoll/internal/view"
	"gaspoll/internal/voucher"
)

type services struct {
	users       *user.Service
	products    *product.Service
	orders      *order.Service
	vouchers    *voucher.Service
	rewards     *reward.Service
	payments    *payment.Service
	mitras      *mitra.Service
	settlements *settlement.Service
	session     *auth.SessionManager
}

func registerRoutes(e *echo.Echo, _ *view.Engine, svc services) {
	e.Static("/assets", "web/assets")

	e.GET("/", func(c echo.Context) error {
		prods, _ := svc.products.ListActive(c.Request().Context())
		return c.Render(http.StatusOK, "home.html", map[string]any{
			"Title":    "Gaspoll",
			"Products": prods,
		})
	})

	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "auth_login.html", map[string]any{
			"Title": "Masuk Driver",
		})
	})

	e.POST("/login", func(c echo.Context) error {
		var in user.LoginInput
		if err := c.Bind(&in); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid payload")
		}
		usr, err := svc.users.Authenticate(c.Request().Context(), in)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "gagal masuk: "+err.Error())
		}
		if err := svc.session.SetUser(c, usr.ID.String()); err != nil {
			return err
		}
		return c.Redirect(http.StatusFound, "/dashboard")
	})

	e.GET("/register", func(c echo.Context) error {
		return c.Render(http.StatusOK, "auth_register.html", map[string]any{
			"Title": "Daftar Driver",
		})
	})

	e.POST("/register", func(c echo.Context) error {
		var in user.RegisterInput
		if err := c.Bind(&in); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid payload")
		}
		if in.Role == "" {
			in.Role = user.RoleDriver
		}
		usr, err := svc.users.Register(c.Request().Context(), in)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		_ = svc.session.SetUser(c, usr.ID.String())
		return c.Redirect(http.StatusFound, "/dashboard")
	})

	e.POST("/logout", func(c echo.Context) error {
		_ = svc.session.Clear(c)
		return c.Redirect(http.StatusFound, "/")
	})

	e.GET("/dashboard", func(c echo.Context) error {
		userID, err := currentUserID(c, svc.session)
		if err != nil {
			return c.Redirect(http.StatusFound, "/login")
		}

		data := map[string]any{
			"Title": "Dashboard Driver",
		}
		orders, _ := svc.orders.ListByUser(c.Request().Context(), userID)
		data["Orders"] = orders

		vouchers, _ := svc.vouchers.ListByUser(c.Request().Context(), userID)
		data["Vouchers"] = vouchers

		progress, _ := svc.rewards.Progress(c.Request().Context(), userID)
		data["Progress"] = progress

		return c.Render(http.StatusOK, "dashboard.html", data)
	})

	e.GET("/partials/products", func(c echo.Context) error {
		prods, _ := svc.products.ListActive(c.Request().Context())
		return c.Render(http.StatusOK, "components/products.html", map[string]any{
			"Products": prods,
		})
	})

	e.POST("/api/orders", func(c echo.Context) error {
		userID, err := currentUserID(c, svc.session)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "login terlebih dahulu")
		}

		var in order.CreateInput
		if err := c.Bind(&in); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "payload invalid")
		}

		ord, invoice, err := svc.orders.Create(c.Request().Context(), userID, in)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"order":   ord,
			"invoice": invoice,
		})
	})

	e.POST("/api/payments/webhook", func(c echo.Context) error {
		// expect reference id in payload
		ref := c.QueryParam("ref")
		if ref == "" {
			return c.NoContent(http.StatusBadRequest)
		}
		orderID, status, err := svc.payments.VerifyWebhook(c.Request().Context(), ref)
		if err != nil || status != payment.StatusPaid {
			return c.NoContent(http.StatusBadRequest)
		}
		if err := svc.orders.MarkPaid(c.Request().Context(), orderID); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	})

	e.POST("/api/vouchers/redeem", func(c echo.Context) error {
		code := c.FormValue("code")
		if code == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "code required")
		}
		partner := c.FormValue("partner")
		if partner == "" {
			partner = "Bengkel Mitra"
		}
		if err := svc.vouchers.Redeem(c.Request().Context(), code, partner); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "redeemed"})
	})

	e.GET("/mitra", func(c echo.Context) error {
		return c.Render(http.StatusOK, "mitra.html", map[string]any{
			"Title": "Dashboard Mitra",
		})
	})

	e.GET("/admin/settlement", func(c echo.Context) error {
		rows, _ := svc.settlements.List(c.Request().Context())
		return c.Render(http.StatusOK, "admin_settlement.html", map[string]any{
			"Title":       "Settlement",
			"Settlements": rows,
		})
	})
}

func currentUserID(c echo.Context, session *auth.SessionManager) (uuid.UUID, error) {
	id, err := session.CurrentUserID(c)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}
