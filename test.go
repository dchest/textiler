package main

import (
	"bytes"
	"fmt"
	"textiler"
)

var passed = 0

func testHtml() {
	passingTests := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	for _, i := range passingTests {
		s := textiler.HtmlTests[i*2]
		expected := []byte(textiler.HtmlTests[i*2+1])
		res := textiler.ToHtml([]byte(s), false, false)
		if !bytes.Equal(res, expected) {
			textiler.ToHtml([]byte(s), false, true)
			fmt.Printf("**Conversion failed!**\n\n'%v'\n\n'%v'\n\n'%v'\n", s, string(expected), string(res))
			return
		}
		passed += 1
	}
}

func testXhtml() {
	passingTests := []int{0, 1, 2}
	for _, i := range passingTests {
		s := textiler.XhtmlTests[i*2]
		expected := []byte(textiler.XhtmlTests[i*2+1])
		res := textiler.ToXhtml([]byte(s), false, false)
		if !bytes.Equal(res, expected) {
			textiler.ToXhtml([]byte(s), false, true)
			fmt.Printf("**Conversion failed!**\n\n'%v'\n\n'%v'\n\n'%v'\n", s, string(expected), string(res))
			return
		}
		passed += 1
	}
}

func main() {
	testXhtml()
	testHtml()
	fmt.Printf("\nPassed %d tests\n", passed)
}
