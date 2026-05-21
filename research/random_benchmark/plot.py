"""Сводка и гистограммы для random_benchmark.

Запуск:
    research/venv/bin/python3 research/random_benchmark/plot.py
"""

import sys
from pathlib import Path

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd

HERE = Path(__file__).parent
CSV_PATH = HERE / "benchmark.csv"
PNG_PATH = HERE / "random_benchmark.png"


def percentiles(series: pd.Series) -> dict:
    return {
        "min":  int(series.min()),
        "p5":   int(np.percentile(series, 5)),
        "p25":  int(np.percentile(series, 25)),
        "med":  int(np.percentile(series, 50)),
        "p75":  int(np.percentile(series, 75)),
        "p95":  int(np.percentile(series, 95)),
        "max":  int(series.max()),
        "mean": float(series.mean()),
    }


def print_summary(df: pd.DataFrame) -> None:
    n = len(df)
    print(f"\n=== Сводка по {n} случайным биективным S-блокам ===\n")
    print(f"{'':6} {'min':>4} {'p5':>4} {'p25':>4} {'med':>4} {'p75':>4} {'p95':>4} {'max':>4} {'mean':>7}")
    for name, col in [
        ("NL",     "nl"),
        ("δ",      "du"),
        ("deg",    "deg"),
        ("AI",     "ai"),
        ("|Фт|",   "fixed"),
        ("|АФт|",  "opposite"),
    ]:
        p = percentiles(df[col])
        print(f"{name:6} {p['min']:>4} {p['p5']:>4} {p['p25']:>4} {p['med']:>4} {p['p75']:>4} {p['p95']:>4} {p['max']:>4} {p['mean']:>7.2f}")

    print()
    no_fixed = int((df["fixed"] == 0).sum())
    no_opp = int((df["opposite"] == 0).sum())
    pass_aes = int(df["passAES"].sum())
    pass_rel = int(df["passRel"].sum())
    pass_res = int(df["passRes"].sum())

    def pct(x: int) -> str:
        return f"{100 * x / n:.1f}%"

    print(f"Без неподвижных точек:    {no_fixed} / {n} ({pct(no_fixed)})   теор. 1/e ≈ 36.8%")
    print(f"Без антиподвижных точек:  {no_opp} / {n} ({pct(no_opp)})")
    print()
    print(f"Под AESLevel:             {pass_aes} / {n} ({pct(pass_aes)})")
    print(f"Под ResearchLevel:        {pass_res} / {n} ({pct(pass_res)})")
    print(f"Под RelaxedLevel:         {pass_rel} / {n} ({pct(pass_rel)})")


def plot_histograms(df: pd.DataFrame, out_path: Path) -> None:
    fig, axes = plt.subplots(2, 2, figsize=(12, 8))

    def draw(ax, col, title, xlabel, aes_value=None, color="steelblue"):
        values = df[col]
        bins = np.arange(values.min(), values.max() + 2)
        ax.hist(values, bins=bins, align="left", rwidth=0.85,
                color=color, edgecolor="black", linewidth=0.5)
        if aes_value is not None:
            ax.axvline(aes_value, color="red", linestyle="--", linewidth=1.5,
                       label=f"AES = {aes_value}")
            ax.legend(loc="best", fontsize=9)
        ax.set_xlabel(xlabel)
        ax.set_ylabel("Число подстановок")
        ax.set_title(f"{title}  (N={len(df)})")
        ax.grid(True, alpha=0.3)

    draw(axes[0, 0], "nl",  "Нелинейность",                   "NL",  aes_value=112, color="steelblue")
    draw(axes[0, 1], "du",  "Дифференциальная равномерность", "δ",   aes_value=4,   color="darkorange")
    draw(axes[1, 0], "deg", "Алгебраическая степень",         "deg",                color="seagreen")
    draw(axes[1, 1], "ai",  "Алгебраическая иммунность",      "AI",                 color="mediumpurple")

    plt.tight_layout()
    plt.savefig(out_path, dpi=120, bbox_inches="tight")
    print(f"\nГистограммы сохранены: {out_path}")


def main() -> None:
    if not CSV_PATH.exists():
        print(f"Не найден {CSV_PATH}. Сначала запусти бенчмарк:")
        print("  go run ./research/random_benchmark/ -n 1000 -out research/random_benchmark/benchmark.csv")
        sys.exit(1)

    df = pd.read_csv(CSV_PATH)
    print_summary(df)
    plot_histograms(df, PNG_PATH)


if __name__ == "__main__":
    main()
