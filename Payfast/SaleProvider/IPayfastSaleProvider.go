package SaleProvider

type IPayfastSaleProvider interface {
	GetSaleFromId(saleId string) ISale
}
