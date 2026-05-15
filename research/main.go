package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/F3dosik/MGSC/generate"
	"github.com/F3dosik/MGSC/sbox"
)

func main() {
	// #1
	// sb, _ := sbox.LoadFromFile("testdata/aes_sbox.json")
	// cf := generate.ClarkJacobStepneyNL(24, 3)
	// state := generate.NewSAState(sb.Table(), cf)

	// fmt.Println("iter,nl,cost")
	// for iter := 0; iter < 1000; iter++ {
	// 	nl := sbox.New(state.Table()).Nonlinearity()
	// 	fmt.Printf("%d,%d,%.0f\n", iter, nl, state.Cost())
	// 	state.RandomSwap()
	// }

	// #2
	// results := generate.ParallelMultiStart(20, 100_000, generate.ClarkJacobStepneyNL(24, 3))
	// for i, r := range results {
	// 	fmt.Printf("run %d: cost=%.0f, NL=%d\n", i, r.Cost, r.NL)
	// }

	// Калибровка начальной температуры
	// cf := generate.ClarkJacobStepneyNL(24, 3)
	// fmt.Println("\nRandom start, 10 calibrations:")
	// for run := 0; run < 10; run++ {
	// 	var table [256]uint8
	// 	for i := range table {
	// 		table[i] = uint8(i)
	// 	}
	// 	rand.Shuffle(len(table), func(i, j int) { table[i], table[j] = table[j], table[i] })

	// 	state := generate.NewSAState(table, cf)
	// 	t0 := generate.EstimateT0(state, 1000, 0.8)
	// 	initialCost := state.Cost()
	// 	fmt.Printf("  initial cost: %.2e, T0: %.0f\n", initialCost, t0)
	// }

	sch := generate.Schedule{
		T0:      2_000_000,
		Alpha:   0.99,
		M:       1000,
		Tmin:    1,
		Etarget: 0,
		MaxIter: 5_000_000,
	}

	// случайный старт
	var table [256]uint8
	for i := range table {
		table[i] = uint8(i)
	}
	rand.Shuffle(len(table), func(i, j int) { table[i], table[j] = table[j], table[i] })

	state := generate.NewSAState(table, generate.ClarkJacobStepneyNL(24, 3))
	res := generate.SA(state, sch)

	nl := sbox.New(res.BestTable).Nonlinearity()
	fmt.Printf("best cost: %.2e, NL: %d, iter: %d, accepted: %d (up: %d)\n",
		res.BestCost, nl, res.TotalIter, res.Accepted, res.AcceptedUp)
}
