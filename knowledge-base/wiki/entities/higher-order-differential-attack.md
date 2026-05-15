---
title: Атака высоких порядков (Higher-order differential)
type: entity
tags: [криптоанализ, s-блоки, дифференциальный, алгебраическая-степень]
sources: [coursework-2, ldc-tutorial-heys]
created: 2026-05-15
updated: 2026-05-15
---

# Атака высоких порядков (Higher-order differential)

> [!todo] Стаб-страница
> Полная теория требует отдельного ingest. Источники в библиографии работ: Lai X. *Higher Order Derivatives and Differential Cryptanalysis* (1994); Knudsen L.R. *Truncated and Higher Order Differentials* // FSE 1995 (LNCS 1008).

## Идея

Обобщение [[differential-cryptanalysis|обычного дифференциального криптоанализа]]. Вместо одной производной по направлению $S'_a$ рассматриваются **производные высоких порядков**:
$$S'_{a_1, a_2, \dots, a_k}(x) = \Delta_{a_1} \Delta_{a_2} \cdots \Delta_{a_k} S(x),$$
где $\Delta_a f(x) = f(x \oplus a) \oplus f(x)$ — производная по направлению.

## Связь со свойством S-блока — [[algebraic-degree|алгебраическая степень]]

**Ключевой факт:** производная порядка $k$ функции степени $d$ имеет степень не выше $d - k$. В частности:
- $k > d$ → производная **тождественно равна нулю**.
- $k = d$ → производная — константа.

**Следствие:** если $\deg(S) = d$, то существует «бесплатная» производная порядка $d + 1$, которая всегда даёт известную фиксированную разность. Это и есть основа атаки: атакующий выбирает $2^{d+1}$ открытых текстов специальной структуры и получает предсказуемое поведение шифра.

## Критерий стойкости

$\deg(S)$ должна быть как можно выше. Для биективных S-блоков 8×8 максимум $\deg = 7$ (достигается AES, SM4, Camellia, Кузнечик — см. [[reference-sboxes-metrics]]).

Целевой порог 2-й курсовой ([[thesis-statement]]): $\deg \geq 7$.

## Связь с другими атаками

Эта атака исторически вытекает из дифференциального криптоанализа, но защищается **другим свойством S-блока** ($\deg$, не $\delta$). Это одна из четырёх атак, на которые ориентируется работа — см. [[attacks-overview]].

В туториале Heys 2002 атака упомянута только в §5 (Advanced Concepts) одним абзацем — формальная теория требует Lai 1994 / Knudsen 1995.

## Источники

- `raw/coursework-2/ldc_tutorial.pdf:30` — упоминание в §5.
- Lai X. Higher Order Derivatives and Differential Cryptanalysis // Communications and Cryptography. Springer, 1994. P. 227–233.
- Knudsen L. R. Truncated and Higher Order Differentials // FSE 1995, LNCS 1008. P. 196–211.
