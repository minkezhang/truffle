// Package table displays a series of Truffle nodes and affords an interface to
// expand to a dedicated view of each node.
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
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle/tui/component/button"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/db/message"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/image"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
	"github.com/minkezhang/truffle/tui/util/node"
	"github.com/minkezhang/truffle/tui/util/node/virtual"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	search_key "github.com/minkezhang/truffle/tui/component/search/key"
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

	// Put is passed along to the PutRequestMessage sent by this UI node.
	// This is checked by this UI node in the PutResponseMessage.
	Put = form.Key{Key: "table-add-link"}
)

type Node struct {
	*focusable.Node

	image          *image.Node
	column         *column.C
	table          table.Model
	selected_index int
	clickable      *clickable.Node
	key            form.Key
	data           []util_node.N

	select_button *button.Node
	add_button    *button.Node
	// link_button *button.Node  // TODO
}

type O struct {
	Prefix         string
	ParentID       string
	Column         *column.C
	Key            form.Key
	Data           []util_node.N
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
			table.WithHeight(11),
			table.WithKeyMap(
				table.KeyMap{
					LineUp:       table.DefaultKeyMap().LineUp,
					LineDown:     table.DefaultKeyMap().LineDown,
					PageUp:       table.DefaultKeyMap().PageUp,
					PageDown:     key.NewBinding(key.WithKeys("f", "pgdn")),
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
	n.select_button = button.New(button.O{
		ParentID: n.ID(),
		Key: form.Key{
			Label: "Select",
			Key:   "table-select",
		},
	})
	n.add_button = button.New(button.O{
		ParentID: n.ID(),
		Key: form.Key{
			Label: "Add",
			Key:   "table-add",
		},
	})
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.clickable.Init(),
		n.image.Init(),
		n.SetValues(n.data),
		n.select_button.Init(),
		n.add_button.Init(),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

type HighlightMessage message.GetNodeRequestMessage

func (n *Node) Value() form.Value[util_node.N] {
	if n.selected_index == -1 {
		return form.Value[util_node.N]{
			Key:   n.key,
			Value: node.N{},
		}
	}
	return form.Value[util_node.N]{
		Key:   n.key,
		Value: n.data[n.selected_index],
	}
}

func (n *Node) do_add_link(node_id string) tea.Cmd {
	// Only add link for single sources.
	if _, ok := n.Value().Value.(virtual.N); !ok {
		return nil
	}

	s, err := n.Value().Value.Virtual()
	if err != nil {
		return func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf("%v: cannot get source: %v", err),
			)
		}
	}

	m := message.PutRequestMessage{
		ID: n.ID(),
		Body: form.Value[source.S]{
			Key:   Put,
			Value: s.WithNodeID(node_id),
		},
	}
	return tea.Sequence(
		func() tea.Msg { return m },
		func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelDebug,
				fmt.Sprintf("%v: add to DB %v", n.ID(), m),
			)
		},
	)
}

func (n *Node) do_select() tea.Cmd {
	m := message.GetNodeRequestMessage{
		ID:   n.ID(),
		Body: n.Value(),
	}
	return tea.Sequence(
		func() tea.Msg { return m },
		func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelDebug,
				fmt.Sprintf("%v: submitted message %v", n.ID(), m),
			)
		},
	)

}

func (n *Node) do_highlight() tea.Cmd {
	if n.selected_index == n.table.Cursor() || len(n.table.Rows()) == 0 {
		return nil
	}

	return tea.Sequence(
		func() tea.Msg {
			n.selected_index = n.table.Cursor()
			return HighlightMessage{
				ID:   n.ID(),
				Body: n.Value(),
			}
		},
		func() tea.Msg {
			m := HighlightMessage{
				ID:   n.ID(),
				Body: n.Value(),
			}
			return errors.ToLogMessage(
				errors.LevelDebug,
				fmt.Sprintf("%v: highlighted message %v", n.ID(), m),
			)
		},
	)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	for _, m := range []tea.Model{
		n.image,
		n.clickable,
		n.select_button,
		n.add_button,
	} {
		_, c = m.Update(msg)
		cmds = append(cmds, c)
	}
	n.table, c = n.table.Update(msg)
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
					cmds = append(cmds, n.do_select())
				}
			}
		}
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
	case button.SubmitMessage:
		switch id := msg.ID; id {
		case n.select_button.ID():
			cmds = append(cmds, n.do_select())
		case n.add_button.ID():
			cmds = append(cmds, n.do_add_link(""))
		}
	case HighlightMessage:
		if msg.ID == n.ID() {
			source, err := msg.Body.Value.Virtual()
			if err != nil {
				cmds = append(cmds, func() tea.Msg {
					return errors.ToErrorMessage(err)
				})
			} else {
				cmds = append(cmds,
					n.image.SetValue(source.PreviewURL()),
					n.add_button.SetIsInvisible(map[epb.SourceAPI]bool{
						epb.SourceAPI_SOURCE_API_NONE:    true,
						epb.SourceAPI_SOURCE_API_TRUFFLE: true,
					}[source.Header().API()]),
				)
			}
		}
	case message.SearchResponseMessage:
		if msg.Body.Key == search_key.Search {
			cmds = append(cmds, n.SetValues(msg.Body.Value))
		}
	case message.PutResponseMessage:
		if msg.Body.Key == Put {
			cmds = append(cmds, n.PutSource(msg.Body.Value.Node, msg.Body.Value.SourceIndex))
		}
	}

	if n.selected_index != n.table.Cursor() && len(n.table.Rows()) > 0 {
		cmds = append(cmds, n.do_highlight())
	}

	return n, tea.Sequence(cmds...)
}

func (n *Node) View() string {
	var buttons []string
	for _, b := range []*button.Node{
		n.select_button,
		n.add_button,
	} {
		if !b.IsInvisible() {
			buttons = append(buttons, b.View())
		}
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		n.image.View(), // image
		lipgloss.JoinVertical( // table and button
			lipgloss.Left,
			lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(
				color_profile.UIForeground[n.FocusState()],
			).Render(
				zone.Mark(
					n.clickable.ID(),
					lipgloss.JoinVertical(
						lipgloss.Left,
						lipgloss.NewStyle().Foreground(
							color_profile.UIForeground[n.FocusState()],
						).Render( // label
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
						lipgloss.JoinVertical( // tabel
							lipgloss.Right,
							lipgloss.JoinHorizontal(
								lipgloss.Top,
								lipgloss.JoinVertical( // border-left
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
							).Render( // line count
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
			lipgloss.JoinHorizontal(lipgloss.Top, buttons...),
		),
	)
}

func (n *Node) to_row(data util_node.N) (table.Row, tea.Cmd) {
	source, err := data.Virtual()
	if err != nil {
		return nil, func() tea.Msg {
			return errors.ToErrorMessage(
				fmt.Errorf("Virtual() returned non-nil error: %v", err),
			)
		}
	}

	return []string{
		util_node.Table.Title(source.Title()),
		util_node.Table.Type(source.Header().Type()),
		util_node.Table.API(source.Header().API()),
		util_node.Table.Status(source.Status()),
		util_node.Table.Score(source.Score()),
	}, nil
}

func (n *Node) PutSource(m node.N, source_index int) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			old_node_index := -1 // Remove from old node if it exists in the table
			old_source_index := -1
			new_node_index := -1 // Add to an existing node if it exists in the table

			s := m.Sources()[source_index]

			for i, o := range n.data {
				if m.Header().ID() != "" && m.Header() == o.Header() {
					new_node_index = i
				}
				for j, t := range o.Sources() {
					if s.Header() == t.Header() {
						old_node_index = i
						old_source_index = j
					}
				}
			}

			if old_node_index >= 0 {
				// Replace old node if it is virtual or being deleted.
				if len(n.data[old_node_index].Sources()) == 1 {
					n.data[old_node_index] = m
					return nil
				}

				// If both old and new node exists, move source.
				sources := append([]source.S{}, n.data[old_node_index].Sources()...)
				n.data[old_node_index] = n.data[old_node_index].(node.N).WithSources(
					append(
						sources[:old_source_index],
						sources[old_source_index+1:]...,
					),
				)
				if new_node_index >= 0 {
					n.data[new_node_index] = m
				} else {
					// If no matching node was found, create a new
					// one in the view.
					n.data = append(
						n.data[:new_node_index],
						append(
							[]util_node.N{m},
							n.data[new_node_index:]...,
						)...,
					)
				}
			}
			return nil
		},
		n.set_values(n.data, false),
	)
}

func (n *Node) SetValues(data []util_node.N) tea.Cmd { return n.set_values(data, true) }

func (n *Node) set_values(data []util_node.N, reset_cursor bool) tea.Cmd {
	var cmds []tea.Cmd
	var rows []table.Row

	cmds = append(
		cmds,
		func() tea.Msg {
			for _, d := range data {
				r, c := n.to_row(d)
				rows = append(rows, r)
				cmds = append(cmds, c)
			}
			n.data = append([]util_node.N{}, data...)
			n.select_button.SetIsInvisible(n.IsInvisible())
			n.add_button.SetIsInvisible(n.IsInvisible())
			n.table.SetRows(rows)
			if reset_cursor {
				n.table.SetCursor(0)
				n.table.GotoTop()
				n.selected_index = -1
				cmds = append(cmds, n.image.SetValue(""))
			}
			return nil
		},
		n.do_highlight(),
	)
	return tea.Batch(cmds...)
}

func (n *Node) IsInvisible() bool { return n.Node.IsInvisible() || len(n.data) == 0 }

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
