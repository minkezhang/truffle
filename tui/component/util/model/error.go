package model

import (
	"errors"

	"github.com/charmbracelet/bubbletea"
)

var (
	// E is the global error model.
	E = Error{}
)

// ErrorMsg may be returned by tea.Cmd in the case of an error.
type ErrorMsg error

type Error struct {
	errors []error
}

func (m *Error) Error() error   { return errors.Join(m.errors...) }
func (m *Error) Append(v error) { m.errors = append(m.errors, v) }

// Init fulfills the tea.Model interface and is unused.
func (m *Error) Init() tea.Cmd { return nil }

func (m *Error) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrorMsg:
		m.Append(error(msg))
	}

	if err := m.Error(); err != nil {
		return m, tea.Quit
	}

	return m, nil
}

// View fulfills the tea.Model interface and is unused.
func (m *Error) View() string { return "" }
