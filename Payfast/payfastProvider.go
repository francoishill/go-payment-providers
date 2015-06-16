package Payfast

import (
	"fmt"
	. "github.com/francoishill/go-payment-providers"
	. "github.com/francoishill/go-payment-providers/Payfast/SaleProvider"
)

type iLocalLogger interface {
	Notice(format string, v ...interface{})
}

type payfastProvider struct {
	logger        iLocalLogger
	saleProvider  IPayfastSaleProvider
	sandboxMode   bool
	merchantId    string
	passPhrase    string
	extractedData *extractedPostITNData
	sale          ISale
}

func CreatePayfastProvider(logger iLocalLogger, saleProvider IPayfastSaleProvider, sandboxMode bool, merchantId, passPhrase string) IPaymentProvider {
	return &payfastProvider{
		logger:       logger,
		saleProvider: saleProvider,
		sandboxMode:  sandboxMode,
		merchantId:   merchantId,
		passPhrase:   passPhrase,
		sale:         nil,
	}
}

func (this *payfastProvider) checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (this *payfastProvider) getRemoteHost() string {
	if this.sandboxMode {
		return fmt.Sprintf("sandbox.payfast.co.za")
	} else {
		return fmt.Sprintf("www.payfast.co.za")
	}
}
