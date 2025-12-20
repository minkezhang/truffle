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
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/util/grid"
)

type Model interface {
	tea.Model

	Column() grid.C
	Focus() (tea.Model, tea.Cmd) // TODO
	Blur() (tea.Model, tea.Cmd)  // TODO
}

type O struct {
	Column grid.C
}

type BaseMsg struct {
	ID string
}

type Base struct {
	id     string // Runtime UUID
	column grid.C
}

func New(o O) *Base {
	return &Base{
		id:     zone.NewPrefix(),
		column: o.Column,
	}
}

// ID is used to uniquely identify a model.
//
// This is useful for passing messages to specific models.
func (m Base) ID() string { return m.id }

// Init fulfills the tea.Model interface and is unused.
func (m Base) Init() tea.Cmd { return nil }

// Update fulfills the tea.Model interface and is unused.
func (m Base) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

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
