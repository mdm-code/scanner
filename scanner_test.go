package scanner

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

type failer struct{}

func (failer) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("err")
}

// Test if non-nil interface value passed to New() passes without errors.
func TestNewScanner(t *testing.T) {
	reader := strings.NewReader("")
	s, err := New(reader)
	if err != nil {
		if errors.Is(err, ErrNilIOReader) {
			t.Fatal("New() returned ErrNilIOReader with a proper reader")
		}
		t.Fatal("New() returned an error with a proper reader")
	}
	t.Logf("new Scanner: %v", s)
}

// Check if a nil interface value causes New() to return ErrNilIOReader error.
func TestErrNilIOReader(t *testing.T) {
	var reader io.Reader
	_, err := New(reader)
	if !errors.Is(err, ErrNilIOReader) {
		t.Errorf("want error: %s; have error: %s", ErrNilIOReader, err)
	}
}

// Test if New() errors out when byte buffer cannot be read from the interface
// value.
func TestNewScannerErrors(t *testing.T) {
	cases := []struct {
		name   string
		reader io.Reader
	}{
		{"nil", nil},
		{"fail", failer{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := New(c.reader)
			if err == nil {
				t.Errorf("New() was expected to return an error")
			}
		})
	}
}

// Test if the string representation of the Token matches the expected format.
func TestTokenString(t *testing.T) {
	cases := []struct {
		rune       rune
		start, end int
	}{
		{'\u0000', 0, 1},
		{'a', 23, 24},
		{'禪', 2, 5},
		{'9', 37, 38},
		{'\uffff', 0, 3},
	}
	for _, c := range cases {
		t.Run(string(c.rune), func(t *testing.T) {
			token := Token{
				Pos{Rune: c.rune, Start: c.start, End: c.end}, nil,
			}
			have := token.String()
			want := fmt.Sprintf("{ %c %d:%d }", c.rune, c.start, c.end)
			if have != want {
				t.Errorf("have %s; want %s", have, want)
			}
		})
	}
}

// Test if Tokens produced by the Scanner are aligned with the expected output.
func TestScannerScan(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  []Token
	}{
		{
			"simple",
			".transaction[].status.柳",
			[]Token{
				{Pos{'.', 0, 1}, nil},
				{Pos{'t', 1, 2}, nil},
				{Pos{'r', 2, 3}, nil},
				{Pos{'a', 3, 4}, nil},
				{Pos{'n', 4, 5}, nil},
				{Pos{'s', 5, 6}, nil},
				{Pos{'a', 6, 7}, nil},
				{Pos{'c', 7, 8}, nil},
				{Pos{'t', 8, 9}, nil},
				{Pos{'i', 9, 10}, nil},
				{Pos{'o', 10, 11}, nil},
				{Pos{'n', 11, 12}, nil},
				{Pos{'[', 12, 13}, nil},
				{Pos{']', 13, 14}, nil},
				{Pos{'.', 14, 15}, nil},
				{Pos{'s', 15, 16}, nil},
				{Pos{'t', 16, 17}, nil},
				{Pos{'a', 17, 18}, nil},
				{Pos{'t', 18, 19}, nil},
				{Pos{'u', 19, 20}, nil},
				{Pos{'s', 20, 21}, nil},
				{Pos{'.', 21, 22}, nil},
				{Pos{'柳', 22, 25}, nil},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := strings.NewReader(c.query)
			s, err := New(r)
			if err != nil {
				t.Fatal("failed to initialize the scanner")
			}
			have := []Token{}
			for s.Scan() {
				t := s.Token()
				t.Buffer = nil // change the buffer pointer to nil
				have = append(have, t)
			}
			if !reflect.DeepEqual(have, c.want) {
				t.Errorf("have %v; want: %v", have, c.want)
			}
		})
	}
}

// Test if the Scanner Scan() method returns false on failure.
func TestScanFailure(t *testing.T) {
	s := Scanner{[]byte(string('\uFFFD')), Pos{'\u0000', 0, 0}}
	if s.Scan() != false {
		t.Error("scan was expected to return false on rune error")
	}
}

// Test if the Goto() method of the Scanner changes the recorded state.
func TestScannerGoto(t *testing.T) {
	cases := []struct {
		name  string
		token Token
	}{
		{`ù`, Token{Pos{'ù', 5, 7}, nil}},
		{`æ`, Token{Pos{'æ', 112, 115}, nil}},
		{`ß`, Token{Pos{'ß', 0, 2}, nil}},
		{`§`, Token{Pos{'§', 14, 16}, nil}},
		{`®`, Token{Pos{'®', 67, 70}, nil}},
	}
	s := Scanner{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.Goto(c.token)
			if c.token.Position() != s.Cursor {
				t.Errorf("have %v; want %v", c.token.Position(), s.Cursor)
			}
		})
	}
}

// Test if the Reset() method of the Scanner resets its cursor back to zero.
func TestScannerReset(t *testing.T) {
	s := Scanner{}
	s.Reset()
	have, want := s.Cursor, Zero
	if have != want {
		t.Errorf("have %v; want %v", have, want)
	}
}

// Test if the Peek() method of the Scanner returns the expected boolean value.
func TestScannerPeek(t *testing.T) {
	cases := []struct {
		name, input, match string
		want               bool
	}{
		{"pass-empty", "input text", "", true},
		{"pass-both-empty", "", "", true},
		{"pass", "input text", "input", true},
		{"fail", "input text", "output", false},
		{"fail-empty", "", "need more text", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := strings.NewReader(c.input)
			s, err := New(r)
			if err != nil {
				t.Fatal("failed to initialize the scanner")
			}
			if have := s.Peek(c.match); have != c.want {
				t.Errorf("have %t; want %t", have, c.want)
			}
		})
	}
}

// Test if ScanAll() method returns the expected slice of scanned tokens.
func TestScannerScanAll(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			"basic",
			".tests[].value",
			[]Token{
				{Pos{'.', 0, 1}, nil},
				{Pos{'t', 1, 2}, nil},
				{Pos{'e', 2, 3}, nil},
				{Pos{'s', 3, 4}, nil},
				{Pos{'t', 4, 5}, nil},
				{Pos{'s', 5, 6}, nil},
				{Pos{'[', 6, 7}, nil},
				{Pos{']', 7, 8}, nil},
				{Pos{'.', 8, 9}, nil},
				{Pos{'v', 9, 10}, nil},
				{Pos{'a', 10, 11}, nil},
				{Pos{'l', 11, 12}, nil},
				{Pos{'u', 12, 13}, nil},
				{Pos{'e', 13, 14}, nil},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := strings.NewReader(c.input)
			s, err := New(r)
			if err != nil {
				t.Fatal("failed to initialize the scanner")
			}
			have, err := s.ScanAll()
			for i := 0; i < len(have); i++ {
				have[i].Buffer = nil
			}
			if err != nil {
				t.Fatal("failed to scan all tokens at once")
			}
			if !reflect.DeepEqual(have, c.want) {
				t.Errorf("have: %v; want: %v", have, c.want)
			}
		})
	}
}
