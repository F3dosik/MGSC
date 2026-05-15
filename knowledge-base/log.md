# Log

Хронологический журнал базы знаний. Append-only. Каждая запись начинается с `## [YYYY-MM-DD] <тип> | <описание>`, чтобы парсилось через `grep "^## \[" log.md`.

Типы: `ingest` (новый источник), `query` (важный вопрос/синтез), `lint` (ревизия), `meta` (изменения схемы/структуры).

---

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
