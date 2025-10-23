package identitydsl

import (
	"testing"
)

func TestLexerStateMethods(t *testing.T) {
	l := &lexer{
		input: "abc\n123",
	}

	// Test next() updates pos and width correctly

	r := l.next()

	if r != 'a' {
		t.Errorf("expected 'a', got %q", r)
	}

	if l.pos <= l.start {
		t.Errorf("pos not advanced after next: start=%d pos=%d", l.start, l.pos)
	}

	if l.width == 0 {
		t.Errorf("width not set after next")
	}

	// Test peek() does not advance pos

	r2 := l.peek()

	if r2 != 'b' {
		t.Errorf("expected peek 'b', got %q", r2)
	}

	if l.pos != 1+l.start {
		t.Errorf("pos advanced after peek: %d", l.pos)
	}

	// Test backup() restores pos

	l.backup()

	if l.pos != l.start {
		t.Errorf("backup did not restore pos to start: %d", l.pos)
	}

	// Test accept() returns true for allowed rune

	l.pos = 0

	if !l.accept("a") {
		t.Errorf("accept did not match 'a'")
	}

	if l.start != 0 {
		t.Errorf("start moved unexpectedly after accept")
	}

	// Test acceptRun() consumes all valid runes
	l.pos = 0

	if !l.acceptRun("abc") {
		t.Errorf("acceptRun failed")
	}

	expectedPos := 3 // "abc"

	if l.pos != expectedPos {
		t.Errorf("acceptRun pos expected %d, got %d", expectedPos, l.pos)
	}

	// Test acceptToLineEnding stops before newline

	l.pos = 0

	l.acceptToLineEnding()

	if l.pos != 3 {
		t.Errorf("acceptToLineEnding expected pos 3, got %d", l.pos)
	}

	// Test acceptString advances pos correctly

	l.pos = 0

	if !l.acceptString("abc") {
		t.Errorf("acceptString failed to match")
	}

	if l.pos != 3 {
		t.Errorf("acceptString did not advance pos correctly: got %d", l.pos)
	}

	// Test value() returns correct substring

	l.start = 0

	l.pos = 3

	val := l.value()

	if val != "abc" {
		t.Errorf("value() expected 'abc', got %q", val)
	}
}
