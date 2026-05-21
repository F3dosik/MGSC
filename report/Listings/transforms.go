package boolfunc

import "slices"

// moebiusInPlace выполняет преобразование Мёбиуса над GF(2) на месте.
// После вызова f содержит АНФ-коэффициенты.
// Сложность: O(n · 2^n), n = log2(len(f)).
func moebiusInPlace(f []uint8, n int) {
	for i := range n {
		step := 1 << i
		blockSize := step << 1
		for blockStart := 0; blockStart < 1<<n; blockStart += blockSize {
			for j := range step {
				a := blockStart + step + j
				f[a] ^= f[a-step]
			}
		}
	}
}

// ANF возвращает копию вектора АНФ-коэффициентов булевой функции.
// Результат кэшируется — повторные вызовы не пересчитывают преобразование.
func (f *BoolFunc) ANF() []uint8 {
	if f.anf == nil {
		anf := make([]uint8, len(f.tv))
		copy(anf, f.tv)
		moebiusInPlace(anf, f.n)
		f.anf = anf
	}
	cp := make([]uint8, len(f.anf))
	copy(cp, f.anf)
	return cp
}

// fwhtInPlace выполняет быстрое преобразование Уолша–Адамара на месте.
// Входной массив sv — знаковый вектор (0→+1, 1→-1), после вызова содержит
// коэффициенты спектра Уолша.
// Сложность: O(n · 2^n), n = log2(len(sv)).
func fwhtInPlace(sv []int32) {
	n := len(sv)
	for h := 1; h < n; h <<= 1 {
		for i := 0; i < n; i += h << 1 {
			for j := i; j < i+h; j++ {
				x := sv[j]
				y := sv[j+h]
				sv[j] = x + y
				sv[j+h] = x - y
			}
		}
	}
}

// computeWalsh вычисляет спектр Уолша и сохраняет в кэш.
// Знаковый вектор — локальная переменная, не хранится как поле.
func (f *BoolFunc) computeWalsh() {
	sv := make([]int32, len(f.tv))
	for i, v := range f.tv {
		sv[i] = 1 - 2*int32(v)
	}
	fwhtInPlace(sv)
	f.walsh = sv
}

// Walsh возвращает копию спектра Уолша–Адамара.
// Результат кэшируется — повторные вызовы не пересчитывают преобразование.
//
// Walsh[a] = sum_{x} (-1)^{f(x) + <a,x>}, где <a,x> — скалярное произведение
// над GF(2). Для сбалансированной функции Walsh[0] == 0.
func (f *BoolFunc) Walsh() []int32 {
	if f.walsh == nil {
		f.computeWalsh()
	}
	return slices.Clone(f.walsh)
}
