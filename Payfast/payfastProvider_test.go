package Payfast

import (
	"crypto/md5"
	"fmt"
	. "github.com/francoishill/go-payment-providers"
	. "github.com/francoishill/go-payment-providers/Payfast/PayfastContext"
	. "github.com/francoishill/go-payment-providers/Payfast/SaleProvider"
	"testing"
)

const (
	tmpRemoteUserAgent = "MobileDummyBrowser1.0"
	tmpRemoteIp        = "100.200.100.200"
	saleId             = "123456"
	tmpMerchantId      = "AABBCC123"
	tmpMerchantKey     = "ZZZXXX"
	payfastPaymentId   = "PFID123"
	paymentStatus      = ""
	buyerEmailAddress  = ""
	firstName          = ""
	lastName           = ""
	amountGross        = "200"
	amountFee          = "2.5"
	amountNet          = "197.5"
)

var logger *tmpLogger = &tmpLogger{}
var validSaleProvider *tmpValidSaleProvider = &tmpValidSaleProvider{}
var invalidSaleProvider *tmpInvalidSaleProvider = &tmpInvalidSaleProvider{}

type tmpLogger struct{}

func (this *tmpLogger) Notice(format string, v ...interface{}) {
	//Do nothing for now
	//fmt.Println("[NOTICE]" + fmt.Sprintf(format, v...))
}

type tmpValidSaleProvider struct{}

func (this *tmpValidSaleProvider) GetSaleFromId(saleId string) ISale { return &tmpValidSale{} }

type tmpInvalidSaleProvider struct{}

func (this *tmpInvalidSaleProvider) GetSaleFromId(saleId string) ISale { return &tmpInvalidSale{} }

type tmpValidSale struct{}

func (this *tmpValidSale) GetId() string              { return "valid-sale-id" }
func (this *tmpValidSale) GetItemName() string        { return "VALID sale name" }
func (this *tmpValidSale) GetItemDescription() string { return "VALID sale description" }
func (this *tmpValidSale) GetAmountGross() float32    { return ValueString(amountGross).ToFloat32() }
func (this *tmpValidSale) IsAuthorized() bool         { return false }
func (this *tmpValidSale) UpdateSaleAsAuthorizedCompleted(payfastPaymentId string, amountFee, amountNet float32) {
}
func (this *tmpValidSale) UpdateSaleAsFailed(payfastPaymentId string)  {}
func (this *tmpValidSale) UpdateSaleAsPending(payfastPaymentId string) {}

type tmpInvalidSale struct{}

func (this *tmpInvalidSale) GetId() string              { return "invalid-sale-id" }
func (this *tmpInvalidSale) GetItemName() string        { return "INVALID sale name" }
func (this *tmpInvalidSale) GetItemDescription() string { return "INVALID sale description" }
func (this *tmpInvalidSale) GetAmountGross() float32 {
	return ValueString(amountGross).ToFloat32() + 0.1
}
func (this *tmpInvalidSale) IsAuthorized() bool { return false }
func (this *tmpInvalidSale) UpdateSaleAsAuthorizedCompleted(payfastPaymentId string, amountFee, amountNet float32) {
}
func (this *tmpInvalidSale) UpdateSaleAsFailed(payfastPaymentId string)  {}
func (this *tmpInvalidSale) UpdateSaleAsPending(payfastPaymentId string) {}

func createPayfastContext(sandboxMode bool, merchantId, merchantKey, passphrase string) *PayfastContext {
	return_url := "http://no_url__testing_mode"
	cancel_url := "http://no_url__testing_mode"
	notify_url := "http://no_url__testing_mode"
	return CreatePayfastContext(sandboxMode, merchantId, merchantKey, passphrase, return_url, cancel_url, notify_url)
}

func TestVerifyValidIPOfRequest(t *testing.T) {
	paymentProvider := CreatePayfastProvider(logger, validSaleProvider, createPayfastContext(true, tmpMerchantId, tmpMerchantKey, "tmp-passphrase1"))

	remoteValidIP := "41.74.179.210"
	remoteInvalidIP := "41.74.179.230"
	DoTestVerifyValidIPOfRequest(t, paymentProvider, remoteValidIP, remoteInvalidIP, "Invalid remote IP '"+remoteInvalidIP+"'")
}

func getRequestPostBodyKeyVals(signatureToInclude, passPhraseToInclude string) SliceOfPostKeyValue {
	keyVals := SliceOfPostKeyValue([]*PostKeyValue{})
	keyVals = append(keyVals, &PostKeyValue{Key: "m_payment_id", Value: ValueString(fmt.Sprintf("%s", saleId))})
	keyVals = append(keyVals, &PostKeyValue{Key: "merchant_id", Value: ValueString(fmt.Sprintf("%s", tmpMerchantId))})
	keyVals = append(keyVals, &PostKeyValue{Key: "pf_payment_id", Value: ValueString(fmt.Sprintf("%s", payfastPaymentId))})
	keyVals = append(keyVals, &PostKeyValue{Key: "payment_status", Value: ValueString(fmt.Sprintf("%s", paymentStatus))})
	keyVals = append(keyVals, &PostKeyValue{Key: "email_address", Value: ValueString(fmt.Sprintf("%s", buyerEmailAddress))})
	keyVals = append(keyVals, &PostKeyValue{Key: "name_first", Value: ValueString(fmt.Sprintf("%s", firstName))})
	keyVals = append(keyVals, &PostKeyValue{Key: "name_last", Value: ValueString(fmt.Sprintf("%s", lastName))})
	keyVals = append(keyVals, &PostKeyValue{Key: "amount_gross", Value: ValueString(fmt.Sprintf("%s", amountGross))})
	keyVals = append(keyVals, &PostKeyValue{Key: "amount_fee", Value: ValueString(fmt.Sprintf("%s", amountFee))})
	keyVals = append(keyVals, &PostKeyValue{Key: "amount_net", Value: ValueString(fmt.Sprintf("%s", amountNet))})
	if signatureToInclude != "" {
		keyVals = append(keyVals, &PostKeyValue{Key: "signature", Value: ValueString(fmt.Sprintf("%s", signatureToInclude))})
	}
	if passPhraseToInclude != "" {
		keyVals = append(keyVals, &PostKeyValue{Key: "passphrase", Value: ValueString(fmt.Sprintf("%s", passPhraseToInclude))})
	}
	return keyVals
}

func localTestVerifySignatureOfPostData(t *testing.T, passPhraseLeaveBlankForNone string) {
	paymentProvider := CreatePayfastProvider(logger, validSaleProvider, createPayfastContext(true, tmpMerchantId, tmpMerchantKey, passPhraseLeaveBlankForNone))

	queryWithPassphraseButNotSignature := getRequestPostBodyKeyVals("", passPhraseLeaveBlankForNone).CombineIntoSingleString(false)
	validSignature := fmt.Sprintf("%x", md5.Sum([]byte(queryWithPassphraseButNotSignature)))
	queryWithoutPassphraseButWithSignature := getRequestPostBodyKeyVals(validSignature, "").CombineIntoSingleString(false)

	validRequestPostBody := []byte(queryWithoutPassphraseButWithSignature)
	invalidRequestPostBody := []byte(queryWithoutPassphraseButWithSignature + "A")

	DoTestVerifySignatureOfPostData(t, paymentProvider, tmpRemoteIp, validRequestPostBody, invalidRequestPostBody, "Invalid Signature")
}

func TestVerifySignatureOfPostData(t *testing.T) {
	passPhrase := "passphrase"
	//Test with passphrase
	localTestVerifySignatureOfPostData(t, passPhrase)

	//Test without passphrase
	localTestVerifySignatureOfPostData(t, "")
}

func getPrePopulatedPayfastPaymentProvider(merchantId string, saleProvider IPayfastSaleProvider) IPaymentProvider {
	paymentProvider := CreatePayfastProvider(logger, saleProvider, createPayfastContext(true, merchantId, tmpMerchantKey, "tmp-passphrase2"))
	queryWithPassphraseButNotSignature := getRequestPostBodyKeyVals("", "tmp-passphrase2").CombineIntoSingleString(false)
	validSignature := fmt.Sprintf("%x", md5.Sum([]byte(queryWithPassphraseButNotSignature)))
	paymentProvider.VerifySignatureOfPostData(getRequestPostBodyKeyVals(validSignature, ""), tmpRemoteIp)
	return paymentProvider
}

func TestVerifySaleDataMatch(t *testing.T) {
	validPaymentProvider := getPrePopulatedPayfastPaymentProvider(tmpMerchantId, validSaleProvider)
	invalidPaymentProvider := getPrePopulatedPayfastPaymentProvider(tmpMerchantId, invalidSaleProvider)
	DoTestVerifySaleDataMatch(t, validPaymentProvider, invalidPaymentProvider, "Gross amount does not match the sale amount")
}

func TestVerifyMerchantDataMatch(t *testing.T) {
	validPaymentProvider := getPrePopulatedPayfastPaymentProvider(tmpMerchantId, validSaleProvider)
	invalidPaymentProvider := getPrePopulatedPayfastPaymentProvider("123abc", validSaleProvider) //Incorrect merchant id
	DoTestVerifyMerchantDataMatch(t, validPaymentProvider, invalidPaymentProvider, "Invalid merchant ID")
}

func TestVerifyFromGatewayTheySentTheRequest(t *testing.T) {
	//TODO: Figure out the best way to test an ITN transaction coming in (to mock it)
	/*validPaymentProvider := getPrePopulatedPayfastPaymentProvider(tmpMerchantId, validSaleProvider)
	invalidPaymentProvider := getPrePopulatedPayfastPaymentProvider(tmpMerchantId, validSaleProvider)
	DoTestVerifyFromGatewayTheySentTheRequest(t, tmpRemoteUserAgent, validPaymentProvider, invalidPaymentProvider, "Data is invalid")*/
}
