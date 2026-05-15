package generate

import (
	"math"
	"math/rand/v2"
	"sync"

	"github.com/F3dosik/MGSC/sbox"
)

// Schedule — параметры расписания SA.
type Schedule struct {
	T0      float64 // начальная температура
	Alpha   float64 // скорость остывания (0.95..0.999)
	M       int     // длина плато (swap'ов на каждой T)
	Tmin    float64 // нижний порог T
	Etarget float64 // целевая cost; ранний выход при E(S*)<= Etarget
	MaxIter int     // hard limit на общее число swap'ов
}

// Result — итог одного запуска SA.
type Result struct {
	BestTable  [256]uint8
	BestCost   float64
	TotalIter  int
	Accepted   int // сколько swap'ов было принято
	AcceptedUp int // из них ухудшающих (полезно для диагностики)
}

// SA выполняет один прогон simulated annealing.
// Меняет state в процессе работы; финальное состояние state может отличаться
// от Result.BestTable.
func SA(state *SAState, sch Schedule) Result {
	best := state.Table()
	bestCost := state.Cost()

	T := sch.T0
	totalIter := 0
	accepted, acceptedUp := 0, 0

	for T > sch.Tmin && totalIter < sch.MaxIter {
		for k := 0; k < sch.M; k++ {
			i, j, oldCost := state.RandomSwap()
			delta := state.Cost() - oldCost

			// accept / reject
			if delta < 0 || rand.Float64() < math.Exp(-delta/T) {
				accepted++
				if delta >= 0 {
					acceptedUp++
				}
				if state.Cost() < bestCost {
					bestCost = state.Cost()
					best = state.Table()
					if bestCost <= sch.Etarget {
						return Result{best, bestCost,
							totalIter + 1, accepted, acceptedUp}
					}
				}
			} else {
				state.Swap(i, j) // откат — swap инволютивен
			}

			totalIter++
			if totalIter >= sch.MaxIter {
				break
			}
		}
		T *= sch.Alpha
	}

	return Result{best, bestCost, totalIter, accepted,
		acceptedUp}
}

func EstimateT0(state *SAState, probes int, p0 float64) float64 {
	var sumDelta float64
	var count int
	for k := 0; k < probes; k++ {
		i, j, oldCost := state.RandomSwap()
		delta := state.Cost() - oldCost
		if delta > 0 {
			sumDelta += delta
			count++
		}
		state.Swap(i, j)
	}
	if count == 0 {
		return 1.0
	}
	mean := sumDelta / float64(count)
	return -mean / math.Log(p0)
}

type RunResultSA struct {
	Table [256]uint8
	Cost  float64
	NL    int
}

func ParallelMultiStartSA(numRuns int, sch Schedule, cf CostFunc) []RunResultSA {
	results := make([]RunResultSA, numRuns)
	var wg sync.WaitGroup

	for idx := 0; idx < numRuns; idx++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// случайный биективный старт
			var table [256]uint8
			for k := range table {
				table[k] = uint8(k)
			}
			rand.Shuffle(len(table), func(a, b int) {
				table[a], table[b] = table[b], table[a]
			})

			state := NewSAState(table, cf)
			res := SA(state, sch)

			results[i] = RunResultSA{
				Table: res.BestTable,
				Cost:  res.BestCost,
				NL:    sbox.New(res.BestTable).Nonlinearity(),
			}
		}(idx)
	}
	wg.Wait()
	return results
}
