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

// MinMaxWalsh минимизирует максимальный Walsh-коэффициент по всем компонентным
// функциям: cost = 256 − 2·NL(S). Эквивалентно максимизации NL без порога —
// градиент есть всегда, включая область NL > 112.
func MinMaxWalsh() CostFunc {
	return func(s [256]uint8) float64 {
		return float64(256 - 2*sbox.New(s).Nonlinearity())
	}
}

// ThresholdNL штрафует только нарушение порога нелинейности:
// cost = max(0, target − NL(S)). Минимум = 0 при NL ≥ target.
func ThresholdNL(target int) CostFunc {
	return func(s [256]uint8) float64 {
		nl := sbox.New(s).Nonlinearity()
		if nl >= target {
			return 0
		}
		return float64(target - nl)
	}
}

func absInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
