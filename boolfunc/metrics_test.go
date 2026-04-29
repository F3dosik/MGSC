package boolfunc

import "testing"

func TestNonlinearity(t *testing.T) {
	tests := []struct {
		name string
		tv   string
		want int
	}{
		{"constant 0 n=2", "0000", 0},
		{"XOR n=2", "0110", 0},
		{"AND n=2", "0001", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := New(tt.tv)
			if got := f.Nonlinearity(); got != tt.want {
				t.Errorf("Nonlinearity() = %d, хотим %d", got, tt.want)
			}
		})
	}

	t.Run("cached", func(t *testing.T) {
		f, _ := New("0001")
		_ = f.Nonlinearity()
		if f.nlin == nil {
			t.Error("nlin не закэшировался после первого вызова")
		}
	})
}

func TestAlgebraicDegree(t *testing.T) {
	tests := []struct {
		name string
		tv   string
		want int
	}{
		{"constant 0", "0000", 0},
		{"x2 only", "0101", 1},
		{"XOR n=2", "0110", 1},
		{"AND n=2", "0001", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := New(tt.tv)
			if got := f.AlgebraicDegree(); got != tt.want {
				t.Errorf("AlgebraicDegree() = %d, хотим %d", got, tt.want)
			}
		})
	}

	t.Run("cached", func(t *testing.T) {
		f, _ := New("0001")
		_ = f.AlgebraicDegree()
		if f.degree == nil {
			t.Error("degree не закэшировался после первого вызова")
		}
	})
}

func TestAlgebraicImmunity_Zero(t *testing.T) {
	// f = 0 везде
	f, err := FromSlice(make([]uint8, 256))
	if err != nil {
		t.Fatal(err)
	}
	if got := f.AlgebraicImmunity(); got != 1 {
		t.Errorf("zero function: got AI=%d, want 1", got)
	}
}

func TestAlgebraicImmunity_One(t *testing.T) {
	// f = 1 везде
	tv := make([]uint8, 256)
	for i := range tv {
		tv[i] = 1
	}
	f, err := FromSlice(tv)
	if err != nil {
		t.Fatal(err)
	}
	if got := f.AlgebraicImmunity(); got != 1 {
		t.Errorf("one function: got AI=%d, want 1", got)
	}
}

func TestAlgebraicImmunity_Linear(t *testing.T) {
	// f = x0
	tv := make([]uint8, 256)
	for x := range tv {
		tv[x] = uint8(x) & 1
	}
	f, err := FromSlice(tv)
	if err != nil {
		t.Fatal(err)
	}
	if got := f.AlgebraicImmunity(); got != 1 {
		t.Errorf("linear x0: got AI=%d, want 1", got)
	}
}

func TestAlgebraicImmunity_Cache(t *testing.T) {
	tv := make([]uint8, 256)
	for x := range tv {
		tv[x] = uint8(x) & 1
	}
	f, err := FromSlice(tv)
	if err != nil {
		t.Fatal(err)
	}
	first := f.AlgebraicImmunity()
	second := f.AlgebraicImmunity()
	if first != second {
		t.Errorf("cache inconsistency: first=%d second=%d", first, second)
	}
}

func TestAlgebraicImmunity_SmallLinear(t *testing.T) {
	// n=2, f = x0: tv = "0101"
	// x=00 -> 0, x=01 -> 1, x=10 -> 0, x=11 -> 1
	f, err := New("0101")
	if err != nil {
		t.Fatal(err)
	}
	if got := f.AlgebraicImmunity(); got != 1 {
		t.Errorf("n=2 linear x0: got AI=%d, want 1", got)
	}
}
