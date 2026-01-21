// Package image loads a bubbletea model which returns a sixel-formatted
// representation of a URL.
//
// TOOD(minkezhang): Implement graphics via Kitty when available
// https://github.com/charmbracelet/bubbletea/issues/163.
//
// TODO(minkezhang): Support animated gifs.
package image

import (
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

func (n *Node) SetValue(v string) tea.Cmd {
	if n.url == v {
		return nil
	}
	n.url = v
	if n.url == "" {
		n.cache = ""
		return nil
	}
	s, err := data(n.url, n.directory, n.width)
	if err != nil {
		return func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelWarn,
				err.Error(),
			)
		}
	}
	n.cache = s
	return nil
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

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return n, nil }

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
