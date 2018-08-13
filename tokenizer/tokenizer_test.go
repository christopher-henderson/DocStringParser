package tokenizer

import (
	"strings"
	"testing"
)

const openDocTest = `
some things
that are
maybe queries
this is a haiku
on mount fuji

/**

*/
`

func TestOpenDoc(t *testing.T) {
	tok := NewTokenizer(strings.NewReader(openDocTest))
	tok.Tokenize()
	if len(tok.tokens) == 0 {
		t.Errorf("Failed to parse open document. Got empty token list")
	}
	switch typ := tok.tokens[0].(type) {
	case OpenDoc:
	default:
		t.Errorf("Incorrect token typed. Got %v want %v", typ, "/**")
	}
}

const openDocTestFalseStart = `
some things
that are
maybe queries
this is a haiku
on mount fuji

/*
Not for us
*/
`

func TestOpenDocFalseStart(t *testing.T) {
	tok := NewTokenizer(strings.NewReader(openDocTestFalseStart))
	tok.Tokenize()
	if len(tok.tokens) != 0 {
		t.Errorf("Failed to parse open document false start. Got %v want %v", tok.tokens, "[]")
	}
}

func TestWholeDocEoF(t *testing.T) {
	tok := NewTokenizer(strings.NewReader(openDocTest))
	tok.Tokenize()
	if len(tok.tokens) != 4 {
		t.Errorf("Document has incorrect number of tokens. Got %v want %v", len(tok.tokens), 4)
	}
	for i, tok := range tok.tokens {
		switch i {
		case 0:
			switch typ := tok.(type) {
			case OpenDoc:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "OpenDoc")
			}
		case 1:
			switch typ := tok.(type) {
			case OpenBlock:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "OpenBlock")
			}
		case 2:
			switch typ := tok.(type) {
			case CloseBlock:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "CloseBlock")
			}
		case 4:
			switch typ := tok.(type) {
			case CloseDoc:
				t.Errorf("%v", tok.Original())
				if tok.Original() != "EOF" {
					t.Errorf("Got wrong document closing delimiter. Got %v want ;", tok.Original())
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "CloseDoc")
			}
		}
	}
}

const wholeDocTestSemi = `
some things
that are
maybe queries
this is a haiku
on mount fuji

/**

*/
SELECT 1 FROM Tests;
`

func TestWholeDocSemicolon(t *testing.T) {
	tok := NewTokenizer(strings.NewReader(wholeDocTestSemi))
	tok.Tokenize()
	if len(tok.tokens) != 4 {
		t.Errorf("Document has incorrect number of tokens. Got %v want %v", len(tok.tokens), 4)
	}
	for i, tok := range tok.tokens {
		switch i {
		case 0:
			switch typ := tok.(type) {
			case OpenDoc:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "OpenDoc")
			}
		case 1:
			switch typ := tok.(type) {
			case OpenBlock:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "OpenBlock")
			}
		case 2:
			switch typ := tok.(type) {
			case CloseBlock:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "CloseBlock")
			}
		case 4:
			switch typ := tok.(type) {
			case CloseDoc:
				if tok.Original() != ";" {
					t.Errorf("Got wrong document closing delimiter. Got %v want ;", tok.Original())
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "CloseDoc")
			}
		}
	}
}

var annotations = [...]string{"title ",
	"description\n",
	"parameter\t",
	"table\r"}

func TestBuildAnnotationName(t *testing.T) {
	for _, annotation := range annotations {
		tok := NewTokenizer(strings.NewReader(annotation))
		ann, err := tok.buildAnnotationName()
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}
		if ann != strings.TrimSpace(annotation) {
			t.Errorf("Wrong annotation name. Got %v want %v", ann, strings.TrimSpace(annotation))
		}
	}
}

var bareWords = [...]string{"number ",
	"state\n",
	"dentist\t",
	"alationnaut\r"}

func TestTokenizeBareWord(t *testing.T) {
	for _, bareWord := range bareWords {
		tok := NewTokenizer(strings.NewReader(bareWord))
		err := tok.tokenizeBareWord()
		if err != nil {
			t.Errorf("Unexpected err %v", err)
		}
		if len(tok.tokens) == 0 {
			t.Error("Got zero tokens, expected one bareword token")
		}
		ann := tok.tokens[0].Original()
		if ann != strings.TrimSpace(bareWord) {
			t.Errorf("Wrong bare word name. Got %v want %v", ann, strings.TrimSpace(bareWord))
		}
	}
}

const buildText = `this is a cool bit of text"`

func TestBuildText(t *testing.T) {
	tok := NewTokenizer(strings.NewReader(buildText))
	err := tok.tokenizeText()
	if err != nil {
		t.Errorf("Got unexpected error %v", err)
	}
	if len(tok.tokens) != 1 {
		t.Errorf("Got wrong number of tokens. Got %v want %v", len(tok.tokens), 1)
	}
	switch typ := tok.tokens[0].(type) {
	case Text:
		if typ.Original() != "this is a cool bit of text" {
			t.Errorf("Incorrect text. Got '%v' want '%v'", typ.Original(), "this is a cool bit of text")
		}
	default:
		t.Errorf("Wrong type of token. Got %v want %v", typ, "Text")
	}
}

const wholeDoc = `
some things
that are
maybe queries
this is a haiku
on mount fuji

/**
@title "HEY NOW, YOU'RE A ROCKSTAR"
@description "GET YOUR SHOW ON, GET PAAAAAAAID"
@param you "the rockstar" "the person who should get their show on"

@table {
	@title "All that glitters"
	@description "Is gold"
	@column name "rockstar name" "The name of all of the rockstars"
}
*/
SELECT 1 FROM Tests;
`

func TestWholeDocWithTitle(t *testing.T) {
	tok := NewTokenizer(strings.NewReader(wholeDoc))
	tok.Tokenize()
	if len(tok.tokens) != 21 {
		t.Errorf("Document has incorrect number of tokens. Got %v want %v", len(tok.tokens), 12)
	}
	for i, tok := range tok.tokens {
		switch i {
		case 0:
			switch typ := tok.(type) {
			case OpenDoc:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "OpenDoc")
			}
		case 1:
			switch typ := tok.(type) {
			case OpenBlock:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "OpenBlock")
			}
		case 2:
			switch typ := tok.(type) {
			case Title:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Title")
			}
		case 3:
			switch typ := tok.(type) {
			case Text:
				if typ.Original() != "HEY NOW, YOU'RE A ROCKSTAR" {
					t.Errorf("Incorrect title text. Got '%v' want '%v'", typ.Original(), "HEY NOW, YOU'RE A ROCKSTAR")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Text")
			}
		case 4:
			switch typ := tok.(type) {
			case Desc:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Title")
			}
		case 5:
			switch typ := tok.(type) {
			case Text:
				if typ.Original() != "GET YOUR SHOW ON, GET PAAAAAAAID" {
					t.Errorf("Incorrect title text. Got '%v' want '%v'", typ.Original(), "GET YOUR SHOW ON, GET PAAAAAAAID")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Text")
			}
		case 6:
			switch typ := tok.(type) {
			case Param:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Param")
			}
		case 7:
			switch typ := tok.(type) {
			case BareWord:
				if typ.Original() != "you" {
					t.Errorf("Incorrect bareword text. Got '%v' want '%v'", typ.Original(), "you")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Text")
			}
		case 8:
			switch typ := tok.(type) {
			case Text:
				if typ.Original() != "the rockstar" {
					t.Errorf("Incorrect param text. Got '%v' want '%v'", typ.Original(), "the rockstar")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Text")
			}
		case 9:
			switch typ := tok.(type) {
			case Text:
				if typ.Original() != "the person who should get their show on" {
					t.Errorf("Incorrect param text. Got '%v' want '%v'", typ.Original(), "the person who should get their show on")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Text")
			}
		case 10:
			switch typ := tok.(type) {
			case Table:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Table")
			}
		case 11:
			switch typ := tok.(type) {
			case Title:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Title")
			}
		case 12:
			switch typ := tok.(type) {
			case Text:
				if typ.Original() != "All that glitters" {
					t.Errorf("Incorrect title text. Got '%v' want '%v'", typ.Original(), "All that glitters")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Text")
			}
		case 13:
			switch typ := tok.(type) {
			case Desc:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Desc")
			}
		case 14:
			switch typ := tok.(type) {
			case Text:
				if typ.Original() != "Is gold" {
					t.Errorf("Incorrect description text. Got '%v' want '%v'", typ.Original(), "Is gold")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Text")
			}
		case 15:
			switch typ := tok.(type) {
			case Column:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Column")
			}
		case 16:
			switch typ := tok.(type) {
			case BareWord:
				if typ.Original() != "name" {
					t.Errorf("Incorrect title text. Got '%v' want '%v'", typ.Original(), "name")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Bareword")
			}
		case 17:
			switch typ := tok.(type) {
			case Text:
				if typ.Original() != "rockstar name" {
					t.Errorf("Incorrect title text. Got '%v' want '%v'", typ.Original(), "rockstar name")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Text")
			}
		case 18:
			switch typ := tok.(type) {
			case Text:
				if typ.Original() != "The name of all of the rockstars" {
					t.Errorf("Incorrect title text. Got '%v' want '%v'", typ.Original(), "The name of all of the rockstars")
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "Text")
			}
		case 19:
			switch typ := tok.(type) {
			case CloseBlock:
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "CloseBlock")
			}
		case 20:
			switch typ := tok.(type) {
			case CloseDoc:
				if tok.Original() != ";" {
					t.Errorf("Got wrong document closing delimiter. Got %v want ;", tok.Original())
				}
			default:
				t.Errorf("Got wrong token type at index %v. Got %v want %v", i, typ, "CloseDoc")
			}
		}
	}
}
