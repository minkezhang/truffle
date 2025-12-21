// Package model is a Truffle-specific tea.Model interface with some rendering
// bolt-ons.
//
// Specifically, model.Model
//
//  1. must always ensure the rendered output is within the bounding box
//  2. will quit if the error field is set -- this may be set in
//
// TODO(minkezhang): Change back to New()
package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/util/grid"
)

func ToCommand(msg tea.Msg) tea.Cmd { return func() tea.Msg { return msg } }

// BlurMsg is published by a specific node which instructs the parent node to
// advance the internal tab index. The ID included in the message is the
// originating node.
//
// Nodes that need to manage a tab index must handle this message.
//
// Upon node deletion, the node must publish a BlurMsg with IsEnd = false. This
// will allow the parent node to appropriately handle this event.
type BlurMsg struct {
	BaseMsg
	IsEnd bool
}

// FocusMsg is published by a specific node takes control due to a mouse event.
// All other nodes will stop handling keyboard input (except for root which
// handles Ctrl + C). The ID included in the message is the originating node.
//
// All nodes which handles user input in some way must handle this message.
type FocusMsg struct {
	BaseMsg
	Index int
}

type O struct {
	Column grid.C
	NTabs  int
}

type BaseMsg struct {
	ID string
}

type Base struct {
	id     string // Runtime UUID
	column grid.C
	focus  bool
	index  int
	tabs   int
}

func New(o O) *Base {
	m := Make(o)
	return &m
}

func Make(o O) Base {
	return Base{
		id:     zone.NewPrefix(),
		column: o.Column,
		focus:  false,
		index:  0,
		tabs:   o.NTabs,
	}
}

func (m Base) Focus() bool      { return m.focus }
func (m *Base) SetFocus(v bool) { m.focus = v }
func (m Base) Index() int       { return m.index }
func (m Base) NTabs() int       { return m.tabs }
func (m Base) SetNTabs(v int)   { m.tabs = v }

func (m *Base) SetIndex(v int) int {
	i := m.index
	m.index = v
	return i
}

// ID is used to uniquely identify a model.
//
// This is useful for passing messages to specific models.
func (m Base) ID() string { return m.id }

// Init fulfills the tea.Model interface and is unused.
func (m Base) Init() tea.Cmd { return nil }

// Update fulfills the tea.Model interface and is unused.
func (m Base) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FocusMsg:
		m.SetFocus(m.ID() == msg.ID)
		if m.Focus() {
			if msg.Index >= m.NTabs() {
				m.SetIndex(msg.Index % (m.NTabs()))
				// TODO
			} else if msg.Index < 0 {
				m.SetIndex(msg.Index + m.NTabs())
			} else {
				m.SetIndex(msg.Index)
			}
		}
		return m, nil
	case BlurMsg:
		m.SetFocus(m.ID() != msg.ID)
	case tea.KeyMsg:
		if !m.Focus() {
			break
		}
		switch t := msg.Type; t {
		case tea.KeyShiftTab:
			dst := m.Index() - 1
			if dst < 0 {
				return m, ToCommand(BlurMsg{
					BaseMsg: BaseMsg{
						ID: m.ID(),
					},
					IsEnd: false,
				})
			} else {
				m.SetIndex(dst)
				return m, ToCommand(FocusMsg{
					BaseMsg: BaseMsg{
						ID: m.ID(),
					},
					Index: m.Index(),
				})
			}
		case tea.KeyTab:
			dst := m.Index() + 1
			if dst >= m.NTabs() {
				return m, ToCommand(BlurMsg{
					BaseMsg: BaseMsg{
						ID: m.ID(),
					},
					IsEnd: false,
				})
			} else {
				m.SetIndex(dst)
				return m, ToCommand(FocusMsg{
					BaseMsg: BaseMsg{
						ID: m.ID(),
					},
					Index: m.Index(),
				})
			}
		}
	}
	return m, nil
}

// View fulfills the tea.Model interface and is unused.
func (m Base) View() string { return "" }

func (m Base) Column() grid.C { return m.column }

// ValidateOrDie will be called from parent Model.View() functions.
//
// Example:
//
//	func (m *M) View() string {
//	  return m.RenderOrDie(lipgloss.NewStyle().Render(...))
//	}
func (m Base) RenderOrDie(s string) string {
	if err := check(s, m.column); err != nil {
		E.Append(err)
		return ""
	}
	return s
}

func axis(w int) string {
	return strings.Repeat("|----:----", w/10) + []string{
		"",
		"|",
		"|-",
		"|--",
		"|---",
		"|----",
		"|----:",
		"|----:-",
		"|----:--",
		"|----:---",
	}[w%10]
}

func check(s string, c grid.C) error {
	if w := lipgloss.Width(s); w > c.Width() {
		return fmt.Errorf(
			"rendered string exceeded bounding box: %d > %d\n"+
				"```\n"+
				"%s\n"+
				"%s\n"+
				"```\n",
			w,
			c.Content,
			axis(w),
			strings.ReplaceAll(s, " ", "."),
		)
	}
	return nil
}
