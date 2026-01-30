package footer

import (
	"cmp"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/timer"
	"github.com/minkezhang/truffle/tui/util/color_profile"
)

var (
	float_duration = map[errors.Level]time.Duration{
		errors.LevelDebug: time.Second,
		errors.LevelInfo:  2 * time.Second,
		errors.LevelWarn:  5 * time.Second,
		errors.LevelError: 10 * time.Second,
	}
)

type O struct {
	Column   *column.C
	MinLevel errors.Level
}

func New(o O) *Node {
	n := &Node{
		Node:      focusable.New("footer", "", 0),
		column:    o.Column,
		min_level: o.MinLevel,
	}
	n.timer = timer.New(timer.O{
		ParentID: n.ID(),
	})
	return n
}

type Node struct {
	*focusable.Node

	column    *column.C
	q         []errors.LogMessage // FIFO by level
	timer     *timer.Node
	run_id    string
	min_level errors.Level
}

func (n *Node) Init() tea.Cmd {
	return func() tea.Msg {
		return base.RegisterMessage{
			Node: n,
		}
	}
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case errors.LogMessage:
		if msg.L < n.min_level {
			break
		}
		var v errors.LogMessage
		if len(n.q) > 0 {
			v = n.q[0]
		}

		n.q = append(n.q, msg)
		slices.SortFunc(n.q, func(a, b errors.LogMessage) int {
			if l := cmp.Compare(b.L, a.L); l != 0 {
				return l
			}
			return cmp.Compare(a.T.Unix(), b.T.Unix()) // FIFO
		})

		if !n.timer.IsRunning() || (v != errors.LogMessage{}) && n.q[0].L > v.L {
			cmds = append(cmds, n.timer.Start(float_duration[n.q[0].L]))
		}
	case timer.TimerMessage:
		if n.timer.ID() == msg.ID {
			switch msg.Type {
			case timer.TimerMessageStart:
				n.run_id = msg.RunID
			case timer.TimerMessageStop:
				if len(n.q) > 0 {
					n.q = n.q[1:]
				}
				if len(n.q) > 0 {
					cmds = append(cmds, n.timer.Start(float_duration[n.q[0].L]))
				}
			}
		}
	}

	_, c := n.timer.Update(msg)
	cmds = append(cmds, c)

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	var v errors.LogMessage
	if len(n.q) > 0 {
		v = n.q[0]
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
