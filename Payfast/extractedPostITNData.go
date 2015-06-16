package Payfast

type extractedPostITNData struct {
	SaleId                     string  //m_payment_id
	MerchantId                 string  //merchant_id
	PayfastPaymentId           string  //pf_payment_id
	PaymentStatus              string  //payment_status
	BuyerEmailAddress          string  //email_address
	FirstName                  string  //name_first
	LastName                   string  //name_last
	AmountGross                float32 //amount_gross
	AmountFee                  float32 //amount_fee
	AmountNet                  float32 //amount_net
	ParamStringForRemoteVerify string  //This string should be all key-val pairs except signature and passphrase
}
