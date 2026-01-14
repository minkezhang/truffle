package textarea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
)

type Node struct {
	*focusable.Node

	max_width int
	input     textarea.Model
	clickable *clickable.Node
	key       form.Key
}

type O struct {
	Prefix      string
	ParentID    string
	Width       int
	Height      int
	Placeholder string
	Label       string
	Value       form.Value[string]
}

func New(o O) *Node {
	t := textarea.New()
	t.SetWidth(o.Width)
	t.SetHeight(o.Height)
	t.Placeholder = o.Placeholder
	t.SetValue(o.Value.Value)
	t.ShowLineNumbers = false
	t.FocusedStyle.Prompt = t.BlurredStyle.Prompt.Foreground(color_profile.UIForeground[types.FocusStateActive])
	t.BlurredStyle.Prompt = t.BlurredStyle.Prompt.Foreground(color_profile.UIForeground[types.FocusStateNone])

	n := &Node{
		Node:      focusable.New(o.Prefix, o.ParentID, 1),
		input:     t,
		max_width: o.Width,
		key:       o.Value.Key,
	}
	n.clickable = clickable.New(n.ID())
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Batch(
		tea.Sequence(
			n.clickable.Init(),
			func() tea.Msg {
				return base.RegisterMessage{
					Node: n,
				}
			},
		),
		// Workaround -- it appears BlurredStyle.Label is not applied
		// until textarea.Blur() is explicitly called.
		tea.Sequence(
			n.OnFocus(0),
			n.OnBlur(),
		),
	)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.MouseMsg:
			if msg.Button == tea.MouseButtonWheelUp {
				n.input, c = n.input.Update(tea.KeyMsg{Type: tea.KeyUp})
				cmds = append(cmds, c)
			} else if msg.Button == tea.MouseButtonWheelDown {
				n.input, c = n.input.Update(tea.KeyMsg{Type: tea.KeyDown})
				cmds = append(cmds, c)
			}
		case tea.KeyMsg:
			n.input, c = n.input.Update(msg)
			cmds = append(cmds, c)
		}
		n.input.Cursor, c = n.input.Cursor.Update(msg)
		cmds = append(cmds, c)
	} else {
		switch msg := msg.(type) {
		case clickable.Click:
			if msg.ID == n.clickable.ID() {
				cmds = append(cmds, func() tea.Msg {
					return types.FocusMessage{
						ID:    n.ID(),
						Index: 0,
					}
				})
			}
		}
		_, c = n.clickable.Update(msg)
		cmds = append(cmds, c)
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return zone.Mark(
		n.clickable.ID(),
		lipgloss.NewStyle().Border(
			lipgloss.NormalBorder(), false, false, true, false,
		).BorderForeground(
			color_profile.UIForeground[n.FocusState()],
		).Width(n.max_width).MaxWidth(n.max_width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				fmt.Sprintf(
					"%v %v",
					n.key.Label,
					lipgloss.NewStyle().Foreground(color_profile.UIForeground[n.FocusState()]).MarginBottom(1).Render(
						strings.Repeat("─", n.max_width-len(n.key.Label)-1),
					),
				),
				n.input.View(),
				lipgloss.NewStyle().Foreground(color_profile.SupplementaryText).Render(
					fmt.Sprintf("line %d / %d", n.input.Line()+1, n.input.LineCount()),
				),
			),
		),
	)
}

func (n *Node) OnFocus(i int) tea.Cmd {
	cmds := []tea.Cmd{n.Node.OnFocus(i)}
	if n.FocusState() == types.FocusStateActive {
		cmds = append(cmds,
			n.input.Focus(),
			n.input.Cursor.Focus(),
		)
	}
	return tea.Sequence(cmds...)
}

func (n *Node) OnBlur() tea.Cmd {
	n.input.Blur()
	n.input.Cursor.Blur()
	return tea.Sequence(
		n.Node.OnBlur(),
	)
}
