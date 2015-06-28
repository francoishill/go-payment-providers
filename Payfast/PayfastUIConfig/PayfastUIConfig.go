package PayfastUIConfig

import (
	"crypto/md5"
	"fmt"
	"strings"

	. "github.com/francoishill/go-payment-providers/Payfast/PayfastContext"
	. "github.com/francoishill/go-payment-providers/Payfast/SaleProvider"
)

type PayfastUIConfig struct {
	PostUrl   string
	FieldList []*keyValPair
}

func CreatePayfastUIConfig(payfastContext *PayfastContext, sale ISale, authContext *authorizedContext) *PayfastUIConfig {
	config := &PayfastUIConfig{
		PostUrl: payfastContext.GetAuthServiceUrl(),
		FieldList: []*keyValPair{
			&keyValPair{Key: "merchant_id", Value: payfastContext.GetMerchantId()},
			&keyValPair{Key: "merchant_key", Value: payfastContext.GetMerchantKey()},
			&keyValPair{Key: "return_url", Value: payfastContext.GetReturnUrl(sale.GetId())},
			&keyValPair{Key: "cancel_url", Value: payfastContext.GetCancelUrl(sale.GetId())},
			&keyValPair{Key: "notify_url", Value: payfastContext.GetNotifyUrl()},
			&keyValPair{Key: "name_first", Value: authContext.FirstName},
			&keyValPair{Key: "name_last", Value: authContext.LastName},
			&keyValPair{Key: "email_address", Value: authContext.Email},
			&keyValPair{Key: "m_payment_id", Value: fmt.Sprintf("%s", sale.GetId())},
			&keyValPair{Key: "amount", Value: fmt.Sprintf("%.2f", sale.GetAmountGross())},
			&keyValPair{Key: "item_name", Value: sale.GetItemName()},
			&keyValPair{Key: "item_description", Value: sale.GetItemDescription()},
			/*&keyValPair{Key: "custom_int1", Value: ""}, //We can have custom_int1..5 and custom_str1..5
			  &keyValPair{Key: "custom_str1", Value: ""},*/
			//email_confirmation, true or false wheter to send an email to the MERCHANT, the buyer will always receive an email
			//confirmation_address, send the confirmation email to the BUYER
		},
	}

	config.sanitizeAndEscapeFieldList()

	passPhrase := payfastContext.GetPassphrase()
	config.RefreshSignature(passPhrase)

	return config
}

const cALLOWED_CHARS_IN_FIELD_VALUE = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_:/@. "

func (this *PayfastUIConfig) sanitizeAndEscapeFieldList() {
	for _, field := range this.FieldList {
		if strings.EqualFold(field.Key, "merchant_id") ||
			strings.EqualFold(field.Key, "merchant_key") ||
			strings.EqualFold(field.Key, "return_url") ||
			strings.EqualFold(field.Key, "cancel_url") ||
			strings.EqualFold(field.Key, "notify_url") ||
			strings.EqualFold(field.Key, "email_address") ||
			strings.EqualFold(field.Key, "m_payment_id") ||
			strings.EqualFold(field.Key, "amount") {
			continue
		}
		origVal := field.Value
		cleanedVal := ""
		for _, char := range origVal {
			charAsStr := string(char)
			if strings.Contains(cALLOWED_CHARS_IN_FIELD_VALUE, charAsStr) {
				cleanedVal += charAsStr
			}
		}
		field.Value = strings.Trim(cleanedVal, " ")
	}
}

func (this *PayfastUIConfig) RefreshSignature(passPhrase string) {
	sliceOfKeyValPairs := []string{}
	for _, keyValPair := range this.FieldList {
		sliceOfKeyValPairs = keyValPair.appendKeyValToStringSlice(sliceOfKeyValPairs)
	}
	//Do not append the passphrase to FieldList too
	if strings.Trim(passPhrase, " ") != "" {
		sliceOfKeyValPairs = appendKeyValToStringSlice(sliceOfKeyValPairs, "passphrase", passPhrase)
	}

	getStr := strings.Join(sliceOfKeyValPairs, "&")

	signatureField := &keyValPair{Key: "signature", Value: fmt.Sprintf("%x", md5.Sum([]byte(getStr)))}
	//Must append the signature to the FieldList
	this.FieldList = append(this.FieldList, signatureField)
}
