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
	Layout grid.L
}

func New(o O) *M {
	return &M{
		book:   o.Book,
		layout: o.Layout,
	}
}

type M struct {
	layout grid.L

	book *book.M
}
