package book

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"
	"github.com/minkezhang/truffle/tui/util/grid"
)

var (
	_ tea.Model = &M{}
)

type O struct {
	Book   *book.M
	Column grid.C
}

func New(o O) *M {
	return &M{
		book:   o.Book,
		column: o.Column,
	}
}

type M struct {
	column grid.C

	book *book.M
}
