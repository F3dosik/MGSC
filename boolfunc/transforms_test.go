package boolfunc

import (
	"slices"
	"testing"
)

// Эталонные пары tv → anf для ручной проверки:
//
// f(x) = x1 XOR x2, n=2
//   tv  = [0,1,1,0]  (f(00)=0, f(01)=1, f(10)=1, f(11)=0)
//   anf = [0,1,1,0]  — мономы: x1 + x2  (без x1·x2)
//
// f(x) = x1 AND x2, n=2
//   tv  = [0,0,0,1]  (f(00)=0, f(01)=0, f(10)=0, f(11)=1)
//   anf = [0,0,0,1]  — моном: x1·x2
//
// f(x) = 1 (константа), n=2
//   tv  = [1,1,1,1]
//   anf = [1,0,0,0]  — только свободный член
//
// f(x) = x1, n=2
//   tv  = [0,0,1,1]
//   anf = [0,0,1,0]  — моном: x1

func TestANF(t *testing.T) {
	tests := []struct {
		name    string
		tv      string
		wantANF []uint8
	}{
		{
			name:    "XOR n=2",
			tv:      "0110",
			wantANF: []uint8{0, 1, 1, 0},
		},
		{
			name:    "AND n=2",
			tv:      "0001",
			wantANF: []uint8{0, 0, 0, 1},
		},
		{
			name:    "constant 1 n=2",
			tv:      "1111",
			wantANF: []uint8{1, 0, 0, 0},
		},
		{
			name:    "constant 0 n=2",
			tv:      "0000",
			wantANF: []uint8{0, 0, 0, 0},
		},
		{
			name:    "x1 only n=2",
			tv:      "0011",
			wantANF: []uint8{0, 0, 1, 0},
		},
		{
			name:    "x2 only n=2",
			tv:      "0101",
			wantANF: []uint8{0, 1, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(tt.tv)
			if err != nil {
				t.Fatalf("New(%q): %v", tt.tv, err)
			}
			got := f.ANF()
			if len(got) != len(tt.wantANF) {
				t.Fatalf("len(ANF()) = %d, хотим %d", len(got), len(tt.wantANF))
			}
			for i := range got {
				if got[i] != tt.wantANF[i] {
					t.Errorf("ANF()[%d] = %d, хотим %d", i, got[i], tt.wantANF[i])
				}
			}
		})
	}

	t.Run("involution", func(t *testing.T) {
		// Преобразование Мёбиуса — инволюция: M(M(tv)) == tv
		f, _ := New("01101001")
		anf := f.ANF()
		g, _ := FromSlice(anf)
		tvBack := g.ANF()
		original := f.TV()
		for i := range original {
			if tvBack[i] != original[i] {
				t.Errorf("M(M(tv))[%d] = %d, хотим %d", i, tvBack[i], original[i])
			}
		}
	})

	t.Run("cached", func(t *testing.T) {
		f, _ := New("0110")
		_ = f.ANF()
		if f.anf == nil {
			t.Error("anf не закэшировался после первого вызова")
		}
	})

	t.Run("isolates cache", func(t *testing.T) {
		f, _ := New("0110")
		anf := f.ANF()
		anf[0] ^= 1 // меняем возвращённую копию
		if f.anf[0] != f.ANF()[0] {
			t.Error("ANF() не защищает кэш от внешних изменений")
		}
	})
}

// Эталонные значения спектра Уолша для ручной проверки:
//
// f = 0 (константа), n=2: sv = [1,1,1,1] -> W = [4,0,0,0]
// f = 1 (константа), n=2: sv = [-1,-1,-1,-1] -> W = [-4,0,0,0]
// f = x1 XOR x2, n=2:     sv = [1,-1,-1,1]  -> W = [0,0,0,4]
// f сбалансирована -> W[0] == 0

func TestWalsh(t *testing.T) {
	tests := []struct {
		name      string
		tv        string
		wantWalsh []int32
	}{
		{
			// sv = [1,1,1,1], W[0]=4, W[1..3]=0
			name:      "constant 0",
			tv:        "0000",
			wantWalsh: []int32{4, 0, 0, 0},
		},
		{
			// sv = [-1,-1,-1,-1], W[0]=-4
			name:      "constant 1",
			tv:        "1111",
			wantWalsh: []int32{-4, 0, 0, 0},
		},
		{
			// f=x2 (младший бит): tv=[0,1,0,1], sv=[1,-1,1,-1]
			// W[0]=0, W[1]=4, W[2]=0, W[3]=0
			name:      "x2 only",
			tv:        "0101",
			wantWalsh: []int32{0, 4, 0, 0},
		},
		{
			// f=x1 XOR x2: tv=[0,1,1,0], sv=[1,-1,-1,1]
			// W[0]=0, W[1]=0, W[2]=0, W[3]=4
			name:      "XOR n=2",
			tv:        "0110",
			wantWalsh: []int32{0, 0, 0, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(tt.tv)
			if err != nil {
				t.Fatalf("New(%q): %v", tt.tv, err)
			}
			got := f.Walsh()
			if !slices.Equal(got, tt.wantWalsh) {
				t.Errorf("Walsh() = %v, хотим %v", got, tt.wantWalsh)
			}
		})
	}

	t.Run("balanced has W[0]=0", func(t *testing.T) {
		// Любая сбалансированная функция имеет W[0] = 0
		f, _ := New("0110")
		if f.Walsh()[0] != 0 {
			t.Errorf("сбалансированная функция: Walsh[0] = %d, хотим 0", f.Walsh()[0])
		}
	})

	t.Run("cached", func(t *testing.T) {
		f, _ := New("0110")
		_ = f.Walsh()
		if f.walsh == nil {
			t.Error("walsh не закэшировался после первого вызова")
		}
	})

	t.Run("isolates cache", func(t *testing.T) {
		f, _ := New("0110")
		w := f.Walsh()
		w[0] = 999
		if f.walsh[0] == 999 {
			t.Error("Walsh() не защищает кэш от внешних изменений")
		}
	})
}
