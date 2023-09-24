/*
Package scanner is a custom text scanner implementation. It has the same
idiomatic Go scanner programming interface, and it lets the client to freely
navigate the buffer. The scanner is also capable of peeking ahead of the
cursor. Read runes are rendered as tokens with additional information on their
position in the buffer.

Usage

	package main

	import (
		"bufio"
		"fmt"
		"log"
		"os"

		"github.com/mdm-code/scanner"
	)

	func main() {
		r := bufio.NewReader(os.Stdin)
		s, err := scanner.New(r)
		if err != nil {
			log.Fatalln(err)
		}
		var ts []scanner.Token
		for s.Scan() {
			t := s.Token()
			ts = append(ts, t)
		}
		fmt.Println(ts)
	}
*/
package scanner
