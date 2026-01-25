package component_node

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/db/message"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/node/source/edit"
	"github.com/minkezhang/truffle/tui/component/node/source/view"
	"github.com/minkezhang/truffle/tui/component/tablist"
	"github.com/minkezhang/truffle/tui/util/form"
	"github.com/minkezhang/truffle/tui/util/node"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	util_view "github.com/minkezhang/truffle/tui/util/node/view"
)

type O struct {
	ParentID       string
	Column         *column.C
	Node           util_node.N
	CacheDirectory string
}

type Node struct {
	*focusable.Node

	column  *column.C
	node    util_node.N
	source  *view.Node
	edit    *edit.Node
	is_edit bool
	tablist *tablist.Node
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("node", o.ParentID, 0),
		column: o.Column,
		node:   o.Node,
	}
	n.source = view.New(view.O{
		ParentID:       n.ID(),
		Column:         o.Column,
		CacheDirectory: o.CacheDirectory,
	})
	n.tablist = tablist.New(tablist.O{
		ParentID: n.ID(),
		Column:   o.Column,
		Key:      form.Key{"", "tab-select"},
	})
	n.edit = edit.New(edit.O{
		ParentID: n.ID(),
		Column:   o.Column,
	})
	return n
}

func (n *Node) SetValue(v util_node.N) tea.Cmd {
	if v == nil {
		return nil
	}
	n.node = v
	s, err := v.Virtual()
	if err != nil {
		return func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf("%v: cannot get a merged source component: %v", err),
			)
		}
	}

	var values []tablist.Tab
	truffle_index := -1
	if _, ok := v.(node.N); ok {
		values = append(values, tablist.Tab{
			Type:  tablist.TabTypeVirtual,
			Label: "⌂",
		})
	}
	for i, s := range v.Sources() {
		if s.Header().API() == epb.SourceAPI_SOURCE_API_TRUFFLE {
			truffle_index = i
		}
		values = append(values, tablist.Tab{
			Type:  tablist.TabTypeSource,
			Key:   i,
			Label: util_view.R(s).API(),
		})
	}
	values = append(values, tablist.Tab{
		Type:  tablist.TabTypeEdit,
		Label: "+",
	})

	cmds := []tea.Cmd{
		tea.Sequence(
			n.tablist.SetValue(values),

			func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: setting tablist value", n.ID()),
				)
			},
		),
		n.source.SetValue(s),
	}

	if truffle_index == -1 {
		s = source.Make(&dpb.Source{
			Header: &dpb.SourceHeader{
				Type: n.node.Header().Type(),
				Api:  epb.SourceAPI_SOURCE_API_TRUFFLE,
			},
		}).WithNodeID(n.node.Header().ID())
	}

	cmds = append(cmds, n.edit.SetValue(s))
	return tea.Batch(cmds...)
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.source.Init(),
		n.tablist.Init(),
		n.edit.Init(),
		n.SetValue(n.node),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	_, c = n.source.Update(msg)
	cmds = append(cmds, c)

	_, c = n.tablist.Update(msg)
	cmds = append(cmds, c)

	_, c = n.edit.Update(msg)
	cmds = append(cmds, c)

	switch msg := msg.(type) {
	case message.GetNodeResponseMessage:
		cmds = append(
			cmds,
			n.SetValue(msg.Body.Value),
			func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: received GetNodeResponseMessage: %v", n.ID(), msg),
				)
			},
		)
	case message.PutResponseMessage:
		cmds = append(
			cmds,
			n.SetValue(msg.Body.Value.Node),
			func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: received PutResponseMessage: %v", n.ID(), msg),
				)
			},
		)
	case tablist.HighlightMessage:
		if msg.ID == n.tablist.ID() {
			n.is_edit = msg.Value.Value.Type == tablist.TabTypeEdit
			cmds = append(
				cmds,
				n.edit.SetIsInvisible(!n.is_edit),
				n.source.SetIsInvisible(n.is_edit),
			)
			cmds = append(cmds, func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: recieved highlight message: %v", n.ID(), msg),
				)
			})
			switch msg.Value.Value.Type {
			case tablist.TabTypeVirtual:
				source, err := n.node.Virtual()
				if err != nil {
					cmds = append(cmds, func() tea.Msg {
						return errors.ToLogMessage(
							errors.LevelWarn,
							fmt.Sprintf("%v: Virtual() returned error: %v", err),
						)
					})
				} else {
					cmds = append(cmds, n.source.SetValue(source))
				}
			case tablist.TabTypeSource:
				cmds = append(cmds, n.source.SetValue(n.node.Sources()[msg.Value.Value.Key]))
			case tablist.TabTypeEdit:
				// Cannot edit remote sources
				if _, ok := n.node.(node.N); !ok {
					break
				}
				var s source.S
				var is_exists bool
				for _, _s := range n.node.Sources() {
					if _s.Header().API() == epb.SourceAPI_SOURCE_API_TRUFFLE {
						if is_exists {
							cmds = append(cmds, func() tea.Msg {
								return errors.ToLogMessage(
									errors.LevelWarn,
									fmt.Sprintf("%v: multiple Truffle sources found: %v", n.ID(), _s.Header()),
								)
							})
						} else {
							s = _s
							is_exists = true
						}
					}
				}
				if !is_exists {
					s = source.Make(&dpb.Source{
						Header: &dpb.SourceHeader{
							Type: n.node.Header().Type(),
							Api:  epb.SourceAPI_SOURCE_API_TRUFFLE,
						},
					}).WithNodeID(n.ID())
				}
				cmds = append(cmds,
					n.edit.SetValue(s),
					func() tea.Msg {
						return errors.ToLogMessage(
							errors.LevelWarn,
							fmt.Sprintf("%v: unimplemented edit source", n.ID()),
						)
					},
				)
			}
		}
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	if n.node == nil {
		return ""
	}
	return n.column.RenderOrDie(n.column.Style().Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(n.tablist.View()),
			map[bool]string{
				false: n.source.View(),
				true:  n.edit.View(),
			}[n.is_edit],
		),
	))
}
