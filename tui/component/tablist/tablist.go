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
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
)

type TabType int

const (
	TabTypeNone TabType = iota
	TabTypeVirtual
	TabTypeSource
	TabTypeEdit
)

type Tab struct {
	Type  TabType
	Key   int
	Label string
}

type O struct {
	ParentID string
	Column   *column.C
	Key      form.Key
}

type Node struct {
	*focusable.Node

	column          *column.C
	tabs            []*clickable.Node
	key             form.Key
	values          []Tab
	viewport        viewport.Model
	viewport_offset int
}

func New(o O) *Node {
	return &Node{
		Node:     focusable.New("tablist", o.ParentID, 0),
		column:   o.Column,
		key:      o.Key,
		viewport: viewport.New(o.Column.Content(), 3),
	}
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.SetValue(n.values),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) SetValue(vs []Tab) tea.Cmd {
	var tabs []*clickable.Node

	cmds := []tea.Cmd{
		func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelDebug,
				fmt.Sprintf("%v: calling tablist.SetValue: %v, i = %v", n.ID(), vs, n.FocusIndex()),
			)
		},
	}

	for i := len(n.tabs); i < len(vs); i++ {
		t := clickable.New(n.ID())
		tabs = append(tabs, t)
		cmds = append(cmds, t.Init())
	}
	n.values = append([]Tab{}, vs...)
	n.tabs = append(n.tabs, tabs...)
	n.viewport.SetContent(n.render())

	cmds = append(
		cmds,
		tea.Sequence(
			n.SetNElements(len(vs)),
			func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: after setNelements: i = %v", n.ID(), vs, n.FocusIndex()),
				)
			},
			func() tea.Msg {
				var c tea.Cmd
				if n.FocusIndex() >= n.NElements() {
					c = n.SetFocusIndex(n.NElements() - 1)
				} else if n.FocusIndex() < 0 && n.NElements() > 0 {
					c = n.SetFocusIndex(0)
				}
				cmds = append(cmds, tea.Sequence(
					c,
					n.update_content(),
				))
				return nil
			},
			n.do_highlight(),
		),
	)

	return tea.Batch(cmds...)
}

func (n *Node) Value() form.Value[Tab] {
	if n.FocusIndex() < 0 { return form.Value[Tab]{} }
	return form.Value[Tab]{
		Key:   n.key,
		Value: n.values[n.FocusIndex()],
	}
}

type HighlightMessage struct {
	ID    string
	Value form.Value[Tab]
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	for _, t := range n.tabs {
		_, c = t.Update(msg)
		cmds = append(cmds, c)
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
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) do_highlight() tea.Cmd {
	return func() tea.Msg {
		return HighlightMessage{
			ID:    n.ID(),
			Value: n.Value(),
		}
	}
}

func (n *Node) update_content() tea.Cmd {
	n.viewport.SetContent(n.render())
	return func() tea.Msg {
		tab_widths := make([]int, len(n.values))
		tab_start := make([]int, len(n.values))
		tab_end := make([]int, len(n.values))

		for i := range len(n.values) {
			tab_widths[i] = lipgloss.Width(n.values[i].Label) + /* padding */ 2 + /* borders */ 2
			if i > 0 {
				tab_start[i] = tab_end[i-1]
			}
			tab_end[i] = tab_start[i] + tab_widths[i]
		}

		forward_buffer := 0
		backward_buffer := 0
		if n.FocusIndex() > 2 {
			backward_buffer = tab_widths[n.FocusIndex()-2] + tab_widths[n.FocusIndex()-1]
		} else if n.FocusIndex() > 1 {
			backward_buffer = tab_widths[n.FocusIndex()-1]
		}
		if n.FocusIndex() < len(n.values)-2 {
			forward_buffer = tab_widths[n.FocusIndex()+2] + tab_widths[n.FocusIndex()+1]
		} else if n.FocusIndex() < len(n.values)-1 {
			forward_buffer = tab_widths[n.FocusIndex()+1]
		}
		if n.FocusIndex() <= 0 {
			n.viewport_offset = 0
		} else if tab_start[n.FocusIndex()]-backward_buffer < n.viewport_offset {
			n.viewport_offset = tab_start[n.FocusIndex()] - backward_buffer
		} else if tab_end[n.FocusIndex()]+forward_buffer > n.viewport_offset+n.column.Content() {
			n.viewport_offset = tab_end[n.FocusIndex()] - n.column.Content() + forward_buffer
		}
		n.viewport.SetXOffset(n.viewport_offset)
		return nil
	}
}

func (n *Node) OnFocus(i int) tea.Cmd {
	return tea.Sequence(
		n.Node.OnFocus(i),
		n.update_content(),
		n.do_highlight(),
	)
}

func (n *Node) OnBlur() tea.Cmd {
	return tea.Sequence(
		n.Node.OnBlur(),
		n.update_content(),
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
	for i, v := range n.values {
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
			).Render(v.Label),
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
