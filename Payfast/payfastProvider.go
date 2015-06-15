package Payfast

import (
	"crypto/md5"
	"fmt"
	"github.com/francoishill/go-payment-providers"
	"strings"
)

type iLocalLogger interface {
	Notice(format string, v ...interface{})
}

type payfastProvider struct {
	logger        iLocalLogger
	remoteIP      string
	passPhrase    string
	extractedData *extractedPostITNData
}

func CreatePayfastProvider(logger iLocalLogger, remoteIP, passPhrase string) gopp.IPaymentProvider {
	return &payfastProvider{
		logger:     logger,
		remoteIP:   remoteIP,
		passPhrase: passPhrase,
	}
}

func (this *payfastProvider) GetValidHostNames() []string {
	return []string{
		"www.payfast.co.za",
		"w1w.payfast.co.za",
		"w2w.payfast.co.za",
		"sandbox.payfast.co.za",
	}
}

func (this *payfastProvider) VerifySignatureOfPostData(postDataInCorrectOrder gopp.SliceOfPostKeyValues) {
	var saleId int64
	var merchantId string
	var payfastPaymentId string
	var paymentStatus string
	var buyerEmailAddress string
	var firstName string
	var lastName string
	var amountGross float32
	var amountFee float32
	var amountNet float32

	itnDataLogStr := ""
	for ind, keyVal := range postDataInCorrectOrder {
		if ind > 0 {
			itnDataLogStr += "&"
		}
		itnDataLogStr += fmt.Sprintf("%s=%s", keyVal.Key, keyVal.Value)
	}
	this.logger.Notice("Starting to validate ITN of remote IP '%s' and data: %s", this.remoteIP, itnDataLogStr)

	var receivedSignature gopp.ValueString = ""
	keyValsExcludingSignature := gopp.SliceOfPostKeyValues([]*gopp.PostKeyValues{})
	for _, keyVal := range postDataInCorrectOrder {
		keyLowerCase := strings.ToLower(keyVal.Key)
		if keyLowerCase == "signature" {
			receivedSignature = keyVal.Value
			continue
		}
		if keyLowerCase == "m_payment_id" {
			saleId = keyVal.Value.ToInt64()
		}
		if keyLowerCase == "merchant_id" {
			merchantId = keyVal.Value.ToString()
		}
		if keyLowerCase == "pf_payment_id" {
			payfastPaymentId = keyVal.Value.ToString()
		}
		if keyLowerCase == "payment_status" {
			paymentStatus = keyVal.Value.ToString()
		}
		if keyLowerCase == "email_address" {
			buyerEmailAddress = keyVal.Value.ToString()
		}
		if keyLowerCase == "name_first" {
			firstName = keyVal.Value.ToString()
		}
		if keyLowerCase == "name_last" {
			lastName = keyVal.Value.ToString()
		}
		if keyLowerCase == "amount_gross" {
			amountGross = keyVal.Value.ToFloat32()
		}
		if keyLowerCase == "amount_fee" {
			amountFee = keyVal.Value.ToFloat32()
		}
		if keyLowerCase == "amount_net" {
			amountNet = keyVal.Value.ToFloat32()
		}
		keyValsExcludingSignature = append(keyValsExcludingSignature, keyVal)
	}

	//Before we add the passphrase
	queryWithoutPassphraseOrSignature := keyValsExcludingSignature.CombineIntoSingleString(false)

	passPhrase := this.passPhrase
	if passPhrase != "" {
		keyValsExcludingSignature = append(keyValsExcludingSignature, &gopp.PostKeyValues{Key: "passphrase", Value: gopp.ValueString(passPhrase)})
	}

	tmpReceivedQueryString := keyValsExcludingSignature.CombineIntoSingleString(false)
	expectedSignature := gopp.ValueString(fmt.Sprintf("%x", md5.Sum([]byte(tmpReceivedQueryString))))
	if expectedSignature != receivedSignature {
		panic("Invalid Signature")
	}

	this.extractedData = &extractedPostITNData{
		SaleId:                     saleId,
		MerchantId:                 merchantId,
		PayfastPaymentId:           payfastPaymentId,
		PaymentStatus:              paymentStatus,
		BuyerEmailAddress:          buyerEmailAddress,
		FirstName:                  firstName,
		LastName:                   lastName,
		AmountGross:                amountGross,
		AmountFee:                  amountFee,
		AmountNet:                  amountNet,
		ParamStringForRemoteVerify: queryWithoutPassphraseOrSignature, ////This string should be all key-val pairs except signature and passphrase
	}
}
