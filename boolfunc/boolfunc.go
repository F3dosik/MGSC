// Package boolfunc реализует булевы функции от произвольного числа переменных.
// Основное применение — анализ компонентных функций S-блоков.
package boolfunc

import (
	"fmt"
	"math/bits"
	"math/rand/v2"
)

// BoolFunc представляет булеву функцию от n переменных.
//
// Вектор истинности tv имеет длину 2^n, где tv[i] — значение функции
// на входе i. Все тяжёлые вычисления кэшируются.
type BoolFunc struct {
	n   int      // количество переменных
	tv  []uint8  // вектор истинности
	tvp []uint64 // TV упакованный побитово

	weight  *int    // вес Хэмминга
	anf     []uint8 // АНФ-коэффициенты
	degree  *int    // алгебраическая степень
	walsh   []int32 // спектр Уолша–Адамара
	nlin    *int    // нелинейность
	bestAff []uint8 // TV лучшего аффинного приближения
	ai      *int    // алгербраическая иммунность
}

// ── Конструкторы ─────────────────────────────────────────────────────────────
func New(s string) (*BoolFunc, error) {
	size := len(s)
	if err := validateSize(size); err != nil {
		return nil, err
	}
	tv := make([]uint8, size)
	for i, ch := range s {
		switch ch {
		case '0':
		case '1':
			tv[i] = 1
		default:
			return nil, fmt.Errorf("boolfunc: недопустимый символ %q на позиции %d", ch, i)
		}
	}
	return newUnchecked(tv), nil
}

// FromSlice создаёт BoolFunc из слайса.
func FromSlice(tv []uint8) (*BoolFunc, error) {
	if err := validateSize(len(tv)); err != nil {
		return nil, err
	}
	for i, v := range tv {
		if v > 1 {
			return nil, fmt.Errorf("boolfunc: значение %d на позиции %d не является булевым", v, i)
		}
	}
	cp := make([]uint8, len(tv))
	copy(cp, tv)
	return newUnchecked(cp), nil
}

// Random возвращает случайную булеву функцию от n переменных.
func Random(n int) (*BoolFunc, error) {
	if err := validateN(n); err != nil {
		return nil, err
	}
	size := 1 << n
	tv := make([]uint8, size)
	for i := range tv {
		tv[i] = uint8(rand.IntN(2))
	}
	return newUnchecked(tv), nil
}

// RandomBalanced возвращает случайную сбалансированную булеву функцию от n переменных
// (ровно 2^(n-1) единиц в векторе истинности).
func RandomBalanced(n int) (*BoolFunc, error) {
	if err := validateN(n); err != nil {
		return nil, err
	}
	size := 1 << n
	tv := make([]uint8, size)
	// Ставим ровно half единиц в начало, затем перемешиваем
	for i := range size / 2 {
		tv[i] = 1
	}
	rand.Shuffle(size, func(i, j int) {
		tv[i], tv[j] = tv[j], tv[i]
	})
	return newUnchecked(tv), nil
}

// newUnchecked — внутренний конструктор без валидации.
// Принимает уже проверенный и скопированный слайс.
func newUnchecked(tv []uint8) *BoolFunc {
	return &BoolFunc{
		n:   bits.Len(uint(len(tv))) - 1, // log2(len(tv))
		tv:  tv,
		tvp: packBits(tv),
	}
}

// ── Геттеры ───────────────────────────────────────────────────────────

// N возвращает количество переменных.
func (f *BoolFunc) N() int { return f.n }

// Size возвращает размер вектора истинности.
func (f *BoolFunc) Size() int { return 1 << f.n }

// TV возвращает копию вектора истинности.
func (f *BoolFunc) TV() []uint8 {
	cp := make([]uint8, len(f.tv))
	copy(cp, f.tv)
	return cp
}

// At возвращает значение функции на входе x.
func (f *BoolFunc) At(x int) uint8 {
	return f.tv[x]
}

// IsBalanced возвращает true, если функция сбалансирована.
func (f *BoolFunc) IsBalanced() bool {
	return f.Weight() == 1<<(f.n-1)
}

// Weight возвращает вес Хэмминга — количество единиц в TV.
func (f *BoolFunc) Weight() int {
	if f.weight == nil {
		w := 0
		for _, word := range f.tvp {
			w += bits.OnesCount64(word)
		}
		f.weight = &w
	}
	return *f.weight
}

// ── Вспомогательные функции ───────────────────────────────────────────────────

// packBits упаковывает []uint8 ∈ {0,1} в []uint64.
// Бит i записывается в слово i/64, позиция i%64.
// Порядок внутри слова — little-endian
func packBits(tv []uint8) []uint64 {
	size := len(tv)
	words := (size + 63) / 64 // ceil(a/b) = (a + b - 1) / b
	packed := make([]uint64, words)
	for i, v := range tv {
		if v == 1 {
			packed[i/64] |= 1 << (i % 64)
		}
	}
	return packed
}

// validateSize проверяет что size > 0 и является степенью двойки.
func validateSize(size int) error {
	if size <= 0 || (size&(size-1)) != 0 {
		return fmt.Errorf("boolfunc: длина вектора истинности должна быть степенью двойки, получено %d", size)
	}
	return nil
}

// validateN проверяет что n > 0.
func validateN(n int) error {
	if n <= 0 {
		return fmt.Errorf("boolfunc: n должно быть > 0, получено %d", n)
	}
	return nil
}
