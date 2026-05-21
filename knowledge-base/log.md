# Log

Хронологический журнал базы знаний. Append-only. Каждая запись начинается с `## [YYYY-MM-DD] <тип> | <описание>`, чтобы парсилось через `grep "^## \[" log.md`.

Типы: `ingest` (новый источник), `query` (важный вопрос/синтез), `lint` (ревизия), `meta` (изменения схемы/структуры).

---

## [2026-05-21] meta | E-10: SA из AES; оптимизация nlFast; финальная картина серии
- E-10: 16 параллельных SA из AES (MinMaxWalsh, T₀≈18) → NL=112 все 16 запусков, BestTable=AES. 64M swap-попыток, ни одного NL>112. AES — изолированный локальный оптимум.
- Обнаружена критическая проблема производительности: `sbox.New(s).Nonlinearity()` вызывал `computeMetrics()`, который вычислял Walsh И ANF (Möbius) для 255 компонент на каждый swap → двойная работа + огромное давление GC.
- Исправления: (1) `sbox/analysis.go` — разделены `computeNonlinearity()` и `computeAlgebraicDeg()`, `Nonlinearity()` больше не вычисляет ANF; (2) `generate/cost.go` — добавлен `nlFast()`: Walsh-преобразование без аллокаций (стековый массив), MinMaxWalsh и ThresholdNL переведены на `nlFast`.
- Финальный вывод серии E-01..E-10: два изолированных бассейна — NL=100 (из random) и NL=112 (из AES). Переход невозможен через случайные swap за разумное число итераций.
- Обновлена [[sa-experiments]] (E-10 завершён, добавлена итоговая таблица).

## [2026-05-21] meta | Реализован CLI инструмента (main.go с verify/analyze/generate)
- Реализован `main.go` в корне — единый CLI с тремя подкомандами. Используется только стандартная библиотека (пакет `flag`).
- `verify [-profile aes|relaxed|research] [-json] <file.json>` — оборачивает `verify.Verify`, exit=0 при PASS, exit=1 при FAIL.
- `analyze [-json] <file.json>` — печатает все 4 метрики + число неподвижных/антиподвижных точек, без оценки по порогам.
- `generate [флаги SA] [-out file.json]` — параметризованный запуск SA: `-cost minmax|threshold|cjs`, `-runs N` для мультистарта, автокалибровка T₀ через `EstimateT0` (если `-t0 0`). После завершения печатает результаты `verify` под AESLevel и RelaxedLevel.
- Добавлен `sbox.SaveToFile(sb, path, name)` — симметрично `LoadFromFile`, формат матрица 16×16. Round-trip JSON проверен smoke-тестом: `generate -out` → `verify` читает тот же файл.
- `go test ./...` зелёный по всем пакетам. Закрыто в [[backlog]]: CLI-задача.

## [2026-05-21] meta | Реализован модуль verify; найдена особенность SM4 (1 fixed point)
- Реализован `verify/verify.go`: тип `Threshold` с предустановленными профилями `AESLevel` (NL≥112, δ≤4, deg≥7, AI≥3, без фикс. и антиподвижных точек), `RelaxedLevel` (уровень «Кузнечика»: NL≥100, δ≤8), `ResearchLevel` (промежуточный, NL≥104, δ≤6). Функция `Verify(s, t) Report` возвращает структурированный отчёт с метриками и списком нарушений. Метод `Report.Summary()` для CLI-вывода.
- 8 тестов в `verify/verify_test.go`, все проходят (`go test ./...` зелёный): AES+Camellia принимаются строгим AESLevel; Кузнечик отвергается под AESLevel (нарушения по NL=100 и δ=8) и принимается под RelaxedLevel; identity (S(x)=x) отвергается через 256 фикс. точек; не-биективная подстановка отказывается сразу, метрики не вычисляются.
- **Содержательная находка:** SM4 имеет ровно одну неподвижную точку — конструктивная особенность шифра, отдельный тест `TestAESLevel_SM4_HasOneFixedPoint` фиксирует это явно. AES и Camellia таких точек не имеют. Под AESLevel без `NoFixedPoints` SM4 проходит.
- Материал для гл. 7.2 «Валидация на эталонных S-блоках»: таблица AES/Camellia/SM4/Кузнечик × AESLevel/RelaxedLevel — содержательнее простого «совпало с литературой».
- Закрыто в [[backlog]]: реализация `verify/verify.go`; косвенно — задача про прогон аналитики на эталонных S-блоках (теперь автоматический тест).

## [2026-05-21] meta | Завершены E-07..E-09c; зафиксирован потолок SA NL=100
- E-07: HC из SA(NL=100) → NL=100 без улучшений за 500K итераций. NL=100 — истинный локальный оптимум.
- E-08: multi-start N=16, MinMaxWalsh → NL=100 все 16 запусков. Нулевая дисперсия. Ландшафт имеет один глубокий аттрактор.
- E-09: CJS(X=36, R=4, T₀=20000 из MDPI 2022) → NL=88-96. T₀ несовместим с нашим масштабом cost (3×10¹⁰).
- E-09b: то же, калиброванный T₀=1.76×10⁸, α=0.95 → NL=96. Мало итераций (298K).
- E-09c: то же, α=0.99, M=5000 (fair vs E-05) → NL=96. **cost(AES) = cost(random) = 3.07×10¹⁰** — ландшафт инвертирован: AES не является минимумом CJS(X=36, R=4).
- **Итог серии:** одиночный SA не выходит за NL=100 ни с одной из 5 опробованных cost-функций. Лучшая — MinMaxWalsh (NL=100, нулевая дисперсия). Параметры MDPI 2022 не воспроизводимы нашей реализацией.
- Обновлена [[sa-experiments]] (E-07..E-09c завершены, сводная таблица, E-10 заблокирован).

## [2026-05-18] meta | Создана структура LaTeX-проекта 2-й курсовой в `report/`
- Скопирован полный шаблон Stulk3 из `knowledge-base/raw/coursework-1` (main.tex, Settings/packages.tex, Settings/format.tex, Settings/listings.tex, Settings/FiraCode-Regular.otf, macros.tex, tocpage.tex, refs.tex, bibpage.tex, TitlePages/). Компилятор — XeLaTeX (polyglossia, Times New Roman 14pt, ГОСТ-2021).
- Контент перенесён из `draft/coursework-2/content.tex` (1701 строка, 6 глав). `draft/coursework-2/` оставлен как backup.
- Точечные правки шаблона: в `Settings/format.tex` переименован `\newtheorem{theorem}{Утверждение}` → `Теорема`, добавлен `\newtheorem{statement}[theorem]{Утверждение}` (теперь оба окружения доступны). В `macros.tex` добавлено пользовательское окружение `\newenvironment{algorithm}[1]` (использовалось в content.tex для псевдокода HC/SA).
- В `content.tex` перед `СПИСКОМ ИСТОЧНИКОВ` вставлены каркасы: гл. 7 «Экспериментальные результаты» (7.1–7.7 + ВЫВОДЫ — план зафиксирован в memory `project-chapter-7-experiments`), ЗАКЛЮЧЕНИЕ, СПИСОК СОКРАЩЁННЫХ ОБОЗНАЧЕНИЙ. После списка источников — ПРИЛОЖЕНИЕ А (под `\codefromfile{...}` для Go-модулей).
- Прежний `report/` (мой pdflatex-каркас прошлой сессии) удалён.

## [2026-05-17] meta | Создан журнал экспериментов sa-experiments
- Создана [[sa-experiments]] — сквозной журнал всех SA-экспериментов с E-01 по E-10.
- Завершённые: E-01 (HC baseline), E-02 (CJS X=24 баг), E-03 (CJS X=0 неправильный ландшафт), E-04 (MinMaxWalsh короткое расписание), E-05 (MinMaxWalsh полное), E-06 (ThresholdNL полное). Лучший результат одиночного SA: NL=100.
- Запланированные: E-07 (HC из SA), E-08 (multi-start N=16), E-09 (CJS X=36 R=4 из MDPI 2022), E-10 (двухфазный SA+HC — главная цель).
- Обновлены [[index]].

## [2026-05-17] meta | Страница sa-for-sboxes; завершение структуры теории SA+HC

- Создана [[sa-for-sboxes]] — специфика SA для биективных 8×8 S-блоков: мутация swap (инволютивность, бесплатный откат), пространство поиска $256!$, анализ cost-функций (CJS vs threshold), сравнительная таблица NL из литературы, двухфазный SA+HC (путь к NL=112), параметры MDPI 2022 (X=36, R=4, α=0.95, success 56.4% для NL=104).
- Расширена [[hill-climbing]] — добавлена секция «HC как вторая фаза в SA+HC»: схема двухфазного алгоритма, алгоритм HC по NL (инволютивный откат), объяснение почему стартовая точка NL=104 ближе к бассейну NL=112.
- Обновлён [[index]] — добавлена [[sa-for-sboxes]].
- Обновлён [[backlog]]: закрыты задачи `generate/annealing.go` и базовые cost-функции; добавлена задача «[СЛЕДУЮЩИЙ ШАГ] Реализовать двухфазный SA+HC».
- Итог сессии: single SA достигает NL=100 независимо от cost-функции (MinMaxWalsh, ThresholdNL, CJS). Барьер принципиальный. Путь к NL=112 — SA+HC по CJS 2005 (Table III).

## [2026-05-15] meta | Создан каркас базы знаний
- Создана структура: `raw/`, `raw/assets/`, `wiki/entities/`.
- Созданы служебные файлы: [[CLAUDE]], [[index]], [[log]].
- Тематика: криптография — булевы функции и S-боксы (контекст курсовой MGSC).
- Категории wiki: пока только `entities/`. По мере роста можно добавить `concepts/`, `sources/`, `comparisons/`.

## [2026-05-15] meta | Multi-start и эмпирические наблюдения SA
- Создана [[multi-start]] — концепция multi-start метаэвристик, формула $1 - (1-p)^N$ для best-of-N, parallelization pattern в Go, сравнение с parallel tempering.
- Создана [[sa-empirical-observations]] — фиксация результатов калибровки: cost stability ~2.7e+08 для random, $T_0 \sim 1.5 \cdot 10^6$ из random vs $1.5 \cdot 10^7$ из AES, объяснение «локальный минимум даёт большие $|\Delta E|$», рекомендованные стартовые параметры SA.
- Пользователь реализовал `ParallelMultiStartSA` в `generate/annealing.go` по предложенному паттерну (sync.WaitGroup, math/rand/v2). Задача закрыта в [[backlog]].

## [2026-05-15] meta | Страница t0-calibration (калибровка начальной температуры)
- Создана [[t0-calibration]] — детальный разбор: зачем нужна калибровка, концепция acceptance rate $p_0$, вывод формулы $T_0 = -\langle\Delta E\rangle / \ln p_0$, алгоритм с инволютивным откатом, edge cases, альтернативы.
- Обновлён [[index]].

## [2026-05-17] meta | Эксперименты MinMaxWalsh и ThresholdNL; диагностика расписания SA
- Реализованы две cost-функции в `generate/cost.go`: `MinMaxWalsh()` (cost=256−2·NL) и `ThresholdNL(target)` (cost=max(0,target−NL)).
- Добавлен `OnTick` callback в `Schedule` для диагностики прогресса SA без изменения алгоритма.
- Эксперимент 1 (Tmin=1, M=1000): ~290K итераций, NL=100 — выявлено, что расписание заканчивается слишком быстро из-за малого T0 (целочисленная cost).
- Эксперимент 2 (Tmin=0.01, M=5000): ~3.7M итераций, NL=100 — SA замерзает при T<0.5, полезная работа прекращается на ~1.5M итераций.
- Установлено: NL=100 — типичный локальный оптимум для одиночного SA с биективных 8×8. Соответствует CJS 2005 (NL=102 из 200 запусков — редкость).
- Обновлена [[sa-empirical-observations]] — добавлен раздел с результатами и анализом замерзания.
- Следующий шаг: multi-start (N≥16) или SA+HC.

## [2026-05-16] meta | Страница threshold-based fitness
- Создана [[threshold-fitness]] — подробный разбор: идея max(0, порог−метрика), формула из тезисов, анализ ландшафта (F=0 в AES), проблема масштаба весов $w_i$, порядок вычисления с early exit, сравнительная таблица CJS vs threshold-fitness.
- Обновлён [[index]].

## [2026-05-16] query | CJS(X=0) ландшафт vs NL=112 — эмпирическое исследование
- Запущен SA с CJS(X=0, R=3), откалиброванным T0=8.54e5 (EstimateT0, p0=0.8), Alpha=0.99, M=1000.
- Результат: застрял на NL=98, cost 4.05e8 после 1.36M итераций. Улучшение только на 5%.
- Установлена структурная причина: CJS(X=0) = Σ|W|³ при фиксированном Σ|W|²=65536 (Парсеваль) оптимизируется на равномерном Walsh-спектре (~16), а не на разреженном (AES: 64×32, 191×0). cost(AES) ≈ 5.35e8 > cost(random) ≈ 4.28e8 — SA движется **от** AES, а не к нему.
- Вывод: для достижения NL=112 нужна threshold-based fitness F=max(0,112−NL) или двухфазный CJS+HC.
- Обновлена [[fitness-function-sbox]] — добавлена секция "Почему CJS не ведёт к NL=112".
- Параллельно: обнаружена ошибка X=24 в `generate/cost.go` (таргетировал Walsh≈24 вместо ≈0, та же структурная проблема). Исправлено в коде.

## [2026-05-15] meta | Страница SA (теория, термины) + hill-climbing стаб
- Создана [[simulated-annealing]] — общая теория SA: физическая метафора, термины (состояние, окрестность, температура, Metropolis, schedule, best-seen), алгоритм псевдокодом, гиперпараметры с эмпирическими правилами выбора, сравнение с HC.
- Создан стаб [[hill-climbing]] со связью с SA через $T=0$ и его ролью baseline в проекте.
- Источники: Kirkpatrick 1983 (основа SA), Russell-Norvig (HC).
- Специфика SA для S-блоков (мутация swap, cost CJS) — остаётся в [[fitness-function-sbox]] и в будущей странице sa-for-sboxes.
- Обновлён [[index]]; добавлено пользователем feedback-правило в memory: длинные объяснения с математикой писать в wiki, не в чат.

## [2026-05-15] ingest | ldc_tutorial.pdf (Heys 2002) → теория LC и DC
- Прочитан целиком (33 стр., 2 захода): Howard M. Heys "A Tutorial on Linear and Differential Cryptanalysis" (Memorial University of Newfoundland, 2002).
- Из tutorial взяли только S-блочную часть: §3.1 (linear probability bias), §3.2 (Piling-Up Lemma), §3.3 (LAT свойства), §3.6 (active S-boxes, $N_L \approx 1/\varepsilon^2$, linear hulls); §4.1 (idea DC), §4.2 (DDT свойства), §4.5 ($N_D \approx c/p_D$, differentials); §5 (advanced — упоминания higher-order, truncated, impossible).
- Пропущено: §2 (toy SPN), §3.4-3.5 и §4.3-4.4 (constructing/extracting на шифре).
- Создано 4 новые страницы:
  - [[differential-cryptanalysis]] — полная страница DC: идея, дифференциал, active S-boxes, complexity, расширения.
  - [[higher-order-differential-attack]] — стаб с минимумом теории (полный ingest Lai 1994 / Knudsen 1995 в backlog).
  - [[algebraic-attack]] — стаб (полный ingest Courtois-Meier 2003 в backlog).
  - [[attacks-overview]] ⭐ — **центральная страница**: сводная таблица 4 атак → 4 метрик → 4 свойств БФ, двойственности LC↔DC и higher-order↔algebraic, нюанс про $\delta$.
- Расширены: [[linear-cryptanalysis]] (piling-up, active S-boxes, complexity, linear hulls), [[linear-approximation-table]] (свойства таблицы), [[differential-uniformity]] (свойства DDT, связь с DC), [[algebraic-degree]] (связь с higher-order differential), [[algebraic-immunity]] (связь с алгебраической атакой).
- Обновлены: [[index]] (секция «Криптоанализ S-блоков»; ⭐ для attacks-overview).
- Понижена в [[backlog]] задача про ingest полной теории DC — базовый материал теперь есть.

## [2026-05-15] ingest | 2-я курсовая (content.tex + thesis_report_main.tex) → 12 страниц wiki
- Прочитаны: `content.tex` (469 стр., обзор 6 шифров + линейный криптоанализ), `thesis_report_main.tex` (167 стр., тезисный доклад).
- Очистка: удалены `tutorial.tex`, `README.md`, `Listings/`, `Images/`, `Settings/`, `TitlePages/`, `macros.tex`, `titlepage.pdf`, `tocpage.tex`, `bibpage.tex`, `refs.tex`. Оставлены `content.tex`, `thesis_report_main.tex`, `main.tex`, `refs.bib`. **Файл `ldc_tutorial.pdf` оставлен** — не идентифицирован, ждём пояснения пользователя.
- Создано 12 страниц:
  - Мета: [[thesis-statement]], [[fitness-function-sbox]].
  - Линейный криптоанализ: [[linear-cryptanalysis]], [[linear-approximation-table]] (с теоремой $L = 2^{n-1} - N$).
  - S-блок как вектор. функция: [[sbox-as-boolean]] (с фактом «биективность → сбалансированность компонент»).
  - Эталонные S-блоки: [[aes-sbox]], [[kuznechik-sbox]], [[serpent-sbox]], [[twofish-sbox]], [[sm4-sbox]], [[camellia-sbox]], [[reference-sboxes-metrics]] (сводная таблица).
- Обновлены: [[nonlinearity]] (TODO заменён на факт $N=116$ для биективных 8×8 со ссылкой на Carlet 2025), [[walsh-hadamard]] (связь с LAT), [[directional-derivative]] (связь с линейным криптоанализом).
- В [[backlog]] закрыты задачи: aes/kuznechik/sm4/camellia-sbox, sbox-as-boolean, LAT, thesis-statement, fitness, reference-sboxes-metrics. Hill climbing убран из приоритетов — фокус только SA.
- В [[index]] добавлены секции «Мета», «Линейный криптоанализ», «Эталонные S-блоки», «Генерация».
- Обновлена память: [[project-research-focus]] — фокус сужен до SA, добавлен таргет $AI \geq 3$.

## [2026-05-15] meta | Заход 2 — мост между wiki и Go-кодом + 2 новые страницы
- Прочитаны: `boolfunc/boolfunc.go`, `boolfunc/metrics.go`, `boolfunc/transforms.go`, `boolfunc/annihilator.go`, `sbox/ddt.go`, `sbox/algebraic.go`.
- Обнаружено: код реализует **алгебраическую иммунность** (`boolfunc/metrics.go:122`) и **дифференциальную равномерность** (`sbox/ddt.go:37`), хотя в 1-й курсовой эти темы формально не покрыты — «код опережает теорию».
- Созданы 2 новые страницы:
  - [[algebraic-immunity]] — определение, граница $\lceil n/2 \rceil$, алгоритм Мейера-Пасалича через гауссов ранг.
  - [[differential-uniformity]] — DDT, $\delta$, минимально достижимое значение для биективных S-блоков 8×8.
- Во все 12 страниц добавлены/уточнены секции «Реализация» со ссылками на конкретные функции и строки в Go-коде.
- Обновлён [[index]] (+2 страницы).
- В [[backlog]] задачи по AI и DDT понижены: код есть, нужен только ingest теоретического источника для отчёта.

## [2026-05-15] ingest | 1-я курсовая (content.tex) → 10 страниц wiki/entities/
- Прочитан `raw/coursework-1/content.tex` целиком (824 стр.).
- Создано 10 страниц с теорией:
  - Основы: [[boolean-function-basics]].
  - Формы представления: [[algebraic-normal-form]].
  - Преобразования: [[walsh-hadamard]], [[mobius-transform]].
  - Метрики: [[balancedness]], [[nonlinearity]], [[algebraic-degree]], [[correlation-immunity]].
  - Алгебраические атаки: [[annihilator]].
  - Дифференциальные свойства: [[directional-derivative]] (мост к будущему [[differential-uniformity]]).
- Не перенесено: вводный материал §1.1 (общая модель криптосистемы) и §3 (описание Python-класса BoolFunc) — нерелевантно текущему фокусу.
- Обнаруженные пробелы (записаны в [[backlog]]): алгебраическая иммунность, дифференциальная равномерность и DDT, специфика S-блоков (компонентные функции, агрегация), бент-функции, LAT, HC/SA — потребуют дополнительных источников.
- Обновлён [[index]].

## [2026-05-15] meta | Добавлена 1-я курсовая в raw/coursework-1/
- Пользователь скопировал исходники LaTeX-курсовой в `raw/coursework-1/`.
- Удалена обвязка LaTeX-шаблона (tutorial.tex, README.md, Settings/, TitlePages/, macros.tex, Listings/Main.java, titlepage.pdf, tocpage.tex, bibpage.tex, refs.tex).
- Осталось содержательное: `content.tex` (824 стр.), `main.tex`, `coursework-1.pdf`, `Images/*.png` (4 файла), `Listings/main.py`, `refs.bib`.
- Следующий шаг — ingest: извлечь теорию в `wiki/entities/` и библиографию в [[backlog]] как источники для дальнейшего чтения.

## [2026-05-15] meta | Зафиксирована роль LLM-научрука и фокус 2-й курсовой
- Создан корневой `/CLAUDE.md` с описанием роли (включение по команде) и фразами-триггерами.
- Уточнён фокус: **генерация S-блоков** (hill climbing, simulated annealing).
- Целевой ориентир: уровень AES (NL=112, DU=4, deg=7).
- Создан `backlog.md` с начальным списком задач (теория + код + эксперименты).
- Контекст текущего кода: `boolfunc/` и `sbox/` готовы и покрыты тестами; `generate/` и `verify/` — пустые заготовки.
