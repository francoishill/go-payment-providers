package payfast

import (
	"fmt"
	"strings"
)

//Config holds the configuration for Payfast payments
type Config struct {
	sandboxMode bool
	merchantID  string
	merchantKey string
	passphrase  string
	returnURL   string
	cancelURL   string
	notifyURL   string
}

//NewConfig creates the new config instance
func NewConfig(sandboxMode bool, merchantID, merchantKey, passphrase, returnURL, cancelURL, notifyURL string) *Config {
	return &Config{
		sandboxMode: sandboxMode,
		merchantID:  merchantID,
		merchantKey: merchantKey,
		passphrase:  passphrase,
		returnURL:   strings.TrimRight(returnURL, "/"),
		cancelURL:   strings.TrimRight(cancelURL, "/"),
		notifyURL:   notifyURL,
	}
}

//IsSandboxmode returns the value of the corresponding private field
func (c *Config) IsSandboxmode() bool { return c.sandboxMode }

//GetMerchantID returns the value of the corresponding private field
func (c *Config) GetMerchantID() string { return c.merchantID }

//GetMerchantKey returns the value of the corresponding private field
func (c *Config) GetMerchantKey() string { return c.merchantKey }

//GetPassphrase returns the value of the corresponding private field
func (c *Config) GetPassphrase() string { return c.passphrase }

//GetReturnURL returns a new URL string for the given SaleID
func (c *Config) GetReturnURL(saleID string) string { return fmt.Sprintf("%s/%s", c.returnURL, saleID) }

//GetCancelURL  returns a new URL string for the given SaleID
func (c *Config) GetCancelURL(saleID string) string { return fmt.Sprintf("%s/%s", c.cancelURL, saleID) }

//GetNotifyURL returns the value of the corresponding private field
func (c *Config) GetNotifyURL() string { return c.notifyURL }

//GetAuthServiceURL returns an URL string based on whether sandbox mode is enabled
func (c *Config) GetAuthServiceURL() string {
	if c.sandboxMode {
		return fmt.Sprintf("https://%s/eng/process", c.GetRemoteHost())
	}

	return fmt.Sprintf("https://%s/eng/process", c.GetRemoteHost())
}

//GetRemoteHost gets the remote host string based on whether sandbox mode is enabled
func (c *Config) GetRemoteHost() string {
	if c.sandboxMode {
		return fmt.Sprintf("sandbox.payfast.co.za")
	}

	return fmt.Sprintf("www.payfast.co.za")
}
