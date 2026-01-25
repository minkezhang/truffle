package edit

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle/tui/component/button"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/textinput"
	"github.com/minkezhang/truffle/tui/util/form"
)

const (
	image_width = 50
)

type O struct {
	ParentID string
	Column   *column.C
}

type t struct {
	Title        *textinput.Node
	Localization *textinput.Node
}

type Node struct {
	*focusable.Node

	column *column.C
	source source.S

	titles        []t
	image         *textinput.Node
	button_unlink *button.Node
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("edit-node", o.ParentID, 0),
		column: o.Column,
	}
	n.button_unlink = button.New(button.O{
		n.ID(),
		form.Key{
			Label: "Unlink",
			Key:   "edit-unlink",
		},
	})
	n.image = textinput.New(textinput.O{
		Prefix:      "edit-source-image",
		ParentID:    n.ID(),
		Width:       image_width,
		Placeholder: "https://cdn.myanimelist.net/images/anime/1015/138006.jpg",
		Prompt:      "> ",
		Value: form.Value[string]{
			Key: form.Key{
				Label: "Image URL",
				Key:   "edit-source-preview-url",
			},
		},
	})
	return n
}

func (n *Node) SetValue(v source.S) tea.Cmd {
	var cmds []tea.Cmd

	n.source = v
	// Additional title for adding.
	for i := len(n.titles); i <= len(v.Titles()); i++ {
		_t := t{
			Title: textinput.New(textinput.O{
				Prefix:      "edit-source-title",
				ParentID:    n.ID(),
				Width:       40,
				Placeholder: "Sousou no Frieren",
				Prompt:      "> ",
				Value: form.Value[string]{
					Key: form.Key{
						Label: "Title",
						Key:   fmt.Sprintf("edit-source-title-%d", i),
					},
				},
			}),
			Localization: textinput.New(textinput.O{
				Prefix:      "edit-source-title",
				ParentID:    n.ID(),
				Width:       9,
				Placeholder: "en",
				Value: form.Value[string]{
					Key: form.Key{
						Label: "Locale",
						Key:   fmt.Sprintf("edit-source-localization-%d", i),
					},
				},
			}),
		}
		n.titles = append(n.titles, _t)
		cmds = append(cmds, tea.Sequence(_t.Title.Init(), _t.Localization.Init()))
	}
	for i, _t := range n.titles {
		is_empty := i >= len(v.Titles())
		is_invisible := i > len(v.Titles())
		cmds = append(
			cmds,
			_t.Title.SetIsInvisible(is_invisible),
			_t.Localization.SetIsInvisible(is_invisible),
		)
		if is_empty {
			cmds = append(
				cmds,
				_t.Title.SetValue(""),
				_t.Localization.SetValue(""),
			)
		} else {
			cmds = append(
				cmds,
				_t.Title.SetValue(v.Titles()[i].Title()),
				_t.Localization.SetValue(v.Titles()[i].Localization()),
			)
		}
	}
	cmds = append(
		cmds,
		n.image.SetValue(v.PreviewURL()),
	)
	return tea.Batch(cmds...)
}

func (n *Node) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, _t := range n.titles {
		cmds = append(
			cmds,
			_t.Title.Init(),
			_t.Localization.Init(),
		)
	}

	cmds = append(
		cmds,
		n.button_unlink.Init(),
		n.image.Init(),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
	return tea.Sequence(cmds...)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	for _, _t := range n.titles {
		_, c = _t.Title.Update(msg)
		cmds = append(cmds, c)
		_, c = _t.Localization.Update(msg)
		cmds = append(cmds, c)
	}

	_, c = n.image.Update(msg)
	cmds = append(cmds, c)

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	var parts []string
	for _, _t := range n.titles {
		if !_t.Title.IsInvisible() {
			parts = append(
				parts,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					lipgloss.NewStyle().Margin(0, 1, 0, 0).Render(
						_t.Title.View(),
					),
					_t.Localization.View(),
				),
			)
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
