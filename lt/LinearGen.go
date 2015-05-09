package lt

import "math"
const (
	a = 16807 // primitive root
)

type LinearGen struct {
	seed uint32
}

func (gen *LinearGen) getSeed() uint32 {
	return gen.seed
}

func (gen *LinearGen) setSeed(seed uint32) {
	gen.seed = seed
}


func (gen *LinearGen) nextInt() uint32 {
	gen.seed = uint32(a * uint64(gen.seed) % math.MaxInt32)
	return gen.seed
}

func (gen *LinearGen) nextFloat() float64 {
	return float64(gen.nextInt())/float64(math.MaxInt32)
}

