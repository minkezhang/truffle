package model

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/components/directory"
)

type O struct {
	Parent       string
	IDPrefix     string
	FocusHandler directory.FocusHandler
}

func New(o O) *M {
	prefix := o.IDPrefix
	if prefix == "" {
		prefix = "none"
	}
	return &M{
		parent: o.Parent,
		id:     fmt.Sprintf("%v:%v", o.IDPrefix, zone.NewPrefix()),
		h:      o.FocusHandler,
	}
}

type M struct {
	parent      string
	id          string
	n_elements  int // Non-children elements
	index       int
	focus_state directory.FocusState

	h directory.FocusHandler
}

func (m *M) Parent() string                   { return m.parent }
func (m *M) ID() string                       { return m.id }
func (m *M) FocusState() directory.FocusState { return m.focus_state }
func (m *M) Init() tea.Cmd                    { return func() tea.Msg { return directory.RegisterNode{Node: m} } }
func (m *M) View() string                     { return "" }
func (m *M) CurrentIndex() int                { return m.index }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *M) NElements() int {
	if m.h != nil {
		return m.h.NElements()
	}
	return 0
}

func (m *M) OnFocus(i int) tea.Cmd {
	m.focus_state = directory.FocusStateActive
	if m.h != nil {
		return m.h.OnFocus(i)
	}
	// Default tab advance: increment index and raise EOF if > n_elements
	return nil
}

func (m *M) OnBlur() tea.Cmd {
	m.focus_state = directory.FocusStateNone
	if m.h != nil {
		return m.h.OnBlur()
	}
	// Default blur message: no-op
	return nil
}
