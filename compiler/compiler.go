package compiler

import (
	"errors"

	"github.com/christopher-henderson/DocStringParser/tokenizer"
)

type QueryDoc struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Params      []Param `json:"params"`
	Output      Table   `json:"output"`
}

func NewQueryDocs() (q QueryDoc) {
	q.Params = make([]Param, 0)
	return q
}

type Table struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Columns     []Column `json:"columns"`
}

func NewTable() Table {
	return Table{Columns: make([]Column, 0)}
}

type Param struct {
	ProperName  string `json:"properName"`
	Blurb       string `json:"blurb"`
	Description string `json:"description"`
}

type Column struct {
	ProperName  string `json:"properName"`
	Blurb       string `json:"blurb"`
	Description string `json:"description"`
}

type Compiler struct {
	Tokens  []tokenizer.Tokener
	docList []QueryDoc
	state   int
}

func NewCompiler(tokens []tokenizer.Tokener) Compiler {
	return Compiler{Tokens: tokens, docList: make([]QueryDoc, 0), state: 0}
}

func Compile(tokens []tokenizer.Tokener) ([]QueryDoc, error) {
	c := Compiler{Tokens: tokens}
	return c.compile()
}

func (c *Compiler) next() (tokenizer.Tokener, error) {
	if c.state >= len(c.Tokens) {
		return nil, errors.New("no")
	}
	defer func() {
		c.state += 1
	}()
	return c.Tokens[c.state], nil
}

func (c *Compiler) compile() ([]QueryDoc, error) {
	for doc, err := c.next(); err == nil; doc, err = c.next() {
		switch doc.(type) {
		case tokenizer.OpenDoc:
			qdoc, err := c.compileDoc()
			if err != nil {
				return nil, err
			}
			c.docList = append(c.docList, qdoc)
		default:
			return nil, errors.New("Unexpected type at top leve.")
		}
	}
	return c.docList, nil
}

func (c *Compiler) compileDoc() (QueryDoc, error) {
	q := NewQueryDocs()
	for t, err := c.next(); err == nil; t, err = c.next() {
		switch t.(type) {
		case tokenizer.Title:
			title, err := c.next()
			if err != nil {
				return q, err
			}
			switch title.(type) {
			case tokenizer.Text:
				q.Title = title.Original()
			default:
				return q, errors.New("Bad parse tree")
			}
		case tokenizer.Desc:
			desc, err := c.next()
			if err != nil {
				return q, err
			}
			switch desc.(type) {
			case tokenizer.Text:
				q.Description = desc.Original()
			default:
				return q, errors.New("Bad parse tree")
			}
		case tokenizer.Param:
			p, err := c.compileParam()
			if err != nil {
				return q, err
			}
			q.Params = append(q.Params, p)
		case tokenizer.Table:
			t, err := c.compileTable()
			if err != nil {
				return q, err
			}
			q.Output = t
		case tokenizer.CloseDoc:
			return q, nil
		}
	}
	return q, nil
}

func (c *Compiler) compileParam() (p Param, err error) {
	name, err := c.next()
	if err != nil {
		return
	}
	switch name.(type) {
	case tokenizer.BareWord:
		p.ProperName = name.Original()
	}
	blurb, err := c.next()
	if err != nil {
		return
	}
	switch blurb.(type) {
	case tokenizer.Text:
		p.Blurb = blurb.Original()
	}
	description, err := c.next()
	if err != nil {
		return
	}
	switch description.(type) {
	case tokenizer.Text:
		p.Description = description.Original()
	}

	return
}

func (c *Compiler) compileTable() (Table, error) {
	table := NewTable()
	var err error
	for t, err := c.next(); err == nil; t, err = c.next() {
		switch t.(type) {
		case tokenizer.Title:
			title, err := c.next()
			if err != nil {
				return table, err
			}
			switch title.(type) {
			case tokenizer.Text:
				table.Title = title.Original()
			default:
				return table, errors.New("Bad parse tree")
			}
		case tokenizer.Desc:
			desc, err := c.next()
			if err != nil {
				return table, err
			}
			switch desc.(type) {
			case tokenizer.Text:
				table.Description = desc.Original()
			default:
				return table, errors.New("Bad parse tree")
			}
		case tokenizer.Column:
			c, err := c.compileColumn()
			if err != nil {
				return table, nil
			}
			table.Columns = append(table.Columns, c)
		default:
			c.state -= 1
			return table, nil
		}
	}
	return table, err
}

func (c *Compiler) compileColumn() (col Column, err error) {
	name, err := c.next()
	if err != nil {
		return
	}
	switch name.(type) {
	case tokenizer.BareWord:
		col.ProperName = name.Original()
	}
	blurb, err := c.next()
	if err != nil {
		return
	}
	switch blurb.(type) {
	case tokenizer.Text:
		col.Blurb = blurb.Original()
	}
	description, err := c.next()
	if err != nil {
		return
	}
	switch description.(type) {
	case tokenizer.Text:
		col.Description = blurb.Original()
	}

	return
}
