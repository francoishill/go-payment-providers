# go-payment-providers
Payment providers for golang

# Travis build
[![Build Status](https://travis-ci.org/francoishill/go-payment-providers.svg?branch=master)](https://travis-ci.org/francoishill/go-payment-providers)

## Goals

- Have a generic platform to handle payment provider ITN (instant transaction notification), including:
  + Verify data (using signatures)
  + Verify IP address
  + Verify sale data
  + Verify merchant data
  + Ask provider if they sent the request
  + Check if the sale was not already processed
- Have support for multiple payment providers, including:
  + [ ] Paypal
  + [x] Payfast
  + [ ] 2Checkout
  + More...
- Other additional "helper" packages for simplifying payment workflows

