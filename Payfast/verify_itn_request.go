package payfast

import (
	"crypto/md5"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/gogits/gogs/pkg/httplib"
	"github.com/pkg/errors"
)

var (
	errSaleAlreadyAuthorized = errors.New("Sale is already authorized")
)

//VerifyITNRequest verifies an Instant Transaction Notification
func VerifyITNRequest(config *Config, remoteIP, remoteUserAgent string, requestPostBody []byte, sale ITNSale) (PaymentStatus, error) {
	if err := verifyValidIPOfRequest(remoteIP); err != nil {
		return "", errors.Wrap(err, "Failed to verify valid request")
	}

	verifiedResult, err := verifySignatureOfPostData(config, requestPostBody, remoteIP)
	if err != nil {
		return "", errors.Wrap(err, "Failed to verify signature of post")
	}

	if err := verifySaleDataMatch(verifiedResult, sale); err != nil {
		return "", errors.Wrap(err, "Failed to verify sale data match")
	}

	if err := verifyMerchantID(config, verifiedResult); err != nil {
		return "", errors.Wrap(err, "Failed to verify merchant data match")
	}

	if err := verifyFromGatewayTheySentTheRequest(config, verifiedResult, remoteUserAgent); err != nil {
		return "", errors.Wrap(err, "Failed to verify that Payfast sent the request")
	}

	if sale.AlreadyAuthorized() {
		return "", errSaleAlreadyAuthorized
	}

	return verifiedResult.PaymentStatus, nil
}

func verifyValidIPOfRequest(remoteIP string) error {
	validHostNames := ValidHostNames

	validIPaddresses := []string{}
	for _, host := range validHostNames {
		hostIPs, err := net.LookupIP(host)
		if err != nil {
			return errors.Wrapf(err, "Unable to lookup host IP of host '%s'", host)
		}
		for _, ip := range hostIPs {
			validIPaddresses = append(validIPaddresses, ip.String())
		}
	}

	ipIsValid := false
	for _, validIP := range validIPaddresses {
		if remoteIP == validIP {
			ipIsValid = true
			break
		}
	}
	if !ipIsValid {
		return errors.New(fmt.Sprintf("Invalid remote IP '%s'", remoteIP))
	}

	return nil
}

func verifySignatureOfPostData(config *Config, requestPostBody []byte, remoteIP string) (*verifySignatureResult, error) {
	postDataInCorrectOrder := readOrderedKeyValuePairsFromPostBody(requestPostBody)

	var saleID string
	var merchantID string
	var payfastPaymentID string
	var paymentStatus string
	var buyerEmailAddress string
	var firstName string
	var lastName string
	var amountGross float64
	var amountFee float64
	var amountNet float64

	itnDataLogStr := ""
	for ind, keyVal := range postDataInCorrectOrder {
		if ind > 0 {
			itnDataLogStr += "&"
		}
		itnDataLogStr += fmt.Sprintf("%s=%s", keyVal.Key, keyVal.Value)
	}

	//TODO: fix logging
	// this.logger.Notice("Starting to validate ITN of remote IP '%s' and data: %s", remoteIP, itnDataLogStr)

	var receivedSignature string
	keyValsExcludingSignature := postKeyValueSlice([]*postKeyValue{})
	for _, keyVal := range postDataInCorrectOrder {
		keyLowerCase := strings.ToLower(keyVal.Key)
		if keyLowerCase == "signature" {
			receivedSignature = keyVal.Value
			continue
		}
		if keyLowerCase == "m_payment_id" {
			saleID = keyVal.Value
		}
		if keyLowerCase == "merchant_id" {
			merchantID = keyVal.Value
		}
		if keyLowerCase == "pf_payment_id" {
			payfastPaymentID = keyVal.Value
		}
		if keyLowerCase == "payment_status" {
			paymentStatus = keyVal.Value
		}
		if keyLowerCase == "email_address" {
			buyerEmailAddress = keyVal.Value
		}
		if keyLowerCase == "name_first" {
			firstName = keyVal.Value
		}
		if keyLowerCase == "name_last" {
			lastName = keyVal.Value
		}
		if keyLowerCase == "amount_gross" {
			floatVal, err := strconv.ParseFloat(keyVal.Value, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse amount_gross '%s' as float", keyVal.Value)
			}
			amountGross = floatVal
		}
		if keyLowerCase == "amount_fee" {
			floatVal, err := strconv.ParseFloat(keyVal.Value, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse amount_fee '%s' as float", keyVal.Value)
			}
			amountFee = floatVal
		}
		if keyLowerCase == "amount_net" {
			floatVal, err := strconv.ParseFloat(keyVal.Value, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse amount_net '%s' as float", keyVal.Value)
			}
			amountNet = floatVal
		}
		keyValsExcludingSignature = append(keyValsExcludingSignature, keyVal)
	}

	//Before we add the passphrase
	queryWithoutPassphraseOrSignature := keyValsExcludingSignature.Combine(false)

	passPhrase := config.GetPassphrase()
	if passPhrase != "" {
		keyValsExcludingSignature = append(keyValsExcludingSignature, &postKeyValue{Key: "passphrase", Value: passPhrase})
	}

	tmpReceivedQueryString := keyValsExcludingSignature.Combine(false)
	expectedSignature := fmt.Sprintf("%x", md5.Sum([]byte(tmpReceivedQueryString)))
	if expectedSignature != receivedSignature {
		return nil, errors.New("Invalid Signature")
	}

	//TODO: fix unused
	tmp := map[string]interface{}{
		"SaleID":                     saleID,
		"MerchantID":                 merchantID,
		"PayfastPaymentID":           payfastPaymentID,
		"PaymentStatus":              paymentStatus,
		"BuyerEmailAddress":          buyerEmailAddress,
		"FirstName":                  firstName,
		"LastName":                   lastName,
		"AmountGross":                amountGross,
		"AmountFee":                  amountFee,
		"AmountNet":                  amountNet,
		"ParamStringForRemoteVerify": queryWithoutPassphraseOrSignature, ////This string should be all key-val pairs except signature and passphrase
	}
	tmp = tmp

	result := &verifySignatureResult{
		ParamStringForRemoteVerify: queryWithoutPassphraseOrSignature,
		AmountGross:                amountGross,
		PaymentStatus:              PaymentStatus(paymentStatus),
	}

	return result, nil
}

func verifySaleDataMatch(verifyResult *verifySignatureResult, sale ITNSale) error {
	saleAmountGross := sale.AmountGross()
	requestAmountGross := verifyResult.AmountGross

	diffAbs := math.Abs(saleAmountGross - requestAmountGross)

	if saleAmountGross == requestAmountGross { // shortcut, handles infinities
		return nil
	}

	//TODO: should this not be customizable
	epsilon := 0.001
	if diffAbs < math.Abs(epsilon) {
		return nil
	}

	return errors.New("Gross amount does not match the sale amount")
}

func verifyMerchantID(config *Config, verifyResult *verifySignatureResult) error {
	if config.GetMerchantID() != verifyResult.MerchantID {
		return errors.New("Invalid merchant ID")
	}
	return nil
}

func verifyFromGatewayTheySentTheRequest(config *Config, verifiedResult *verifySignatureResult, remoteUserAgent string) error {
	host := config.GetRemoteHost()
	url := fmt.Sprintf("https://%s/eng/query/validate", host)

	postBodyBytes := []byte(verifiedResult.ParamStringForRemoteVerify)

	request := httplib.Post(url).
		Header("Host", host).
		SetUserAgent(remoteUserAgent).
		Header("Content-Type", "application/x-www-form-urlencoded").
		Header("Content-Length", fmt.Sprintf("%d", len(postBodyBytes))).
		Body(postBodyBytes)
	// request = request.SetTimeout(connectTimeout, readWriteTimeout)
	//TODO: Proxies

	responseString, err := request.String()
	if err != nil {
		return errors.Wrap(err, "Failed to get response")
	}

	if strings.ToUpper(responseString) != "VALID" {
		return errors.New("Data is invalid")
	}

	return nil
}
