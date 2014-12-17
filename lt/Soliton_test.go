package lt

import "testing"

type testCase struct {
	arr []float64; 
	arg float64;
	res uint;
}

func TestBinSearch(t *testing.T) {
	steps := func(n int) []float64 {
		ret := make([]float64, n+1)
		for i := 0; i <= n; i++ { 
			ret[i] = float64(i)/float64(n)
		}
		return ret
	}

	cases := []testCase {
		{steps(1), 0.0, 0},
		{steps(1), 0.5, 0},
		{steps(2), 0.0, 0},
		{steps(2), 0.49, 0},
		{steps(2), 0.5, 1},
		{steps(2), 0.99, 1},
		{steps(3), 0.2, 0},
		{steps(3), 0.50, 1},
		{steps(3), 0.80, 2},
		{steps(10), 0.01, 0},
		{steps(10), 0.14, 1},
		{steps(10), 0.31, 3},
		{steps(10), 0.81, 8},
		{steps(10), 0.90, 9},
		{steps(25), 0.10, 2},
		{steps(25), 0.14, 3},
		{steps(25), 0.25, 6},
		{steps(25), 0.50, 12},
		{steps(25), 0.52, 13},
	}
	for _, c := range(cases) {
		ans := binSearch(c.arr, c.arg)
		if ans != c.res {
			t.Error("binSearch(", c.arr, ",",c.arg,") =", ans, ", want", c.res)
		}
	}
}

