package asasn

const (
	maskTagClass    = 0b11000000
	maskConstructed = 0b00100000
	maskTagValue    = 0b00011111
	maskLoNibble    = 0b1111
	maskHiNibble    = maskLoNibble << 4
)

// DerReader - base type
type DerReader struct {
	buffer []byte
}

// NewReader - create reader
func NewReader(buf []byte) *DerReader {

	return &DerReader{buffer: buf}
}

func (r *DerReader) next() int {

	if len(r.buffer) > 0 {

		value := int(r.buffer[0])
		r.buffer = r.buffer[1:]
		return value
	}
	return -1
}

func (r *DerReader) readTag() (int, bool, int) {

	tagClass := r.next()
	var tagValue int = int(tagClass & maskTagValue)
	tagConstructed := (tagClass & maskConstructed) != 0
	tagClass = tagClass & maskTagClass

	if tagValue == 31 {
		tagValue = r.readBase128()
	}

	return tagClass, tagConstructed, tagValue
}

func (r *DerReader) readBase128() int {
	value := 0

	for true {
		byteValue := r.next()
		value |= (int)(byteValue & 0b01111111)
		if 0 == (byteValue & 0b10000000) {
			break
		}
		value <<= 7
	}
	return value
}
