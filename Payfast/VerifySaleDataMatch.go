package Payfast

import (
	"math"
)

func (this *payfastProvider) VerifySaleDataMatch() {
	sale := this.saleProvider.GetSaleFromId(this.extractedData.SaleId)

	saleActualGrossAmount := sale.GetAmountGross()
	incomingRequestGrossAmount := this.extractedData.AmountGross

	if !checkFloat32sAreNearlyEqual_Using_AbsoluteError(saleActualGrossAmount, incomingRequestGrossAmount, 0.001) {
		panic("Gross amount does not match the sale amount")
	}

	this.sale = sale
}

func checkFloat32sAreNearlyEqual_Using_RelativeError(a, b, epsilon float32) bool {
	absA := math.Abs(float64(a))
	absB := math.Abs(float64(b))
	diffAbs := math.Abs(float64(a - b))

	if a == b { // shortcut, handles infinities
		return true
	} else if a == 0 || b == 0 || diffAbs < math.SmallestNonzeroFloat32 {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return diffAbs < (float64(epsilon) * float64(math.SmallestNonzeroFloat32))
	} else { // use relative error
		return diffAbs/(absA+absB) < float64(epsilon)
	}
}

func checkFloat32sAreNearlyEqual_Using_AbsoluteError(a, b, epsilon float32) bool {
	diffAbs := math.Abs(float64(a - b))

	if a == b { // shortcut, handles infinities
		return true
	} else {
		return diffAbs < math.Abs(float64(epsilon))
	}
}
