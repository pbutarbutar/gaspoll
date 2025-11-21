package server

import (
	"bytes"
	"crypto/sha1"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"gaspoll/internal/constant"
	"gaspoll/internal/entity"
)

type IPLBill struct {
	House  string
	Period string
	Amount int
}

func (a *App) showLogin(c echo.Context) error {
	return a.render(c, "login.html", map[string]interface{}{
		"title": "Masuk / Daftar",
	})
}

func (a *App) handleLogin(c echo.Context) error {
	phone := c.FormValue("phone")
	name := c.FormValue("name")
	vehicle := c.FormValue("vehicle")

	user, err := a.auth.Login(phone)
	if err != nil {
		user, err = a.auth.Register(name, phone, vehicle)
		if err != nil {
			return c.String(http.StatusBadRequest, "Gagal login/daftar: "+err.Error())
		}
	}

	sess, _ := a.session.Get(c.Request(), constant.SessionName)
	sess.Values["uid"] = user.ID
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/dashboard")
}

func (a *App) handleRegister(c echo.Context) error {
	name := c.FormValue("name")
	phone := c.FormValue("phone")
	vehicle := c.FormValue("vehicle")

	user, err := a.auth.Register(name, phone, vehicle)
	if err != nil {
		return c.String(http.StatusBadRequest, "Gagal daftar: "+err.Error())
	}
	sess, _ := a.session.Get(c.Request(), constant.SessionName)
	sess.Values["uid"] = user.ID
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, "/dashboard")
}

func (a *App) handleHome(c echo.Context) error {
	mitras, _ := a.mitra.List()
	products, _ := a.product.List()

	if len(products) == 0 {
		products = []*entity.Product{
			{ID: "p-gj", Name: "Topup Gojek 50k", Category: "Topup Driver", Provider: "Gojek", Price: 52000, Description: "Dummy topup Gojek"},
			{ID: "p-gr", Name: "Topup Grab 50k", Category: "Topup Driver", Provider: "Grab", Price: 52500, Description: "Dummy topup Grab"},
			{ID: "p-mx", Name: "Topup Maxim 50k", Category: "Topup Driver", Provider: "Maxim", Price: 53000, Description: "Dummy topup Maxim"},
			{ID: "p-id", Name: "Topup InDrive 50k", Category: "Topup Driver", Provider: "InDrive", Price: 53500, Description: "Dummy topup InDrive"},
			{ID: "p-pls", Name: "Pulsa 25k", Category: "Pulsa", Provider: "Telkomsel", Price: 26000, Description: "Dummy pulsa"},
			{ID: "p-data", Name: "Paket Data 20GB", Category: "Paket Internet", Provider: "Indosat", Price: 75000, Description: "Dummy data"},
		}
	}
	if len(mitras) == 0 {
		mitras = []*entity.Mitra{
			// Depok (6)
			{ID: "m-1", Name: "Oto Bengkel", Address: "Jl. Melati No. 12", Location: "Depok", Promo: "Diskon oli 20%", Services: []string{"Ganti oli", "Tune up"}, Distance: "1.2 km"},
			{ID: "m-2", Name: "Garasi Juara", Address: "Jl. Pahlawan No. 45", Location: "Depok", Promo: "Voucher servis ringan", Services: []string{"Spooring", "Balancing"}, Distance: "3 km"},
			{ID: "m-3", Name: "Sentosa Motor", Address: "Jl. Sudirman No. 10", Location: "Depok", Promo: "Diskon filter udara", Services: []string{"Filter udara", "Tune up"}, Distance: "4.5 km"},
			{ID: "m-4", Name: "Lenteng Motor", Address: "Jl. Lenteng Agung Raya", Location: "Depok", Promo: "Cek rem gratis", Services: []string{"Cek rem", "Ganti kampas"}, Distance: "2.9 km"},
			{ID: "m-5", Name: "Margonda Service", Address: "Jl. Margonda Raya No. 88", Location: "Depok", Promo: "Diskon servis 15%", Services: []string{"Servis ringan", "Ganti oli"}, Distance: "3.4 km"},
			{ID: "m-6", Name: "Sawangan Garage", Address: "Jl. Raya Sawangan", Location: "Depok", Promo: "Voucher 25k", Services: []string{"Tune up", "Ganti ban"}, Distance: "5 km"},

			// Cikarang (6)
			{ID: "m-7", Name: "Jaya Auto", Address: "Jl. Panjang No. 88", Location: "Cikarang", Promo: "Voucher 25k", Services: []string{"Spooring", "Balancing"}, Distance: "6 km"},
			{ID: "m-8", Name: "Green Garage", Address: "Jl. Cemara No. 2", Location: "Cikarang", Promo: "Paket oli hemat", Services: []string{"Ganti oli", "Cuci motor"}, Distance: "2.8 km"},
			{ID: "m-9", Name: "Delta Service", Address: "Jl. Industri No. 5", Location: "Cikarang", Promo: "Diskon tune up 10%", Services: []string{"Tune up", "Ganti busi"}, Distance: "4.2 km"},
			{ID: "m-10", Name: "Mega Auto", Address: "Jl. Boulevard Cikarang", Location: "Cikarang", Promo: "Cek motor gratis", Services: []string{"Cek motor", "Ganti oli"}, Distance: "3.7 km"},
			{ID: "m-11", Name: "C-Tech Workshop", Address: "Jl. Kenari No. 3", Location: "Cikarang", Promo: "Diskon kampas rem 12%", Services: []string{"Rem", "Suspensi"}, Distance: "5.3 km"},
			{ID: "m-12", Name: "Benua Motor", Address: "JL. Anggrek Industri", Location: "Cikarang", Promo: "Gratis nitrogen ban", Services: []string{"Isi nitrogen", "Tambal ban"}, Distance: "3 km"},

			// Bekasi (6)
			{ID: "m-13", Name: "Bengkel Cepat", Address: "Jl. Raya Karet", Location: "Bekasi", Promo: "Gratis cek motor", Services: []string{"Cek motor", "Ganti ban"}, Distance: "5 km"},
			{ID: "m-14", Name: "Bengkel Berkah", Address: "Jl. Anggrek No. 9", Location: "Bekasi", Promo: "Diskon ban 10%", Services: []string{"Ganti ban", "Tambal ban"}, Distance: "2.1 km"},
			{ID: "m-15", Name: "Patriot Garage", Address: "Jl. Ahmad Yani Bekasi", Location: "Bekasi", Promo: "Voucher servis ringan", Services: []string{"Servis ringan", "Ganti oli"}, Distance: "2.9 km"},
			{ID: "m-16", Name: "Harapan Motor", Address: "Jl. Harapan Indah", Location: "Bekasi", Promo: "Diskon 15% spooring", Services: []string{"Spooring", "Balancing"}, Distance: "4.6 km"},
			{ID: "m-17", Name: "Galaxy Service", Address: "Jl. Galaxy Raya", Location: "Bekasi", Promo: "Gratis konsultasi mesin", Services: []string{"Tune up", "Cek mesin"}, Distance: "3.3 km"},
			{ID: "m-18", Name: "Summarecon Auto", Address: "Jl. Boulevard Selatan", Location: "Bekasi", Promo: "Diskon filter & oli 10%", Services: []string{"Ganti oli", "Filter udara"}, Distance: "5.5 km"},

			// Bogor (6)
			{ID: "m-19", Name: "Laju Motor", Address: "Jl. Kramat No. 3", Location: "Bogor", Promo: "Gratis cek rem", Services: []string{"Cek rem", "Tune up"}, Distance: "4 km"},
			{ID: "m-20", Name: "Prima Service", Address: "Jl. Kenanga No. 7", Location: "Bogor", Promo: "Diskon servis 15%", Services: []string{"Servis ringan", "Ganti oli"}, Distance: "3.6 km"},
			{ID: "m-21", Name: "Pajajaran Garage", Address: "Jl. Pajajaran No. 10", Location: "Bogor", Promo: "Voucher 30k", Services: []string{"Tune up", "Ganti busi"}, Distance: "2.4 km"},
			{ID: "m-22", Name: "Puncak Tech", Address: "Jl. Raya Puncak KM 10", Location: "Bogor", Promo: "Diskon rem 12%", Services: []string{"Rem", "Suspensi"}, Distance: "7 km"},
			{ID: "m-23", Name: "Taman Safari Auto", Address: "Jl. Safari", Location: "Bogor", Promo: "Gratis cek ban", Services: []string{"Ban", "Cuci motor"}, Distance: "8.5 km"},
			{ID: "m-24", Name: "Dramaga Motor", Address: "Jl. Dramaga Raya", Location: "Bogor", Promo: "Diskon oli 10%", Services: []string{"Ganti oli", "Filter udara"}, Distance: "5.2 km"},
		}
	}

	topups := map[string]*entity.Product{}
	for _, p := range products {
		if strings.EqualFold(p.Category, "Topup Driver") {
			topups[strings.ToLower(p.Provider)] = p
		}
	}

	if len(topups) == 0 {
		topups["gojek"] = &entity.Product{Name: "Topup Gojek 50k", Price: 52000, Description: "Dummy topup Gojek"}
		topups["grab"] = &entity.Product{Name: "Topup Grab 50k", Price: 52500, Description: "Dummy topup Grab"}
		topups["maxim"] = &entity.Product{Name: "Topup Maxim 50k", Price: 53000, Description: "Dummy topup Maxim"}
		topups["indrive"] = &entity.Product{Name: "Topup InDrive 50k", Price: 53500, Description: "Dummy topup InDrive"}
	}

	topupCards := []map[string]interface{}{
		{"label": "Gojek", "product": topups["gojek"]},
		{"label": "Grab", "product": topups["grab"]},
		{"label": "Maxim", "product": topups["maxim"]},
		{"label": "InDrive", "product": topups["indrive"]},
	}

	locationOrder := []string{"Depok", "Cikarang", "Bekasi", "Bogor"}
	groupedMitras := map[string][]*entity.Mitra{}
	for _, m := range mitras {
		loc := m.Location
		if loc == "" {
			loc = "Lainnya"
		}
		groupedMitras[loc] = append(groupedMitras[loc], m)
	}

	mitraLocations := []struct {
		Name  string
		Items []*entity.Mitra
	}{}
	for _, loc := range locationOrder {
		if items, ok := groupedMitras[loc]; ok {
			mitraLocations = append(mitraLocations, struct {
				Name  string
				Items []*entity.Mitra
			}{Name: loc, Items: items})
			delete(groupedMitras, loc)
		}
	}
	for loc, items := range groupedMitras {
		mitraLocations = append(mitraLocations, struct {
			Name  string
			Items []*entity.Mitra
		}{Name: loc, Items: items})
	}

	return a.render(c, "home.html", map[string]interface{}{
		"title":          "Gaspoll - Topup Driver & Reward",
		"mitras":         mitras,
		"mitraLocations": mitraLocations,
		"products":       products,
		"topupCards":     topupCards,
	})
}

func (a *App) listProducts(c echo.Context) error {
	products, _ := a.product.List()
	return a.render(c, "products.html", map[string]interface{}{
		"title":    "Produk Digital",
		"products": products,
	})
}

func (a *App) handleCheckout(c echo.Context) error {
	user := c.Get(constant.ContextUserKey).(*entity.User)
	productID := c.FormValue("product_id")
	method := c.FormValue("payment_method")
	if method == "" {
		method = "QRIS"
	}

	order, instr, err := a.order.Checkout(user.ID, productID, method)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	refreshedUser, _ := a.repo.User.FindUserByID(user.ID)
	payload := map[string]interface{}{
		"order":       order,
		"instruction": instr,
		"user":        refreshedUser,
	}

	if c.Request().Header.Get("HX-Request") != "" {
		return c.Render(http.StatusOK, "order_result", payload)
	}
	return a.render(c, "order_result_page.html", payload)
}

func (a *App) handleDashboard(c echo.Context) error {
	user := c.Get(constant.ContextUserKey).(*entity.User)
	orders, _ := a.order.ListByUser(user.ID)
	vouchers, _ := a.voucher.ListByUser(user.ID)
	progress := a.reward.Progress(user)

	return a.render(c, "dashboard.html", map[string]interface{}{
		"title":    "Dashboard Driver",
		"user":     user,
		"orders":   orders,
		"vouchers": vouchers,
		"progress": progress,
	})
}

func (a *App) renderVoucherQR(c echo.Context) error {
	id := c.Param("id")
	v, err := a.voucher.Vouchers.FindVoucherByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "voucher tidak ditemukan")
	}
	img := makeSimpleQR(v.QRData)
	return c.Blob(http.StatusOK, "image/png", img)
}

func makeSimpleQR(data string) []byte {
	const size = 240
	const cell = 12
	grid := size / cell
	hash := sha1.Sum([]byte(data))

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	bg := color.RGBA{12, 18, 35, 255}
	fg := color.RGBA{52, 211, 153, 255}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, bg)
		}
	}

	idx := 0
	for y := 0; y < grid; y++ {
		for x := 0; x < grid; x++ {
			bit := hash[(idx/8)%len(hash)]>>(uint(idx%8))&1 == 1
			if bit {
				for yy := 0; yy < cell; yy++ {
					for xx := 0; xx < cell; xx++ {
						img.Set(x*cell+xx, y*cell+yy, fg)
					}
				}
			}
			idx++
		}
	}

	buf := bytes.NewBuffer(nil)
	_ = png.Encode(buf, img)
	return buf.Bytes()
}

func (a *App) showRedeem(c echo.Context) error {
	return a.render(c, "redeem.html", map[string]interface{}{
		"title": "Redeem Voucher Mitra",
	})
}

func (a *App) showMitra(c echo.Context) error {
	id := c.Param("id")
	m, err := a.mitra.Find(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Mitra tidak ditemukan")
	}
	return a.render(c, "mitra.html", map[string]interface{}{
		"title": "Promo " + m.Name,
		"mitra": m,
	})
}

func (a *App) showMitraRedeem(c echo.Context) error {
	id := c.Param("id")
	m, err := a.mitra.Find(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Mitra tidak ditemukan")
	}
	// Render the existing redeem template but include mitra info to prefill UI
	return a.render(c, "redeem.html", map[string]interface{}{
		"title": "Redeem Voucher Mitra",
		"mitra": m,
	})
}

// showIPL renders mock unpaid IPL bills for a given perumahan and house number.
func (a *App) showIPL(c echo.Context) error {
	perumahan := c.QueryParam("perumahan")
	nomor := c.QueryParam("nomor")
	if perumahan == "" && nomor == "" {
		// render empty form result page
		return a.render(c, "ipl_list.html", map[string]interface{}{
			"title": "Tagihan IPL",
		})
	}

	// create 3 mock unpaid bills for previous months
	now := time.Now()
	bills := []IPLBill{}
	for i := 0; i < 3; i++ {
		t := now.AddDate(0, -i-1, 0)
		period := t.Format("Jan 2006")
		amount := 75000 + i*10000
		bills = append(bills, IPLBill{House: nomor, Period: period, Amount: amount})
	}

	return a.render(c, "ipl_list.html", map[string]interface{}{
		"title":     "Tagihan IPL",
		"Perumahan": perumahan,
		"Nomor":     nomor,
		"Bills":     bills,
	})
}

// handleIPLPAY simulates payment for a period. For now, it returns a simple confirmation.
func (a *App) handleIPLPAY(c echo.Context) error {
	perumahan := c.FormValue("perumahan")
	nomor := c.FormValue("nomor")
	period := c.FormValue("period")

	// In a real app, you'd call payment & mark bill paid. Here we just show confirmation.
	return a.render(c, "ipl_list.html", map[string]interface{}{
		"title":     "Pembayaran IPL",
		"Perumahan": perumahan,
		"Nomor":     nomor,
		"Bills":     []IPLBill{},
		"message":   "Pembayaran untuk periode " + period + " berhasil.",
	})
}

func (a *App) handleRedeem(c echo.Context) error {
	code := strings.TrimSpace(c.FormValue("voucher_id"))
	v, err := a.voucher.Redeem(code)
	if err != nil {
		return c.String(http.StatusBadRequest, "Redeem gagal: "+err.Error())
	}
	return a.render(c, "redeem.html", map[string]interface{}{
		"title":   "Redeem Voucher Mitra",
		"success": true,
		"voucher": v,
	})
}

func (a *App) showSettlement(c echo.Context) error {
	settlements, _ := a.settlement.List()
	return a.render(c, "settlement.html", map[string]interface{}{
		"title":       "Settlement Admin",
		"settlements": settlements,
	})
}
