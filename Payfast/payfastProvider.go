package Payfast

import (
	. "github.com/francoishill/go-payment-providers"
	. "github.com/francoishill/go-payment-providers/Payfast/PayfastContext"
	. "github.com/francoishill/go-payment-providers/Payfast/SaleProvider"
)

type iLocalLogger interface {
	Notice(format string, v ...interface{})
}

type payfastProvider struct {
	logger         iLocalLogger
	saleProvider   IPayfastSaleProvider
	payfastContext *PayfastContext
	extractedData  *extractedPostITNData
	sale           ISale
}

func CreatePayfastProvider(logger iLocalLogger, saleProvider IPayfastSaleProvider, payfastContext *PayfastContext) IPaymentProvider {
	return &payfastProvider{
		logger:         logger,
		saleProvider:   saleProvider,
		payfastContext: payfastContext,
		sale:           nil,
	}
}

func (this *payfastProvider) checkError(err error) {
	if err != nil {
		panic(err)
	}
}
