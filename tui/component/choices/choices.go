package choices

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
)

type O struct {
	ParentID string
	Width    int
	Height   int
	Key      form.Key
	Choices  []form.Key
}

type Node struct {
	*focusable.Node

	max_width   int
	clickable   *clickable.Node
	key         form.Key
	choices     []form.Key
	index       int
	height      int
	start_index int
}

func New(o O) *Node {
	n := &Node{
		Node:      focusable.New("choices", o.ParentID, 1),
		max_width: o.Width,
		key:       o.Key,
		choices:   o.Choices,
		height:    o.Height,
	}
	n.clickable = clickable.New(n.ID())
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.clickable.Init(),
		n.SetValue(n.choices, 0),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) SetValue(vs []form.Key, index int) tea.Cmd {
	n.choices = vs
	n.index = index
	n.start_index = 0
	return nil
}

type HighlightMessage struct {
	ID    string
	Value form.Value[form.Key]
}

func (n *Node) do_scroll(is_up bool) tea.Cmd {
	if is_up {
		if n.index > 0 {
			n.index -= 1
			if n.start_index > n.index {
				n.start_index -= 1
			}
		}
	} else {
		if n.index < len(n.choices)-1 {
			n.index += 1
			if end_index := n.start_index + n.height; end_index <= n.index {
				n.start_index += 1
			}
		}
	}
	return func() tea.Msg {
		return HighlightMessage{
			ID:    n.ID(),
			Value: n.Value(),
		}
	}
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.MouseMsg:
			if msg.Button == tea.MouseButtonWheelDown {
				cmds = append(cmds, n.do_scroll(false))
			}
			if msg.Button == tea.MouseButtonWheelUp {
				cmds = append(cmds, n.do_scroll(true))
			}
		case tea.KeyMsg:
			if msg.Type == tea.KeyDown {
				cmds = append(cmds, n.do_scroll(false))
			}
			if msg.Type == tea.KeyUp {
				cmds = append(cmds, n.do_scroll(true))
			}
			if msg.Type == tea.KeySpace || msg.Type == tea.KeyEnter || msg.Type == tea.KeyEsc {
				cmds = append(cmds, func() tea.Msg {
					return types.EOFMessage{
						ID:     n.ID(),
						IsHead: true,
					}
				})
			}
		case clickable.Click:
			if msg.ID == n.clickable.ID() {
				cmds = append(cmds, func() tea.Msg {
					return types.EOFMessage{
						ID:     n.ID(),
						IsHead: true,
					}
				})
			}
		}
	}

	_, c = n.clickable.Update(msg)
	cmds = append(cmds, c)

	return n, tea.Batch(cmds...)
}

func (n *Node) Value() form.Value[form.Key] {
	if n.index >= len(n.choices) {
		return form.Value[form.Key]{
			Key:   n.key,
			Value: form.Key{},
		}
	}
	return form.Value[form.Key]{
		Key:   n.key,
		Value: n.choices[n.index],
	}
}

func (n *Node) View() string {
	var rows []string
	end_index := n.start_index + n.height
	if end_index >= len(n.choices) {
		end_index = len(n.choices)
	}

	for i, c := range n.choices[n.start_index:end_index] {
		w := n.max_width - /* padding-left */ 1 - lipgloss.Width("┃")
		style := lipgloss.NewStyle().Width(w).MaxWidth(w).Inline(true)
		row := lipgloss.NewStyle().PaddingLeft(1).Foreground(
			map[bool]lipgloss.TerminalColor{
				false: color_profile.UIForeground[n.FocusState()],
				true:  color_profile.ForegroundInverted,
			}[i+n.start_index == n.index],
		).Background(
			map[bool]lipgloss.TerminalColor{
				true:  color_profile.UIForeground[n.FocusState()],
				false: color_profile.ForegroundInverted,
			}[i+n.start_index == n.index],
		).Render(
			style.Render(ansi.Truncate(c.Label, w, "…")),
		)
		rows = append(rows, fmt.Sprintf("%v%v", lipgloss.NewStyle().Foreground(
			color_profile.UIForeground[n.FocusState()],
		).Render("┃"), row))
	}

	parts := []string{
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
	}

	if len(n.choices) > n.height {
		parts = append(parts,
			lipgloss.NewStyle().Foreground(color_profile.SupplementaryText).Render(
				fmt.Sprintf("%d / %d", n.index+1, len(n.choices)),
			),
		)
	}

	return zone.Mark(
		n.clickable.ID(),
		lipgloss.NewStyle().Width(n.max_width).Render(
			lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(
				color_profile.UIForeground[n.FocusState()],
			).Render(
				lipgloss.JoinVertical(
					lipgloss.Right,
					parts...,
				),
			),
		),
	)

}

func (n *Node) OnBlur() tea.Cmd {
	return tea.Sequence(n.Node.OnBlur(), n.SetIsInvisible(true))
}
