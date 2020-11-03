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

// AsBCD -
func (b *Buffer) AsBCD(out Output, len int) {

	for len > 0 {
		next := b.NextByte()
		out.WriteByte(hexDigit[next&hiNibble>>4])
		lo := next & loNibble
		if lo != 15 {
			out.WriteByte(hexDigit[lo])
		}
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

// AsRBCD -
func (b *Buffer) AsRBCD(out Output, len int) {

	for len > 0 {
		next := b.NextByte()
		out.WriteByte(hexDigit[next&loNibble])
		hi := next & hiNibble >> 4
		if hi != 15 {
			out.WriteByte(hexDigit[hi])
		}
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

// AsTimestamp -
func (b *Buffer) AsTimestamp(out Output, len int) {

	b.AsHex(out, 3) // date
	out.WriteByte(' ')
	b.AsHex(out, 3) // time
	out.WriteByte(' ')
	b.NextByte()
	b.AsHex(out, len-7)
}

// AsPLMN -
func (b *Buffer) AsPLMN(out Output, len int) {

	b.AsRBCD(out, 2)
	out.WriteByte('-')
	b.AsRBCD(out, len-2)
}

// AsBool -
func (b *Buffer) AsBool(out Output, len int) {

	if 0 == b.buffer[0] {

		out.WriteByte('T')
	} else {

		out.WriteByte('F')
	}

	b.skipBytes(len)
}
