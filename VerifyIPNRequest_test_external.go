package gopp

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func DoTestVerifyValidIPOfRequest(t *testing.T, paymentProvider IPaymentProvider, remoteValidIP, remoteInvalidIP string) {
	Convey("Verifying IPN request", t, func() {
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
			"Invalid remote IP '"+remoteInvalidIP+"'",
		)
	})
}

func DoTestVerifySignatureOfPostData(t *testing.T, paymentProvider IPaymentProvider, remoteValidIP, remoteInvalidIP string) {
	Convey("Verifying IPN request", t, func() {
		handler := createIpnHandler(paymentProvider)

		So(
			func() {
				handler.verifySignatureOfPostData(remoteValidIP)
			},
			ShouldNotPanic,
		)

		So(
			func() {
				handler.verifySignatureOfPostData(remoteInvalidIP)
			},
			ShouldPanicWith,
			"Invalid remote IP '"+remoteInvalidIP+"'",
		)
	})
}
