package payfast

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/gogits/gogs/pkg/httplib"
	"github.com/pkg/errors"
)

var (
	errSaleAlreadyProcessed = errors.New("Sale is already processed")
)

//VerifyITNRequest verifies an Instant Transaction Notification
func VerifyITNRequest(config *Config, remoteIP, remoteUserAgent string, bodyReader io.Reader, saleProvider SaleProvider) (PaymentStatus, error) {
	verifications := &itnVerifications{
		config:          config,
		remoteIP:        remoteIP,
		remoteUserAgent: remoteUserAgent,
		bodyReader:      bodyReader,
		saleProvider:    saleProvider,
	}

	if err := verifications.verifyValidIPOfRequest(); err != nil {
		return "", errors.Wrap(err, "Failed to verify valid request")
	}

	if err := verifications.verifySignatureOfPostData(); err != nil {
		return "", errors.Wrap(err, "Failed to verify signature of post")
	}

	if err := verifications.verifySaleDataMatch(); err != nil {
		return "", errors.Wrap(err, "Failed to verify sale data match")
	}

	if err := verifications.verifyMerchantID(); err != nil {
		return "", errors.Wrap(err, "Failed to verify merchant data match")
	}

	if err := verifications.verifyFromGatewayTheySentTheRequest(); err != nil {
		return "", errors.Wrap(err, "Failed to verify that Payfast sent the request")
	}

	if verifications.verifiedRequest.ActualSale.AlreadyProcessed() {
		return "", errSaleAlreadyProcessed
	}

	return verifications.verifiedRequest.PaymentStatus, nil
}

type itnVerifications struct {
	config          *Config
	remoteIP        string
	remoteUserAgent string
	bodyReader      io.Reader
	saleProvider    SaleProvider

	verifiedRequest *verifiedRequest
}

func (i *itnVerifications) verifyValidIPOfRequest() error {
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
		if i.remoteIP == validIP {
			ipIsValid = true
			break
		}
	}
	if !ipIsValid {
		return fmt.Errorf("Invalid remote IP '%s'", i.remoteIP)
	}

	return nil
}

func (i *itnVerifications) verifySignatureOfPostData() error {
	bodyBytes, err := ioutil.ReadAll(i.bodyReader)
	if err != nil {
		return errors.Wrap(err, "failed to read body bytes")
	}
	postDataInCorrectOrder := readOrderedKeyValuePairsFromPostBody(bodyBytes)

	var (
		saleID            string
		merchantID        string
		payfastPaymentID  string
		paymentStatus     string
		buyerEmailAddress string
		firstName         string
		lastName          string
		amountGross       float64
		amountFee         float64
		amountNet         float64
	)

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
				return errors.Wrapf(err, "Failed to parse amount_gross '%s' as float", keyVal.Value)
			}
			amountGross = floatVal
		}
		if keyLowerCase == "amount_fee" {
			floatVal, err := strconv.ParseFloat(keyVal.Value, 64)
			if err != nil {
				return errors.Wrapf(err, "Failed to parse amount_fee '%s' as float", keyVal.Value)
			}
			amountFee = floatVal
		}
		if keyLowerCase == "amount_net" {
			floatVal, err := strconv.ParseFloat(keyVal.Value, 64)
			if err != nil {
				return errors.Wrapf(err, "Failed to parse amount_net '%s' as float", keyVal.Value)
			}
			amountNet = floatVal
		}
		keyValsExcludingSignature = append(keyValsExcludingSignature, keyVal)
	}

	//TODO: add to avoid golang compile errors (unused variables)
	payfastPaymentID = payfastPaymentID
	buyerEmailAddress = buyerEmailAddress
	firstName = firstName
	lastName = lastName
	amountFee = amountFee
	amountNet = amountNet

	//Before we add the passphrase
	queryWithoutPassphraseOrSignature := keyValsExcludingSignature.Combine(false)

	passPhrase := i.config.GetPassphrase()
	if passPhrase != "" {
		keyValsExcludingSignature = append(keyValsExcludingSignature, &postKeyValue{Key: "passphrase", Value: passPhrase})
	}

	tmpReceivedQueryString := keyValsExcludingSignature.Combine(false)
	expectedSignature := fmt.Sprintf("%x", md5.Sum([]byte(tmpReceivedQueryString)))
	if expectedSignature != receivedSignature {
		return errors.New("Invalid Signature")
	}

	sale, err := i.saleProvider.GetByID(saleID)
	if err != nil {
		return errors.Wrapf(err, "Unable to get Sale (from Provider) by ID '%s'", saleID)
	}

	i.verifiedRequest = &verifiedRequest{
		ParamStringForRemoteVerify: queryWithoutPassphraseOrSignature,
		MerchantID:                 merchantID,
		AmountGross:                amountGross,
		PaymentStatus:              PaymentStatus(paymentStatus),
		ActualSale:                 sale,
	}

	return nil
}

func (i *itnVerifications) verifySaleDataMatch() error {
	saleAmountGross := i.verifiedRequest.ActualSale.AmountGross()
	requestAmountGross := i.verifiedRequest.AmountGross

	diffAbs := math.Abs(saleAmountGross - requestAmountGross)

	if saleAmountGross == requestAmountGross { // shortcut, handles infinities
		return nil
	}

	//TODO: should this not be customizable
	epsilon := 0.001
	if diffAbs < math.Abs(epsilon) {
		return nil
	}

	return fmt.Errorf("Gross amount (%v) does not match the sale amount (%v)", requestAmountGross, saleAmountGross)
}

func (i *itnVerifications) verifyMerchantID() error {
	if i.config.GetMerchantID() != i.verifiedRequest.MerchantID {
		return fmt.Errorf("Merchant ID (%s) is not valid", i.verifiedRequest.MerchantID)
	}
	return nil
}

func (i *itnVerifications) verifyFromGatewayTheySentTheRequest() error {
	host := i.config.GetRemoteHost()
	url := fmt.Sprintf("https://%s/eng/query/validate", host)

	postBodyBytes := []byte(i.verifiedRequest.ParamStringForRemoteVerify)

	request := httplib.Post(url).
		Header("Host", host).
		SetUserAgent(i.remoteUserAgent).
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
