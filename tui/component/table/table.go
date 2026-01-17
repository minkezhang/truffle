package table

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/image"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
)

const (
	image_width = 26
)

var (
	styles = map[types.FocusState]table.Styles{
		types.FocusStateNone: table.Styles{
			Header: table.DefaultStyles().Header.Foreground(
				color_profile.UIForeground[types.FocusStateNone],
			).Border(
				lipgloss.NormalBorder(), false, false, true, false,
			).BorderForeground(
				color_profile.UIForeground[types.FocusStateNone],
			),
			Cell: lipgloss.NewStyle().Padding(0, 1),
			Selected: lipgloss.NewStyle().Foreground(
				color_profile.ForegroundNormal,
			).Background(
				color_profile.BackgroundNegligible,
			),
		},
		types.FocusStateActive: table.Styles{
			Header: table.DefaultStyles().Header.Foreground(
				color_profile.UIForeground[types.FocusStateActive],
			).Border(
				lipgloss.NormalBorder(), false, false, true, false,
			).BorderForeground(
				color_profile.UIForeground[types.FocusStateActive],
			),
			Cell: lipgloss.NewStyle().Padding(0, 1),
			Selected: lipgloss.NewStyle().Foreground(
				color_profile.ForegroundInverted,
			).Background(
				color_profile.ForegroundNormal,
			),
		},
	}
)

type Node struct {
	*focusable.Node

	image          *image.Node
	column         *column.C
	table          table.Model
	selected_index int
	clickable      *clickable.Node
	key            form.Key
	data           []node.N
}

type O struct {
	Prefix         string
	ParentID       string
	Column         *column.C
	Key            form.Key
	Data           []node.N
	CacheDirectory string
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New(o.Prefix, o.ParentID, 1),
		column: o.Column,
		key:    o.Key,
		table: table.New(
			table.WithColumns([]table.Column{
				table.Column{
					Title: "Title",
					Width: o.Column.Content() - /* other columns */ 45 - /* padding */ 10 - /* border-left */ 1 - /* image */ image_width,
				},
				table.Column{
					Title: "Media", // e.g. "Light Novel"
					Width: 15,
				},
				table.Column{
					Title: "API",
					Width: 15,
				},
				table.Column{
					Title: "Status",
					Width: 10,
				},
				table.Column{
					Title: "Score",
					Width: 5,
				},
			}),
			table.WithFocused(false),
			table.WithHeight(14),
			table.WithKeyMap(
				table.KeyMap{
					LineUp:       table.DefaultKeyMap().LineUp,
					LineDown:     table.DefaultKeyMap().LineDown,
					PageUp:       table.DefaultKeyMap().PageUp,
					PageDown:     key.NewBinding(key.WithKeys("f", "pgdn")), // table.DefaultKeyMap().PageDown,
					HalfPageUp:   table.DefaultKeyMap().HalfPageUp,
					HalfPageDown: table.DefaultKeyMap().HalfPageDown,
					GotoTop:      table.DefaultKeyMap().GotoTop,
					GotoBottom:   table.DefaultKeyMap().GotoBottom,
				},
			),
		),
		data: o.Data,
	}
	n.table.SetStyles(styles[n.FocusState()])
	n.clickable = clickable.New(n.ID())
	n.image = image.New(image.O{
		ParentID:       n.ID(),
		URL:            "",
		Width:          image_width,
		CacheDirectory: o.CacheDirectory,
	})
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.clickable.Init(),
		n.image.Init(),
		n.SetValues(n.data),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

type HighlightMessage struct {
	ID    string
	Value form.Value[node.N]
}

type SelectMessage HighlightMessage

func (n *Node) Value() form.Value[node.N] {
	if n.selected_index == -1 {
		return form.Value[node.N]{
			Key:   n.key,
			Value: node.N{},
		}
	}
	return form.Value[node.N]{
		Key:   n.key,
		Value: n.data[n.selected_index],
	}
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	_, c = n.image.Update(msg)
	cmds = append(cmds, c)
	_, c = n.clickable.Update(msg)
	cmds = append(cmds, c)

	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.MouseMsg:
			if msg.Button == tea.MouseButtonWheelUp {
				n.table.MoveUp(1)
			}
			if msg.Button == tea.MouseButtonWheelDown {
				n.table.MoveDown(1)
			}
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				if n.selected_index >= 0 {
					m := SelectMessage{
						ID:    n.ID(),
						Value: n.Value(),
					}
					cmds = append(cmds, tea.Sequence(
						func() tea.Msg { return m },
						func() tea.Msg {
							return errors.ToLogMessage(
								errors.LevelDebug,
								fmt.Sprintf("%v: submitted message %v", n.ID(), m),
							)
						},
					))
				}
			}
		}
		n.table, c = n.table.Update(msg)
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
	}

	switch msg := msg.(type) {
	case HighlightMessage:
		if msg.ID == n.ID() {
			source, err := msg.Value.Value.Virtual()
			if err != nil {
				cmds = append(cmds, func() tea.Msg {
					return errors.ToErrorMessage(err)
				})
			} else {
				cmds = append(cmds, n.image.SetURL(source.PreviewURL()))
			}
		}
	}

	if (n.selected_index == -1 && len(n.table.Rows()) > 0) || (n.selected_index >= 0 && n.selected_index != n.table.Cursor()) {
		n.selected_index = n.table.Cursor()
		m := HighlightMessage{
			ID:    n.ID(),
			Value: n.Value(),
		}
		cmds = append(cmds, tea.Sequence(
			func() tea.Msg { return m },
			func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: highlighted message %v", n.ID(), m),
				)
			},
		))
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return zone.Mark(
		n.clickable.ID(),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			n.image.View(),
			lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(
				color_profile.UIForeground[n.FocusState()],
			).Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
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
							strings.Repeat("─", n.column.Width()-len(n.key.Label)-3),
						),
					),
					lipgloss.JoinVertical(
						lipgloss.Right,
						lipgloss.JoinHorizontal(
							lipgloss.Top,
							lipgloss.JoinVertical(
								lipgloss.Left,
								lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(
									color_profile.UIForeground[n.FocusState()],
								).Render(" "),
								lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false, false, false, true).BorderForeground(
									color_profile.UIForeground[n.FocusState()],
								).Render(strings.Repeat("\n", n.table.Height()-1)),
							),
							n.table.View(),
						),
						lipgloss.NewStyle().Foreground(
							color_profile.BackgroundNegligible,
						).Render(
							fmt.Sprintf(
								"%d / %d",
								map[bool]int{
									true:  n.table.Cursor() + 1,
									false: 0,
								}[len(n.table.Rows()) > 0],
								len(n.table.Rows()),
							),
						),
					),
				),
			),
		),
	)
}

func (n *Node) to_row(data node.N) (table.Row, tea.Cmd) {
	source, err := data.Virtual()
	if err != nil {
		return nil, func() tea.Msg {
			return errors.ToErrorMessage(
				fmt.Errorf("Virtual() returned non-nil error: %v", err),
			)
		}
	}

	score := (source.Score() + 10) / 20
	return []string{
		source.Title().Title(),
		source.Header().Type().String(),
		source.Header().API().String(),
		source.Status().String(),
		strings.Repeat("★", score) + strings.Repeat("☆", 5-score),
	}, nil
}

func (n *Node) SetValues(data []node.N) tea.Cmd {
	var cmds []tea.Cmd
	var rows []table.Row
	for _, d := range data {
		r, c := n.to_row(d)
		rows = append(rows, r)
		cmds = append(cmds, c)
	}
	n.table.SetRows(rows)
	n.selected_index = -1
	return tea.Batch(cmds...)
}

func (n *Node) OnFocus(i int) tea.Cmd {
	n.table.SetStyles(styles[types.FocusStateActive])
	n.table.Focus()
	return n.Node.OnFocus(i)
}

func (n *Node) OnBlur() tea.Cmd {
	n.table.SetStyles(styles[types.FocusStateNone])
	n.table.Blur()
	return n.Node.OnBlur()
}
