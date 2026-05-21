// Бенчмарк случайных биективных подстановок 8×8.
//
// Запуск:
//
//	go run ./research/random_benchmark -n 1000 -out research/random_benchmark/benchmark.csv
//
// Программа генерирует N случайных биективных S-блоков через rand.Shuffle,
// вычисляет на каждой четыре криптографические метрики и две структурные
// характеристики, сохраняет всё в CSV и печатает сводку процентилей.
// Доля прохождения профилей verify даёт baseline для главы 7 курсовой.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/F3dosik/MGSC/sbox"
	"github.com/F3dosik/MGSC/verify"
)

type record struct {
	idx     int
	nl      int
	du      int
	deg     int
	ai      int
	nFixed  int
	nOpp    int
	passAES bool
	passRel bool
	passRes bool
}

func main() {
	n := flag.Int("n", 1000, "число случайных подстановок")
	seed := flag.Uint64("seed", 42, "seed для воспроизводимости")
	workers := flag.Int("workers", runtime.NumCPU(), "число горутин-исполнителей")
	out := flag.String("out", "benchmark.csv", "выходной CSV-файл")
	flag.Parse()

	fmt.Fprintf(os.Stderr, "Бенчмарк: N=%d, workers=%d, seed=%d\n", *n, *workers, *seed)

	jobs := make(chan int, *n)
	results := make(chan record, *n)

	var wg sync.WaitGroup
	for w := 0; w < *workers; w++ {
		wg.Add(1)
		// У каждого воркера свой детерминированный PCG-источник.
		go func(workerID int) {
			defer wg.Done()
			r := rand.New(rand.NewPCG(*seed, uint64(workerID+1)))
			for idx := range jobs {
				table := randomTable(r)
				sb := sbox.New(table)
				results <- record{
					idx:     idx,
					nl:      sb.Nonlinearity(),
					du:      sb.DiffUniformity(),
					deg:     sb.AlgebraicDegree(),
					ai:      sb.AlgebraicImmunity(),
					nFixed:  len(sb.FixedPoints()),
					nOpp:    len(sb.OppositeFixedPoints()),
					passAES: verify.Verify(sb, verify.AESLevel).Passed,
					passRel: verify.Verify(sb, verify.RelaxedLevel).Passed,
					passRes: verify.Verify(sb, verify.ResearchLevel).Passed,
				}
			}
		}(w)
	}

	go func() {
		for i := 0; i < *n; i++ {
			jobs <- i
		}
		close(jobs)
	}()
	go func() {
		wg.Wait()
		close(results)
	}()

	all := make([]record, 0, *n)
	start := time.Now()
	for r := range results {
		all = append(all, r)
		if len(all)%100 == 0 {
			fmt.Fprintf(os.Stderr, "\r%d/%d", len(all), *n)
		}
	}
	fmt.Fprintf(os.Stderr, "\rзавершено за %s\n", time.Since(start))

	sort.Slice(all, func(i, j int) bool { return all[i].idx < all[j].idx })
	if err := writeCSV(*out, all); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка записи CSV:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "CSV сохранён:", *out)

	printSummary(all)
}

func randomTable(r *rand.Rand) [256]uint8 {
	var t [256]uint8
	for i := range t {
		t[i] = uint8(i)
	}
	r.Shuffle(len(t), func(a, b int) { t[a], t[b] = t[b], t[a] })
	return t
}

func writeCSV(path string, records []record) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"idx", "nl", "du", "deg", "ai", "fixed", "opposite", "passAES", "passRel", "passRes"}); err != nil {
		return err
	}
	for _, r := range records {
		row := []string{
			strconv.Itoa(r.idx),
			strconv.Itoa(r.nl),
			strconv.Itoa(r.du),
			strconv.Itoa(r.deg),
			strconv.Itoa(r.ai),
			strconv.Itoa(r.nFixed),
			strconv.Itoa(r.nOpp),
			strconv.FormatBool(r.passAES),
			strconv.FormatBool(r.passRel),
			strconv.FormatBool(r.passRes),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func printSummary(records []record) {
	n := len(records)

	nls := make([]int, n)
	dus := make([]int, n)
	degs := make([]int, n)
	ais := make([]int, n)
	fixed := make([]int, n)
	opp := make([]int, n)
	noFixed, noOpp := 0, 0
	passAES, passRel, passRes := 0, 0, 0
	for i, r := range records {
		nls[i] = r.nl
		dus[i] = r.du
		degs[i] = r.deg
		ais[i] = r.ai
		fixed[i] = r.nFixed
		opp[i] = r.nOpp
		if r.nFixed == 0 {
			noFixed++
		}
		if r.nOpp == 0 {
			noOpp++
		}
		if r.passAES {
			passAES++
		}
		if r.passRel {
			passRel++
		}
		if r.passRes {
			passRes++
		}
	}

	fmt.Println()
	fmt.Println("=== Сводка по", n, "случайным биективным S-блокам ===")
	printStats("NL    ", nls)
	printStats("δ     ", dus)
	printStats("deg   ", degs)
	printStats("AI    ", ais)
	printStats("|F.т.|", fixed)
	printStats("|АФт.|", opp)

	fmt.Println()
	pct := func(x int) float64 { return 100 * float64(x) / float64(n) }
	fmt.Printf("Без неподвижных точек:    %d / %d (%.1f%%)   теор. 1/e ≈ 36.8%%\n", noFixed, n, pct(noFixed))
	fmt.Printf("Без антиподвижных точек:  %d / %d (%.1f%%)\n", noOpp, n, pct(noOpp))
	fmt.Println()
	fmt.Printf("Под AESLevel:             %d / %d (%.1f%%)\n", passAES, n, pct(passAES))
	fmt.Printf("Под ResearchLevel:        %d / %d (%.1f%%)\n", passRes, n, pct(passRes))
	fmt.Printf("Под RelaxedLevel:         %d / %d (%.1f%%)\n", passRel, n, pct(passRel))
}

func printStats(name string, xs []int) {
	sorted := append([]int(nil), xs...)
	sort.Ints(sorted)
	n := len(sorted)
	mean := 0.0
	for _, v := range xs {
		mean += float64(v)
	}
	mean /= float64(n)
	p := func(q float64) int {
		i := int(q * float64(n-1))
		return sorted[i]
	}
	fmt.Printf("%s: min=%d  p5=%d  p25=%d  med=%d  p75=%d  p95=%d  max=%d  mean=%.2f\n",
		name, sorted[0], p(0.05), p(0.25), p(0.50), p(0.75), p(0.95), sorted[n-1], mean)
}
