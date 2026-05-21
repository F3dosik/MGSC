// Package main — CLI инструмента mgsc.
//
// Использование:
//
//	mgsc verify   [-profile aes|relaxed|research] [-json] <file.json>
//	mgsc analyze  [-json] <file.json>
//	mgsc generate [флаги SA] [-out file.json] [-json]
//
// Флаги указываются перед позиционными аргументами (поведение пакета flag).
//
// Подкоманды независимы; общий флаг -json переключает вывод в структурированный JSON.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/F3dosik/MGSC/generate"
	"github.com/F3dosik/MGSC/sbox"
	"github.com/F3dosik/MGSC/verify"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "verify":
		cmdVerify(args)
	case "analyze":
		cmdAnalyze(args)
	case "generate":
		cmdGenerate(args)
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "неизвестная подкоманда: %s\n\n", cmd)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `mgsc — инструмент анализа и генерации S-блоков 8×8.

Подкоманды:
  verify    [флаги] <file.json>   проверить S-блок по профилю порогов
  analyze   [флаги] <file.json>   вычислить и напечатать все метрики
  generate  [флаги]               сгенерировать S-блок методом SA

Флаги указываются ПЕРЕД позиционными аргументами.
Используйте «mgsc <команда> -h» для списка флагов конкретной подкоманды.`)
}

// ── verify ──────────────────────────────────────────────────────────────────

func cmdVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	profile := fs.String("profile", "aes", "профиль порогов: aes | relaxed | research")
	asJSON := fs.Bool("json", false, "вывод в JSON")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: mgsc verify [-profile aes|relaxed|research] [-json] <file.json>")
		fs.PrintDefaults()
	}
	_ = fs.Parse(args)

	pos := fs.Args()
	if len(pos) != 1 {
		fs.Usage()
		os.Exit(2)
	}

	sb, err := sbox.LoadFromFile(pos[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	t, err := resolveProfile(*profile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	r := verify.Verify(sb, t)
	if *asJSON {
		printJSON(r)
	} else {
		fmt.Println(r.Summary())
		printVerifyDetails(r)
	}
	if !r.Passed {
		os.Exit(1)
	}
}

func printVerifyDetails(r verify.Report) {
	if !r.Bijective {
		return
	}
	if n := len(r.FixedPoints); n > 0 {
		fmt.Printf("  Неподвижных точек:    %d\n", n)
	}
	if n := len(r.OppositeFixed); n > 0 {
		fmt.Printf("  Антиподвижных точек:  %d\n", n)
	}
	if !r.Passed {
		fmt.Println("  Нарушения:")
		for _, f := range r.Failures {
			fmt.Printf("    - %s\n", f)
		}
	}
}

// ── analyze ─────────────────────────────────────────────────────────────────

func cmdAnalyze(args []string) {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "вывод в JSON")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: mgsc analyze [-json] <file.json>")
		fs.PrintDefaults()
	}
	_ = fs.Parse(args)

	pos := fs.Args()
	if len(pos) != 1 {
		fs.Usage()
		os.Exit(2)
	}

	sb, err := sbox.LoadFromFile(pos[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	printAnalysis(sb, *asJSON)
}

type analysisReport struct {
	Bijective     bool    `json:"bijective"`
	NL            int     `json:"NL,omitempty"`
	DU            int     `json:"DU,omitempty"`
	Deg           int     `json:"deg,omitempty"`
	AI            int     `json:"AI,omitempty"`
	FixedPoints   []uint8 `json:"fixed_points,omitempty"`
	OppositeFixed []uint8 `json:"opposite_fixed,omitempty"`
}

func printAnalysis(sb *sbox.SBox, asJSON bool) {
	rep := analysisReport{Bijective: sb.IsBijective()}
	if rep.Bijective {
		rep.NL = sb.Nonlinearity()
		rep.DU = sb.DiffUniformity()
		rep.Deg = sb.AlgebraicDegree()
		rep.AI = sb.AlgebraicImmunity()
		rep.FixedPoints = sb.FixedPoints()
		rep.OppositeFixed = sb.OppositeFixedPoints()
	}

	if asJSON {
		printJSON(rep)
		return
	}

	if !rep.Bijective {
		fmt.Println("Подстановка не биективна — метрики не вычисляются.")
		return
	}
	fmt.Printf("Биективна:              да\n")
	fmt.Printf("Нелинейность (NL):      %d\n", rep.NL)
	fmt.Printf("Дифф. равномерность δ:  %d\n", rep.DU)
	fmt.Printf("Алгебр. степень:        %d\n", rep.Deg)
	fmt.Printf("Алгебр. иммунность:     %d\n", rep.AI)
	fmt.Printf("Неподвижных точек:      %d\n", len(rep.FixedPoints))
	fmt.Printf("Антиподвижных точек:    %d\n", len(rep.OppositeFixed))
}

// ── generate ────────────────────────────────────────────────────────────────

func cmdGenerate(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	costName := fs.String("cost", "minmax", "функция штрафов: minmax | threshold | cjs")
	target := fs.Int("target", 112, "целевая нелинейность (для threshold)")
	x := fs.Int("x", 0, "параметр X функции CJS")
	r := fs.Int("r", 3, "параметр R функции CJS")
	alpha := fs.Float64("alpha", 0.99, "скорость остывания")
	mLen := fs.Int("m", 5000, "длина плато (swap'ов на каждой температуре)")
	tmin := fs.Float64("tmin", 0.01, "нижний порог температуры")
	maxIter := fs.Int("max-iter", 10_000_000, "верхний предел числа итераций")
	t0 := fs.Float64("t0", 0, "начальная температура (0 = автокалибровка)")
	p0 := fs.Float64("p0", 0.8, "целевая доля принятых ухудшений при калибровке T0")
	probes := fs.Int("probes", 1000, "число проб при калибровке T0")
	runs := fs.Int("runs", 1, "число параллельных запусков (мультистарт)")
	out := fs.String("out", "", "сохранить лучший S-блок в файл JSON")
	asJSON := fs.Bool("json", false, "вывод в JSON")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: mgsc generate [флаги]")
		fs.PrintDefaults()
	}
	_ = fs.Parse(args)

	cf, err := resolveCost(*costName, *target, *x, *r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Автокалибровка T0, если не задано явно.
	t0Used := *t0
	if t0Used == 0 {
		probe := randomSBoxTable()
		probeState := generate.NewSAState(probe, cf)
		t0Used = generate.EstimateT0(probeState, *probes, *p0)
	}

	sch := generate.Schedule{
		T0:      t0Used,
		Alpha:   *alpha,
		M:       *mLen,
		Tmin:    *tmin,
		Etarget: 0,
		MaxIter: *maxIter,
	}

	best := runSA(*runs, sch, cf)
	sb := sbox.New(best.Table)

	if *out != "" {
		if err := sbox.SaveToFile(sb, *out, fmt.Sprintf("generated-NL%d", best.NL)); err != nil {
			fmt.Fprintln(os.Stderr, "ошибка сохранения:", err)
			os.Exit(1)
		}
	}

	if *asJSON {
		printJSON(map[string]any{
			"cost_function":  *costName,
			"runs":           *runs,
			"t0_used":        t0Used,
			"best_NL":        best.NL,
			"best_cost":      best.Cost,
			"verify_aes":     verify.Verify(sb, verify.AESLevel),
			"verify_relaxed": verify.Verify(sb, verify.RelaxedLevel),
			"out":            *out,
		})
		return
	}

	fmt.Printf("Cost-функция:     %s\n", *costName)
	fmt.Printf("Запусков:         %d\n", *runs)
	fmt.Printf("T0 использован:   %g\n", t0Used)
	fmt.Println()
	fmt.Printf("Лучший результат: NL=%d, cost=%g\n", best.NL, best.Cost)
	fmt.Println()
	fmt.Println(verify.Verify(sb, verify.AESLevel).Summary())
	fmt.Println(verify.Verify(sb, verify.RelaxedLevel).Summary())
	if *out != "" {
		fmt.Printf("\nСохранено в %s\n", *out)
	}
}

func runSA(numRuns int, sch generate.Schedule, cf generate.CostFunc) generate.RunResultSA {
	if numRuns > 1 {
		results := generate.ParallelMultiStartSA(numRuns, sch, cf)
		best := results[0]
		for _, res := range results[1:] {
			if res.NL > best.NL || (res.NL == best.NL && res.Cost < best.Cost) {
				best = res
			}
		}
		return best
	}

	table := randomSBoxTable()
	state := generate.NewSAState(table, cf)
	res := generate.SA(state, sch)
	return generate.RunResultSA{
		Table: res.BestTable,
		Cost:  res.BestCost,
		NL:    sbox.New(res.BestTable).Nonlinearity(),
	}
}

// ── helpers ─────────────────────────────────────────────────────────────────

func resolveProfile(name string) (verify.Threshold, error) {
	switch name {
	case "aes":
		return verify.AESLevel, nil
	case "relaxed":
		return verify.RelaxedLevel, nil
	case "research":
		return verify.ResearchLevel, nil
	default:
		return verify.Threshold{}, fmt.Errorf("неизвестный профиль: %s (доступны: aes, relaxed, research)", name)
	}
}

func resolveCost(name string, target, x, r int) (generate.CostFunc, error) {
	switch name {
	case "minmax":
		return generate.MinMaxWalsh(), nil
	case "threshold":
		return generate.ThresholdNL(target), nil
	case "cjs":
		return generate.ClarkJacobStepneyNL(x, r), nil
	default:
		return nil, fmt.Errorf("неизвестная cost-функция: %s (доступны: minmax, threshold, cjs)", name)
	}
}

func randomSBoxTable() [256]uint8 {
	var table [256]uint8
	for k := range table {
		table[k] = uint8(k)
	}
	rand.Shuffle(len(table), func(a, b int) {
		table[a], table[b] = table[b], table[a]
	})
	return table
}

func printJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка JSON-вывода:", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
