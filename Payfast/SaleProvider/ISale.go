package SaleProvider

type ISale interface {
	GetId() string
	GetItemName() string
	GetItemDescription() string
	GetAmountGross() float32
	IsAuthorized() bool
	UpdateSaleAsAuthorizedCompleted(payfastPaymentId string, amountFee, amountNet float32)
	UpdateSaleAsFailed(payfastPaymentId string)
	UpdateSaleAsPending(payfastPaymentId string)
}
