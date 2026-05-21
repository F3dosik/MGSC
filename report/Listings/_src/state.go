package generate

import "math/rand/v2"

// SAState — состояние симулированного отжига.
// Хранит текущую подстановку и закэшированное значение cost.
// Метод Swap идемпотентно меняет state; отмена = повторный Swap тех же индексов.
type SAState struct {
	s    [256]uint8
	cost float64
	cf   CostFunc
}

// CostFunc — что-то, что превращает подстановку в число.
// Чем меньше — тем лучше.
type CostFunc func([256]uint8) float64

// NewSAState инициализирует state с начальной подстановкой и cost-функцией.
// Сразу считает cost.
func NewSAState(initial [256]uint8, cf CostFunc) *SAState {
	st := &SAState{s: initial, cf: cf}
	st.cost = cf(initial)
	return st
}

// Cost — текущая cost (кэширована).
func (st *SAState) Cost() float64 { return st.cost }

// Table — копия текущей подстановки.
func (st *SAState) Table() [256]uint8 { return st.s }

// Swap меняет местами s[i] и s[j] и пересчитывает cost.
// Возвращает старую cost.
func (st *SAState) Swap(i, j int) (oldCost float64) {
	oldCost = st.cost
	st.s[i], st.s[j] = st.s[j], st.s[i]
	st.cost = st.cf(st.s)
	return oldCost
}

// RandomSwap выбирает две различные позиции и делает Swap.
// Возвращает индексы и старую cost.
func (st *SAState) RandomSwap() (i, j int, oldCost float64) {
	i = int(rand.UintN(256))
	j = int(rand.UintN(256))
	for j == i {
		j = int(rand.UintN(256))
	}
	return i, j, st.Swap(i, j)
}
