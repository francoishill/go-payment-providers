package payfast

//Sale is just an interface that holds the basic details of a Sale required by Payfast interaction
type Sale interface {
	AmountGross() float64
	AlreadyProcessed() bool
}
