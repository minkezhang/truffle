package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/db"
	"github.com/minkezhang/truffle/tui/component/db/message"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/form"
	"github.com/minkezhang/truffle/tui/util/node"
	"github.com/minkezhang/truffle/tui/util/node/virtual"
	"github.com/minkezhang/truffle/tui/util/search"
	"google.golang.org/protobuf/encoding/prototext"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
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
	case message.SearchRequestMessage:
		cmds = append(cmds, do_search(n.context, n.db, msg))
	case message.PutRequestMessage:
		cmds = append(
			cmds,
			tea.Sequence(
				func() tea.Msg {
					buf, _ := prototext.Marshal(msg.Body.Value.PB())
					return errors.ToLogMessage(
						errors.LevelInfo,
						fmt.Sprintf("%v: got PutRequestMessage\n%v", n.ID(), string(buf)),
					)
				},
				do_put(n.context, n.db, msg),
			),
		)
	case message.PutResponseMessage:
		cmds = append(cmds, func() tea.Msg {
			buf, _ := prototext.Marshal(msg.Body.Value.Node.PB())
			parts := []string{string(buf)}
			for _, s := range msg.Body.Value.Node.Sources() {
				buf, _ = prototext.Marshal(s.PB())
				parts = append(parts, string(buf))
			}
			return errors.ToLogMessage(
				errors.LevelInfo,
				fmt.Sprintf("%v: received PutResponseMessage:\n%v", n.ID(), strings.Join(parts, "\n")),
			)
		})
	case message.GetNodeRequestMessage:
		cmds = append(cmds, do_get_node(n.context, n.db, msg))
	case message.GetRequestMessage:
		cmds = append(cmds, do_get(n.context, n.db, msg))
	}

	return n, tea.Batch(cmds...)
}

func do_get(ctx context.Context, _db *db.DB, msg message.GetRequestMessage) tea.Cmd {
	return func() tea.Msg {
		s, err := _db.Get(ctx, msg.Body.Value, option.Remote(false))
		if err != nil {
			buf, _ := prototext.Marshal(msg.Body.Value.PB())
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf("db.Get(%v) returned non-nil error: %v", string(buf), err.Error()),
			)
		}
		var n util_node.N
		if s.NodeID() != "" {
			n, err = _db.GetNode(ctx, node.Make(&dpb.Node{
				Header: &dpb.NodeHeader{
					Id:   s.NodeID(),
					Type: s.Header().Type(),
				},
			}).Header(), option.Remote(false))
			if err != nil {
				buf, _ := prototext.Marshal(msg.Body.Value.PB())
				return errors.ToLogMessage(
					errors.LevelWarn,
					fmt.Sprintf("db.Get(%v) returned non-nil error: %v", string(buf), err.Error()),
				)
			}
		} else {
			n = virtual.Make(s.PB())
		}

		source_index := -1
		for i, t := range n.Sources() {
			if t.Header() == msg.Body.Value {
				source_index = i
			}
		}
		return message.GetResponseMessage{
			ID: msg.ID,
			Body: form.Value[message.GetResponseBody]{
				Key: msg.Body.Key,
				Value: message.GetResponseBody{
					Node:        n,
					SourceIndex: source_index,
				},
			},
		}
	}
}

func do_get_node(ctx context.Context, _db *db.DB, msg message.GetNodeRequestMessage) tea.Cmd {
	return func() tea.Msg {
		n, err := _db.GetNode(ctx, msg.Body.Value.Header(), option.Remote(false))
		if err != nil {
			buf, _ := prototext.Marshal(msg.Body.Value.Header().PB())
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf("db.GetNode(%v) returned non-nil error: %v", string(buf), err.Error()),
			)
		}
		return message.GetNodeResponseMessage{
			ID: msg.ID,
			Body: form.Value[util_node.N]{
				Key:   msg.Body.Key,
				Value: n,
			},
		}
	}
}

func do_put(ctx context.Context, _db *db.DB, msg message.PutRequestMessage) tea.Cmd {
	return func() tea.Msg {
		h, err := _db.Put(ctx, msg.Body.Value)
		if err != nil {
			buf, _ := prototext.Marshal(msg.Body.Value.Header().PB())
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf("db.Put(%v) returned non-nil error: %v", string(buf), err.Error()),
			)
		}

		s, err := _db.Get(ctx, h, option.Remote(false))
		if err != nil {
			buf, _ := prototext.Marshal(h.PB())
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf("db.Get(%v) returned non-nil error: %v", string(buf), err.Error()),
			)
		}

		n, err := _db.GetNode(
			ctx,
			node.Make(&dpb.Node{
				Header: &dpb.NodeHeader{
					Id:   s.NodeID(),
					Type: s.Header().Type(),
				},
			}).Header(),
			option.Remote(false),
		)
		if err != nil {
			buf, _ := prototext.Marshal(&dpb.NodeHeader{
				Id:   s.NodeID(),
				Type: s.Header().Type(),
			})
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf("db.GetNode(%v) returned non-nil error: %v", string(buf), err.Error()),
			)
		}

		source_index := -1
		for i, s := range n.Sources() {
			if s.Header() == h {
				source_index = i
			}
		}

		return message.PutResponseMessage{
			ID: msg.ID,
			Body: form.Value[message.PutResponseBody]{
				Key: msg.Body.Key,
				Value: message.PutResponseBody{
					Node:        n,
					SourceIndex: source_index,
				},
			},
		}
	}
}

func do_search(ctx context.Context, _db *db.DB, msg message.SearchRequestMessage) tea.Cmd {
	return func() tea.Msg {
		var types []epb.SourceType
		for t, ok := range msg.Body.Value.SourceTypes {
			if ok {
				types = append(types, t)
			}
		}

		_opts := []option.O{option.Remote(true)}
		for _, o := range msg.Body.Value.Options {
			_opts = append(_opts, o)
		}

		opts := map[epb.SourceAPI][]option.O{}
		for api, ok := range msg.Body.Value.APIs {
			if ok {
				opts[api] = append([]option.O{}, _opts...)
			}
		}

		results, err := util_search.Search(ctx, _db, msg.Body.Value.Query, opts, types)
		if err != nil {
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf(
					"db.Search(\"%v\") returned non-nil error: %v",
					msg.Body.Value.Query,
					err.Error(),
				),
			)
		}

		return message.SearchResponseMessage{
			ID: msg.ID,
			Body: form.Value[[]util_node.N]{
				Key:   msg.Body.Key,
				Value: results,
			},
		}
	}
}

func (n *Node) View() string { return "" }
