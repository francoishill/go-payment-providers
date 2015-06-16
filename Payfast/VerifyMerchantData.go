package Payfast

func (this *payfastProvider) VerifyMerchantData() {
	if this.merchantId != this.extractedData.MerchantId {
		panic("Invalid merchant ID")
	}
}
