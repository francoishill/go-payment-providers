package payfast

import (
	"strings"
)

//PaymentStatus is just a convenience wrapper around the string value
type PaymentStatus string

//IsComplete checks if it is complete
func (p PaymentStatus) IsComplete() bool { return strings.EqualFold(string(p), "COMPLETE") }

//IsFailed checks if it is failed
func (p PaymentStatus) IsFailed() bool { return strings.EqualFold(string(p), "FAILED") }

//IsPending checks if it is pending
func (p PaymentStatus) IsPending() bool { return strings.EqualFold(string(p), "PENDING") }
