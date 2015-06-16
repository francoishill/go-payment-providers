package Payfast

func (this *payfastProvider) CheckSaleAlreadyAuthorized() bool {
	return this.sale.IsAuthorized()
}
