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

func (b CodedBlock) Pack() []byte {
	header := make([]byte, 12)
	NETWORK_BYTEORDER.PutUint32(header, b.fileSize)
	NETWORK_BYTEORDER.PutUint32(header[4:], b.blockSize)
	NETWORK_BYTEORDER.PutUint32(header[8:], b.seed)
	return append(header, b.data...)
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

	b.data = make([]byte, b.blockSize)
	_, err = r.Read(b.data) // todo: handle partial read
	if err != nil {
		return
	}
	return
}
