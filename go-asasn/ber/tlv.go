package ber

import (
	"strconv"
	"strings"
)

const indent = "                                                "

var tagCls = [4]byte{'U', 'A', 'C', 'P'}

// ReadTLV - top level
func (b *Buffer) ReadTLV(out Output) {
	b.readOneTLV(out, "", 0)
}

func tagStr(cls int, value int) string {

	sb := strings.Builder{}
	sb.WriteByte('.')
	sb.WriteByte(tagCls[cls])
	sb.WriteString(strconv.FormatInt(int64(value), 10))
	return sb.String()
}

func (b *Buffer) readOneTLV(out Output, path string, depth int) {
	cls, constructed, value := b.ReadTag()
	len := b.ReadLen()

	path += tagStr(cls, value)

	name := path
	fmter := (*Buffer).AsHex

	if f, ok := FieldByPath(path); ok {
		name = f.name
		fmter = f.format
	}

	out.WriteString(indent[:depth])
	out.WriteString(name)
	out.WriteByte(':')
	out.WriteByte(' ')

	valueReader := b.SubBuffer(len)
	if constructed {
		out.WriteByte('\n')
		valueReader.readAllTLVs(out, path, depth+2)
	} else {
		fmter(valueReader, out, len)
		out.WriteByte('\n')
	}
}

func (b *Buffer) readAllTLVs(out Output, path string, depth int) {

	for b.Available() {

		b.readOneTLV(out, path, depth)
	}
}
