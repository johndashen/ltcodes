package lt

import "bytes"
import "testing"

func TestCodedBlock(t *testing.T) {
	block := CodedBlock{
		fileSize: 30, 
		blockSize: 20, 
		seed: 91532, 
		data: []byte("abcdefghijklmnopqrst"),
	}
	
	buf := bytes.NewBuffer(block.Pack())
	readBlock := ReadBlockFrom(buf)
	if block.fileSize != readBlock.fileSize || 
		block.blockSize != readBlock.blockSize || 
		block.seed != readBlock.seed ||
		!bytes.Equal(block.data, readBlock.data) {
		t.Error("unread/read =", readBlock, ", expected", block)
	} 
}

