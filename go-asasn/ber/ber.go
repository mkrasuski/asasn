package ber

// Buffer - buffer
type Buffer struct {
	buffer []byte
}

// NewBuffer -
func NewBuffer(buf []byte) *Buffer {
	return &Buffer{buffer: buf}
}

// NextByte -
func (b *Buffer) NextByte() int {
	if len(b.buffer) > 0 {
		ret := b.buffer[0]
		b.buffer = b.buffer[1:]
		return int(ret)
	}
	return -1
}

// Available -
func (b *Buffer) Available() bool {
	return len(b.buffer) > 0
}

// ReadTag - read tag
func (b *Buffer) ReadTag() (int, bool, int) {

	next := b.NextByte()

	cls := (next & 0b11000000) >> 6
	constructed := (next & 0b00100000) > 0
	value := (next & 0b00011111)

	if value == 31 {
		value = 0
		for true {
			value <<= 7
			next = b.NextByte()
			value |= (next & 0b01111111)
			if next < 128 {
				break
			}
		}
	}
	return cls, constructed, value
}

// ReadLen -
func (b *Buffer) ReadLen() int {

	next := b.NextByte()
	if next < 128 {
		return next
	}
	next -= 128
	value := 0
	for next > 0 {
		value <<= 8
		value |= b.NextByte()
		next--
	}
	return value
}

// SubBuffer -
func (b *Buffer) SubBuffer(len int) *Buffer {
	sub := b.buffer[:len]
	b.buffer = b.buffer[len:]
	return &Buffer{buffer: sub}
}

func (b *Buffer) skipBytes(len int) {
	b.buffer = b.buffer[len:]
}
