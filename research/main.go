package main

import (
	"fmt"
	"math/rand/v2"
	"sort"

	"github.com/F3dosik/MGSC/generate"
	"github.com/F3dosik/MGSC/sbox"
)

func randomTable() [256]uint8 {
	var table [256]uint8
	for i := range table {
		table[i] = uint8(i)
	}
	rand.Shuffle(len(table), func(i, j int) { table[i], table[j] = table[j], table[i] })
	return table
}

func runHC(label string, startTable [256]uint8, cf generate.CostFunc, maxIter int) int {
	state := generate.NewSAState(startTable, cf)
	bestNL := sbox.New(startTable).Nonlinearity()
	for iter := 0; iter < maxIter; iter++ {
		i, j, oldCost := state.RandomSwap()
		if state.Cost() >= oldCost {
			state.Swap(i, j)
		}
		if nl := sbox.New(state.Table()).Nonlinearity(); nl > bestNL {
			bestNL = nl
		}
	}
	fmt.Printf("HC (%s): NL=%d\n", label, bestNL)
	return bestNL
}

func main() {
	const numRuns = 16
	cf := generate.MinMaxWalsh()

	// Калибруем T0 один раз (стабильно для random start)
	calibState := generate.NewSAState(randomTable(), cf)
	t0 := generate.EstimateT0(calibState, 1000, 0.8)
	fmt.Printf("T0=%.2e\n\n", t0)

	sch := generate.Schedule{
		T0:      t0,
		Alpha:   0.99,
		M:       5000,
		Tmin:    0.01,
		Etarget: 0,
		MaxIter: 10_000_000,
	}

	fmt.Printf("=== E-08: multi-start SA, N=%d ===\n\n", numRuns)
	results := generate.ParallelMultiStartSA(numRuns, sch, cf)

	// Распределение NL
	nls := make([]int, numRuns)
	for i, r := range results {
		nls[i] = r.NL
	}
	sort.Ints(nls)

	nlCounts := map[int]int{}
	for _, nl := range nls {
		nlCounts[nl]++
	}
	fmt.Println("NL distribution:")
	for nl := nls[0]; nl <= nls[len(nls)-1]; nl += 2 {
		if c, ok := nlCounts[nl]; ok {
			fmt.Printf("  NL=%d: %d/%d\n", nl, c, numRuns)
		}
	}

	// Лучший результат
	best := results[0]
	for _, r := range results[1:] {
		if r.NL > best.NL {
			best = r
		}
	}
	fmt.Printf("\nBest SA: NL=%d\n\n", best.NL)

	// HC из лучшего
	runHC("из best SA", best.Table, cf, 500_000)
}
