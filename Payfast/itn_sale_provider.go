package payfast

//SaleProvider is a simple provider to get the Sale by its ID
type SaleProvider interface {
	GetByID(id string) (Sale, error)
}
