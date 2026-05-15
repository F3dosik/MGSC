package main

import (
	"fmt"

	"github.com/F3dosik/MGSC/generate"
)

func main() {
	// sb, _ := sbox.LoadFromFile("testdata/aes_sbox.json")
	// cf := generate.ClarkJacobStepneyNL(24, 3)
	// state := generate.NewSAState(sb.Table(), cf)

	// fmt.Println("iter,nl,cost")
	// for iter := 0; iter < 1000; iter++ {
	// 	nl := sbox.New(state.Table()).Nonlinearity()
	// 	fmt.Printf("%d,%d,%.0f\n", iter, nl, state.Cost())
	// 	state.RandomSwap()
	// }

	results := generate.ParallelMultiStart(20, 100_000, generate.ClarkJacobStepneyNL(24, 3))
	for i, r := range results {
		fmt.Printf("run %d: cost=%.0f, NL=%d\n", i, r.Cost, r.NL)
	}

}
