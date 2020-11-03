package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"plk/asasn/ber"
	"strings"
)

func main() {

	if buf, err := ioutil.ReadFile(os.Args[1]); err == nil {

		ber := ber.NewBuffer(buf)
		out := &strings.Builder{}

		out.Grow(100_000_000)
		for ber.Available() {
			ber.ReadTLV(out)
		}
		fmt.Println(out.String())

		return
	}
	fmt.Fprintln(os.Stderr, "Failure")

}
