package domain

type PaymentMethod string

const (
	BANK_TRANSFER   PaymentMethod = "BANK_TRANSFER"
	VIRTUAL_ACCOUNT PaymentMethod = "VIRTUAL_ACCOUNT"
	CREDIT_CARD     PaymentMethod = "CREDIT_CARD"
	GOPAY           PaymentMethod = "GOPAY"
)

func (p PaymentMethod) IsValid() bool {
	switch p {
	case BANK_TRANSFER, VIRTUAL_ACCOUNT, CREDIT_CARD, GOPAY:
		return true
	default:
		return false
	}
}
