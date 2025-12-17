// Package model is a Truffle-specific tea.Model interface with some rendering
// bolt-ons.
//
// Specifically, model.Model
//
//  1. must always ensure the rendered output is within the bounding box
//  2. will quit if the error field is set -- this may be set in
package model

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/util/grid"
)

type Model interface {
	tea.Model
	Column() grid.C
	Focus() (tea.Model, tea.Cmd)
	Blur() (tea.Model, tea.Cmd)
}

type O struct {
	Column grid.C
}

type Base struct {
	column grid.C
}

func New(o O) *Base {
	return &Base{
		column: o.Column,
	}
}

// Init fulfills the tea.Model interface and is unused.
func (m *Base) Init() tea.Cmd { return nil }

// Update fulfills the tea.Model interface and is unused.
func (m *Base) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

// View fulfills the tea.Model interface and is unused.
func (m *Base) View() string { return "" }

func (m *Base) Column() grid.C { return m.column }

// ValidateOrDie will be called from parent Model.View() functions.
//
// Example:
//
//	func (m *M) View() string {
//	  return m.RenderOrDie(lipgloss.NewStyle().Render(...))
//	}
func (m *Base) RenderOrDie(s string) string {
	if w := lipgloss.Width(s); w > m.column.Content {
		E.Append(fmt.Errorf("rendered string exceeded bounding box: %d > %d", w, m.column.Content))
		return ""
	}
	return s
}
