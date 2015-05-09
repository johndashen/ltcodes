package lt

import (
	"bytes"
	"fmt"
	"math"
)

type decodedBlock []byte
type mixedBlock struct {
	data []byte
	mix []uint32
}

type cleanBlock struct{
	data decodedBlock;
	i uint32
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
func (x *mixedBlock) removeFromMix(idx uint32) bool {
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
		panic("xoring block that is not contained in the mixed block")
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

func NewDecoder(firstBlock CodedBlock) *BlockDecoder {
	nb := uint32(math.Ceil(float64(firstBlock.fileSize)/float64(firstBlock.blockSize - BLOCK_HEADER_SIZE)))
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
	return &em
}

func (dec *BlockDecoder) BlockSize() uint32 {
	return dec.blockSize
}

func (dec *BlockDecoder) FileSize() uint32 {
	return dec.fileSize
}

func (dec *BlockDecoder) addToDone(mb mixedBlock) {
	if len(mb.mix) != 1 {
		panic("shouldn't call addToDone")
		return
	} else if _, ok := dec.doneBlocks[mb.mix[0]]; ok {
		return
	} else {
		dec.doneBlocks[mb.mix[0]] = mb.data
		dec.nLeft--
		dec.isDone[mb.mix[0]] = true
	}
}

func (dec *BlockDecoder) Validate(block CodedBlock) bool {
	if dec.blockSize == 0 { //decoder not initialized
		return true
	} else if dec.blockSize != block.blockSize || dec.fileSize != block.fileSize {
		return false
	} 
	return true
}

// precondition: validate first or risk panic based on sizes
func (dec *BlockDecoder) Include(block CodedBlock) {
	if dec.blockSize == 0 { //decoder not initialized
		*dec = *NewDecoder(block)
		return
	}

	if dec.blockSize != block.blockSize || dec.fileSize != block.fileSize {
		panic(fmt.Sprintf("(block, file)size not matching: decoder = (%d,%d) incoming = (%d, %d)", 
			dec.blockSize, dec.fileSize, block.blockSize, block.fileSize))
	} 

	dec.planner.rando.setSeed(block.seed)
	blockList, seed := dec.planner.NextBlockList()

	if seed != block.seed {
		panic(fmt.Sprintf("seed not matching: decoder = %d, incoming = %d", seed, block.seed))
	} 

	// the new block
	mb := mixedBlock{
		data: block.data,
		mix: blockList,
	}
		
	ptr := 0
	for ptr < len(mb.mix) { // blockList { 
		b := mb.mix[ptr]
		if dec.isDone[b] { // block b in the mix is already uncoded
			if doneBlock, ok := dec.doneBlocks[b]; !ok {
				panic("done block not in map") 
			} else {
				mb.xorBlock(doneBlock, b)
			}
		} else { // block b is not encoded
			ptr++
		}
	}
	if len(mb.mix) == 0 {
		if !bytes.Equal(mb.data, make([]byte, dec.blockSize - BLOCK_HEADER_SIZE)) {
			panic("completely unmixed block is not zero")
		}
	} else if len(mb.mix) == 1 {
		dec.addToDone(mb) // if this block is cleaned
		dec.reduceOther(mb.data, mb.mix[0])
	} else {
		dec.stack = append(dec.stack, mb) // add to dirty block
	}
}

func (dec BlockDecoder) AttemptDone() (done bool, data []byte) {
	if len(dec.doneBlocks) == 0 || dec.nLeft > 0 {
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
	if uint32(len(data)) < dec.fileSize {
		panic("not enough data for the file")
	}
	if uint32(len(data)) > dec.fileSize {
		data = data[:dec.fileSize]
	}
	return
}

func (dec *BlockDecoder) reduceOther(clean decodedBlock, idx uint32) {
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
			}

			if len(dirtyBlock.mix) == 1 { // block is now clean
				if block, in := dec.doneBlocks[dirtyBlock.mix[0]]; in { // it was clean already
					if !bytes.Equal(block, dirtyBlock.data) {
						panic ("block mismatch")
					}
				} else { // new clean block
					dec.addToDone(dirtyBlock)
					alsoCleaned = append(alsoCleaned, 
						cleanBlock{
							data:dirtyBlock.data, 
							i:dirtyBlock.mix[0]})
				}
			} else { // block is still dirty
				newStack = append(newStack, dirtyBlock)
			}
		}
		dec.stack = newStack
	}
}
