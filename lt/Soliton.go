package lt

import "math"

const (
	c = 0.1
	delta = 0.5
)

type Soliton struct {
	k uint32 // number of blocks
	cdf []float64
}

func NewSoliton(nblocks uint32) (sol Soliton) {
	sol.k = nblocks
	sol.cdf = pdf2cdf(robustPDF(sol.k))
	return sol
}

func (sol Soliton) generate(rando RandGen) uint {
	newDbl := rando.nextFloat()
	return binSearch(sol.cdf, newDbl) + 1
}

// find the largest index smaller than num
// guarantee: num in [arr[0], arr[-1])
func binSearch(arr []float64, num float64) uint {
	if len(arr) == 1 {
		return 0
	} else {
		mid := uint((len(arr) - 1) / 2)  
		if arr[mid] <= num && num < arr[mid+1] {
			return mid
		} else if arr[mid] > num {
			return binSearch(arr[:mid+1], num)
		} else { // arr[mid+1] <= num 
			return binSearch(arr[mid+1:], num) + mid + 1
		}
	}
}

func intMin(a uint32, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func idealPDF(k uint32) []float64 {
	pdf := make([]float64, k)
	pdf[0] = 1.0 / float64(k)
	for d := uint32(2); d <= k; d++ {
		pdf[d-1] = 1.0 / float64(d * (d - 1))
	}
	return pdf
}

func robustPDF(k uint32) (pdf []float64) {
	s := c * math.Log(float64(k) / delta) * math.Sqrt(float64(k))
	pdf = idealPDF(k)

	zeroPt := uint32(math.Floor(float64(k) / s))

	var runSum float64 = 1.0
	var term float64
	for i := uint32(1); i < intMin(zeroPt, k); i++ {
		term = s / float64(k) * (1.0 / float64(i) )
		pdf[i-1] += term
		runSum   += term
	}

	if zeroPt < k {
		term = s / float64(k) * math.Log(s / delta)
		pdf[zeroPt - 1] += term
		runSum += term
	}

	for i, _ := range(pdf) {
		pdf[i] /= runSum
	}
	return
}

func pdf2cdf(pdf []float64) []float64 {
	cdf := make([]float64, len(pdf)+1)
	cdf[0] = 0.0
	for i, pdfVal := range(pdf) {
		cdf[i+1] = cdf[i] + pdfVal
	}
	return cdf
}
