package sbox

import (
	"testing"
)

func TestDiffUniformity(t *testing.T) {
	tests := []struct {
		file string
		want int
	}{
		{"aes_sbox.json", 4},
		{"kuznechik_sbox.json", 8},
		{"sm4_sbox.json", 4},
		{"camellia_sbox.json", 4},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			sb := loadSBox(t, tc.file)
			if got := sb.DiffUniformity(); got != tc.want {
				t.Errorf("DiffUniformity() = %d, хотим %d", got, tc.want)
			}
		})
	}

	t.Run("cached", func(t *testing.T) {
		sb := loadSBox(t, "aes_sbox.json")
		_ = sb.DiffUniformity()
		if sb.diffUniformity == nil {
			t.Error("diffUniformity не закэшировалась после первого вызова")
		}
	})
}

func TestDDT_properties(t *testing.T) {
	sb := New(aesSBox)
	ddt := sb.DDT()

	t.Run("row a=0 is zero", func(t *testing.T) {
		for b := range 256 {
			if ddt[0][b] != 0 {
				t.Errorf("DDT[0][%d] = %d, хотим 0 (строка не заполняется)", b, ddt[0][b])
			}
		}
	})

	t.Run("each row a!=0 sums to 256", func(t *testing.T) {
		// для биективного S-блока каждая строка a≠0 суммируется в 256
		for a := 1; a < 256; a++ {
			sum := 0
			for b := range 256 {
				sum += int(ddt[a][b])
			}
			if sum != 256 {
				t.Errorf("DDT[%d] сумма = %d, хотим 256", a, sum)
			}
		}
	})

	t.Run("all values are even", func(t *testing.T) {
		// для биективного S-блока все значения DDT чётны
		for a := 1; a < 256; a++ {
			for b := range 256 {
				if ddt[a][b]%2 != 0 {
					t.Errorf("DDT[%d][%d] = %d нечётное", a, b, ddt[a][b])
				}
			}
		}
	})

	t.Run("cached", func(t *testing.T) {
		d1 := sb.DDT()
		d2 := sb.DDT()
		if d1 != d2 {
			t.Error("DDT должен возвращать закэшированный указатель")
		}
	})
}
