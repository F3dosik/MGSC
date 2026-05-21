// Package sbox реализует анализ криптографических свойств S-блоков 8×8.
package sbox

import (
	"fmt"

	"github.com/F3dosik/MGSC/boolfunc"
)

// SBox представляет S-блок 8×8 — биективную подстановку над {0,...,255}.
type SBox struct {
	s [256]uint8 // таблица подстановки

	// ── Ленивый кэш ─────────────────────────────────────────────────────────
	components        [255]*boolfunc.BoolFunc // компонентные функции, индекс f_b = b-1, b ∈ {1..255}
	ddt               *[256][256]uint8        // таблица разностных распределений
	nonlinearity      *int                    // минимальная нелинейность по всем компонентным функциям
	algebraicDeg      *int                    // максимальная алгебраическая степень
	diffUniformity    *int                    // максимум DDT — дифференциальная равномерность
	algebraicImmunity *int                    // минимум AI по всем компонентным функциям f_b

}

// ── Конструкторы ─────────────────────────────────────────────────────────────

// New создаёт SBox из массива подстановки.
// Не проверяет биективность — используй IsBijective() отдельно.
func New(s [256]uint8) *SBox {
	return &SBox{s: s}
}

// FromSlice создаёт SBox из слайса длиной ровно 256.
func FromSlice(s []uint8) (*SBox, error) {
	if len(s) != 256 {
		return nil, fmt.Errorf("sbox: длина таблицы должна быть 256, получено %d", len(s))
	}
	var arr [256]uint8
	copy(arr[:], s)
	return New(arr), nil
}

// ── Базовые геттеры ───────────────────────────────────────────────────────────

// At возвращает значение S-блока на входе x.
func (sb *SBox) At(x uint8) uint8 {
	return sb.s[x]
}

// Table возвращает копию таблицы подстановки.
func (sb *SBox) Table() [256]uint8 {
	return sb.s
}

// ── Базовые проверки ──────────────────────────────────────────────────────────

// IsBijective возвращает true если S-блок является биекцией (перестановкой).
// Проверяет что каждое значение {0,...,255} встречается ровно один раз.
func (sb *SBox) IsBijective() bool {
	var seen [256]bool
	for _, v := range sb.s {
		if seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}

// FixedPoints возвращает список неподвижных точек — входов где S(x) == x.
func (sb *SBox) FixedPoints() []uint8 {
	var points []uint8
	for x, v := range sb.s {
		if v == uint8(x) {
			points = append(points, uint8(x))
		}
	}
	return points
}

// OppositeFixedPoints возвращает список антиподвижных точек — входов где S(x) == ^x.
func (sb *SBox) OppositeFixedPoints() []uint8 {
	var points []uint8
	for x, v := range sb.s {
		if v == ^uint8(x) {
			points = append(points, uint8(x))
		}
	}
	return points
}

// ── Компонентные функции ──────────────────────────────────────────────────────

// Component возвращает компонентную функцию f_b(x) = <b, S(x)> над GF(2),
// где b ∈ {1,...,255} — маска выходных битов.
// Результат кэшируется.
func (sb *SBox) Component(b uint8) *boolfunc.BoolFunc {
	if b == 0 {
		panic("sbox: маска b не может быть 0")
	}
	idx := b - 1
	if sb.components[idx] == nil {
		sb.components[idx] = computeComponent(sb.s, b)
	}
	return sb.components[idx]
}

// computeComponent вычисляет компонентную функцию f_b(x) = <b, S(x)> над GF(2).
func computeComponent(s [256]uint8, b uint8) *boolfunc.BoolFunc {
	tv := make([]uint8, 256)
	for x := range 256 {
		v := b & s[x]
		v ^= v >> 4
		v ^= v >> 2
		v ^= v >> 1
		tv[x] = v & 1
	}
	f, _ := boolfunc.FromSlice(tv)
	return f
}
