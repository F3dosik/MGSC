// Package verify реализует runtime-проверку S-блоков на соответствие
// заранее заданному профилю криптографических порогов.
//
// Назначение модуля — отделить корректность вычисления метрик (см. sbox)
// от вопроса «годен ли конкретный S-блок к применению». Verify принимает
// произвольную подстановку (например, результат генератора) и сравнивает
// её метрики с порогами выбранного профиля; возвращает структурированный
// отчёт со списком нарушений.
package verify

import (
	"fmt"
	"strings"

	"github.com/F3dosik/MGSC/sbox"
)

// Threshold задаёт набор пороговых требований к S-блоку.
// Нулевое значение поля MinNL/MinDeg/MinAI трактуется как «нет нижней
// границы»; нулевое MaxDU — как «нет верхней границы». Для типовых
// профилей см. AESLevel, RelaxedLevel, ResearchLevel.
type Threshold struct {
	Name string // человекочитаемое имя профиля

	MinNL  int // нижняя граница нелинейности
	MaxDU  int // верхняя граница дифференциальной равномерности
	MinDeg int // нижняя граница алгебраической степени
	MinAI  int // нижняя граница алгебраической иммунности

	RequireBijective bool // подстановка обязана быть биекцией
	NoFixedPoints    bool // не должно быть S(x) = x
	NoOppositeFixed  bool // не должно быть S(x) = ^x
}

// AESLevel — промышленный потолок: характеристики S-блоков AES, Camellia, SM4.
var AESLevel = Threshold{
	Name:             "AES",
	MinNL:            112,
	MaxDU:            4,
	MinDeg:           7,
	MinAI:            3,
	RequireBijective: true,
	NoFixedPoints:    true,
	NoOppositeFixed:  true,
}

// RelaxedLevel — уровень «Кузнечика»: то, что достижимо одиночным SA в наших экспериментах.
var RelaxedLevel = Threshold{
	Name:             "Relaxed (Кузнечик)",
	MinNL:            100,
	MaxDU:            8,
	MinDeg:           7,
	MinAI:            3,
	RequireBijective: true,
}

// ResearchLevel — целевой промежуточный профиль, ожидаемый выход двухфазного SA + HC.
var ResearchLevel = Threshold{
	Name:             "Research (SA+HC target)",
	MinNL:            104,
	MaxDU:            6,
	MinDeg:           6,
	MinAI:            3,
	RequireBijective: true,
}

// Report — результат верификации.
type Report struct {
	Profile string // имя профиля порогов

	Bijective     bool
	FixedPoints   []uint8
	OppositeFixed []uint8

	// Криптографические метрики; заполняются только если подстановка биективна.
	NL  int
	DU  int
	Deg int
	AI  int

	Passed   bool     // true, если ни один порог не нарушен
	Failures []string // список нарушений; пустой, если Passed
}

// Verify проверяет S-блок на соответствие профилю.
//
// Если профиль требует биективности и подстановка не биективна, проверка
// останавливается на этом шаге — метрики не вычисляются, поскольку их
// определения завязаны на биективность.
func Verify(s *sbox.SBox, t Threshold) Report {
	r := Report{Profile: t.Name}

	r.Bijective = s.IsBijective()
	if t.RequireBijective && !r.Bijective {
		r.Failures = append(r.Failures, "подстановка не биективна")
		return r
	}

	r.FixedPoints = s.FixedPoints()
	r.OppositeFixed = s.OppositeFixedPoints()
	r.NL = s.Nonlinearity()
	r.DU = s.DiffUniformity()
	r.Deg = s.AlgebraicDegree()
	r.AI = s.AlgebraicImmunity()

	if t.NoFixedPoints && len(r.FixedPoints) > 0 {
		r.Failures = append(r.Failures,
			fmt.Sprintf("найдено %d неподвижных точек", len(r.FixedPoints)))
	}
	if t.NoOppositeFixed && len(r.OppositeFixed) > 0 {
		r.Failures = append(r.Failures,
			fmt.Sprintf("найдено %d антиподвижных точек", len(r.OppositeFixed)))
	}
	if t.MinNL > 0 && r.NL < t.MinNL {
		r.Failures = append(r.Failures,
			fmt.Sprintf("NL=%d < %d", r.NL, t.MinNL))
	}
	if t.MaxDU > 0 && r.DU > t.MaxDU {
		r.Failures = append(r.Failures,
			fmt.Sprintf("δ=%d > %d", r.DU, t.MaxDU))
	}
	if t.MinDeg > 0 && r.Deg < t.MinDeg {
		r.Failures = append(r.Failures,
			fmt.Sprintf("deg=%d < %d", r.Deg, t.MinDeg))
	}
	if t.MinAI > 0 && r.AI < t.MinAI {
		r.Failures = append(r.Failures,
			fmt.Sprintf("AI=%d < %d", r.AI, t.MinAI))
	}

	r.Passed = len(r.Failures) == 0
	return r
}

// Summary возвращает однострочное человекочитаемое резюме отчёта —
// удобно для CLI-вывода и логов.
func (r Report) Summary() string {
	if !r.Bijective {
		return fmt.Sprintf("[%s] FAIL: подстановка не биективна", r.Profile)
	}
	status := "PASS"
	if !r.Passed {
		status = "FAIL"
	}
	body := fmt.Sprintf("[%s] %s: NL=%d δ=%d deg=%d AI=%d",
		r.Profile, status, r.NL, r.DU, r.Deg, r.AI)
	if !r.Passed {
		body += " — " + strings.Join(r.Failures, "; ")
	}
	return body
}
