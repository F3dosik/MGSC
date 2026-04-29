package sbox

import (
	"path/filepath"
	"testing"
)

// Эталонные значения из опубликованных источников:
// AES:  нелинейность = 112, алгебраическая степень = 7  (FIPS 197)

func loadSBox(t *testing.T, name string) *SBox {
	t.Helper()
	path := filepath.Join("..", "testdata", name)
	sb, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("не удалось загрузить %s: %v", name, err)
	}
	return sb
}

func TestNonlinearity(t *testing.T) {
	tests := []struct {
		file string
		want int
	}{
		{"aes_sbox.json", 112},
		{"kuznechik_sbox.json", 100},
		{"sm4_sbox.json", 112},
		{"camellia_sbox.json", 112},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			sb := loadSBox(t, tc.file)
			if got := sb.Nonlinearity(); got != tc.want {
				t.Errorf("Nonlinearity() = %d, хотим %d", got, tc.want)
			}
		})
	}

	t.Run("cached", func(t *testing.T) {
		sb := loadSBox(t, "aes_sbox.json")
		_ = sb.Nonlinearity()
		if sb.nonlinearity == nil {
			t.Error("nonlinearity не закэшировалась после первого вызова")
		}
	})
}

func TestAlgebraicDegree(t *testing.T) {
	tests := []struct {
		file string
		want int
	}{
		{"aes_sbox.json", 7},
		{"kuznechik_sbox.json", 7},
		{"sm4_sbox.json", 7},
		{"camellia_sbox.json", 7},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			sb := loadSBox(t, tc.file)
			if got := sb.AlgebraicDegree(); got != tc.want {
				t.Errorf("AlgebraicDegree() = %d, хотим %d", got, tc.want)
			}
		})
	}

	t.Run("cached", func(t *testing.T) {
		sb := loadSBox(t, "aes_sbox.json")
		_ = sb.AlgebraicDegree()
		if sb.algebraicDeg == nil {
			t.Error("algebraicDeg не закэшировалась после первого вызова")
		}
	})
}

func TestAlgebraicImmunity(t *testing.T) {
	tests := []struct {
		file string
		want int
	}{
		{"aes_sbox.json", 3},
		{"kuznechik_sbox.json", 3},
		{"sm4_sbox.json", 3},
		{"camellia_sbox.json", 3},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			sb := loadSBox(t, tc.file)
			if got := sb.AlgebraicImmunity(); got != tc.want {
				t.Errorf("AlgebraicImmunity() = %d, хотим %d", got, tc.want)
			}
		})
	}
}
