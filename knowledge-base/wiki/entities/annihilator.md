---
title: Аннигиляторы булевой функции
type: entity
tags: [булевы-функции, алгебраические-атаки]
sources: [coursework-1]
created: 2026-05-15
updated: 2026-05-15
---

# Аннигиляторы

## Определение

Аннигилятор функции $f \in \mathcal{F}_n$ — функция $g \in \mathcal{F}_n$ такая, что:
$$f(x) \cdot g(x) = 0 \quad \forall x \in \mathbb{Z}_2^n.$$

Эквивалентно: $\mathrm{supp}(g) \subseteq \overline{\mathrm{supp}(f)}$ — носитель $g$ лежит в нулях $f$.

## Криптографический смысл

Если для $f$ существует аннигилятор низкой [[algebraic-degree|алгебраической степени]] — это **алгебраическая атака** (Courtois, Meier, 2003): шифр сводится к системе уравнений низкой степени, разрешимой за полиномиальное время.

Минимальная степень аннигиляторов $f$ или $f \oplus 1$ — это [[algebraic-immunity|алгебраическая иммунность]] $AI(f)$. Известно: $AI(f) \leq \lceil n/2 \rceil$.

## Реализация

- `boolfunc/annihilator.go:5` — `subsets(n, d)`: генерация всех подмножеств $\{0..n-1\}$ размера $\leq d$ (мономы степени $\leq d$).
- `boolfunc/annihilator.go:26` — `evalMonomial(x, subset)`: значение монома в точке $x$ (битовая маска).
- `boolfunc/metrics.go:93` — `hasAnnihilator(tv, n, d, complement)`: строит матрицу значений мономов в точках носителя $f$ (или $f \oplus 1$), проверяет ранг.
- `boolfunc/metrics.go:59` — `gaussianRank(matrix, cols)`: гауссов ранг матрицы над GF(2), используется как примитив поиска ядра.

Прямо «найти аннигилятор» (вектор коэффициентов $g$) — не реализовано; считается только наличие через $\mathrm{rank} < \mathrm{cols}$. Этого достаточно для вычисления [[algebraic-immunity|AI]].

## Источники

- `raw/coursework-1/content.tex:499–510`.
