package Payfast

import (
	. "github.com/francoishill/go-payment-providers"
	"strings"
)

func (this *payfastProvider) UpdateSale(eventHandler IIPNEventHandler) {
	paymentStatusUppercase := strings.ToUpper(this.extractedData.PaymentStatus)

	switch paymentStatusUppercase {
	case "COMPLETE":
		this.sale.UpdateSaleAsAuthorizedCompleted(this.extractedData.PayfastPaymentId, this.extractedData.AmountFee, this.extractedData.AmountNet)
		eventHandler.SaleSucceeded()
		break
	case "FAILED":
		this.sale.UpdateSaleAsFailed(this.extractedData.PayfastPaymentId)
		eventHandler.SaleFailed()
		break
	case "PENDING":
		this.sale.UpdateSaleAsPending(this.extractedData.PayfastPaymentId)
		eventHandler.SalePending()
		break
	default:
		// Note on payfast site (https://www.payfast.co.za/documentation/itn/): "If unknown status, do nothing (safest course of action)"
		break
	}
}
