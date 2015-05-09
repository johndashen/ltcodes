package lt

import "encoding/binary"
import "io"

type CodedBlock struct {
	fileSize uint32
	blockSize uint32
	seed uint32
	data []byte
}

var (
	NETWORK_BYTEORDER = binary.BigEndian
)
const BLOCK_HEADER_SIZE = 12

func EmptyCodedBlock(fs uint32, bs uint32) CodedBlock {
	return CodedBlock{
		fileSize: fs,
		blockSize: bs,
		seed: 1,
		data: make([]byte, bs - BLOCK_HEADER_SIZE)}
}

func (b CodedBlock) Pack() []byte {
	header := make([]byte, BLOCK_HEADER_SIZE)
	NETWORK_BYTEORDER.PutUint32(header, b.fileSize)
	NETWORK_BYTEORDER.PutUint32(header[4:], b.blockSize)
	NETWORK_BYTEORDER.PutUint32(header[8:], b.seed)
	return append(header, b.data...)
}

func (b CodedBlock) Seed() uint32 {
	return b.seed
}

func ReadBlockFrom(r io.Reader) (b CodedBlock, err error) {
	err = binary.Read(r, NETWORK_BYTEORDER, &b.fileSize)
	if err != nil {
		return 
	}

	err = binary.Read(r, NETWORK_BYTEORDER, &b.blockSize)
	if err != nil {
		return 
	}

	err = binary.Read(r, NETWORK_BYTEORDER, &b.seed)
	if err != nil {
		return
	}

	b.data = make([]byte, b.blockSize - BLOCK_HEADER_SIZE)
	_, err = r.Read(b.data) // todo: handle partial read and other errors
	if err != nil {
		return
	}
	return
}
