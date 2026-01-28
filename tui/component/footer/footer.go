package footer

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
)

type O struct {
	Column *column.C
}

func New(o O) *Node {
	return &Node{
		Node:   focusable.New("footer", "", 0),
		column: o.Column,
	}
}

type Node struct {
	*focusable.Node

	column *column.C
	m      errors.LogMessage
}

func (n *Node) Init() tea.Cmd {
	return func() tea.Msg {
		return base.RegisterMessage{
			Node: n,
		}
	}
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errors.LogMessage:
		n.m = msg
	}
	return n, nil
}

func (n *Node) View() string {
	if n.m == (errors.LogMessage{}) {
		return ""
	}
	return n.column.RenderOrDie(
		n.column.Style().Inline(true).Foreground(
			map[errors.Level]lipgloss.TerminalColor{
				errors.LevelDebug: color_profile.ForegroundNegligible,
				errors.LevelInfo:  color_profile.ForegroundInverted,
				errors.LevelWarn:  color_profile.ForegroundInverted,
				errors.LevelError: color_profile.ForegroundInverted,
			}[n.m.L],
		).Background(
			map[errors.Level]lipgloss.TerminalColor{
				errors.LevelDebug: lipgloss.NoColor{},
				errors.LevelInfo:  color_profile.LogForeground[errors.LevelInfo],
				errors.LevelWarn:  color_profile.LogForeground[errors.LevelWarn],
				errors.LevelError: color_profile.LogForeground[errors.LevelError],
			}[n.m.L],
		).Render(
			runewidth.Truncate(
				strings.ReplaceAll(n.m.V, "\n", " "),
				n.column.Content(),
				"…",
			),
		),
	)
}
