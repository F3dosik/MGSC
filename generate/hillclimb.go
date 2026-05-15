package generate

import (
	"math/rand/v2"
	"sync"

	"github.com/F3dosik/MGSC/sbox"
)

type RunResult struct {
	Table [256]uint8
	Cost  float64
	NL    int
}

// ParallelMultiStart запускает numRuns независимых HC параллельно.
// Возвращает все результаты — выбирать лучший снаружи.
func ParallelMultiStart(numRuns, maxIter int, cf CostFunc) []RunResult {
	results := make([]RunResult, numRuns)
	var wg sync.WaitGroup

	for idx := 0; idx < numRuns; idx++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// случайная биективная подстановка
			var table [256]uint8
			for k := range table {
				table[k] = uint8(k)
			}
			rand.Shuffle(len(table), func(a, b int) {
				table[a], table[b] = table[b], table[a]
			})

			state := NewSAState(table, cf)
			HillClimb(state, maxIter)

			results[i] = RunResult{
				Table: state.Table(),
				Cost:  state.Cost(),
				NL:    sbox.New(state.Table()).Nonlinearity(),
			}
		}(idx)
	}
	wg.Wait()
	return results
}

// HillClimb — простейший локальный поиск. Принимает swap только при улучшении cost.
// Возвращает финальное состояние и число принятых swap'ов.
func HillClimb(state *SAState, maxIter int) (bestTable [256]uint8, bestCost float64) {
	bestTable = state.Table()
	bestCost = state.Cost()
	for iter := 0; iter < maxIter; iter++ {
		i, j, oldCost := state.RandomSwap()
		if state.Cost() >= oldCost {
			state.Swap(i, j)
		}
		if state.Cost() < bestCost {
			bestCost = state.Cost()
			bestTable = state.Table()
		}
	}
	return
}
