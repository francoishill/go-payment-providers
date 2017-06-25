package payfast

type ITNSale interface {
	AmountGross() float64
	AlreadyAuthorized() bool
}
