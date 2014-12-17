package lt

type BlockPlanner struct {
	nblocks uint32
	rando RandGen
	sol Soliton
}

func NewBlockPlanner(nb uint32, seed uint32) BlockPlanner {
	return BlockPlanner{
		nblocks: nb,
		rando: &LinearGen{seed: seed},
		sol: NewSoliton(nb),
	}
}

func (planner BlockPlanner) CurrSeed() uint32 {
	return planner.rando.getSeed()
}

func (planner BlockPlanner) NextBlockList() (blockList []uint32, currSeed uint32) {
	currSeed = planner.rando.getSeed()
	nToCode := planner.sol.generate(planner.rando)
	blockList = make([]uint32, nToCode)
	var n uint

addBlock:
	for n < nToCode {
		nextBlock := planner.rando.nextInt() % planner.nblocks;
		for _, currBlock := range blockList[:n] {
			if currBlock == nextBlock {
				continue addBlock
			}
		}
		blockList[n] = nextBlock
		n++
	}
	return
}
