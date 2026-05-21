#!/usr/bin/env python3
"""Собирает HTML-документ приложения А с подсветкой Go через pygments."""

import subprocess
from pathlib import Path

HERE = Path(__file__).parent
OUT = HERE / "appendix.html"

SECTIONS = [
    ("Модуль булевых функций (boolfunc)", [
        ("Преобразование Мёбиуса и Уолша–Адамара (transforms.go)", "transforms.go"),
        ("Метрики булевых функций (metrics.go)", "metrics.go"),
        ("Поиск аннигилятора (annihilator.go)", "annihilator.go"),
    ]),
    ("Модуль S-блока (sbox)", [
        ("Тип SBox и базовые операции (sbox.go)", "sbox.go"),
        ("Агрегированные метрики через компонентные функции (analysis.go)", "analysis.go"),
        ("Таблица разностного распределения (ddt.go)", "ddt.go"),
        ("Алгебраическая иммунность S-блока (algebraic.go)", "algebraic.go"),
    ]),
    ("Модуль генерации (generate)", [
        ("Состояние SA и инволютивная транспозиция (state.go)", "state.go"),
        ("Функции штрафов (cost.go)", "cost.go"),
        ("Имитация отжига и мультистарт (annealing.go)", "annealing.go"),
        ("Восхождение к вершине (hillclimb.go)", "hillclimb.go"),
    ]),
    ("Модуль верификации (verify)", [
        ("Профили порогов, верификация и отчёт (verify.go)", "verify.go"),
    ]),
    ("Точка входа (main.go)", [
        ("CLI с подкомандами verify, analyze, generate (main.go)", "main.go"),
    ]),
]

INTRO = (
    "В данном приложении приведены ключевые исходные файлы инструмента, "
    "реализующие алгоритмы и структуры данных, описанные в главах 5 и 6. "
    "Файлы расположены в порядке зависимостей: модуль булевых функций, "
    "модуль S-блока, модуль генерации, модуль верификации, точка входа CLI. "
    "Пакеты с тестами и служебная инфраструктура (I/O в JSON, "
    "конструкторы-обёртки) в приложение не вынесены — их описание есть в "
    "тексте главы 6."
)

CSS_STATIC = """
@page { margin: 1.5cm; size: A4; }
body { font-family: 'DejaVu Serif', serif; font-size: 11pt; line-height: 1.4; }
h1 { font-size: 17pt; text-align: center; margin-top: 0; line-height: 1.3; }
h2 { font-size: 13pt; margin-top: 2em;
     border-bottom: 1px solid #ccc; padding-bottom: 0.3em; page-break-after: avoid; }
h3 { font-size: 11pt; margin-top: 1.2em; font-weight: bold; page-break-after: avoid; }
.highlight { background: #f8f8f8; padding: 6px 10px;
             border: 1px solid #ddd; border-radius: 3px; }
.highlight pre { margin: 0; font-family: 'DejaVu Sans Mono', monospace;
                 font-size: 8.5pt; line-height: 1.25;
                 white-space: pre-wrap; word-break: break-word; }
.intro { text-align: justify; margin: 1em 0 2em; }
"""


def pygmentize_fragment(path: Path) -> str:
    return subprocess.run(
        ["pygmentize", "-f", "html", "-O", "linenos=false", "-l", "go", str(path)],
        capture_output=True, text=True, check=True,
    ).stdout


def pygmentize_css() -> str:
    return subprocess.run(
        ["pygmentize", "-f", "html", "-S", "friendly", "-a", ".highlight"],
        capture_output=True, text=True, check=True,
    ).stdout


def build() -> str:
    parts = [
        '<!DOCTYPE html>',
        '<html lang="ru"><head><meta charset="utf-8">',
        '<title>Приложение А — программная реализация</title>',
        '<style>', CSS_STATIC, pygmentize_css(), '</style>',
        '</head><body>',
        '<h1>Приложение А.<br>'
        'Программная реализация инструмента анализа и генерации S-блоков</h1>',
        f'<p class="intro">{INTRO}</p>',
    ]
    for section_title, files in SECTIONS:
        parts.append(f'<h2>{section_title}</h2>')
        for sub_title, filename in files:
            parts.append(f'<h3>{sub_title}</h3>')
            parts.append(pygmentize_fragment(HERE / filename))
    parts.append('</body></html>')
    return '\n'.join(parts)


if __name__ == "__main__":
    OUT.write_text(build(), encoding="utf-8")
    print(f"wrote {OUT}")
