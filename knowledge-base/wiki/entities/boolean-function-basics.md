---
title: Булевы функции — основные понятия
type: entity
tags: [булевы-функции, справочник]
sources: [coursework-1]
created: 2026-05-15
updated: 2026-05-15
---

# Булевы функции — основные понятия

Справочная страница. Формальные детали и доказательства — в `raw/coursework-1/content.tex:68–169`.

## Определение

Булевой функцией от $n$ переменных называется $f: \mathbb{Z}_2^n \to \mathbb{Z}_2$. Множество всех таких функций обозначается $\mathcal{F}_n$, его мощность:
$$|\mathcal{F}_n| = 2^{2^n}.$$

Для $n=8$ имеем $|\mathcal{F}_8| = 2^{256}$.

## Вес и носитель

$$w(f) = |\{x : f(x) = 1\}| = \sum_{x \in \mathbb{Z}_2^n} f(x).$$

$$\mathrm{supp}(f) = \{x : f(x) = 1\}, \quad |\mathrm{supp}(f)| = w(f).$$

## Расстояние Хэмминга

$$d(f, g) = |\{x : f(x) \neq g(x)\}| = w(f \oplus g).$$

## Скалярное произведение, линейные и аффинные функции

Скалярное произведение векторов $a, x \in \mathbb{Z}_2^n$:
$$(a, x) = \bigoplus_{i=1}^n a_i x_i.$$

**Линейная** функция: $g(x) = (a, x)$ для некоторого $a$. Множество — $\mathcal{L}(n)$, мощность $2^n$.

**Аффинная**: $g(x) = a_0 \oplus (a, x)$. Множество — $\mathcal{A}(n) = \mathcal{L}(n) \cup \{1 \oplus \ell : \ell \in \mathcal{L}(n)\}$, мощность $2^{n+1}$.

## Формы представления

- **Таблица истинности** — массив длины $2^n$ значений. Используется как базовый формат.
- **Полярная таблица истинности** — заменяет 0 → +1, 1 → −1; вход для [[walsh-hadamard|ПУА]].
- **[[algebraic-normal-form|АНФ]] / полином Жегалкина** — рабочая форма для анализа алгебраической структуры.
- СДНФ/СКНФ — канонические, но в криптографии не применяются.

## Связанные страницы

**Метрики:** [[balancedness]], [[nonlinearity]], [[algebraic-degree]], [[correlation-immunity]], [[algebraic-immunity]].

**Преобразования:** [[walsh-hadamard]], [[mobius-transform]].

**Дифференциальные свойства:** [[directional-derivative]], [[differential-uniformity]].

## Реализация

- `boolfunc/boolfunc.go:15` — структура `BoolFunc` (`n`, `tv []uint8`, `tvp []uint64`, кэши метрик).
- Конструкторы: `boolfunc.New(string)` (`:30`), `boolfunc.FromSlice([]uint8)` (`:49`), `boolfunc.Random(n)` (`:64`), `boolfunc.RandomBalanced(n)` (`:78`).
- Геттеры: `N()`, `Size()`, `TV()`, `At(x)` (`:107–122`).
- Базовые метрики: `Weight()` (`:130`), `IsBalanced()` (`:125`).

## Источники

- `raw/coursework-1/content.tex:17–169`.
