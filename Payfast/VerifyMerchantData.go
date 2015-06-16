package Payfast

func (this *payfastProvider) VerifyMerchantData() {
	if this.payfastContext.GetMerchantId() != this.extractedData.MerchantId {
		panic("Invalid merchant ID")
	}
}
