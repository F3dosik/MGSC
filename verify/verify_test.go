package verify

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/F3dosik/MGSC/sbox"
)

func loadSBox(t *testing.T, name string) *sbox.SBox {
	t.Helper()
	path := filepath.Join("..", "testdata", name)
	sb, err := sbox.LoadFromFile(path)
	if err != nil {
		t.Fatalf("не удалось загрузить %s: %v", name, err)
	}
	return sb
}

// identitySBox — тождественная подстановка S(x) = x.
// Биективна, но все 256 точек неподвижны.
func identitySBox(t *testing.T) *sbox.SBox {
	t.Helper()
	var arr [256]uint8
	for i := range arr {
		arr[i] = uint8(i)
	}
	return sbox.New(arr)
}

// nonBijectiveSBox — подстановка с повторами: S(0)=S(1)=0, остальное как identity.
func nonBijectiveSBox(t *testing.T) *sbox.SBox {
	t.Helper()
	var arr [256]uint8
	for i := range arr {
		arr[i] = uint8(i)
	}
	arr[1] = 0
	return sbox.New(arr)
}

// ── Профиль AESLevel: эталонные шифры ───────────────────────────────────────

// AES и Camellia спроектированы без неподвижных и антиподвижных точек —
// проходят строгий AESLevel целиком.
func TestAESLevel_AcceptsAESAndCamellia(t *testing.T) {
	for _, name := range []string{"aes_sbox.json", "camellia_sbox.json"} {
		t.Run(name, func(t *testing.T) {
			sb := loadSBox(t, name)
			r := Verify(sb, AESLevel)
			if !r.Passed {
				t.Errorf("ожидалось Passed=true, получено: %s", r.Summary())
			}
		})
	}
}

// SM4 имеет ровно одну неподвижную точку по конструкции шифра.
// Это документированное расхождение с AES/Camellia — фиксируем его явно.
// Под профилем без NoFixedPoints SM4 проходит по всем четырём метрикам.
func TestAESLevel_SM4_HasOneFixedPoint(t *testing.T) {
	sb := loadSBox(t, "sm4_sbox.json")

	rStrict := Verify(sb, AESLevel)
	if rStrict.Passed {
		t.Errorf("под строгим AESLevel SM4 не должен проходить (одна fixed point): %s",
			rStrict.Summary())
	}
	if len(rStrict.FixedPoints) != 1 {
		t.Errorf("ожидалась ровно 1 неподвижная точка для SM4, найдено %d",
			len(rStrict.FixedPoints))
	}

	relaxedFP := AESLevel
	relaxedFP.NoFixedPoints = false
	rRelaxed := Verify(sb, relaxedFP)
	if !rRelaxed.Passed {
		t.Errorf("SM4 обязан проходить AESLevel без NoFixedPoints: %s", rRelaxed.Summary())
	}
}

func TestAESLevel_RejectsKuznechik(t *testing.T) {
	sb := loadSBox(t, "kuznechik_sbox.json")
	r := Verify(sb, AESLevel)
	if r.Passed {
		t.Fatalf("ожидался FAIL для «Кузнечика» под AESLevel, получено PASS: %s", r.Summary())
	}
	if !hasFailure(r, "NL=") {
		t.Errorf("ожидалось нарушение по NL; failures=%v", r.Failures)
	}
	if !hasFailure(r, "δ=") {
		t.Errorf("ожидалось нарушение по δ; failures=%v", r.Failures)
	}
}

// ── Профиль RelaxedLevel ────────────────────────────────────────────────────

func TestRelaxedLevel_AcceptsKuznechik(t *testing.T) {
	sb := loadSBox(t, "kuznechik_sbox.json")
	r := Verify(sb, RelaxedLevel)
	if !r.Passed {
		t.Errorf("ожидалось Passed=true для «Кузнечика» под RelaxedLevel: %s", r.Summary())
	}
}

func TestRelaxedLevel_AcceptsAES(t *testing.T) {
	sb := loadSBox(t, "aes_sbox.json")
	r := Verify(sb, RelaxedLevel)
	if !r.Passed {
		t.Errorf("AES обязан проходить более слабый профиль RelaxedLevel: %s", r.Summary())
	}
}

// ── Структурные нарушения ───────────────────────────────────────────────────

func TestIdentity_FailsByFixedPoints(t *testing.T) {
	sb := identitySBox(t)
	r := Verify(sb, AESLevel)
	if r.Passed {
		t.Fatalf("identity не должна проходить AESLevel: %s", r.Summary())
	}
	if !r.Bijective {
		t.Error("identity всё же биективна — Bijective должно быть true")
	}
	if len(r.FixedPoints) != 256 {
		t.Errorf("ожидалось 256 неподвижных точек, найдено %d", len(r.FixedPoints))
	}
	if !hasFailure(r, "неподвижных точек") {
		t.Errorf("ожидалось нарушение по NoFixedPoints; failures=%v", r.Failures)
	}
}

func TestNonBijective_FailsImmediately(t *testing.T) {
	sb := nonBijectiveSBox(t)
	r := Verify(sb, AESLevel)
	if r.Passed {
		t.Fatalf("не-биективная подстановка не должна проходить: %s", r.Summary())
	}
	if r.Bijective {
		t.Error("Bijective должно быть false")
	}
	// Метрики не должны вычисляться при отказе по биективности.
	if r.NL != 0 || r.DU != 0 || r.Deg != 0 || r.AI != 0 {
		t.Errorf("метрики не должны заполняться при FAIL по биективности: NL=%d δ=%d deg=%d AI=%d",
			r.NL, r.DU, r.Deg, r.AI)
	}
	if len(r.Failures) != 1 || !strings.Contains(r.Failures[0], "биектив") {
		t.Errorf("ожидалось одно нарушение по биективности, получено: %v", r.Failures)
	}
}

// ── Содержательный контент отчёта ───────────────────────────────────────────

func TestReport_FillsMetricsForBijective(t *testing.T) {
	sb := loadSBox(t, "aes_sbox.json")
	r := Verify(sb, AESLevel)
	if r.NL != 112 || r.DU != 4 || r.Deg != 7 || r.AI != 3 {
		t.Errorf("AES: ожидалось NL=112 δ=4 deg=7 AI=3, получено NL=%d δ=%d deg=%d AI=%d",
			r.NL, r.DU, r.Deg, r.AI)
	}
}

func TestSummary_Format(t *testing.T) {
	sb := loadSBox(t, "aes_sbox.json")
	r := Verify(sb, AESLevel)
	s := r.Summary()
	if !strings.Contains(s, "PASS") || !strings.Contains(s, "AES") {
		t.Errorf("Summary должен содержать PASS и имя профиля: %q", s)
	}

	sb = loadSBox(t, "kuznechik_sbox.json")
	r = Verify(sb, AESLevel)
	s = r.Summary()
	if !strings.Contains(s, "FAIL") {
		t.Errorf("Summary для отказавшего отчёта должен содержать FAIL: %q", s)
	}
}

// ── helpers ─────────────────────────────────────────────────────────────────

func hasFailure(r Report, substr string) bool {
	for _, f := range r.Failures {
		if strings.Contains(f, substr) {
			return true
		}
	}
	return false
}
