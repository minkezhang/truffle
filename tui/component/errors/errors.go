package errors

import (
	"errors"
	"log/slog"

	"github.com/charmbracelet/bubbletea"
)

var _program *tea.Program

func SetProgram(p *tea.Program) {
	_program = p
}

func Error(e error) {
	if _program != nil {
		_program.Send(ToErrorMessage(e))
	}
}

func Warn(e error) {
	if _program != nil {
		_program.Send(ToWarnMessage(e))
	}
}

func Info(v string) {
	if _program != nil {
		_program.Send(ToInfoMessage(v))
	}
}

func Debug(v string) {
	if _program != nil {
		_program.Send(ToDebugMessage(v))
	}
}

func ToErrorMessage(e error) ErrorMessage {
	slog.Error(e.Error())
	return ErrorMessage{e: e}
}

func ToWarnMessage(e error) WarnMessage {
	slog.Warn(e.Error())
	return WarnMessage{e: e}
}

func ToInfoMessage(v string) InfoMessage {
	slog.Info(v)
	return InfoMessage(v)
}

func ToDebugMessage(v string) DebugMessage {
	slog.Debug(v)
	return DebugMessage(v)
}

// ErrorMsg may be returned by tea.Cmd in the case of an error.
type ErrorMessage struct {
	e error
}

type WarnMessage ErrorMessage
type InfoMessage string
type DebugMessage InfoMessage

type Node struct {
	errors []error
}

func (n *Node) Init() tea.Cmd { return nil }
func (n *Node) Error() error  { return errors.Join(n.errors...) }
func (n *Node) View() string  { return "" }

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrorMessage:
		n.errors = append(n.errors, msg.e)
	}

	if err := n.Error(); err != nil {
		return n, tea.Quit
	}

	return n, nil
}
