// Package image loads a bubbletea model which returns a sixel-formatted
// representation of a URL.
//
// TOOD(minkezhang): Implement graphics via Kitty when available
// https://github.com/charmbracelet/bubbletea/issues/163.
//
// TODO(minkezhang): Support animated gifs.
package image

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
)

type O struct {
	ParentID       string
	URL            string
	CacheDirectory string
	Width          int
}

type Node struct {
	*focusable.Node

	width     int
	url       string
	filepath  string
	directory string
	cache     string // image
}

func New(o O) *Node {
	return &Node{
		Node:      focusable.New("image", o.ParentID, 0),
		url:       o.URL,
		width:     o.Width,
		directory: o.CacheDirectory,
	}
}

type update_cache_message struct {
	id      string
	payload string
}

func (n *Node) URL() string { return n.url }

func (n *Node) SetValue(v string) tea.Cmd {
	return func() tea.Msg {
		if n.url == v {
			return nil
		}
		n.url = v
		if n.url == "" {
			return update_cache_message{
				id:      n.ID(),
				payload: "",
			}
		}
		s, err := data(n.url, n.directory, n.width)
		if err != nil {
			return errors.ToLogMessage(
				errors.LevelWarn,
				err.Error(),
			)
		}
		return update_cache_message{
			id:      n.ID(),
			payload: s,
		}
	}
}

func (n *Node) Init() tea.Cmd {
	return tea.Batch(
		n.SetValue(n.url),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case update_cache_message:
		if n.ID() == msg.id {
			n.cache = string(msg.payload)
			cmds = append(cmds, func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: updating cache with url = \"%v\"", n.ID(), n.url),
				)
			})
		}
	}
	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return lipgloss.Place(
		n.width,
		// Two vertical pixels per character may leave a pixel
		// unaccounted for.
		height(n.width)/2+1,
		lipgloss.Top,
		lipgloss.Center,
		n.cache,
	)
}
