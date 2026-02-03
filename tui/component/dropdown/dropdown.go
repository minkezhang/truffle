package dropdown

import (
	"fmt"
	"strings"

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
	Prefix   string
	ParentID string
	Width    int
	Key      form.Key
	Choices  []form.Key
	Prompt   string
	Height   int
}

type Node struct {
	*focusable.Node

	max_width   int
	prompt      string
	clickable   *clickable.Node
	key         form.Key
	choices     []form.Key
	index       int
	height      int
	start_index int
}

func New(o O) *Node {
	n := &Node{
		Node:      focusable.New(o.Prefix, o.ParentID, 1),
		max_width: o.Width,
		key:       o.Key,
		choices:   o.Choices,
		prompt:    o.Prompt,
		height:    o.Height,
	}
	n.clickable = clickable.New(n.ID())
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.clickable.Init(),
		n.SetValue(n.choices),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) SetValue(vs []form.Key) tea.Cmd {
	n.choices = vs
	n.index = 0
	n.start_index = 0
	return nil
}

func (n *Node) scroll(is_up bool) tea.Cmd {
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
	return nil
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.MouseMsg:
			if msg.Button == tea.MouseButtonWheelDown {
				cmds = append(cmds, n.scroll(false))
			}
			if msg.Button == tea.MouseButtonWheelUp {
				cmds = append(cmds, n.scroll(true))
			}
		case tea.KeyMsg:
			if msg.Type == tea.KeyDown {
				cmds = append(cmds, n.scroll(false))
			}
			if msg.Type == tea.KeyUp {
				cmds = append(cmds, n.scroll(true))
			}
		}
	} else {
		switch msg := msg.(type) {
		case clickable.Click:
			if msg.ID == n.clickable.ID() {
				cmds = append(cmds, func() tea.Msg {
					return types.FocusMessage{
						ID: n.ID(),
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
	for i, c := range n.choices[n.start_index : n.start_index+n.height] {
		w := n.max_width - /* padding-left */ 1 - lipgloss.Width("┃")
		style := lipgloss.NewStyle().Width(w).MaxWidth(w).Inline(true)
		row := lipgloss.NewStyle().PaddingLeft(1).Foreground(
			map[bool]lipgloss.TerminalColor{
				false: color_profile.UIForeground[types.FocusStateActive],
				true:  color_profile.ForegroundInverted,
			}[i+n.start_index == n.index],
		).Background(
			map[bool]lipgloss.TerminalColor{
				true:  color_profile.UIForeground[types.FocusStateActive],
				false: color_profile.ForegroundInverted,
			}[i+n.start_index == n.index],
		).Render(
			style.Render(ansi.Truncate(c.Label, w, "…")),
		)
		rows = append(rows, fmt.Sprintf("┃%v", row))
	}

	return zone.Mark(
		n.clickable.ID(),
		lipgloss.NewStyle().Width(n.max_width).Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Inline(true).Width(n.max_width).MaxWidth(n.max_width).Render(
					fmt.Sprintf(
						"%v%v",
						lipgloss.NewStyle().Foreground(color_profile.UIForeground[n.FocusState()]).Render(n.prompt),
						ansi.Truncate(
							n.Value().Value.Label,
							n.max_width-lipgloss.Width(n.prompt),
							"…",
						),
					),
				),
				lipgloss.NewStyle().Foreground(color_profile.UIForeground[n.FocusState()]).Render(
					fmt.Sprintf(
						"%v%v%v%v",
						map[bool]string{
							true:  "──",
							false: "─ ",
						}[n.key.Label == ""],
						n.key.Label,
						map[bool]string{
							true:  "─",
							false: " ",
						}[n.key.Label == ""],
						strings.Repeat("─", n.max_width-len(n.key.Label)-3),
					),
				),
				lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).Render(
					lipgloss.JoinVertical(
						lipgloss.Right,
						lipgloss.JoinVertical(
							lipgloss.Left,
							rows...,
						),
						lipgloss.NewStyle().Foreground(color_profile.SupplementaryText).Render(
							fmt.Sprintf("%d / %d", n.index+1, len(n.choices)),
						),
					),
				),
			),
		),
	)
}
