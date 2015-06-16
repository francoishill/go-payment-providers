package gopp

import (
	"fmt"
)

type IIPNEventHandler interface {
	OnError(err error)
	SaleAlreadyPreviouslyAuthorized()
	SaleSucceeded()
	SaleFailed()
	SalePending()
}

func VerifyIPNRequest(paymentProvider IPaymentProvider, eventHandler IIPNEventHandler, remoteIP, remoteUserAgent string, requestPostBody []byte) {
	defer func() {
		if r := recover(); r != nil {
			eventHandler.OnError(fmt.Errorf("%+v", r))
		}
	}()

	handler := createIpnHandler(paymentProvider)
	handler.verifyValidIPOfRequest(remoteIP)
	handler.verifySignatureOfPostData(requestPostBody, remoteIP)
	handler.verifySaleDataMatch()
	handler.verifyMerchantData()
	handler.verifyFromGatewayTheySentTheRequest(remoteUserAgent)
	if handler.checkSaleAlreadyAuthorized() {
		eventHandler.SaleAlreadyPreviouslyAuthorized()
		return
	}
	handler.updateSale(eventHandler)
}
