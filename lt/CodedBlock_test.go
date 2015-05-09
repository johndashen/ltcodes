package lt

import "bytes"
import "testing"

const str = "abcdefghijklmnopqrst"
func TestCodedBlock(t *testing.T) {
	block := CodedBlock{
		fileSize: 30, 
		blockSize: uint32(len(str) + BLOCK_HEADER_SIZE), 
		seed: 91532, 
		data: []byte(str),
	}
	
	buf := bytes.NewBuffer(block.Pack())
	readBlock, err := ReadBlockFrom(buf)
	if err != nil {
		t.Error("Unable to read block, error: ", err)
	} else if block.fileSize != readBlock.fileSize || 
		block.blockSize != readBlock.blockSize || 
		block.seed != readBlock.seed ||
		!bytes.Equal(block.data, readBlock.data) {
		t.Error("unread/read =", readBlock, ", expected", block)
	} 
}

