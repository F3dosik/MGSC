package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/F3dosik/MGSC/generate"
	"github.com/F3dosik/MGSC/sbox"
)

func printNLDistribution(results []generate.RunResultSA) {
	nls := make([]int, len(results))
	for i, r := range results {
		nls[i] = r.NL
	}
	sort.Ints(nls)

	counts := map[int]int{}
	for _, nl := range nls {
		counts[nl]++
	}
	for nl := nls[0]; nl <= nls[len(nls)-1]; nl += 2 {
		if c, ok := counts[nl]; ok {
			fmt.Printf("  NL=%-3d %s %d/%d\n", nl, strings.Repeat("█", c), c, len(results))
		}
	}
}

// multiStartFromTable запускает numRuns SA параллельно, все из одной стартовой таблицы.
func multiStartFromTable(start [256]uint8, sch generate.Schedule, cf generate.CostFunc, numRuns int) []generate.RunResultSA {
	results := make([]generate.RunResultSA, numRuns)
	var wg sync.WaitGroup
	for idx := 0; idx < numRuns; idx++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			state := generate.NewSAState(start, cf)
			res := generate.SA(state, sch)
			results[i] = generate.RunResultSA{
				Table: res.BestTable,
				Cost:  res.BestCost,
				NL:    sbox.New(res.BestTable).Nonlinearity(),
			}
		}(idx)
	}
	wg.Wait()
	return results
}

func main() {
	const numRuns = 16
	cf := generate.MinMaxWalsh()

	aes, err := sbox.LoadFromFile("testdata/aes_sbox.json")
	if err != nil {
		panic(err)
	}
	aesTable := aes.Table()
	fmt.Printf("AES: NL=%d  cost=%.0f\n\n", aes.Nonlinearity(), cf(aesTable))

	// Калибруем T0 из AES — ΔE будут маленькими (AES в глубоком минимуме MinMaxWalsh)
	calibState := generate.NewSAState(aesTable, cf)
	t0 := generate.EstimateT0(calibState, 1000, 0.8)

	sch := generate.Schedule{
		T0:      t0,
		Alpha:   0.99,
		M:       5000,
		Tmin:    0.01,
		Etarget: 0,
		MaxIter: 10_000_000,
	}
	fmt.Printf("T0=%.2e  α=%.2f  M=%d  N=%d\n\n", sch.T0, sch.Alpha, sch.M, numRuns)

	// E-10: multi-start SA из AES
	fmt.Println("=== E-10: multi-start SA из AES ===")
	results := multiStartFromTable(aesTable, sch, cf, numRuns)

	fmt.Println("NL distribution:")
	printNLDistribution(results)

	best := results[0]
	for _, r := range results[1:] {
		if r.NL > best.NL {
			best = r
		}
	}
	fmt.Printf("\nBest NL: %d\n", best.NL)

	// Для сравнения: сколько уникальных S-блоков с NL=112 нашли?
	// (даже если все NL=112 — это разные объекты, не AES)
	atTarget := 0
	for _, r := range results {
		if r.NL >= 112 {
			atTarget++
		}
	}
	fmt.Printf("NL>=112: %d/%d\n", atTarget, numRuns)

	// Проверяем: лучший результат — это сам AES или другой S-блок?
	if best.NL >= 112 {
		bestSB := sbox.New(best.Table)
		if best.Table == aesTable {
			fmt.Println("\nBest table = AES (не сдвинулись)")
		} else {
			fmt.Printf("\nBest table ≠ AES (нашли другой S-блок: NL=%d DU=%d deg=%d)\n",
				bestSB.Nonlinearity(), bestSB.DiffUniformity(), bestSB.AlgebraicDegree())
		}
	}

	// Baseline: SA из random — калибруем T0 отдельно (ΔE из random >> ΔE из AES)
	var randTable [256]uint8
	for i := range randTable {
		randTable[i] = uint8(i)
	}
	randState := generate.NewSAState(randTable, cf)
	t0Random := generate.EstimateT0(randState, 1000, 0.8)
	schRandom := sch
	schRandom.T0 = t0Random

	fmt.Printf("\n=== Baseline: multi-start SA из random (T0=%.2e) ===\n", t0Random)
	randomResults := generate.ParallelMultiStartSA(numRuns, schRandom, cf)
	fmt.Println("NL distribution:")
	printNLDistribution(randomResults)
	bestRandom := randomResults[0]
	for _, r := range randomResults[1:] {
		if r.NL > bestRandom.NL {
			bestRandom = r
		}
	}
	fmt.Printf("\nBest NL (random start): %d\n", bestRandom.NL)

	// Ещё один тест: HC из AES — есть ли NL>112 соседи?
	fmt.Println("\n=== HC из AES: ищем соседей с NL>112 ===")
	hcState := generate.NewSAState(aesTable, cf)
	bestNL := aes.Nonlinearity()
	improvements := 0
	for iter := 0; iter < 500_000; iter++ {
		i, j, oldCost := hcState.RandomSwap()
		if hcState.Cost() >= oldCost {
			hcState.Swap(i, j)
		} else {
			improvements++
			if nl := sbox.New(hcState.Table()).Nonlinearity(); nl > bestNL {
				bestNL = nl
				fmt.Printf("  iter=%d  NL=%d\n", iter, nl)
			}
		}
	}
	if improvements == 0 {
		fmt.Println("  0 улучшений — AES является локальным оптимумом MinMaxWalsh")
	} else {
		fmt.Printf("  Best NL из HC: %d  (улучшений: %d)\n", bestNL, improvements)
	}

}
