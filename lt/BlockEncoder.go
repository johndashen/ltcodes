package lt

import (
	"fmt"
	"io"
	"math"
)

type BlockEncoder struct {
	fileSize uint32
	blockSize uint32 /* size of the data encoded + the header */
	dataSize uint32 /* size of the data alone */

	fileData []byte
	nBlocks uint32
	planner BlockPlanner
}

type uncodedBlock []byte

// datablockSize is the size of the individual data blocks apart from the header 
func NewEncoder(buf io.Reader, size uint64, datablockSize uint32, initSeed uint32) *BlockEncoder {
	// read the stream 
	fileBuf := make([]byte, size)
	var bufPtr uint64
	for bufPtr < size {
		nRead, err := buf.Read(fileBuf[bufPtr:])
		if err != nil {
			fmt.Errorf(err.Error())
			return nil
		}
		bufPtr += uint64(nRead)
	}

	// pad with bytes here
	padNum := size % uint64(datablockSize)
	if padNum != 0 {
		padNum = uint64(datablockSize) - padNum
		pads := make([]byte, padNum)
		fileBuf = append(fileBuf, pads...)
	}

	nb := uint32(math.Ceil(float64(size)/float64(datablockSize)))
	return &BlockEncoder{
		blockSize: datablockSize + BLOCK_HEADER_SIZE,
		dataSize: datablockSize,
		fileSize: uint32(size),

		fileData: fileBuf,
		nBlocks: nb,
		planner: NewBlockPlanner(nb, initSeed),
	}
}

func (enc *BlockEncoder) getBlock(bnum uint32) uncodedBlock {
	startIdx := bnum * enc.dataSize
	return enc.fileData[startIdx:startIdx + enc.dataSize]
} 

func (enc *BlockEncoder) NextCodedBlock() CodedBlock {
	blockList, currSeed := enc.planner.NextBlockList()
	accum_block := uncodedBlock(make([]byte, enc.dataSize))

	for _, blockIdx := range blockList {
		accum_block.xorBlock(enc.getBlock(blockIdx))
	}

	ans := CodedBlock{
		fileSize: enc.fileSize,
		blockSize: enc.blockSize,
		seed: currSeed,
		data: accum_block}
	return ans
}

// assumes equal length, modifies x but not y
func (x uncodedBlock) xorBlock(y uncodedBlock) { //[]byte {
	if len(x) != len(y) {
		panic("xoring unequal length lists")
	}
	for i, yb := range y {
		x[i] ^= yb
	}
}
	

