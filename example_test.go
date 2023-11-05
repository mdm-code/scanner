package scanner_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/mdm-code/scanner"
)

// ExampleScanner_ScanAll shows how to convert text into a list of tokens with
// a single method call to ScanAll() instead of using a for loop to traverse
// the input one token at a time.
func ExampleScanner_ScanAll() {
	in := "Hello!"
	r := strings.NewReader(in)
	s, err := scanner.New(r)
	if err != nil {
		log.Fatal(err)
	}
	ts, err := s.ScanAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ts)
	// Output: [{ H 0:1 } { e 1:2 } { l 2:3 } { l 3:4 } { o 4:5 } { ! 5:6 }]
}

// ExampleScanner_Scan shows how to translate text into a list of tokens with
// the Scanner public API. It combines New, Scan and Token to get a slice of
// tokens matching the provided "Hello\!" input.
func ExampleScanner_Scan() {
	in := "Hello!"
	r := strings.NewReader(in)
	s, err := scanner.New(r)
	if err != nil {
		log.Fatal(err)
	}

	var ts = []scanner.Token{}
	for s.Scan() {
		t := s.Token()
		ts = append(ts, t)
	}
	fmt.Println(ts)
	// Output: [{ H 0:1 } { e 1:2 } { l 2:3 } { l 3:4 } { o 4:5 } { ! 5:6 }]
}

// ExampleScanner_Reset shows how to reset the scanner back to its initial,
// zero state. In the example, tokens produced by the scanner the usual way are
// discarded, and then the scanner gets reset back to its initial state.
func ExampleScanner_Reset() {
	r := strings.NewReader("Hello!")
	s, err := scanner.New(r)
	if err != nil {
		log.Fatal(err)
	}
	var t scanner.Token
	for s.Scan() {
	}
	s.Reset()
	s.Scan()
	t = s.Token()
	fmt.Println(t)
	// Output: { H 0:1 }
}

// ExampleScanner_Goto shows how an already emitted token can be used to move
// the cursor of the scanner back to the position it's pointing at.
func ExampleScanner_Goto() {
	r := strings.NewReader("Hello!")
	s, err := scanner.New(r)
	if err != nil {
		log.Fatal(err)
	}

	var final scanner.Token
	for s.Scan() {
		if curr := s.Token(); curr.Rune == 'e' {
			final = curr
		}
	}
	s.Goto(final)
	fmt.Println(s.Token())
	// Output: { e 1:2 }
}

// ExampleScanner_Peek shows how to peek ahead of the scanner cursor to see
// whether the buffer ahead matches the provided string.
func ExampleScanner_Peek() {
	r := strings.NewReader("There's a match!")
	s, err := scanner.New(r)
	if err != nil {
		log.Fatal(err)
	}
	for s.Scan() {
		if t := s.Token(); t.Rune == 's' {
			break
		}
	}
	result := s.Peek(" a match!")
	fmt.Println(result)
	// Output: true
}
