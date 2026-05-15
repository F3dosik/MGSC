---
title: Multi-start метаэвристики
type: entity
tags: [оптимизация, метаэвристики, параллелизм]
sources: []
created: 2026-05-15
updated: 2026-05-15
---

# Multi-start метаэвристики

**Multi-start** — стратегия: вместо одного запуска оптимизатора делать $N$ независимых запусков с разными стартовыми точками и возвращать **best-of-N**.

---

## Зачем нужно

Все локальные методы оптимизации — [[hill-climbing|HC]], [[simulated-annealing|SA]], градиентный спуск — застревают в **локальных** оптимумах. SA с температурой даёт **некоторую** возможность выбраться из них, но не гарантирует. Один запуск SA — это одна **выборка** из распределения возможных решений; её качество — случайная величина.

Multi-start снижает дисперсию финального результата:
$$P(\text{best-of-N} \leq E^*) = 1 - (1 - P(\text{single} \leq E^*))^N.$$

Если вероятность одного запуска найти решение лучше порога $E^*$ — $p$, то вероятность хотя бы одного из $N$ запусков:
$$P_N = 1 - (1 - p)^N.$$

Для $p = 0.1$ (один запуск находит «хорошее» решение в 10% случаев):
- $N=1$: 10%;
- $N=10$: 65%;
- $N=20$: 88%;
- $N=50$: 99.5%.

То есть для **редко достижимых** целей multi-start критичен.

---

## Параллелизм — линейный буст

Запуски multi-start **полностью независимы**: нет общего состояния, нет коммуникации между ними. Это **embarrassingly parallel** — идеальный случай для горутин/потоков.

На $K$-ядерном CPU $N$ запусков выполняются за время $\lceil N/K \rceil$ одного запуска. Для $N = K = 16$ — за время одного запуска получаешь 16 попыток.

Это самый практичный способ ускорить SA для задач генерации S-блоков.

---

## Реализация в Go

```go
type RunResult struct {
    Table [256]uint8
    Cost  float64
    NL    int
}

func ParallelMultiStartSA(numRuns int, sch Schedule, cf CostFunc) []RunResult {
    results := make([]RunResult, numRuns)
    var wg sync.WaitGroup
    
    for idx := 0; idx < numRuns; idx++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            
            // случайный биективный старт
            var table [256]uint8
            for k := range table { table[k] = uint8(k) }
            rand.Shuffle(len(table), func(a, b int) {
                table[a], table[b] = table[b], table[a]
            })
            
            state := NewSAState(table, cf)
            res := SA(state, sch)
            
            results[i] = RunResult{
                Table: res.BestTable,
                Cost:  res.BestCost,
                NL:    sbox.New(res.BestTable).Nonlinearity(),
            }
        }(idx)
    }
    wg.Wait()
    return results
}
```

После запуска — выбрать лучший:
```go
results := ParallelMultiStartSA(16, sch, cf)
best := results[0]
for _, r := range results[1:] {
    if r.Cost < best.Cost {
        best = r
    }
}
```

---

## Тонкости

### Память

Каждый параллельный запуск держит свой `SAState` со кэшем компонентных функций и Walsh-спектров. Для одного запуска SA ~270KB. Для 16 — ~4.3MB. Для 1000 — 270MB. На современных машинах не проблема.

### PRNG

`math/rand/v2` потокобезопасен на top-level — `rand.Shuffle`, `rand.UintN`, `rand.Float64` можно вызывать из любой горутины. Для **детерминированности** (воспроизводимые эксперименты) нужен `rand.New(rand.NewPCG(seed1, seed2))` в каждой горутине со своим seed'ом.

### Worker pool

Для $N \gg K$ ядер (например $N=1000$, $K=16$) запускать 1000 горутин — overhead. Лучше **worker pool**:

```go
jobs := make(chan int, numRuns)
results := make(chan RunResult, numRuns)

for w := 0; w < runtime.NumCPU(); w++ {
    go func() {
        for range jobs {
            // ... run SA ...
            results <- result
        }
    }()
}
```

Это даёт **строгое** $K$ горутин, исполняющих $N$ задач. Не лучше по wall-time, но экономит память.

---

## Multi-start vs Parallel Tempering

Multi-start — **N независимых** запусков, **best-of-N**.

**Parallel tempering** — продвинутая альтернатива: $N$ SA-запусков с **разными температурами** $T_1 < T_2 < \ldots < T_N$, периодически обмениваются состояниями. Это даёт лучшее качество поиска, но **сложнее**, требует коммуникации между запусками. Для нашей задачи multi-start обычно достаточно; parallel tempering — следующая ступень.

---

## Связанные

- [[simulated-annealing]] — основа метода.
- [[hill-climbing]] — multi-start особенно важен для HC (без температуры застревание неизбежно).
- [[t0-calibration]] — калибровка $T_0$ может быть **общей** для всех запусков, если cost стабильна (см. [[sa-empirical-observations]]).

## Источники

- Martí R., Resende M.G.C., Ribeiro C.C. **Multi-start methods for combinatorial optimization** // European Journal of Operational Research. 2013. Vol. 226, № 1. P. 1–8.
