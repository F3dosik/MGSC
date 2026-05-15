package generate

import (
	"math"

	"github.com/F3dosik/MGSC/sbox"
)

// ClarkJacobStepneyNL — cost по нелинейности из Clark-Jacob-Stepney 2005.
// X — целевой уровень коэффициентов Walsh (24 для NL≈116, 32 для NL≈112).
// R — степень.
func ClarkJacobStepneyNL(X, R int) CostFunc {
	return func(s [256]uint8) float64 {
		sb := sbox.New(s)
		var sum float64
		for b := 1; b < 256; b++ {
			comp := sb.Component(uint8(b))
			walsh := comp.Walsh()
			for w := 1; w < 256; w++ {
				d := math.Abs(float64(absInt32(walsh[w])) - float64(X))
				sum += math.Pow(d, float64(R))
			}
		}
		return sum
	}
}

func absInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
