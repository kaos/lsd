package main

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
	"unicode"
	"unicode/utf8"
)

type lsdLex struct {
	*Request
	log  *log.Entry
	data []byte
	peek rune
	row  int
	col  int
}

func (lex *lsdLex) Error(s string) {
	i, suffix := len(lex.data), ""
	if i > 25 {
		i = 25
		suffix = "..."
	}
	lex.log.Errorf(
		"%d:%d:%s: at \"%s%s\"%s\n",
		lex.row+1, lex.col, s,
		strings.Replace(string(lex.peek), "\n", "\\n", -1),
		strings.Replace(string(lex.data[0:i]), "\n", "\\n", -1),
		suffix,
	)
}

func (lex *lsdLex) Lex(lval *lsdSymType) int {
	if lex.row == 0 && lex.col == 0 {
		// first token on first row is command
		return lex.command(lval)
	}

	for {
		c := lex.Peek()
		switch c {
		case 0:
			return 0
		case '<', '>', '=':
			switch p := lex.drop_peek(); p {
			case '=':
				switch c {
				case '<':
					return lex.symbol(LTEQ)
				case '>':
					return lex.symbol(GTEQ)
				}
			case '~':
				switch c {
				case '=':
					return lex.symbol(EQI)
				}
			}
			return int(c)
		case '!', '~':
			return lex.symbol(c)
		}

		switch true {
		case unicode.In(c, unicode.L):
			if lex.col == 1 {
				return lex.header(lval)
			}
			return lex.ident(lval)
		case unicode.In(c, unicode.N):
			return lex.number_literal(lval)
		case unicode.In(c, unicode.P):
			return lex.punct(lval)
		case !unicode.In(c, unicode.Z, unicode.Cc):
			lex.Error("unexpected input")
		}

		// discard the rest
		lex.peek = 0
	}
}

func (lex *lsdLex) drop_peek() rune {
	lex.peek = 0
	return lex.Peek()
}

// Return the next rune for the lexer.
func (lex *lsdLex) next() rune {
	if lex.peek != 0 {
		r := lex.peek
		lex.peek = 0
		return r
	}
recurse:
	if len(lex.data) == 0 {
		return 0
	}
	c, size := utf8.DecodeRune(lex.data)
	lex.data = lex.data[size:]

	if c == '\n' {
		lex.row++
		lex.col = 0
	} else {
		lex.col++
	}

	if c == utf8.RuneError && size == 1 {
		lex.log.Warn("invalid utf8")
		goto recurse
	}

	if c == '#' && lex.col == 1 {
		for {
			c = lex.next()
			if c == 0 || c == '\n' {
				break
			}
		}
		goto recurse
	}

	return c
}

func (lex *lsdLex) Peek() rune {
	if lex.peek == 0 {
		lex.peek = lex.next()
	}
	return lex.peek
}

func (lex *lsdLex) symbol(c rune) int {
	lex.peek = 0
	return int(c)
}

func (lex *lsdLex) until(cond func(rune) bool) string {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			lex.log.Errorf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	for {
		c := lex.next()
		if c == 0 || cond(c) {
			lex.peek = c
			break
		}
		add(&b, c)
	}
	return b.String()
}

func (lex *lsdLex) command(lval *lsdSymType) int {
	str := lex.until(func(c rune) bool {
		return c == ' '
	})
	switch str {
	case "GET":
		return GET
	default:
		lex.Error("invalid request command")
	}
	return 0
}

func (lex *lsdLex) header(lval *lsdSymType) int {
	lex.ident(lval)
	if lex.peek != ':' {
		lex.Error("expected ':' after header name")
	}
	// eat :
	lex.peek = 0
	return HEADER
}

func (lex *lsdLex) ident(lval *lsdSymType) int {
	lval.str = lex.until(func(c rune) bool {
		return !unicode.In(c, unicode.L, unicode.N, unicode.Pc)
	})
	return IDENT
}

func (lex *lsdLex) punct(lval *lsdSymType) int {
	lval.punct = lex.peek
	lex.peek = 0
	return PUNCT
}

func (lex *lsdLex) number_literal(lval *lsdSymType) int {
	str := lex.until(func(c rune) bool {
		return unicode.In(c, unicode.Z, unicode.Cc)
	})

	i, err := fmt.Sscanf(str, "%v", &lval.num)
	if i != 1 || err != nil || str != fmt.Sprintf("%v", lval.num) {
		lex.ident(lval)
		lval.str = str + lval.str
		return IDENT
	}

	return NUMBER
}
