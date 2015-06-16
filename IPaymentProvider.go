package gopp

type IPaymentProvider interface {
	GetValidHostNames() []string
	VerifySignatureOfPostData(postDataInCorrectOrder SliceOfPostKeyValue, remoteIP string)
	VerifySaleDataMatch()
	VerifyMerchantData()
	VerifyFromGatewayTheySentTheRequest(remoteUserAgent string)
	CheckSaleAlreadyAuthorized() bool
	UpdateSale(eventHandler IIPNEventHandler)
}
