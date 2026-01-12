package errors

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
)

type Level int

func (l Level) ToSLog() slog.Level {
	return map[Level]slog.Level{
		LevelDebug: slog.LevelDebug,
		LevelInfo:  slog.LevelInfo,
		LevelWarn:  slog.LevelWarn,
		LevelError: slog.LevelError,
	}[l]
}

func (l Level) String() string {
	return map[Level]string{
		LevelDebug: "d",
		LevelInfo:  "i",
		LevelWarn:  "w",
		LevelError: "e",
	}[l]
}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var _program *tea.Program

func SetProgram(p *tea.Program) {
	_program = p
}

func Error(e error) {
	if _program != nil {
		_program.Send(ToErrorMessage(e))
		_program.Send(ToLogMessage(LevelError, e.Error()))
	}
}

func Warn(e error) {
	if _program != nil {
		_program.Send(ToLogMessage(LevelError, e.Error()))
	}
}

func Info(v string) {
	if _program != nil {
		_program.Send(ToLogMessage(LevelInfo, v))
	}
}

func Debug(v string) {
	if _program != nil {
		_program.Send(ToLogMessage(LevelDebug, v))
	}
}

func ToErrorMessage(e error) ErrorMessage {
	slog.Error(e.Error())
	return ErrorMessage{e: e}
}

func ToLogMessage(l Level, v string) LogMessage {
	parts := strings.Split(v, "\n")
	for _, p := range parts {
		slog.Log(context.Background(), l.ToSLog(), p)
	}
	return LogMessage{
		V: v,
		L: l,
		T: time.Now(),
	}
}

// ErrorMsg may be returned by tea.Cmd in the case of an error.
type ErrorMessage struct {
	e error
}

type LogMessage struct {
	V string
	L Level
	T time.Time
}

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
