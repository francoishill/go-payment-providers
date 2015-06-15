package gopp

func VerifyIPNRequest(paymentProvider IPaymentProvider, remoteIP string, requestPostBody []byte) {
	handler := createIpnHandler(paymentProvider)
	handler.verifyValidIPOfRequest(remoteIP)
	handler.verifySignatureOfPostData(requestPostBody)
}
