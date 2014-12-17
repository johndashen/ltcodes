package lt

import "bytes"
import "math"

type decodedBlock []byte
type mixedBlock struct {
	data []byte
	mix []uint32
}

func (x mixedBlock) contains(idx uint32) (b bool, mixIdx int) {
	var mixNum uint32
	for mixIdx, mixNum = range x.mix {
		if idx == mixNum {
			b = true
			return b, mixIdx
		}
	}
	return
}

// returns true if removed, false if not found
func (x mixedBlock) removeFromMix(idx uint32) bool {
	var ok bool
	var mixIdx int
	if ok, mixIdx = x.contains(idx); ok {
		x.mix = append(x.mix[0:mixIdx], x.mix[mixIdx+1:]...)
	}
	return ok
}

// blockIdx
func (x *mixedBlock) xorBlock(y []byte, idx uint32) { //[]byte {
	if len(x.data) != len(y) {
		panic("xoring unequal length lists")
	} 
	if !x.removeFromMix(idx) {
		panic("xoring unfound block")
	}
	for i, yb := range y {
		x.data[i] ^= yb
	}
}

type BlockDecoder struct {
	blockSize uint32
	fileSize uint32
	
	fileData []byte
	nBlocks uint32
	planner BlockPlanner
	
	doneBlocks map[uint32]decodedBlock

	isDone []bool
	stack []mixedBlock
	nLeft int
}

func NewDecoder(firstBlock CodedBlock) BlockDecoder {
	nb := uint32(math.Ceil(float64(firstBlock.fileSize)/float64(firstBlock.blockSize)))
	em := BlockDecoder{
		blockSize: firstBlock.blockSize,
		fileSize: firstBlock.fileSize,
		fileData: make([]byte, firstBlock.fileSize),
		nBlocks: nb,
		planner: NewBlockPlanner(nb, firstBlock.seed),
		doneBlocks: make(map[uint32]decodedBlock),
		isDone: make([]bool, nb),
		stack: make([]mixedBlock, 0, nb),
		nLeft: int(nb),
	}
	em.Include(firstBlock)
	return em
}

func (dec *BlockDecoder) Include(block CodedBlock) {
	blockList, seed := dec.planner.NextBlockList()
	if seed != block.seed {
		panic("seed not matching")
	}

	mb := mixedBlock{
		data: block.data,
		mix: blockList,
	}

	for _, b := range blockList {
		if dec.isDone[b] { // block b in the mix is already uncoded
			if doneBlock, ok := dec.doneBlocks[b]; !ok {
				panic("done block not in map") 
			} else {
				mb.xorBlock(doneBlock, b)
			}
		} // block b is not encoded
	}
	if len(mb.mix) == 0 {
		if !bytes.Equal(mb.data, make([]byte, dec.blockSize)) {
			panic("completely unmixed block is not zero")
		}
	} else if len(mb.mix) == 1 {
		dec.doneBlocks[mb.mix[0]] = mb.data // if this block is cleaned
		dec.nLeft--
		dec.reduceOther(mb.data, mb.mix[0])
	} else {
		dec.stack = append(dec.stack, mb) // add to dirty block
	}
}

func (dec BlockDecoder) AttemptDone() (done bool, data []byte) {

	if dec.nLeft > 0 {
		return // false, empty
	} else {
		done = true
		for i := uint32(0); i < dec.nBlocks; i++ {
			if block, ok := dec.doneBlocks[i]; !ok {
				done = false
				return // false, partial
				panic("reported done but not all blocks in map")
			} else {
				data = append(data, block...)
			}
		}
	}
	return
}

func (dec *BlockDecoder) reduceOther(clean decodedBlock, idx uint32) {
	type cleanBlock struct{
		data decodedBlock;
		i uint32}
	alsoCleaned := []cleanBlock{{
		data:clean, 
		i:idx}}
	for len(alsoCleaned) > 0 && dec.nLeft > 0 {
		clean := alsoCleaned[0]
		alsoCleaned = alsoCleaned[1:]
		
		newStack := make([]mixedBlock, 0, len(dec.stack))
		for _, dirtyBlock := range dec.stack {
			if in, _ := dirtyBlock.contains(clean.i); in {
				
				dirtyBlock.xorBlock(clean.data, clean.i)
				if len(dirtyBlock.mix) == 1 { // block is now clean
					if block, in := dec.doneBlocks[dirtyBlock.mix[0]]; in {
						if !bytes.Equal(block, dirtyBlock.data) {
							panic ("block mismatch")
						}
					} else {
						dec.nLeft--

						dec.doneBlocks[dirtyBlock.mix[0]] = dirtyBlock.data
						alsoCleaned = append(alsoCleaned, 
							cleanBlock{
								data:dirtyBlock.data, 
								i:dirtyBlock.mix[0]})
					}
				} else { // block is still dirty
					newStack = append(newStack, dirtyBlock)
				}
			}
		}
		dec.stack = newStack
	}
}
