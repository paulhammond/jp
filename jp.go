// Copyright 2013-2014 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

package jp

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

type dict struct {
	indented   string
	objEmpty   string
	objOpen    string
	objClose   string
	arrEmpty   string
	arrOpen    string
	arrClose   string
	colon      string
	comma      string
	strOpen    string
	strClose   string
	otherSpace string
	otherOpen  string
	otherClose string
	end        string
}

var dicts = map[string]dict{
	"pretty16":  {"\n  ", "\033[0;32m{ }", "\033[0;32m{\n", "\033[0;32m\n}", "\033[0;32m[ ]", "\033[0;32m[\n", "\033[0;32m\n]", "\033[0;32m: ", "\033[0;32m,\n", "\033[0m\"\033[1m", "\033[0m\"", "\033[0;41m \033[0m\033[1;33m", "\033[1;33m", "", "\033[0m\n"},
	"compact16": {"\n", "\033[0;32m{}", "\033[0;32m{", "\033[0;32m}", "\033[0;32m[]", "\033[0;32m[", "\033[0;32m]", "\033[0;32m:", "\033[0;32m,", "\033[0m\"\033[1m", "\033[0m\"", "\033[0;41m \033[0m\033[1;33m", "\033[1;33m", "", "\033[0"},
	"pretty":    {"\n  ", "{ }", "{\n", "\n}", "[ ]", "[\n", "\n]", ": ", ",\n", `"`, `"`, " ", "", "", "\n"},
	"compact":   {"\n", "{}", "{", "}", "[]", "[", "]", ":", ",", `"`, `"`, " ", "", "", ""},
}

func (d dict) indent() dict {
	return dict{
		d.indented,
		strings.Replace(d.objEmpty, "\n", d.indented, 1),
		strings.Replace(d.objOpen, "\n", d.indented, 1),
		strings.Replace(d.objClose, "\n", d.indented, 1),
		strings.Replace(d.arrEmpty, "\n", d.indented, 1),
		strings.Replace(d.arrOpen, "\n", d.indented, 1),
		strings.Replace(d.arrClose, "\n", d.indented, 1),
		strings.Replace(d.colon, "\n", d.indented, 1),
		strings.Replace(d.comma, "\n", d.indented, 1),
		strings.Replace(d.strOpen, "\n", d.indented, 1),
		strings.Replace(d.strClose, "\n", d.indented, 1),
		d.otherSpace,
		strings.Replace(d.otherOpen, "\n", d.indented, 1),
		strings.Replace(d.otherClose, "\n", d.indented, 1),
		d.end,
	}
}

type scanner struct {
	r           *bufio.Reader
	w           *bufio.Writer
	indentSize  int
	indentDicts []dict
	dict        *dict
}

func (s scanner) writeRune(r rune) (e error) {
	_, e = s.w.WriteRune(r)
	return e
}

func (s scanner) writeString(str string) (e error) {
	_, e = s.w.WriteString(str)
	return e
}

func (s *scanner) indent(d int) {
	s.indentSize += d
	if len(s.indentDicts) <= s.indentSize {
		s.indentDicts = append(s.indentDicts, s.indentDicts[len(s.indentDicts)-1].indent())
	}
	s.dict = &s.indentDicts[s.indentSize]
}

func (s scanner) unread() {
	e := s.r.UnreadRune()
	if e != nil {
		panic(e.Error()) // we only ever read runes, so this shouldn't happen
	}
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
	for e == nil {
		r, _, e = s.r.ReadRune()
		if e != nil {
			break
		}

		if r == '"' && last != '\\' {
			e = s.writeString(s.dict.strClose)
			break
		} else if last != '\\' {
			last = r
			e = s.writeRune(r)
		} else {
			last = 0
			e = s.writeRune(r)
		}
	}
	return e
}

func (s scanner) copyOther() (e error) {
	var r rune
	var space bool
	e = s.writeString(s.dict.otherOpen)
	for e == nil {
		r, _, e = s.r.ReadRune()
		if e != nil {
			break
		}
		switch r {
		case '{', '}', '[', ']', ',', ':', '"':
			s.unread()
			e = s.writeString(s.dict.otherClose)
			return e
		default:
			if unicode.IsSpace(r) {
				space = true
				break
			}
			if space {
				space = false
				e = s.writeString(s.dict.otherSpace)
				if e != nil {
					return e
				}
			}
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
				s.unread()
				s.indent(1)
				e = s.writeString(s.dict.objOpen)
			}
		case '}':
			s.indent(-1)
			e = s.writeString(s.dict.objClose)
		case '[':
			r, e = s.read()
			if e != nil {
				break
			}
			if r == ']' {
				e = s.writeString(s.dict.arrEmpty)
			} else {
				s.unread()
				s.indent(1)
				e = s.writeString(s.dict.arrOpen)
			}
		case ']':
			s.indent(-1)
			e = s.writeString(s.dict.arrClose)
		case ',':
			e = s.writeString(s.dict.comma)
		case ':':
			e = s.writeString(s.dict.colon)
		case '"':
			e = s.copyString()
		// todo unicode.ReplacementChar
		default:
			s.unread()
			e = s.copyOther()
		}
	}
	s.writeString(s.dict.end)
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
	indentDicts := []dict{d}
	s := &scanner{
		r:           bufio.NewReader(reader),
		w:           bufio.NewWriter(writer),
		indentDicts: indentDicts,
		dict:        &indentDicts[0],
	}
	return s.expand()
}
