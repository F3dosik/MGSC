package generate

import "github.com/F3dosik/MGSC/sbox"

// nlFast вычисляет нелинейность биективного 8×8 S-блока напрямую через
// Walsh-Hadamard преобразование, без аллокации SBox/BoolFunc объектов.
// Используется в горячем цикле SA для минимизации давления на GC.
func nlFast(s [256]uint8) int {
	var maxW int32
	var w [256]int32
	for b := 1; b < 256; b++ {
		// TV компонентной функции f_b(x) = popcount(b & S[x]) mod 2
		// как ±1 для Walsh-преобразования
		for x := range 256 {
			v := int(b) & int(s[x])
			v ^= v >> 4
			v ^= v >> 2
			v ^= v >> 1
			if v&1 == 0 {
				w[x] = 1
			} else {
				w[x] = -1
			}
		}
		// Walsh-Hadamard butterfly (in-place)
		for h := 1; h < 256; h <<= 1 {
			for i := 0; i < 256; i += h << 1 {
				for j := i; j < i+h; j++ {
					lo, hi := w[j], w[j+h]
					w[j] = lo + hi
					w[j+h] = lo - hi
				}
			}
		}
		// ищем максимум |W_b(ω)| по ω ≠ 0
		for _, v := range w[1:] {
			if v < 0 {
				v = -v
			}
			if v > maxW {
				maxW = v
			}
		}
	}
	return 128 - int(maxW)/2
}

// ClarkJacobStepneyNL — cost по нелинейности из Clark-Jacob-Stepney 2005.
// X — целевой уровень коэффициентов Walsh.
// R — степень (поддерживается 1..4; за пределами диапазона паникует).
func ClarkJacobStepneyNL(X, R int) CostFunc {
	fX := float64(X)
	return func(s [256]uint8) float64 {
		sb := sbox.New(s)
		var sum float64
		for b := 1; b < 256; b++ {
			walsh := sb.Component(uint8(b)).Walsh()
			for w := 1; w < 256; w++ {
				d := float64(absInt32(walsh[w])) - fX
				if d < 0 {
					d = -d
				}
				switch R {
				case 1:
					sum += d
				case 2:
					sum += d * d
				case 3:
					sum += d * d * d
				case 4:
					sum += d * d * d * d
				default:
					panic("ClarkJacobStepneyNL: R must be 1..4")
				}
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
		return float64(256 - 2*nlFast(s))
	}
}

// ThresholdNL штрафует только нарушение порога нелинейности:
// cost = max(0, target − NL(S)). Минимум = 0 при NL ≥ target.
func ThresholdNL(target int) CostFunc {
	return func(s [256]uint8) float64 {
		nl := nlFast(s)
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
