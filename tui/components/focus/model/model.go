package model

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle/tui/components/directory"
)

type M struct {
	parent     string
	id         string
	n_elements int // Non-children elements
	index      int
}

func (m *M) Parent() string { return m.parent }
func (m *M) ID() string     { return m.id }

func (m *M) Init() tea.Cmd { return func() tea.Msg { return directory.RegisterNodeMsg{Node: m} } }
func (m *M) View() string  { return "" }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *M) CurrentIndex() int   { return m.index }
func (m *M) Focus(i int) tea.Cmd { return nil } // Default tab advance: increment index and raise EOF if > n_elements
func (m *M) Blur() tea.Cmd       { return nil } // Default blur message: no-op
