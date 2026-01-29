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
	"google.golang.org/protobuf/encoding/prototext"

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

	column      *column.C
	node        util_node.N
	source      *view.Node
	edit        *edit.Node
	render_type tablist.TabType
	tablist     *tablist.Node
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

	var values []tablist.Tab

	if _, ok := v.(node.N); ok {
		values = append(values, tablist.Tab{
			Type:  tablist.TabTypeVirtual,
			Label: "⌂",
		})
	}
	for i, s := range v.Sources() {
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

	var cmds []tea.Cmd

	var parts []string

	buf, _ := prototext.Marshal(v.PB())

	parts = append(parts, string(buf))

	for _, s := range v.Sources() {
		buf, _ = prototext.Marshal(s.PB())
		parts = append(parts, string(buf))
	}

	cmds = append(cmds,
		tea.Sequence(
			n.tablist.SetValue(values),
			n.do_highlight(n.tablist.Value()),
		),
	)

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

func (n *Node) do_highlight(v form.Value[tablist.Tab]) tea.Cmd {
	var cmds []tea.Cmd

	switch v.Value.Type {
	case tablist.TabTypeVirtual:
		s, err := n.node.Virtual()
		if err != nil {
			cmds = append(cmds, func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelWarn,
					fmt.Sprintf("cannot get virtual node: %v", err),
				)
			})
		} else {
			cmds = append(cmds, n.source.SetValue(s))
		}
	case tablist.TabTypeSource:
		cmds = append(cmds, n.source.SetValue(n.node.Sources()[v.Value.Key]))
	case tablist.TabTypeEdit:
		// Cannot edit remote sources
		//
		// TODO(minkezhang): Add handler for linking virtual source.
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
							fmt.Sprintf("multiple Truffle sources found for node %v", n.node.Header().ID()),
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
			}).WithNodeID(n.node.Header().ID())
		}

		cmds = append(cmds, n.edit.SetValue(s))
	}
	n.render_type = v.Value.Type
	if n.node != nil && n.node.Header().ID() == "" && n.render_type == tablist.TabTypeEdit {
		n.render_type = tablist.TabTypeNone
	}
	cmds = append(
		cmds,
		n.edit.SetIsInvisible(n.render_type != tablist.TabTypeEdit),
		n.source.SetIsInvisible(n.render_type != tablist.TabTypeVirtual && n.render_type != tablist.TabTypeSource),
	)
	return tea.Sequence(cmds...)
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
		cmds = append(cmds, n.SetValue(msg.Body.Value))
	case message.GetResponseMessage:
		cmds = append(cmds, n.SetValue(msg.Body.Value.Node))
	case message.PutResponseMessage:
		cmds = append(cmds, n.SetValue(msg.Body.Value.Node))
	case tablist.HighlightMessage:
		if msg.ID == n.tablist.ID() {
			cmds = append(cmds, n.do_highlight(msg.Value))
		}
	}

	return n, tea.Sequence(cmds...)
}

func (n *Node) View() string {
	if n.node == nil {
		return ""
	}
	return n.column.RenderOrDie(n.column.Style().Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(n.tablist.View()),
			map[tablist.TabType]string{
				tablist.TabTypeSource:  n.source.View(),
				tablist.TabTypeVirtual: n.source.View(),
				tablist.TabTypeEdit:    n.edit.View(),
			}[n.render_type],
		),
	))
}
