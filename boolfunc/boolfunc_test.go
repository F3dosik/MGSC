package boolfunc

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantN   int
		wantErr bool
	}{
		{"valid", "01101001", 3, false},
		{"invalid char", "01102010", 0, true},
		{"not power of 2", "010", 0, true},
		{"empty", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && f.N() != tt.wantN {
				t.Errorf("N() = %d, хотим %d", f.N(), tt.wantN)
			}
		})
	}
}

func TestFromSlice(t *testing.T) {
	tests := []struct {
		name    string
		input   []uint8
		wantN   int
		wantErr bool
	}{
		{"valid", []uint8{1, 0, 1, 0}, 2, false},
		{"invalid value", []uint8{1, 2, 0, 0}, 0, true},
		{"not power of 2", []uint8{1, 0, 1}, 0, true},
		{"empty", []uint8{}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := FromSlice(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && f.N() != tt.wantN {
				t.Errorf("N() = %d, хотим %d", f.N(), tt.wantN)
			}
		})
	}
	t.Run("isolates slice", func(t *testing.T) {
		tv := []uint8{0, 1, 1, 0}
		f, _ := FromSlice(tv)
		tv[0] ^= 1
		if f.At(0) != 0 {
			t.Error("FromSlice не изолирует слайс от внешних изменений")
		}
	})
}

func TestWeight(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"all zeros", "0000", 0},
		{"all ones", "1111", 4},
		{"mixed", "0110", 2},
		{"xor n=3", "01101001", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := New(tt.input)
			if got := f.Weight(); got != tt.want {
				t.Errorf("Weight(%q) = %d, хотим %d", tt.input, got, tt.want)
			}
		})
	}

	t.Run("check cached", func(t *testing.T) {
		f, _ := New("0110")
		_ = f.Weight()
		if f.weight == nil {
			t.Error("weight не закэшировался после первого вызова")
		}
	})
}

func TestIsBalanced(t *testing.T) {
	balanced, _ := New("0110")
	unbalanced, _ := New("0111")
	if !balanced.IsBalanced() {
		t.Error("0110 должна быть сбалансированной")
	}
	if unbalanced.IsBalanced() {
		t.Error("0111 не должна быть сбалансированной")
	}
}

func TestRandomBalanced(t *testing.T) {
	for range 100 {
		f, err := RandomBalanced(4)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if !f.IsBalanced() {
			t.Errorf("RandomBalanced вернула несбалансированную функцию: weight=%d, size=%d",
				f.Weight(), f.Size())
		}
	}
}

func TestValidateSize(t *testing.T) {
	for _, s := range []int{1, 2, 4, 8, 16, 256} {
		if err := validateSize(s); err != nil {
			t.Errorf("validateSize(%d) вернул ошибку: %v", s, err)
		}
	}
	for _, s := range []int{0, -1, 3, 5, 6, 7, 255} {
		if err := validateSize(s); err == nil {
			t.Errorf("validateSize(%d) должен вернуть ошибку", s)
		}
	}
}
