package lt

import "bytes"
import "testing"

func TestBlockEncoder(t *testing.T) {
	blockA := uncodedBlock([]byte{1,2,3,4,5,6,7,8})
	blockB := []byte{12,13,14,15,16,17,18,19}

	blockX := blockA
	blockX.xorBlock(blockB)

	blockRes := []byte{13,15,13,11, 21,23,21,27}
	if !bytes.Equal(blockX, blockRes)  {
		t.Errorf("%x xor %x = %x, expected %x", blockA, blockB, blockX, blockRes)
	}
}
