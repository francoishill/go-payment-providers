package PayfastContext

import (
	"fmt"
	"strings"
)

type PayfastContext struct {
	sandboxMode bool
	merchantId  string
	merchantKey string
	passphrase  string
	returnUrl   string
	cancelUrl   string
	notifyUrl   string
}

func CreatePayfastContext(sandboxMode bool, merchantId, merchantKey, passphrase, return_url, cancel_url, notify_url string) *PayfastContext {
	return &PayfastContext{
		sandboxMode: sandboxMode,
		merchantId:  merchantId,
		merchantKey: merchantKey,
		passphrase:  passphrase,
		returnUrl:   strings.TrimRight(return_url, "/"),
		cancelUrl:   strings.TrimRight(cancel_url, "/"),
		notifyUrl:   notify_url,
	}
}

func (this *PayfastContext) IsSandboxmode() bool {
	return this.sandboxMode
}

func (this *PayfastContext) GetMerchantId() string {
	return this.merchantId
}

func (this *PayfastContext) GetMerchantKey() string {
	return this.merchantKey
}

func (this *PayfastContext) GetPassphrase() string {
	return this.passphrase
}

func (this *PayfastContext) GetReturnUrl(saleId string) string {
	return fmt.Sprintf("%s/%s", this.returnUrl, saleId)
}

func (this *PayfastContext) GetCancelUrl(saleId string) string {
	return fmt.Sprintf("%s/%s", this.cancelUrl, saleId)
}

func (this *PayfastContext) GetNotifyUrl() string {
	return this.notifyUrl
}

func (this *PayfastContext) GetAuthServiceUrl() string {
	if this.sandboxMode {
		return fmt.Sprintf("https://%s/eng/process", this.GetRemoteHost())
	} else {
		return fmt.Sprintf("https://%s/eng/process", this.GetRemoteHost())
	}
}

func (this *PayfastContext) GetRemoteHost() string {
	if this.sandboxMode {
		return fmt.Sprintf("sandbox.payfast.co.za")
	} else {
		return fmt.Sprintf("www.payfast.co.za")
	}
}
