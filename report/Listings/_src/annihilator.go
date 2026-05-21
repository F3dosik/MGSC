package boolfunc

// subsets возвращает все подмножества {0..n-1} размера <= d
// каждое подмножество представлено как []int индексов переменных
func subsets(n, d int) [][]int {
	var result [][]int
	var rec func(start, rem int, current []int)
	rec = func(start, rem int, current []int) {
		// добавляем текущее подмножество (включая пустое — константный моном)
		tmp := make([]int, len(current))
		copy(tmp, current)
		result = append(result, tmp)
		if rem == 0 {
			return
		}
		for i := start; i < n; i++ {
			rec(i+1, rem-1, append(current, i))
		}
	}
	rec(0, d, []int{})
	return result
}

// evalMonomial вычисляет значение монома S в точке x над GF(2)
// x передаётся как uint8 — битовое представление точки
func evalMonomial(x uint8, subset []int) uint8 {
	val := uint8(1)
	for _, i := range subset {
		val &= (x >> i) & 1
	}
	return val
}
