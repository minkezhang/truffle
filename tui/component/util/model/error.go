package model

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/charmbracelet/bubbletea"
)

var (
	// E is the global error model.
	E = &Error{}
)

func ToErrorMsg(e error) ErrorMsg {
	slog.Error(e.Error())
	return ErrorMsg{e: e}
}

func ToWarningMsg(e error) WarningMsg {
	slog.Warn(e.Error())
	return WarningMsg{
		e: fmt.Errorf(
			"WARNING(%v): %v",
			time.Now().Format("15:04:05"),
			e.Error(),
		),
	}
}

func ToNoticeMsg(v string) NoticeMsg {
	slog.Info(v)
	return NoticeMsg(
		fmt.Sprintf(
			"NOTICE(%v): %v",
			time.Now().Format("15:04:05"),
			v,
		),
	)
}

// ErrorMsg may be returned by tea.Cmd in the case of an error.
//
// Callers should use the constructors ToError() instead.
type ErrorMsg struct {
	e error
}

func (e ErrorMsg) Error() string { return e.e.Error() }

type WarningMsg ErrorMsg

func (w WarningMsg) Error() string { return w.e.Error() }

type NoticeMsg string

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
		m.Append(error(msg.e))
	}

	if err := m.Error(); err != nil {
		return m, tea.Quit
	}

	return m, nil
}

// View fulfills the tea.Model interface and is unused.
func (m *Error) View() string { return "" }
