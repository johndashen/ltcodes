package lt

type RandGen interface {
	getSeed() uint32
	setSeed(uint32)
	nextInt() uint32
	nextFloat() float64	
}
