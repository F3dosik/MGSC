package sbox

// computeNonlinearity вычисляет минимальную нелинейность по Walsh-спектру.
func (sb *SBox) computeNonlinearity() {
	minNL := int(^uint(0) >> 1)
	for b := 1; b < 256; b++ {
		if nl := sb.Component(uint8(b)).Nonlinearity(); nl < minNL {
			minNL = nl
		}
	}
	sb.nonlinearity = &minNL
}

// computeAlgebraicDeg вычисляет максимальную алгебраическую степень через АНФ.
func (sb *SBox) computeAlgebraicDeg() {
	maxDeg := 0
	for b := 1; b < 256; b++ {
		if deg := sb.Component(uint8(b)).AlgebraicDegree(); deg > maxDeg {
			maxDeg = deg
		}
	}
	sb.algebraicDeg = &maxDeg
}

// Nonlinearity возвращает нелинейность S-блока —
// минимальную нелинейность среди всех 255 компонентных функций.
func (sb *SBox) Nonlinearity() int {
	if sb.nonlinearity == nil {
		sb.computeNonlinearity()
	}
	return *sb.nonlinearity
}

// AlgebraicDegree возвращает алгебраическую степень S-блока —
// максимальную степень среди всех 255 компонентных функций.
func (sb *SBox) AlgebraicDegree() int {
	if sb.algebraicDeg == nil {
		sb.computeAlgebraicDeg()
	}
	return *sb.algebraicDeg
}
