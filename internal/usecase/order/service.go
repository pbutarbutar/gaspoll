package order

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"gaspoll/internal/entity"
	"gaspoll/internal/repository"
	"gaspoll/internal/usecase/reward"
)

type Service struct {
	Users    repository.UserRepository
	Products repository.ProductRepository
	Orders   repository.OrderRepository
	Vouchers repository.VoucherRepository
	Mitras   repository.MitraRepository
	Rewards  reward.Service
}

func (s Service) Checkout(userID, productID, method string) (*entity.Order, *entity.PaymentInstruction, error) {
	user, err := s.Users.FindUserByID(userID)
	if err != nil {
		return nil, nil, errors.New("user tidak ditemukan")
	}
	product, err := s.Products.FindProductByID(productID)
	if err != nil {
		return nil, nil, errors.New("produk tidak ditemukan")
	}

	order := &entity.Order{
		ID:            fmt.Sprintf("ord-%d", time.Now().UnixNano()),
		UserID:        user.ID,
		ProductID:     product.ID,
		ProductName:   product.Name,
		Amount:        product.Price,
		PaymentMethod: method,
		Status:        "waiting_payment",
		WhatsAppLink:  fmt.Sprintf("https://wa.me/6281234567890?text=%s", url.QueryEscape("Halo admin, ini bukti bayar untuk "+product.Name)),
		RewardEarned:  1,
		CreatedAt:     time.Now(),
	}

	payment := buildPaymentInstruction(method, order)

	s.Rewards.AddPoints(user, order.RewardEarned)
	issueVoucher := s.Rewards.ShouldIssueVoucher(user)
	if issueVoucher {
		s.Rewards.ConsumeForVoucher(user)
		voucher := s.makeVoucher(user.ID)
		s.Vouchers.Save(voucher)
	}

	if err := s.Users.Save(user); err != nil {
		return nil, nil, err
	}
	if err := s.Orders.Save(order); err != nil {
		return nil, nil, err
	}

	return order, payment, nil
}

func (s Service) ListByUser(userID string) ([]*entity.Order, error) {
	return s.Orders.ListByUser(userID)
}

func (s Service) makeVoucher(userID string) *entity.Voucher {
	mitras, _ := s.Mitras.ListMitras()
	target := ""
	if len(mitras) > 0 {
		target = mitras[0].ID
	}
	now := time.Now()
	return &entity.Voucher{
		ID:          fmt.Sprintf("v-%d", now.UnixNano()),
		Code:        fmt.Sprintf("VCHR-%d", now.Unix()),
		UserID:      userID,
		MitraID:     target,
		Title:       "Voucher Servis Bengkel Mitra",
		Description: "Tunjukkan QR ke bengkel mitra, khusus driver Gaspoll",
		PointsUsed:  5,
		IsRedeemed:  false,
		ExpiredAt:   now.Add(30 * 24 * time.Hour),
		QRData:      fmt.Sprintf("voucher:%d:%s", now.Unix(), userID),
	}
}

func buildPaymentInstruction(method string, order *entity.Order) *entity.PaymentInstruction {
	switch method {
	case "VA":
		return &entity.PaymentInstruction{
			Method: "Virtual Account",
			Code:   "8808" + order.ID[len(order.ID)-4:],
			Steps: []string{
				"Buka aplikasi bank pilihan",
				"Pilih pembayaran VA",
				"Masukkan kode di atas dan konfirmasi",
				"Simpan bukti untuk verifikasi WA",
			},
		}
	case "E-Wallet":
		return &entity.PaymentInstruction{
			Method: "E-Wallet",
			Code:   "DANA/SPPAY: " + order.ID[len(order.ID)-4:],
			Steps: []string{
				"Buka Dana/ShopeePay/OVO",
				"Pilih transfer ke pengguna",
				"Masukkan kode dan nominal sesuai tagihan",
			},
		}
	default:
		return &entity.PaymentInstruction{
			Method:   "QRIS",
			QRString: "00020101021226680012IDCOBAQRIS0002015204599953033605802ID5908GASPOLL6007JAKARTA",
			Steps: []string{
				"Simpan/buka QR lalu bayar via app bank/e-wallet",
				"Setelah bayar, klik tombol buka WhatsApp untuk konfirmasi",
			},
		}
	}
}
