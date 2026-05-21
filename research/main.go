package main

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"

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

func main() {
	const numRuns = 16
	cf := generate.ClarkJacobStepneyNL(36, 4)

	// T0 калибруем под нашу реализацию — T0=20000 из MDPI 2022 несовместим
	// с нашим масштабом cost (3e10 vs их ожидаемые единицы тысяч)
	calibState := generate.NewSAState(randomTable(), cf)
	t0 := generate.EstimateT0(calibState, 1000, 0.8)

	sch := generate.Schedule{
		T0:      t0,
		Alpha:   0.99, // замедляем: α=0.95 давало только ~298K итераций
		M:       5000, // как в E-05 для fair-сравнения с MinMaxWalsh (~3.7M итераций)
		Tmin:    0.01,
		Etarget: 0,
		MaxIter: 50_000_000,
	}

	fmt.Println("=== E-09c: CJS(X=36, R=4, α=0.99, M=5000) — fair сравнение с E-05 ===")
	fmt.Printf("T0=%.2e  α=%.2f  M=%d  N=%d\n\n", sch.T0, sch.Alpha, sch.M, numRuns)

	results := generate.ParallelMultiStartSA(numRuns, sch, cf)

	fmt.Println("NL distribution:")
	printNLDistribution(results)

	best := results[0]
	for _, r := range results[1:] {
		if r.NL > best.NL {
			best = r
		}
	}
	fmt.Printf("\nBest NL: %d\n", best.NL)

	// Сколько запусков дали NL >= 104?
	above104 := 0
	for _, r := range results {
		if r.NL >= 104 {
			above104++
		}
	}
	fmt.Printf("NL>=104: %d/%d (литература: ~56%%)\n", above104, numRuns)

	// CJS cost для AES S-блока как ориентир
	aes, err := sbox.LoadFromFile("testdata/aes_sbox.json")
	if err != nil {
		panic(err)
	}

	probe := generate.NewSAState(aes.Table(), cf)
	fmt.Printf("\nCJS cost для AES: %.2e\n", probe.Cost())
}
