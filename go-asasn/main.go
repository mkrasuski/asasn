package main

import (
	"fmt"
	"io/ioutil"
)

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

// Next - next byte or -1  
func (r *DerReader) Next() int {

	if len(r.buffer) > 0 {

		value := int(r.buffer[0])
		r.buffer = r.buffer[1:]
		return value
	}
	return -1
}

func (r *DerReader) readTag() (int, bool, int) {

	tagClass := r.Next()
	var tagValue int = int(tagClass & maskTagValue)
	tagConstructed := (tagClass & maskConstructed) != 0
	tagClass = tagClass & maskTagClass

	if tagValue == 31 {
		tagValue = r.readBase128()
	}

	return (tagClass >> 6), tagConstructed, tagValue
}

func (r *DerReader) readBase128() int {
	value := 0

	for true {
		byteValue := r.Next()
		value |= (int)(byteValue & 0b01111111)
		if 0 == (byteValue & 0b10000000) {
			break
		}
		value <<= 7
	}
	return value
}

func (r *DerReader) readLen() int {

	byteValue := r.Next()

	if byteValue > 128 {
		count := byteValue - 128
		value := 0
		for count > 0 {
			value <<= 8
			byteValue = r.Next()
			value |= byteValue
			count--
		}
		return value
	}

	return int(byteValue)
}

func main() {

	fmt.Println("Start")

	if buf, err := ioutil.ReadFile("/home/mkr/SAMPLE"); err == nil {
		fmt.Printf("file size %d\n", len(buf))
		r := NewReader(buf)
		cls, cons, val := r.readTag()
		l := r.readLen()
		fmt.Printf("%b %v %v [%d]\n", cls, cons, val, l)
	} else {

		fmt.Println("failed.")
	}

}
