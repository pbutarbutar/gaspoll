package memory

import (
	"errors"
	"sync"
	"time"

	"gaspoll/internal/entity"
	"gaspoll/internal/repository"
)

type Store struct {
	mu          sync.RWMutex
	users       map[string]*entity.User
	products    map[string]*entity.Product
	orders      map[string]*entity.Order
	vouchers    map[string]*entity.Voucher
	mitras      map[string]*entity.Mitra
	settlements map[string]*entity.Settlement
}

func NewSeededRepository() *repository.Repository {
	store := &Store{
		users:       map[string]*entity.User{},
		products:    map[string]*entity.Product{},
		orders:      map[string]*entity.Order{},
		vouchers:    map[string]*entity.Voucher{},
		mitras:      map[string]*entity.Mitra{},
		settlements: map[string]*entity.Settlement{},
	}
	store.seed()

	return &repository.Repository{
		User:       userRepo{store},
		Product:    productRepo{store},
		Order:      orderRepo{store},
		Voucher:    voucherRepo{store},
		Mitra:      mitraRepo{store},
		Settlement: settlementRepo{store},
	}
}

func (s *Store) seed() {
	now := time.Now()
	driver := &entity.User{
		ID:          "user-1",
		Name:        "Budi Driver",
		Phone:       "081234567890",
		VehicleType: "Motor",
		Points:      2,
		RewardGoal:  5,
		CreatedAt:   now.Add(-24 * time.Hour),
	}

	products := []*entity.Product{
		{ID: "p-gj-20", Name: "Topup Gojek 20k", Category: "Topup Driver", Provider: "Gojek", Price: 21000, Description: "Saldo Gojek driver instan", Tags: []string{"topup", "gojek"}},
		{ID: "p-grab-50", Name: "Topup Grab 50k", Category: "Topup Driver", Provider: "Grab", Price: 50500, Description: "Saldo Grab driver, ready 24 jam", Tags: []string{"topup", "grab"}},
		{ID: "p-maxim-25", Name: "Topup Maxim 25k", Category: "Topup Driver", Provider: "Maxim", Price: 26000, Description: "Topup Maxim driver promo", Tags: []string{"topup", "maxim"}},
		{ID: "p-ind-30", Name: "Topup InDrive 30k", Category: "Topup Driver", Provider: "InDrive", Price: 31000, Description: "Voucher saldo InDrive driver", Tags: []string{"topup", "indrive"}},
		{ID: "p-pulsa-25", Name: "Pulsa 25k", Category: "Pulsa", Provider: "Telkomsel", Price: 26000, Description: "Pulsa reguler", Tags: []string{"pulsa"}},
		{ID: "p-data-20", Name: "Paket Data 20GB", Category: "Paket Internet", Provider: "Indosat", Price: 72000, Description: "Paket data 30 hari", Tags: []string{"data", "internet"}},
	}

	mitras := []*entity.Mitra{
		// Depok
		{ID: "m-depok-oto", Name: "Oto Bengkel", Address: "Jl. Melati No. 12", Location: "Depok", Promo: "Diskon oli 20%", Services: []string{"Ganti oli", "Tune up"}, Distance: "1.2 km", BasePrice: 120000, DiscountPercent: 20, PromoDetail: "Promo ini bisa didapatkan dengan syarat top driver dalam 1 minggu 10X.", Packages: []entity.Package{
			{ID: "pkg-1", Name: "Service + Ganti Oli", Description: "Service lengkap + penggantian oli mesin standar.", BasePrice: 120000, DiscountPercent: 20},
			{ID: "pkg-2", Name: "Ganti Kampas Rem", Description: "Penggantian kampas rem depan atau belakang sesuai kebutuhan.", BasePrice: 150000, DiscountPercent: 10},
			{ID: "pkg-3", Name: "Ganti 1 Ban Motor", Description: "Penggantian 1 buah ban motor termasuk pemasangan dan penyeimbangan.", BasePrice: 200000, DiscountPercent: 15},
			{ID: "pkg-4", Name: "Ganti 2 Ban Motor", Description: "Penggantian 2 ban motor (diskon paket).", BasePrice: 380000, DiscountPercent: 18},
		}},
		{ID: "m-depok-garasi", Name: "Garasi Juara", Address: "Jl. Pahlawan No. 45", Location: "Depok", Promo: "Voucher servis ringan", Services: []string{"Spooring", "Balancing"}, Distance: "3 km", BasePrice: 85000, DiscountPercent: 0, PromoDetail: "Tunjukkan voucher untuk servis ringan gratis (syarat berlaku)."},
		{ID: "m-depok-sentosa", Name: "Sentosa Motor", Address: "Jl. Sudirman No. 10", Location: "Depok", Promo: "Diskon filter udara", Services: []string{"Filter udara", "Tune up"}, Distance: "4.5 km", BasePrice: 90000, DiscountPercent: 10, PromoDetail: "Diskon filter udara, tidak termasuk pemasangan khusus."},
		{ID: "m-depok-lenteng", Name: "Lenteng Motor", Address: "Jl. Lenteng Agung Raya", Location: "Depok", Promo: "Cek rem gratis", Services: []string{"Cek rem", "Ganti kampas"}, Distance: "2.9 km", BasePrice: 0, DiscountPercent: 100, PromoDetail: "Cek rem gratis untuk semua tipe motor."},
		{ID: "m-depok-margonda", Name: "Margonda Service", Address: "Jl. Margonda Raya No. 88", Location: "Depok", Promo: "Diskon servis 15%", Services: []string{"Servis ringan", "Ganti oli"}, Distance: "3.4 km", BasePrice: 150000, DiscountPercent: 15, PromoDetail: "Diskon berlaku untuk servis paket hemat."},
		{ID: "m-depok-sawangan", Name: "Sawangan Garage", Address: "Jl. Raya Sawangan", Location: "Depok", Promo: "Voucher 25k", Services: []string{"Tune up", "Ganti ban"}, Distance: "5 km", BasePrice: 230000, DiscountPercent: 10, PromoDetail: "Voucher potongan langsung 25k untuk layanan ban."},

		// Cikarang
		{ID: "m-cikarang-jaya", Name: "Jaya Auto", Address: "Jl. Panjang No. 88", Location: "Cikarang", Promo: "Voucher 25k", Services: []string{"Spooring", "Balancing"}, Distance: "6 km", BasePrice: 250000, DiscountPercent: 10, PromoDetail: "Voucher potongan langsung sebesar 25k untuk layanan tertentu."},
		{ID: "m-cikarang-green", Name: "Green Garage", Address: "Jl. Cemara No. 2", Location: "Cikarang", Promo: "Paket oli hemat", Services: []string{"Ganti oli", "Cuci motor"}, Distance: "2.8 km", BasePrice: 100000, DiscountPercent: 15, PromoDetail: "Paket hemat untuk ganti oli + cuci motor."},
		{ID: "m-cikarang-delta", Name: "Delta Service", Address: "Jl. Industri No. 5", Location: "Cikarang", Promo: "Diskon tune up 10%", Services: []string{"Tune up", "Ganti busi"}, Distance: "4.2 km", BasePrice: 180000, DiscountPercent: 10, PromoDetail: "Diskon berlaku untuk servis reguler."},
		{ID: "m-cikarang-mega", Name: "Mega Auto", Address: "Jl. Boulevard Cikarang", Location: "Cikarang", Promo: "Cek motor gratis", Services: []string{"Cek motor", "Ganti oli"}, Distance: "3.7 km", BasePrice: 0, DiscountPercent: 100, PromoDetail: "Cek motor gratis tanpa syarat pembelian."},
		{ID: "m-cikarang-tech", Name: "C-Tech Workshop", Address: "Jl. Kenari No. 3", Location: "Cikarang", Promo: "Diskon kampas rem 12%", Services: []string{"Rem", "Suspensi"}, Distance: "5.3 km", BasePrice: 210000, DiscountPercent: 12, PromoDetail: "Diskon kampas rem untuk semua tipe motor."},
		{ID: "m-cikarang-benua", Name: "Benua Motor", Address: "JL. Anggrek Industri", Location: "Cikarang", Promo: "Gratis nitrogen ban", Services: []string{"Isi nitrogen", "Tambal ban"}, Distance: "3 km", BasePrice: 50000, DiscountPercent: 0, PromoDetail: "Isi nitrogen gratis untuk driver Gaspoll."},

		// Bekasi
		{ID: "m-bekasi-cepat", Name: "Bengkel Cepat", Address: "Jl. Raya Karet", Location: "Bekasi", Promo: "Gratis cek motor", Services: []string{"Cek motor", "Ganti ban"}, Distance: "5 km", BasePrice: 0, DiscountPercent: 100, PromoDetail: "Gratis cek motor, hanya berlaku untuk 1 kali per pengguna."},
		{ID: "m-bekasi-berkah", Name: "Bengkel Berkah", Address: "Jl. Anggrek No. 9", Location: "Bekasi", Promo: "Diskon ban 10%", Services: []string{"Ganti ban", "Tambal ban"}, Distance: "2.1 km", BasePrice: 300000, DiscountPercent: 10, PromoDetail: "Diskon otomatis untuk tipe ban tertentu."},
		{ID: "m-bekasi-patriot", Name: "Patriot Garage", Address: "Jl. Ahmad Yani Bekasi", Location: "Bekasi", Promo: "Voucher servis ringan", Services: []string{"Servis ringan", "Ganti oli"}, Distance: "2.9 km", BasePrice: 130000, DiscountPercent: 12, PromoDetail: "Voucher berlaku untuk servis ringan paket A."},
		{ID: "m-bekasi-harapan", Name: "Harapan Motor", Address: "Jl. Harapan Indah", Location: "Bekasi", Promo: "Diskon 15% spooring", Services: []string{"Spooring", "Balancing"}, Distance: "4.6 km", BasePrice: 260000, DiscountPercent: 15, PromoDetail: "Diskon spooring untuk mobil dan motor matik."},
		{ID: "m-bekasi-galaxy", Name: "Galaxy Service", Address: "Jl. Galaxy Raya", Location: "Bekasi", Promo: "Gratis konsultasi mesin", Services: []string{"Tune up", "Cek mesin"}, Distance: "3.3 km", BasePrice: 0, DiscountPercent: 100, PromoDetail: "Gratis pengecekan awal mesin."},
		{ID: "m-bekasi-summa", Name: "Summarecon Auto", Address: "Jl. Boulevard Selatan", Location: "Bekasi", Promo: "Diskon filter & oli 10%", Services: []string{"Ganti oli", "Filter udara"}, Distance: "5.5 km", BasePrice: 200000, DiscountPercent: 10, PromoDetail: "Diskon kombinasi filter + oli."},

		// Bogor
		{ID: "m-bogor-laju", Name: "Laju Motor", Address: "Jl. Kramat No. 3", Location: "Bogor", Promo: "Gratis cek rem", Services: []string{"Cek rem", "Tune up"}, Distance: "4 km", BasePrice: 0, DiscountPercent: 100, PromoDetail: "Cek rem gratis, tanpa pembelian komponen."},
		{ID: "m-bogor-prima", Name: "Prima Service", Address: "Jl. Kenanga No. 7", Location: "Bogor", Promo: "Diskon servis 15%", Services: []string{"Servis ringan", "Ganti oli"}, Distance: "3.6 km", BasePrice: 150000, DiscountPercent: 15, PromoDetail: "Diskon berlaku untuk layanan servis ringan."},
		{ID: "m-bogor-pajajaran", Name: "Pajajaran Garage", Address: "Jl. Pajajaran No. 10", Location: "Bogor", Promo: "Voucher 30k", Services: []string{"Tune up", "Ganti busi"}, Distance: "2.4 km", BasePrice: 175000, DiscountPercent: 10, PromoDetail: "Voucher 30k untuk paket tune up."},
		{ID: "m-bogor-puncak", Name: "Puncak Tech", Address: "Jl. Raya Puncak KM 10", Location: "Bogor", Promo: "Diskon rem 12%", Services: []string{"Rem", "Suspensi"}, Distance: "7 km", BasePrice: 220000, DiscountPercent: 12, PromoDetail: "Diskon rem untuk motor beban berat."},
		{ID: "m-bogor-safari", Name: "Taman Safari Auto", Address: "Jl. Safari", Location: "Bogor", Promo: "Gratis cek ban", Services: []string{"Ban", "Cuci motor"}, Distance: "8.5 km", BasePrice: 0, DiscountPercent: 100, PromoDetail: "Gratis cek ban + tekanan angin."},
		{ID: "m-bogor-dramaga", Name: "Dramaga Motor", Address: "Jl. Dramaga Raya", Location: "Bogor", Promo: "Diskon oli 10%", Services: []string{"Ganti oli", "Filter udara"}, Distance: "5.2 km", BasePrice: 110000, DiscountPercent: 10, PromoDetail: "Diskon oli untuk semua produk standar."},
	}

	order := &entity.Order{
		ID:            "ord-1001",
		UserID:        driver.ID,
		ProductID:     products[0].ID,
		ProductName:   products[0].Name,
		Amount:        products[0].Price,
		PaymentMethod: "QRIS",
		Status:        "settled",
		WhatsAppLink:  "https://wa.me/6281234567890?text=Gaspoll%20Order%201001",
		RewardEarned:  1,
		CreatedAt:     now.Add(-12 * time.Hour),
	}

	voucher := &entity.Voucher{
		ID:          "v-7001",
		Code:        "VCHR-7001",
		UserID:      driver.ID,
		MitraID:     mitras[0].ID,
		Title:       "Voucher Servis Ringan",
		Description: "Tunjukkan ke bengkel mitra, berlaku sampai akhir bulan",
		PointsUsed:  2,
		IsRedeemed:  false,
		ExpiredAt:   now.Add(7 * 24 * time.Hour),
		QRData:      "voucher:v-7001",
	}

	s.users[driver.ID] = driver
	for _, p := range products {
		s.products[p.ID] = p
	}
	// default package set applied to mitras without explicit packages
	defaultPackages := []entity.Package{
		{ID: "pkg-1", Name: "Service + Ganti Oli", Description: "Service lengkap + penggantian oli mesin standar.", BasePrice: 120000, DiscountPercent: 20},
		{ID: "pkg-2", Name: "Ganti Kampas Rem", Description: "Penggantian kampas rem depan atau belakang sesuai kebutuhan.", BasePrice: 150000, DiscountPercent: 10},
		{ID: "pkg-3", Name: "Ganti 1 Ban Motor", Description: "Penggantian 1 buah ban motor termasuk pemasangan dan penyeimbangan.", BasePrice: 200000, DiscountPercent: 15},
		{ID: "pkg-4", Name: "Ganti 2 Ban Motor", Description: "Penggantian 2 ban motor (diskon paket).", BasePrice: 380000, DiscountPercent: 18},
	}

	for _, m := range mitras {
		if len(m.Packages) == 0 {
			// copy default packages so slices are independent per mitra
			m.Packages = make([]entity.Package, len(defaultPackages))
			copy(m.Packages, defaultPackages)
		}
		s.mitras[m.ID] = m
	}
	s.orders[order.ID] = order
	s.vouchers[voucher.ID] = voucher
}

type userRepo struct{ store *Store }
type productRepo struct{ store *Store }
type orderRepo struct{ store *Store }
type voucherRepo struct{ store *Store }
type mitraRepo struct{ store *Store }
type settlementRepo struct{ store *Store }

// UserRepository
func (r userRepo) ListUsers() ([]*entity.User, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	result := make([]*entity.User, 0, len(r.store.users))
	for _, u := range r.store.users {
		cp := *u
		result = append(result, &cp)
	}
	return result, nil
}

func (r userRepo) FindUserByID(id string) (*entity.User, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	user, ok := r.store.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	cp := *user
	return &cp, nil
}

func (r userRepo) FindUserByPhone(phone string) (*entity.User, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	for _, u := range r.store.users {
		if u.Phone == phone {
			cp := *u
			return &cp, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r userRepo) Save(user *entity.User) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	r.store.users[user.ID] = &entity.User{
		ID:           user.ID,
		Name:         user.Name,
		Phone:        user.Phone,
		VehicleType:  user.VehicleType,
		Points:       user.Points,
		RewardGoal:   user.RewardGoal,
		RewardStarts: user.RewardStarts,
		CreatedAt:    user.CreatedAt,
	}
	return nil
}

// ProductRepository
func (r productRepo) ListProducts() ([]*entity.Product, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	result := make([]*entity.Product, 0, len(r.store.products))
	for _, p := range r.store.products {
		cp := *p
		result = append(result, &cp)
	}
	return result, nil
}

func (r productRepo) FindProductByID(id string) (*entity.Product, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	p, ok := r.store.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	cp := *p
	return &cp, nil
}

// OrderRepository
func (r orderRepo) ListByUser(userID string) ([]*entity.Order, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	result := []*entity.Order{}
	for _, o := range r.store.orders {
		if o.UserID == userID {
			cp := *o
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (r orderRepo) Save(order *entity.Order) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	r.store.orders[order.ID] = &entity.Order{
		ID:             order.ID,
		UserID:         order.UserID,
		ProductID:      order.ProductID,
		ProductName:    order.ProductName,
		Amount:         order.Amount,
		PaymentMethod:  order.PaymentMethod,
		Status:         order.Status,
		WhatsAppLink:   order.WhatsAppLink,
		PaymentCode:    order.PaymentCode,
		PaymentPayload: order.PaymentPayload,
		RewardEarned:   order.RewardEarned,
		CreatedAt:      order.CreatedAt,
	}
	return nil
}

// VoucherRepository
func (r voucherRepo) ListByUser(userID string) ([]*entity.Voucher, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	result := []*entity.Voucher{}
	for _, v := range r.store.vouchers {
		if v.UserID == userID {
			cp := *v
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (r voucherRepo) FindVoucherByID(id string) (*entity.Voucher, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	v, ok := r.store.vouchers[id]
	if !ok {
		return nil, errors.New("voucher not found")
	}
	cp := *v
	return &cp, nil
}

func (r voucherRepo) Save(voucher *entity.Voucher) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	r.store.vouchers[voucher.ID] = &entity.Voucher{
		ID:          voucher.ID,
		Code:        voucher.Code,
		UserID:      voucher.UserID,
		MitraID:     voucher.MitraID,
		Title:       voucher.Title,
		Description: voucher.Description,
		PointsUsed:  voucher.PointsUsed,
		IsRedeemed:  voucher.IsRedeemed,
		RedeemedAt:  voucher.RedeemedAt,
		ExpiredAt:   voucher.ExpiredAt,
		QRData:      voucher.QRData,
	}
	return nil
}

// MitraRepository
func (r mitraRepo) ListMitras() ([]*entity.Mitra, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	result := []*entity.Mitra{}
	for _, m := range r.store.mitras {
		cp := *m
		result = append(result, &cp)
	}
	return result, nil
}

func (r mitraRepo) FindMitraByID(id string) (*entity.Mitra, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	m, ok := r.store.mitras[id]
	if !ok {
		return nil, errors.New("mitra not found")
	}
	cp := *m
	return &cp, nil
}

// SettlementRepository
func (r settlementRepo) ListSettlements() ([]*entity.Settlement, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	result := []*entity.Settlement{}
	for _, st := range r.store.settlements {
		cp := *st
		result = append(result, &cp)
	}
	return result, nil
}

func (r settlementRepo) Save(settlement *entity.Settlement) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	r.store.settlements[settlement.ID] = &entity.Settlement{
		ID:           settlement.ID,
		MitraID:      settlement.MitraID,
		Period:       settlement.Period,
		TotalVoucher: settlement.TotalVoucher,
		Amount:       settlement.Amount,
		Status:       settlement.Status,
	}
	return nil
}
