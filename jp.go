package jp

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

type dict struct {
	objEmpty string
	objOpen  string
	objClose string
	arrEmpty string
	arrOpen  string
	arrClose string
	colon    string
	comma    string
	strOpen  string
	strClose string
}

var dicts = map[string]dict{
	"pretty":  {"{ }", "{\n", "\n}", "[ ]", "[\n", "\n]", ": ", ",\n", `"`, `"`},
	"compact": {"{}", "{", "}", "[]", "[", "]", ":", ",", `"`, `"`},
}

type scanner struct {
	r      *bufio.Reader
	w      *bufio.Writer
	indent int
	dict   dict
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

func (s scanner) read() (r rune, e error) {
	for e == nil {
		r, _, e = s.r.ReadRune()
		if !unicode.IsSpace(r) {
			break
		}
	}
	return r, e
}

func (s scanner) copyString() (e error) {
	var r rune
	var last rune
	e = s.writeString(s.dict.strOpen)
loop:
	for e == nil {
		r, _, e = s.r.ReadRune()
		if e != nil {
			break
		}

		switch r {
		case '"':
			if last == '\\' {
				e = s.writeRune(r)
			} else {
				e = s.writeString(s.dict.strClose)
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
		r, e = s.read()
		if e != nil {
			break
		}
		switch r {
		case '{':
			r, e = s.read()
			if e != nil {
				break
			}
			if r == '}' {
				e = s.writeString(s.dict.objEmpty)
			} else {
				e = s.r.UnreadRune()
				if e != nil {
					break // this really shouldn't happen
				}
				s.indent++
				e = s.writeIndented(s.dict.objOpen)
			}
		case '}':
			s.indent--
			e = s.writeIndented(s.dict.objClose)
		case '[':
			r, e = s.read()
			if e != nil {
				break
			}
			if r == ']' {
				e = s.writeString(s.dict.arrEmpty)
			} else {
				e = s.r.UnreadRune()
				if e != nil {
					break // this really shouldn't happen
				}
				s.indent++
				e = s.writeIndented(s.dict.arrOpen)
			}
		case ']':
			s.indent--
			e = s.writeIndented(s.dict.arrClose)
		case ',':
			e = s.writeIndented(s.dict.comma)
		case ':':
			e = s.writeString(s.dict.colon)
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

func Expand(reader io.Reader, writer io.Writer, format string) error {
	d, ok := dicts[format]
	if !ok {
		return errors.New("unknown format")
	}
	s := &scanner{
		r:    bufio.NewReader(reader),
		w:    bufio.NewWriter(writer),
		dict: d,
	}
	return s.expand()
}
