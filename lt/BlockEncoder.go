package lt

import "fmt"
import "math"
import "os"

type BlockEncoder struct {
	blockSize uint32
	fileSize uint32
	fileData []byte
	nBlocks uint32
	planner BlockPlanner
}

type uncodedBlock []byte

func NewEncoder(filename string, bs uint32, initSeed uint32) BlockEncoder {
	var em BlockEncoder
	stats, err := os.Lstat(filename)
	if err != nil {
		fmt.Errorf(err.Error())
		return em
	}

	fSize := stats.Size()
	f, err := os.Open(filename)
	if err != nil {
		fmt.Errorf(err.Error())
		return em
	}
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}()
	
	fileBuf := make([]byte, fSize)
	var bufPtr uint64
	for bufPtr < uint64(fSize) {
		nRead, err := f.Read(fileBuf[bufPtr:])
		if err != nil {
			fmt.Errorf(err.Error())
			return em
		}
		bufPtr += uint64(nRead)
	}
	// pad with bytes here
	padNum := fSize % int64(bs)
	if padNum != 0 {
		padNum = int64(bs) - padNum
		pads := make([]byte, padNum)
		fileBuf = append(fileBuf, pads...)
	}

	nb := uint32(math.Ceil(float64(fSize)/float64(bs)))
	return BlockEncoder{
		blockSize: bs,
		fileSize: uint32(fSize),
		fileData: fileBuf,
		nBlocks: nb,
		planner: NewBlockPlanner(nb, initSeed),
	}
}

func (enc BlockEncoder) getBlock(bnum uint32) uncodedBlock {
	startIdx := bnum * enc.blockSize
	return enc.fileData[startIdx:startIdx + enc.blockSize]
} 

func (enc BlockEncoder) NextCodedBlock() CodedBlock {
	blockList, currSeed := enc.planner.NextBlockList()
	accum_block := uncodedBlock(make([]byte, enc.blockSize))
	copy(accum_block, enc.getBlock(blockList[0]))

	if len(blockList) > 1 {
		for _, blockIdx := range blockList[1:] {
//			accum_block = 
			accum_block.xorBlock(enc.getBlock(blockIdx))
		}
	}

	return CodedBlock{
		fileSize: enc.fileSize,
		blockSize: enc.blockSize,
		seed: currSeed,
		data: accum_block}
}

// assumes equal length, modifies x but not y
func (x uncodedBlock) xorBlock(y uncodedBlock) { //[]byte {
	if len(x) != len(y) {
		panic("xoring unequal length lists")
	}
	for i, yb := range y {
		x[i] ^= yb
	}
//	return x
}
	

