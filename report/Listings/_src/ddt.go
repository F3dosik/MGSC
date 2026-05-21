package sbox

// computeDDT вычисляет таблицу разностных распределений.
// DDT[a][b] = |{x | S(x⊕a) ⊕ S(x) = b}|, a ∈ {1..255}, b ∈ {0..255}.
// Строка a=0 не заполняется — тривиальна и не используется при вычислении
// дифференциальной равномерности.
// Сложность: O(2^(2n)) = O(65536) для n=8.
func (sb *SBox) computeDDT() {
	var ddt [256][256]uint8
	maxVal := 0

	for a := 1; a < 256; a++ {
		for x := 0; x < 256; x++ {
			b := sb.s[x^a] ^ sb.s[x]
			ddt[a][b]++
			if int(ddt[a][b]) > maxVal {
				maxVal = int(ddt[a][b])
			}
		}
	}

	sb.ddt = &ddt
	sb.diffUniformity = &maxVal
}

// DDT возвращает указатель на таблицу разностных распределений.
// Результат кэшируется.
func (sb *SBox) DDT() *[256][256]uint8 {
	if sb.ddt == nil {
		sb.computeDDT()
	}
	return sb.ddt
}

// DiffUniformity возвращает дифференциальную равномерность —
// максимальное значение DDT[a][b] для a ≠ 0.
func (sb *SBox) DiffUniformity() int {
	if sb.diffUniformity == nil {
		sb.computeDDT()
	}
	return *sb.diffUniformity
}
