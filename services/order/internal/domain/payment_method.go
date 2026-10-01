package domain

type PaymentMethod string

const (
	PaymentTypeBankTransfer PaymentMethod = "bank_transfer"
	PaymentTypeVirtualAccount PaymentMethod = "virtual_account"
	PaymentTypeCreditCard PaymentMethod = "credit_card"
	PaymentTypeGopay PaymentMethod = "gopay"
)

func (p PaymentMethod) IsValid() bool {
	switch p {
	case PaymentTypeBankTransfer, PaymentTypeVirtualAccount, PaymentTypeCreditCard, PaymentTypeGopay:
		return true
	default:
		return false
	}
}
