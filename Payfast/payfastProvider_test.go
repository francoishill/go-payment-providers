package Payfast

import (
	"fmt"
	. "github.com/francoishill/go-payment-providers"
	"testing"
)

type tmpLogger struct{}

func (this *tmpLogger) Notice(format string, v ...interface{}) {
	fmt.Println("[NOTICE]" + fmt.Sprintf(format, v...))
}

func TestPayfastVerifyValidIPOfRequest(t *testing.T) {
	logger := &tmpLogger{}
	tmpRemoteIP := "41.74.179.210"
	passPhrase := "passphrase"

	paymentProvider := CreatePayfastProvider(logger, tmpRemoteIP, passPhrase)
	remoteValidIP := "41.74.179.210"
	remoteInvalidIP := "41.74.179.230"
	DoTestVerifyValidIPOfRequest(t, paymentProvider, remoteValidIP, remoteInvalidIP)
}
