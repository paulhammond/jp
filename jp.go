package jp

import (
	"bufio"
	"io"
	"strings"
)

type scanner struct {
	r      *bufio.Reader
	w      *bufio.Writer
	indent int
}

func (s scanner) writeRune(r rune) (e error) {
	_, e = s.w.WriteRune(r)
	return e
}

func (s scanner) writeString(str string) (e error) {
	_, e = s.w.WriteString(str)
	return e
}

func (s scanner) writeIndented(str string) (e error) {
	indent := "\n" + strings.Repeat("  ", s.indent)
	return s.writeString(strings.Replace(str, "\n", indent, 1))
}

func (s scanner) readRune() (r rune, e error) {
	r, _, e = s.r.ReadRune()
	return r, e
}

func (s scanner) copyString() (e error) {
	var r rune
	var last rune
	e = s.writeString(`"`)
loop:
	for e == nil {
		r, e = s.readRune()
		if e != nil {
			break
		}

		switch r {
		case '"':
			if last == '\\' {
				e = s.writeRune(r)
			} else {
				e = s.writeString(`"`)
				break loop
			}
		default:
			last = r
			e = s.writeRune(r)
		}
	}
	return e
}

func (s scanner) expand() (e error) {
	var r rune
	for e == nil {
		r, e = s.readRune()
		if e != nil {
			break
		}
		switch r {
		case '{':
			r, e = s.readRune()
			if e != nil {
				break
			}
			if r == '}' {
				e = s.writeString("{ }")
			} else {
				e = s.r.UnreadRune()
				if e != nil {
					break // this really shouldn't happen
				}
				s.indent++
				e = s.writeIndented("{\n")
			}
		case '}':
			s.indent--
			e = s.writeIndented("\n}")
		case '[':
			r, e = s.readRune()
			if e != nil {
				break
			}
			if r == ']' {
				e = s.writeString("[ ]")
			} else {
				e = s.r.UnreadRune()
				if e != nil {
					break // this really shouldn't happen
				}
				s.indent++
				e = s.writeIndented("[\n")
			}
		case ']':
			s.indent--
			e = s.writeIndented("\n]")
		case ',':
			e = s.writeIndented(",\n")
		case ':':
			e = s.writeString(": ")
		case '"':
			e = s.copyString()
		// todo unicode.ReplacementChar
		default:
			e = s.writeRune(r)
		}
	}
	s.w.Flush()
	if e == io.EOF {
		return nil
	}
	return e
}

func Expand(reader io.Reader, writer io.Writer) error {
	s := &scanner{
		r: bufio.NewReader(reader),
		w: bufio.NewWriter(writer),
	}
	return s.expand()
}
