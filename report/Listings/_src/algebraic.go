package sbox

import "math"

// AlgebraicImmunity возвращает минимальную AI по всем компонентным функциям
func (sb *SBox) AlgebraicImmunity() int {
	minAI := math.MaxInt32
	for b := 1; b <= 255; b++ {
		ai := sb.Component(uint8(b)).AlgebraicImmunity()
		if ai < minAI {
			minAI = ai
		}
		// ранний выход — хуже 1 не бывает
		if minAI == 1 {
			return 1
		}
	}
	return minAI
}
