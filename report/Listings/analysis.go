package sbox

// computeMetrics вычисляет агрегированные метрики за один проход
// по всем 255 компонентным функциям.
// Заполняет: nonlinearity, algebraicDeg.
func (sb *SBox) computeMetrics() {
	minNL := int(^uint(0) >> 1) // MaxInt
	maxDeg := 0

	for b := 1; b < 256; b++ {
		f := sb.Component(uint8(b))

		if nl := f.Nonlinearity(); nl < minNL {
			minNL = nl
		}
		if deg := f.AlgebraicDegree(); deg > maxDeg {
			maxDeg = deg
		}
	}

	sb.nonlinearity = &minNL
	sb.algebraicDeg = &maxDeg
}

// Nonlinearity возвращает нелинейность S-блока —
// минимальную нелинейность среди всех 255 компонентных функций.
func (sb *SBox) Nonlinearity() int {
	if sb.nonlinearity == nil {
		sb.computeMetrics()
	}
	return *sb.nonlinearity
}

// AlgebraicDegree возвращает алгебраическую степень S-блока —
// максимальную степень среди всех 255 компонентных функций.
func (sb *SBox) AlgebraicDegree() int {
	if sb.algebraicDeg == nil {
		sb.computeMetrics()
	}
	return *sb.algebraicDeg
}
