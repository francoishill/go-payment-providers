package gopp

import (
	"net"
)

func createIpnHandler(paymentProvider IPaymentProvider) *ipnHandler {
	return &ipnHandler{
		paymentProvider: paymentProvider,
	}
}

type ipnHandler struct {
	paymentProvider IPaymentProvider
}

func (this *ipnHandler) verifyValidIPOfRequest(remoteIP string) {
	validHostNames := this.paymentProvider.GetValidHostNames()

	validIPaddresses := []string{}
	for _, host := range validHostNames {
		hostIPs, err := net.LookupIP(host)
		if err != nil {
			panic("Unable to lookup host IP: " + err.Error())
		}
		for _, ip := range hostIPs {
			validIPaddresses = append(validIPaddresses, ip.String())
		}
	}

	ipIsValid := false
	for _, validIP := range validIPaddresses {
		if remoteIP == validIP {
			ipIsValid = true
			break
		}
	}
	if !ipIsValid {
		panic("Invalid remote IP '" + remoteIP + "'")
	}
}

func (this *ipnHandler) verifySignatureOfPostData(requestPostBody []byte, remoteIp string) {
	allKeyValuePairsInCorrectOrder := readKeyValuePairsInCorrectOrderFromPostBody(requestPostBody)
	this.paymentProvider.VerifySignatureOfPostData(allKeyValuePairsInCorrectOrder, remoteIp)
}

func (this *ipnHandler) verifySaleDataMatch() {
	this.paymentProvider.VerifySaleDataMatch()
}

func (this *ipnHandler) verifyMerchantData() {
	this.paymentProvider.VerifyMerchantData()
}

func (this *ipnHandler) verifyFromGatewayTheySentTheRequest(remoteUserAgent string) {
	this.paymentProvider.VerifyFromGatewayTheySentTheRequest(remoteUserAgent)
}

func (this *ipnHandler) checkSaleAlreadyAuthorized() bool {
	return this.paymentProvider.CheckSaleAlreadyAuthorized()
}

func (this *ipnHandler) updateSale(eventHandler IIPNEventHandler) {
	this.paymentProvider.UpdateSale(eventHandler)
}
