package gopp

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func DoTestVerifyValidIPOfRequest(t *testing.T, paymentProvider IPaymentProvider, remoteValidIP, remoteInvalidIP string, invalidExpectedPanic interface{}) {
	Convey("verifyValidIPOfRequest", t, func() {
		handler := createIpnHandler(paymentProvider)

		So(
			func() {
				handler.verifyValidIPOfRequest(remoteValidIP)
			},
			ShouldNotPanic,
		)

		So(
			func() {
				handler.verifyValidIPOfRequest(remoteInvalidIP)
			},
			ShouldPanicWith,
			invalidExpectedPanic,
		)
	})
}

func DoTestVerifySignatureOfPostData(t *testing.T, paymentProvider IPaymentProvider, remoteIp string, validRequestPostBody, invalidRequestPostBody []byte, invalidExpectedPanic interface{}) {
	Convey("verifySignatureOfPostData", t, func() {
		handler := createIpnHandler(paymentProvider)

		So(
			func() {
				handler.verifySignatureOfPostData(validRequestPostBody, remoteIp)
			},
			ShouldNotPanic,
		)

		So(
			func() {
				handler.verifySignatureOfPostData(invalidRequestPostBody, remoteIp)
			},
			ShouldPanicWith,
			invalidExpectedPanic,
		)
	})
}

func DoTestVerifySaleDataMatch(t *testing.T, validPaymentProvider, invalidPaymentProvider IPaymentProvider, invalidExpectedPanic interface{}) {
	Convey("verifySaleDataMatch", t, func() {
		Convey("Payment provider with valid sale data", func() {
			So(
				func() {
					handler := createIpnHandler(validPaymentProvider)
					handler.verifySaleDataMatch()
				},
				ShouldNotPanic,
			)
		})

		Convey("Payment provider with invalid sale data", func() {
			So(
				func() {
					handler := createIpnHandler(invalidPaymentProvider)
					handler.verifySaleDataMatch()
				},
				ShouldPanicWith,
				invalidExpectedPanic,
			)
		})
	})
}

func DoTestVerifyMerchantDataMatch(t *testing.T, validPaymentProvider, invalidPaymentProvider IPaymentProvider, invalidExpectedPanic interface{}) {
	Convey("verifySaleDataMatch", t, func() {
		Convey("Payment provider with valid Merchant ID", func() {
			So(
				func() {
					handler := createIpnHandler(validPaymentProvider)
					handler.verifyMerchantData()
				},
				ShouldNotPanic,
			)
		})
		Convey("Payment provider with invalid Merchant ID", func() {
			So(
				func() {
					handler := createIpnHandler(invalidPaymentProvider)
					handler.verifyMerchantData()
				},
				ShouldPanicWith,
				invalidExpectedPanic,
			)
		})
	})
}

func DoTestVerifyFromGatewayTheySentTheRequest(t *testing.T, remoteUserAgent string, validPaymentProvider, invalidPaymentProvider IPaymentProvider, invalidExpectedPanic interface{}) {
	Convey("verifySaleDataMatch", t, func() {
		Convey("Payment provider returning valid from provider", func() {
			So(
				func() {
					handler := createIpnHandler(validPaymentProvider)
					handler.verifyFromGatewayTheySentTheRequest(remoteUserAgent)
				},
				ShouldNotPanic,
			)
		})
		Convey("Payment provider failing to return valid from provider", func() {
			So(
				func() {
					handler := createIpnHandler(invalidPaymentProvider)
					handler.verifyFromGatewayTheySentTheRequest(remoteUserAgent)
				},
				ShouldPanicWith,
				invalidExpectedPanic,
			)
		})
	})
}
