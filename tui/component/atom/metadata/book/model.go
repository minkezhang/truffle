package book

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"
)

var (
	_ tea.Model = &M{}
)

type O struct {
	Book  *book.M
	Width int
}

func New(o O) *M {
	return &M{
		book:  o.Book,
		width: o.Width,
	}
}

type M struct {
	book  *book.M
	hover bool
	width int
}
