package boolfunc

import "math/bits"

// computeNonlinearity вычисляет нелинейность за один проход по спектру Уолша
// и сохраняет результат в кэш.
func (f *BoolFunc) computeNonlinearity() {
	walsh := f.Walsh()

	maxAbs := int32(0)
	for _, w := range walsh {
		if w > maxAbs {
			maxAbs = w
		} else if -w > maxAbs {
			maxAbs = -w
		}
	}

	nl := (1 << (f.n - 1)) - int(maxAbs>>1)
	f.nlin = &nl
}

// Nonlinearity возвращает нелинейность булевой функции —
// минимальное расстояние Хэмминга до множества аффинных функций.
//
// Формула: NL(f) = 2^(n-1) - max|W[a]| / 2
func (f *BoolFunc) Nonlinearity() int {
	if f.nlin == nil {
		f.computeNonlinearity()
	}
	return *f.nlin
}

// AlgebraicDegree возвращает алгебраическую степень булевой функции —
// максимальный вес Хэмминга среди индексов ненулевых коэффициентов АНФ.
func (f *BoolFunc) AlgebraicDegree() int {
	if f.degree != nil {
		return *f.degree
	}

	anf := f.ANF()
	deg := 0
	for i, coeff := range anf {
		if coeff == 1 {
			d := bits.OnesCount(uint(i))
			if d > deg {
				deg = d
			}
		}
	}

	f.degree = &deg
	return deg
}

// gaussianRank вычисляет ранг матрицы над GF(2)
// matrix[i] — i-я строка, упакованная побитово
// cols — количество столбцов
func gaussianRank(matrix []uint64, cols int) int {
	rows := make([]uint64, len(matrix))
	copy(rows, matrix)

	rank := 0
	for col := 0; col < cols; col++ {
		// ищем опорную строку с единицей в позиции col
		pivot := -1
		for row := rank; row < len(rows); row++ {
			if (rows[row]>>col)&1 == 1 {
				pivot = row
				break
			}
		}
		if pivot == -1 {
			continue // свободная переменная
		}

		// переставляем опорную строку наверх
		rows[rank], rows[pivot] = rows[pivot], rows[rank]

		// обнуляем col во всех остальных строках
		for row := 0; row < len(rows); row++ {
			if row != rank && (rows[row]>>col)&1 == 1 {
				rows[row] ^= rows[rank]
			}
		}
		rank++
	}
	return rank
}

// hasAnnihilator проверяет существование аннулятора степени <= d
// для функции tv либо для (1 xor tv) в зависимости от флага complement
func hasAnnihilator(tv []uint8, n, d int, complement bool) bool {
	monomials := subsets(n, d)
	cols := len(monomials)

	// собираем строки матрицы — точки носителя
	var rows []uint64
	for x := 0; x < (1 << n); x++ {
		fx := tv[x]
		if complement {
			fx ^= 1
		}
		if fx != 1 {
			continue
		}
		// строим строку: значения всех мономов в точке x
		var row uint64
		for j, mono := range monomials {
			if evalMonomial(uint8(x), mono) == 1 {
				row |= 1 << j
			}
		}
		rows = append(rows, row)
	}

	// ненулевой аннулятор существует если rank < cols
	return gaussianRank(rows, cols) < cols
}

// AlgebraicImmunity возвращает алгебраическую иммунность булевой функции
func (f *BoolFunc) AlgebraicImmunity() int {
	if f.ai != nil {
		return *f.ai
	}
	maxDeg := (f.n + 1) / 2
	result := maxDeg
	for d := 1; d <= maxDeg; d++ {
		if hasAnnihilator(f.tv, f.n, d, false) ||
			hasAnnihilator(f.tv, f.n, d, true) {
			result = d
			break
		}
	}
	f.ai = &result
	return result
}
