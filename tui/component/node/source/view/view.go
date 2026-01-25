package view

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/image"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/node/view"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

const (
	image_width = 50
)

type O struct {
	ParentID       string
	Column         *column.C
	CacheDirectory string
}

type Node struct {
	*focusable.Node

	column *column.C
	image  *image.Node
	source source.S
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("source-view", o.ParentID, 0),
		column: o.Column,
	}
	n.image = image.New(image.O{
		ParentID:       n.ID(),
		URL:            "",
		Width:          image_width,
		CacheDirectory: o.CacheDirectory,
	})
	return n
}

func (n *Node) SetValue(v source.S) tea.Cmd {
	n.source = v
	return n.image.SetValue(v.PreviewURL())
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.image.Init(),
		n.SetValue(n.source),
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

	_, c = n.image.Update(msg)
	cmds = append(cmds, c)

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	parts := []string{
		fmt.Sprintf(
			"%v%v%v",
			view.WithHeader("Type", view.R(n.source).Type(), view.RenderHeaderModeNone),
			lipgloss.NewStyle().Foreground(color_profile.ForegroundNegligible).Render(" > "),
			view.WithHeader("ID", view.R(n.source).ID(), view.RenderHeaderModeNone),
		),
		view.R(n.source).Title(),
		view.WithHeader("Score", view.R(n.source).Score(), view.RenderHeaderModeNone),
		view.WithHeader("Status", view.R(n.source).Status(), view.RenderHeaderModeInline),
		view.WithHeader("Genres", view.R(n.source).Genres(), view.RenderHeaderModeIndent),
	}

	if map[epb.SourceType]bool{
		epb.SourceType_SOURCE_TYPE_BOOK:             true,
		epb.SourceType_SOURCE_TYPE_BOOK_MANGA:       true,
		epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL: true,
	}[n.source.Header().Type()] {
		parts = append(
			parts,
			view.WithHeader("Authors", view.R(n.source).Authors(), view.RenderHeaderModeIndent),
			view.WithHeader("Illustrators", view.R(n.source).Illustrators(), view.RenderHeaderModeIndent),
		)
	}

	if map[epb.SourceType]bool{
		epb.SourceType_SOURCE_TYPE_SERIES_ANIME: true,
		epb.SourceType_SOURCE_TYPE_MOVIE_ANIME:  true,
	}[n.source.Header().Type()] {
		parts = append(
			parts,
			view.WithHeader("Studio", view.R(n.source).Studios(), view.RenderHeaderModeIndent),
		)
	}

	if map[epb.SourceType]bool{
		epb.SourceType_SOURCE_TYPE_SERIES_ANIME: true,
	}[n.source.Header().Type()] {
		parts = append(
			parts,
			view.WithHeader("Seasons", view.R(n.source).Seasons(), view.RenderHeaderModeIndent),
		)
	}

	parts = append(
		parts,
		view.WithHeader("Synopsis", view.R(n.source).Synopsis(), view.RenderHeaderModeNone),
		view.WithHeader("Notes", view.R(n.source).Notes(), view.RenderHeaderModeSpacer),
	)

	for i, p := range parts {
		if i != len(parts)-1 {
			parts[i] = lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(p)
		}
	}

	return n.column.RenderOrDie(
		n.column.Style().Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				n.image.View(),
				n.column.WithWidth(n.column.Width()-image_width).WithMargin(0, 0, 0, 1).Style().Render(
					lipgloss.JoinVertical(
						lipgloss.Left,
						parts...,
					),
				),
			),
		),
	)
}
