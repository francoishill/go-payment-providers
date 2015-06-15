package gopp

type IPaymentProvider interface {
	GetValidHostNames() []string
	VerifySignatureOfPostData(postDataInCorrectOrder SliceOfPostKeyValues)
}
