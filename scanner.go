package scanner

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

// ErrNilIOReader indicates that the parameter passed to an attribute of the
// inteface type io.Reader has a nil value.
var ErrNilIOReader error = errors.New("provided io.Reader is nil")

// ErrRuneError says that UTF-8 Unicode replacement character was encountered
// by the Scanner.
var ErrRuneError error = errors.New("Unicode replacement character found")

// Zero represents the initial state of the Scanner with the cursor pointing at
// the start of the byte buffer.
var Zero = Pos{Rune: '\u0000', Start: 0, End: 0}

// Pos carries information about the position of the rune in the byte buffer.
type Pos struct {
	Rune       rune
	Start, End int
}

// Token represents a single rune read from the byte buffer.
type Token struct {
	Pos
	Buffer *[]byte
}

// Scanner encapsulates the logic of scanning runes from a text file. Its
// instance is stateful and unsafe to use across multiple threads.
type Scanner struct {
	Buffer []byte
	Errors []error
	Cursor Pos
}

// New creates an instance of the Scanner in its initial state.
func New(r io.Reader) (*Scanner, error) {
	if r == nil {
		return nil, ErrNilIOReader
	}
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	s := Scanner{
		Cursor: Zero,
		Buffer: buf,
	}
	return &s, nil
}

// Position returns the position of the recorded character in the byte buffer.
func (t Token) Position() Pos {
	return Pos{
		Rune:  t.Rune,
		Start: t.Start,
		End:   t.End,
	}
}

// String returns a text representation of the Pos.
func (p Pos) String() string {
	repr := fmt.Sprintf("{ %c %d:%d }", p.Rune, p.Start, p.End)
	return repr
}

// Reset puts the Scanner back in its initial state with the cursor pointing at
// the start of the byte buffer and clears all the recored scanner errors.
func (s *Scanner) Reset() {
	s.Cursor = Zero
	s.Errors = nil
}

// Goto moves the cursor of the Scanner to the position of the t Token.
func (s *Scanner) Goto(t Token) { s.Cursor = t.Position() }

// Token returns the Token currently pointed at by the cursor of the Scanner.
func (s *Scanner) Token() Token {
	token := Token{
		Pos:    s.Cursor,
		Buffer: &s.Buffer,
	}
	return token
}

// Scan advances the cursor of the Scanner by a single UTF-8 encoded Unicode
// character. The method returns a boolean value so that is can be used
// idiomatically the same way other scanners in the standard Go library are
// used.
func (s *Scanner) Scan() bool {
	if s.Cursor.End >= len(s.Buffer) {
		return false
	}
	r, size := rune(s.Buffer[s.Cursor.End]), 1
	if r >= utf8.RuneSelf {
		r, size = utf8.DecodeRune(s.Buffer[s.Cursor.End:])
		if r == utf8.RuneError {
			s.Errors = append(s.Errors, ErrRuneError)
			return false
		}
	}
	s.Cursor = Pos{r, s.Cursor.End, s.Cursor.End + size}
	return true
}

// Peek reports whether the v string matches the byte buffer from the position
// currently pointed at by the cursor. It returns true if there is a match.
// It returns false either if there is no match or the provided v string goes
// beyond the length of the buffer. It does not advance the Scanner.
func (s *Scanner) Peek(v string) bool {
	if len(v)+s.Cursor.End > len(s.Buffer) {
		return false
	}
	start, end := s.Cursor.End, s.Cursor.End+len(v)
	if string(s.Buffer[start:end]) == v {
		return true
	}
	return false
}

// ScanAll scans all Tokens representing UTF-8 encoded Unicode characters from
// the byte buffer underlying the Scanner.
func (s *Scanner) ScanAll() ([]Token, bool) {
	result := make([]Token, s.Cursor.End)
	for s.Scan() {
		t := s.Token()
		result = append(result, t)
	}
	if s.Errored() {
		return result, false
	}
	return result, true
}

// Errored reports if the Scanner encountered errors while scanning the
// underlying byte buffer.
func (s *Scanner) Errored() bool {
	if len(s.Errors) > 0 {
		return true
	}
	return false
}
