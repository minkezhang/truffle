package footer

import (
	"cmp"
	"slices"
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
		q:      map[errors.Level][]errors.LogMessage{},
	}
}

type Node struct {
	*focusable.Node

	column *column.C
	q      map[errors.Level][]errors.LogMessage
}

func (n *Node) Init() tea.Cmd {
	return func() tea.Msg {
		return base.RegisterMessage{
			Node: n,
		}
	}
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) { // TODO(minkezhang): Add timer to expire messages.
	case errors.LogMessage:
		n.q[msg.L] = append(n.q[msg.L], msg)
		slices.SortFunc(n.q[msg.L], func(a, b errors.LogMessage) int {
			return cmp.Compare(a.T.Unix(), b.T.Unix()) // FIFO
		})
	}
	return n, nil
}

func (n *Node) View() string {
	var v errors.LogMessage
	for _, l := range []errors.Level{
		errors.LevelError,
		errors.LevelWarn,
		errors.LevelInfo,
		errors.LevelDebug,
	} {
		if ms := n.q[l]; len(ms) > 0 {
			v = ms[0]
			break
		}
	}

	if v == (errors.LogMessage{}) {
		return ""
	}

	return n.column.RenderOrDie(
		n.column.Style().Inline(true).Foreground(
			map[errors.Level]lipgloss.TerminalColor{
				errors.LevelDebug: color_profile.ForegroundNegligible,
				errors.LevelInfo:  color_profile.ForegroundInverted,
				errors.LevelWarn:  color_profile.ForegroundInverted,
				errors.LevelError: color_profile.ForegroundInverted,
			}[v.L],
		).Background(
			map[errors.Level]lipgloss.TerminalColor{
				errors.LevelDebug: lipgloss.NoColor{},
				errors.LevelInfo:  color_profile.LogForeground[errors.LevelInfo],
				errors.LevelWarn:  color_profile.LogForeground[errors.LevelWarn],
				errors.LevelError: color_profile.LogForeground[errors.LevelError],
			}[v.L],
		).Render(
			runewidth.Truncate(
				strings.ReplaceAll(v.V, "\n", " "),
				n.column.Content(),
				"…",
			),
		),
	)
}
