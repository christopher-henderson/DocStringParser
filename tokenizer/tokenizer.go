package tokenizer

/*
Known bugs:
	/*X where X is o * else causes a bad state since the reader isn't putting X  and * back
*/

import (
	"bufio"
	"io"
	"strings"
)

type Tokenizer struct {
	src    io.RuneScanner
	tokens []Tokener
}

func (t *Tokenizer) Tokens() []Tokener {
	return t.tokens
}

type (
	Tokener interface {
		Original() string
		SetOriginal(string)
	}

	Token struct {
		original string
	}

	OpenDoc struct {
		Token
	}

	CloseDoc struct {
		Token
	}

	OpenBlock struct {
		Token
	}

	CloseBlock struct {
		Token
	}

	At struct {
		Token
	}

	Title struct {
		Token
	}

	Desc struct {
		Token
	}

	Param struct {
		Token
	}

	Column struct {
		Token
	}

	Table struct {
		Token
	}

	Text struct {
		Token
	}

	BareWord struct {
		Token
	}
)

func (t Token) Original() string {
	return t.original
}

func (t Token) SetOriginal(s string) {
	t.original = s
}

func NewTokenizer(src io.Reader) *Tokenizer {
	return &Tokenizer{src: bufio.NewReader(src), tokens: make([]Tokener, 0)}
}

func (t *Tokenizer) Peek() (s string, err error) {
	s, err = t.Read()
	t.src.UnreadRune()
	return
}

func (t *Tokenizer) Read() (s string, err error) {
	r, _, err := t.src.ReadRune()
	s = string(r)
	return
}

func (t *Tokenizer) Tokenize() error {
	for {
		c, err := t.Read()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch c {
		case "/":
			err := t.attemptDoc()
			if err != nil {
				return err
			}
		}
	}
}

func (t *Tokenizer) attemptDoc() error {
	peek, err := t.Peek()
	switch err {
	case nil:
		break
	case io.EOF:
		return nil
	default:
		return err
	}
	switch peek {
	case "*":
		t.Read()
		peek, err := t.Peek()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch peek {
		case "*":
			t.Read()
			t.tokens = append(t.tokens, OpenDoc{Token{"/**"}})
			t.tokens = append(t.tokens, OpenBlock{Token{"/**"}})
			return t.tokenizeDoc()
		}
	}
	return nil
}

func (t *Tokenizer) tokenizeDoc() error {
	err := t.tokenizeBlock()
	if err != nil {
		return err
	}
	for {
		c, err := t.Read()
		switch err {
		case nil:
			break
		case io.EOF:
			t.tokens = append(t.tokens, CloseDoc{Token{"EOF"}})
			return io.EOF
		default:
			return err
		}
		switch c {
		case ";":
			t.tokens = append(t.tokens, CloseDoc{Token{";"}})
			return nil
		case "/":
			err := t.attemptBlock()
			if err != nil {
				return err
			}

		}

	}
}

func (t *Tokenizer) attemptBlock() error {
	peek, err := t.Peek()
	switch err {
	case nil:
		break
	case io.EOF:
		return nil
	default:
		return err
	}
	switch peek {
	case "*":
		t.Read()
		peek, err := t.Peek()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch peek {
		case "*":
			t.Read()
			t.tokens = append(t.tokens, OpenBlock{Token{"/**"}})
			return t.tokenizeBlock()
		}
	}
	return nil
}

func (t *Tokenizer) tokenizeBlock() error {
	for {
		c, err := t.Read()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch c {
		case "*":
			peek, err := t.Peek()
			switch err {
			case nil:
				break
			case io.EOF:
				return nil
			default:
				return err
			}
			switch peek {
			case "/":
				t.Read()
				t.tokens = append(t.tokens, CloseBlock{Token{"*/"}})
				return nil
			}
		case "@":
			err := t.tokenizeAnnotation()
			if err != nil {
				return err
			}
		}
	}
}

func (t *Tokenizer) tokenizeAnnotation() error {
	annotation, err := t.buildAnnotationName()
	if err != nil {
		return err
	}
	switch annotation {
	// @TODO magic strings. evil.
	case "title":
		t.tokens = append(t.tokens, Title{Token{annotation}})
		err := t.consumeUntilQuote()
		if err != nil {
			return err
		}
		return t.tokenizeText()
	case "description":
		t.tokens = append(t.tokens, Desc{Token{annotation}})
		err := t.consumeUntilQuote()
		if err != nil {
			return err
		}
		return t.tokenizeText()
	case "param":
		t.tokens = append(t.tokens, Param{Token{annotation}})
		t.tokenizeParamColumnContents()
	case "column":
		t.tokens = append(t.tokens, Column{Token{annotation}})
		t.tokenizeParamColumnContents()
	case "table":
		t.tokens = append(t.tokens, Table{Token{annotation}})
		t.tokenizeTable()
	}
	return nil
}

func (t *Tokenizer) buildAnnotationName() (string, error) {
	b := strings.Builder{}
	for {
		c, err := t.Read()
		switch err {
		case nil:
			break
		case io.EOF:
			return "", nil
		default:
			return "", err
		}
		switch c {
		case "\n", "\t", "\r", " ":
			return b.String(), nil
		default:
			b.WriteString(c)
		}
	}
	return "", nil
}

func (t *Tokenizer) tokenizeTable() error {
	err := t.consumeUntilLBracket()
	if err != nil {
		return err
	}
	for {
		c, err := t.Read()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch c {
		case "}":
			return nil
		case "@":
			err := t.tokenizeAnnotation()
			if err != nil {
				return err
			}
		}
	}
}

func (t *Tokenizer) tokenizeText() error {
	b := strings.Builder{}
	for {
		c, err := t.Read()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch c {
		case "\"":
			t.tokens = append(t.tokens, Text{Token{b.String()}})
			return nil
		default:
			b.WriteString(c)
		}
	}
}

func (t *Tokenizer) tokenizeBareWord() error {
	b := strings.Builder{}
	for {
		c, err := t.Read()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch c {
		case " ", "\n", "\t", "\r":
			t.tokens = append(t.tokens, BareWord{Token{b.String()}})
			return nil
		default:
			b.WriteString(c)
		}
	}
}

func (t *Tokenizer) tokenizeParamColumnContents() error {
	err := t.consumeSpaces()
	if err != nil {
		return err
	}
	err = t.tokenizeBareWord()
	if err != nil {
		return err
	}
	err = t.consumeUntilQuote()
	if err != nil {
		return err
	}
	err = t.tokenizeText()
	if err != nil {
		return err
	}
	err = t.consumeUntilQuote()
	if err != nil {
		return err
	}
	err = t.tokenizeText()
	if err != nil {
		return err
	}
	return nil
}

func (t *Tokenizer) consumeUntilQuote() error {
	for {
		c, err := t.Read()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch c {
		case "\"":
			return nil
		}
	}
}

func (t *Tokenizer) consumeUntilLBracket() error {
	for {
		c, err := t.Read()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch c {
		case "{":
			return nil
		}
	}
}

func (t *Tokenizer) consumeSpaces() error {
	for {
		c, err := t.Peek()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		switch c {
		case " ", "\n", "\t", "\r":
			t.Read()
		default:
			return nil
		}
	}
}
