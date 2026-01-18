package db

import (
	"context"

	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle-api/db"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/search"
	"github.com/minkezhang/truffle/tui/component/table"
	"github.com/minkezhang/truffle/tui/util/node"
	"github.com/minkezhang/truffle/tui/util/search"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type Node struct {
	*focusable.Node
	db      *db.DB
	context context.Context
}

type O struct {
	ParentID string
	DB       *db.DB
}

func New(ctx context.Context, o O) *Node {
	return &Node{
		Node:    focusable.New("db-connection", o.ParentID, 0),
		db:      o.DB,
		context: ctx,
	}
}

func (n *Node) Init() tea.Cmd {
	return func() tea.Msg {
		return base.RegisterMessage{
			Node: n,
		}
	}
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case search.SubmitSearchMessage:
		cmds = append(cmds, do_search(n.context, n.db, msg))
	case table.AddLinkMessage:
		cmds = append(cmds, do_add_link(n.context, n.db, msg))
	}

	return n, tea.Batch(cmds...)
}

type SearchResultMessage struct {
	Results []util_node.N
}

type AddLinkResultMessage struct {
	Result source.S
}

func do_add_link(ctx context.Context, _db *db.DB, msg table.AddLinkMessage) tea.Cmd {
	return func() tea.Msg {
		h, err := _db.Put(ctx, msg.Value.Value.WithNodeID(msg.NodeID))
		if err != nil {
			return errors.ToLogMessage(
				errors.LevelWarn,
				err.Error(),
			)
		}

		s, err := _db.Get(ctx, h, option.Remote(false))
		if err != nil {
			return errors.ToLogMessage(
				errors.LevelWarn,
				err.Error(),
			)
		}

		return AddLinkResultMessage{
			Result: s,
		}
	}
}

func do_search(ctx context.Context, _db *db.DB, msg search.SubmitSearchMessage) tea.Cmd {
	return func() tea.Msg {
		var types []epb.SourceType
		for _, t := range msg.SourceTypes.Value {
			if t.Value {
				if v, ok := epb.SourceType_value[t.Key.Key]; ok {
					types = append(types, epb.SourceType(v))
				}
			}
		}

		_opts := []option.O{option.Remote(true)}
		for _, o := range msg.Options.Value {
			if o.Value {
				_opts = append(_opts, map[string]option.O{
					"options-nsfw": option.NSFW(true),
				}[o.Key.Key])
			}
		}

		opts := map[epb.SourceAPI][]option.O{}
		for _, api := range msg.APIs.Value {
			if api.Value {
				if v, ok := epb.SourceAPI_value[api.Key.Key]; ok {
					opts[epb.SourceAPI(v)] = append([]option.O{}, _opts...)
				}
			}
		}

		results, err := util_search.Search(ctx, _db, msg.Query.Value, opts, types)
		if err != nil {
			return errors.ToLogMessage(
				errors.LevelWarn,
				err.Error(),
			)
		}

		return SearchResultMessage{
			Results: results,
		}
	}
}

func (n *Node) View() string { return "" }
