package dropdown

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/choices"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
)

type O struct {
	Prefix   string
	ParentID string
	Column   *column.C
	Width    int
	Key      form.Key
	Choices  []form.Key
	Prompt   string
	Height   int
}

type Node struct {
	*focusable.Node

	column    *column.C
	prompt    string
	clickable *clickable.Node
	choices   *choices.Node
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New(o.Prefix, o.ParentID, 1),
		column: o.Column,
		prompt: o.Prompt,
	}
	n.clickable = clickable.New(n.ID())
	n.choices = choices.New(choices.O{
		ParentID: n.ID(),
		Width:    o.Column.Content(),
		Height:   o.Height,
		Key:      o.Key,
		Choices:  o.Choices,
	})
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.clickable.Init(),
		n.choices.Init(),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) SetValue(vs []form.Key, index int) tea.Cmd { return n.choices.SetValue(vs, index) }
func (n *Node) Value() form.Value[form.Key]               { return n.choices.Value() }

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	if n.FocusState() != types.FocusStateActive {
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
	_, c = n.choices.Update(msg)
	cmds = append(cmds, c)

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return n.column.RenderOrDie(
		lipgloss.JoinVertical(
			lipgloss.Left,
			zone.Mark(
				n.clickable.ID(),
				n.column.Style().Render(
					lipgloss.JoinVertical(
						lipgloss.Left,
						n.column.Style().MaxWidth(n.column.Width()).Inline(true).Render(
							fmt.Sprintf(
								"%v%v",
								lipgloss.NewStyle().Foreground(color_profile.UIForeground[n.FocusState()]).Render(n.prompt),
								ansi.Truncate(
									n.Value().Value.Label,
									n.column.Content()-lipgloss.Width(n.prompt),
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
								}[n.Value().Key.Label == ""],
								n.Value().Key.Label,
								map[bool]string{
									true:  "─",
									false: " ",
								}[n.Value().Key.Label == ""],
								strings.Repeat("─", n.column.Content()-len(n.Value().Key.Label)-3),
							),
						),
					),
				),
			),
			n.choices.View(),
		),
	)
}
