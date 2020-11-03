package ber

import "strconv"

// Output -
type Output interface {
	Write([]byte) (int, error)
	WriteString(string) (int, error)
	WriteByte(byte) error
}

const (
	hexDigit = "0123456789ABCDEF"
	loNibble = 0x0F
	hiNibble = 0xF0
)

// AsHex -
func (b *Buffer) AsHex(out Output, len int) {

	for len > 0 {
		next := b.NextByte()
		out.WriteByte(hexDigit[next&hiNibble>>4])
		out.WriteByte(hexDigit[next&loNibble])
		len--
	}
}

// AsRHex -
func (b *Buffer) AsRHex(out Output, len int) {

	for len > 0 {
		next := b.NextByte()
		out.WriteByte(hexDigit[next&loNibble])
		out.WriteByte(hexDigit[next&hiNibble>>4])
		len--
	}
}

// AsString -
func (b *Buffer) AsString(out Output, len int) {
	str := b.buffer[:len]
	b.buffer = b.buffer[len:]
	out.Write(str)
}

// AsInt -
func (b *Buffer) AsInt(out Output, len int) {
	val := b.readInt(len)
	out.WriteString(strconv.FormatInt(val, 10))
}

func (b *Buffer) readInt(len int) int64 {
	val := int64(0)
	for len > 0 {
		next := b.NextByte()
		val <<= 8
		val |= int64(next)
		len--
	}
	return val
}

// AsIPAddress -
func (b *Buffer) AsIPAddress(out Output, len int) {

	out.WriteString(strconv.FormatInt(int64(b.NextByte()), 10))
	out.WriteByte('.')
	out.WriteString(strconv.FormatInt(int64(b.NextByte()), 10))
	out.WriteByte('.')
	out.WriteString(strconv.FormatInt(int64(b.NextByte()), 10))
	out.WriteByte('.')
	out.WriteString(strconv.FormatInt(int64(b.NextByte()), 10))
}
