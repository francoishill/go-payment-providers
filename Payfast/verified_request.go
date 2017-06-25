package payfast

type verifiedRequest struct {
	ParamStringForRemoteVerify string
	MerchantID                 string
	AmountGross                float64
	PaymentStatus              PaymentStatus

	ActualSale Sale
}
