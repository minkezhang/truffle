package tablist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/util/color_profile"
)

type O struct {
	ParentID string
	Column   *column.C
}

type Node struct {
	*focusable.Node

	column   *column.C
	tabs     []*clickable.Node
	labels   []string
	viewport viewport.Model
}

func New(o O) *Node {
	return &Node{
		Node:     focusable.New("tablist", o.ParentID, 0),
		column:   o.Column,
		viewport: viewport.New(o.Column.Content(), 3),
	}
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.SetValue(n.labels),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) SetValue(vs []string) tea.Cmd {
	var tabs []*clickable.Node

	var cmds []tea.Cmd
	for i := len(n.tabs); i < len(vs); i++ {
		t := clickable.New(n.ID())
		tabs = append(tabs, t)
		cmds = append(cmds, t.Init())
	}
	n.labels = append([]string{}, vs...)
	n.tabs = append(n.tabs, tabs...)
	n.viewport.SetContent(n.render())

	cmds = append(
		cmds,
		n.SetNElements(len(vs)),
	)

	return tea.Batch(cmds...)
}

type HighlightMessage struct {
	ID    string
	Index int
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	for _, t := range n.tabs {
		_, c = t.Update(msg)
		cmds = append(cmds, c)
	}

	if n.FocusState() == types.FocusStateActive {
		if msg, ok := msg.(tea.KeyMsg); ok {
			n.viewport.SetHorizontalStep(5)
			n.viewport, c = n.viewport.Update(msg)
			cmds = append(cmds, c, func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: viewport processing key message %v", n.ID(), msg),
				)
			})
		}
	}

	switch msg := msg.(type) {
	case clickable.Click:
		for i, t := range n.tabs {
			if msg.ID == t.ID() && (n.FocusState() == types.FocusStateNone || (n.FocusState() == types.FocusStateActive && n.FocusIndex() != i)) {
				cmds = append(cmds, func() tea.Msg {
					return types.FocusMessage{
						ID:    n.ID(),
						Index: i,
					}
				})
			}
		}
	case types.FocusMessage:
		if msg.ID == n.ID() {
			cmds = append(cmds, func() tea.Msg {
				return HighlightMessage{
					ID:    n.ID(),
					Index: msg.Index,
				}
			})
		}
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) set_content(s string) tea.Cmd {
	n.viewport.SetContent(s)
	return nil
}

func (n *Node) OnFocus(i int) tea.Cmd {
	return tea.Sequence(
		n.Node.OnFocus(i),
		n.set_content(n.render()),
	)
}

func (n *Node) OnBlur() tea.Cmd {
	return tea.Sequence(
		n.Node.OnBlur(),
		n.set_content(n.render()),
	)
}

var (
	selected_border = lipgloss.Border{
		Top:      lipgloss.RoundedBorder().Top,
		Left:     lipgloss.RoundedBorder().Left,
		Right:    lipgloss.RoundedBorder().Right,
		TopLeft:  lipgloss.RoundedBorder().TopLeft,
		TopRight: lipgloss.RoundedBorder().TopRight,
	}
)

func (n *Node) render() string {
	var parts []string
	var border []string
	for i, l := range n.labels {
		p := zone.Mark(
			n.tabs[i].ID(),
			lipgloss.NewStyle().Border(
				lipgloss.RoundedBorder(), true, true, false, true,
			).BorderForeground(
				map[bool]lipgloss.Color{
					false: color_profile.UIForeground[types.FocusStateNone],
					true:  color_profile.UIForeground[types.FocusStateActive],
				}[i == n.FocusIndex() && n.FocusState() == types.FocusStateActive],
			).Padding(0, 1).Foreground(
				map[bool]lipgloss.Color{
					false: color_profile.UIForeground[types.FocusStateNone],
					true:  color_profile.UIForeground[types.FocusStateActive],
				}[i == n.FocusIndex() && n.FocusState() == types.FocusStateActive],
			).Render(l),
		)
		parts = append(parts, p)
		if i == n.FocusIndex() {
			border = append(border, fmt.Sprintf("╯%v╰", strings.Repeat(" ", lipgloss.Width(p)-2)))
		} else {
			border = append(border, strings.Repeat("─", lipgloss.Width(p)))
		}
	}
	remainder := n.column.Content() - lipgloss.Width(strings.Join(border, ""))
	if remainder < 0 {
		remainder = 0
	}
	border = append(border, strings.Repeat("─", remainder))
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, parts...),
		lipgloss.NewStyle().Foreground(
			color_profile.UIForeground[n.FocusState()],
		).Render(strings.Join(border, "")),
	)
}

func (n *Node) View() string {
	return n.viewport.View()
}
