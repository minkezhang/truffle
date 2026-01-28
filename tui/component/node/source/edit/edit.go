package edit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle/tui/component/button"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/db/message"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/textinput"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
	"github.com/minkezhang/truffle/tui/util/node/view"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
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

	titles []t
	image  *textinput.Node
	score  *textinput.Node
	// synopsis
	// notes
	genres *textinput.Node
	// status
	studios      *textinput.Node
	seasons      *textinput.Node
	authors      *textinput.Node
	illustrators *textinput.Node
	// last updated
	save_button *button.Node
}

var (
	key = form.Key{"", "edit-source"}
)

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("edit-node", o.ParentID, 0),
		column: o.Column,
	}
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
	n.score = textinput.New(textinput.O{
		Prefix:      "edit-source-score",
		ParentID:    n.ID(),
		Width:       9,
		Placeholder: "98",
		Prompt:      "> ",
		Value: form.Value[string]{
			Key: form.Key{
				Label: "Score",
				Key:   "edit-source-score",
			},
		},
	})
	n.genres = textinput.New(textinput.O{
		Prefix:      "edit-source-genres",
		ParentID:    n.ID(),
		Width:       50,
		Placeholder: "Adventure, Drama, Fantasy, Shounen",
		Prompt:      "> ",
		Value: form.Value[string]{
			Key: form.Key{
				Label: "Genres",
				Key:   "edit-source-genres",
			},
		},
	})
	n.studios = textinput.New(textinput.O{
		Prefix:      "edit-source-studios",
		ParentID:    n.ID(),
		Width:       50,
		Placeholder: "Madhouse",
		Prompt:      "> ",
		Value: form.Value[string]{
			Key: form.Key{
				Label: "Studio",
				Key:   "edit-source-studios",
			},
		},
	})
	n.seasons = textinput.New(textinput.O{
		Prefix:      "edit-source-seasons",
		ParentID:    n.ID(),
		Width:       50,
		Placeholder: "Fall 2023",
		Prompt:      "> ",
		Value: form.Value[string]{
			Key: form.Key{
				Label: "Seasons",
				Key:   "edit-source-seasons",
			},
		},
	})
	n.authors = textinput.New(textinput.O{
		Prefix:      "edit-source-authors",
		ParentID:    n.ID(),
		Width:       50,
		Placeholder: "Kanehito Yamada",
		Prompt:      "> ",
		Value: form.Value[string]{
			Key: form.Key{
				Label: "Authors",
				Key:   "edit-source-authors",
			},
		},
	})
	n.illustrators = textinput.New(textinput.O{
		Prefix:      "edit-source-illustrators",
		ParentID:    n.ID(),
		Width:       50,
		Placeholder: "Tsukasa Abe",
		Prompt:      "> ",
		Value: form.Value[string]{
			Key: form.Key{
				Label: "Illustrators",
				Key:   "edit-source-illustrators",
			},
		},
	})
	n.save_button = button.New(button.O{
		n.ID(),
		form.Key{
			Label: "Save",
			Key:   "edit-save",
		},
	})
	return n
}

func (n *Node) SetIsInvisible(v bool) tea.Cmd {
	type i interface {
		SetIsInvisible(v bool) tea.Cmd
		IsInvisible() bool
	}

	var cmds []tea.Cmd
	children := []i{n.Node}
	for _, _t := range n.titles {
		children = append(children,
			_t.Title,
			_t.Localization,
		)
	}
	for _, c := range children {
		cmds = append(cmds, c.SetIsInvisible(v || c.IsInvisible()))
	}

	for _, c := range []i{
		n.image,
		n.score,
		n.genres,
		n.save_button,
	} {
		cmds = append(cmds, c.SetIsInvisible(v))
	}

	_t := n.source.Header().Type()

	cmds = append(
		cmds,
		n.studios.SetIsInvisible(v || !map[epb.SourceType]bool{
			epb.SourceType_SOURCE_TYPE_SERIES_ANIME: true,
			epb.SourceType_SOURCE_TYPE_MOVIE_ANIME:  true,
		}[_t]),
		n.seasons.SetIsInvisible(v || !map[epb.SourceType]bool{
			epb.SourceType_SOURCE_TYPE_SERIES_ANIME: true,
		}[_t]),
		n.authors.SetIsInvisible(v || !map[epb.SourceType]bool{
			epb.SourceType_SOURCE_TYPE_BOOK:             true,
			epb.SourceType_SOURCE_TYPE_BOOK_MANGA:       true,
			epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL: true,
		}[_t]),
		n.illustrators.SetIsInvisible(v || !map[epb.SourceType]bool{
			epb.SourceType_SOURCE_TYPE_BOOK:             true,
			epb.SourceType_SOURCE_TYPE_BOOK_MANGA:       true,
			epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL: true,
		}[_t]),
	)

	return tea.Batch(cmds...)
}

func (n *Node) SetValue(v source.S) tea.Cmd {
	var cmds []tea.Cmd

	n.source = v
	// Additional title for adding.
	var _cmds []tea.Cmd
	var _titles []t
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
				Prefix:      "edit-source-localization",
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
		_titles = append(_titles, _t)
		_cmds = append(
			_cmds,
			_t.Title.Init(),
			_t.Localization.Init(),
		)
	}
	n.titles = append(n.titles, _titles...)

	// Reorder children tab order.
	children := []string{
		n.image.ID(),
	}
	for _, _t := range n.titles {
		children = append(children, _t.Title.ID(), _t.Localization.ID())
	}
	children = append(
		children,
		n.score.ID(),
		n.genres.ID(),
		n.studios.ID(),
		n.seasons.ID(),
		n.authors.ID(),
		n.illustrators.ID(),
		n.save_button.ID(),
	)

	_cmds = append(_cmds, func() tea.Msg {
		return base.PutChildrenMessage{
			ID:       n.ID(),
			Children: children,
		}
	})

	cmds = append(cmds, tea.Sequence(_cmds...))

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
		n.score.SetValue(fmt.Sprintf("%d", v.Score())),
		n.genres.SetValue(strings.Join(v.Genres(), ", ")),
		n.studios.SetValue(strings.Join(v.Studios(), ", ")),
		n.seasons.SetValue(strings.Join(v.Seasons(), ", ")),
		n.authors.SetValue(strings.Join(v.Authors(), ", ")),
		n.illustrators.SetValue(strings.Join(v.Illustrators(), ", ")),
	)

	return tea.Batch(cmds...)
}

func (n *Node) Init() tea.Cmd {
	cmds := []tea.Cmd{
		n.image.Init(),
	}

	for _, _t := range n.titles {
		cmds = append(
			cmds,
			_t.Title.Init(),
			_t.Localization.Init(),
		)
	}

	for _, m := range []tea.Model{
		n.score,
		n.genres,
		n.studios,
		n.seasons,
		n.authors,
		n.illustrators,
		n.save_button,
	} {
		cmds = append(cmds, m.Init())
	}
	cmds = append(
		cmds,
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
	return tea.Sequence(cmds...)
}

func (n *Node) do_save() tea.Cmd {
	return func() tea.Msg {
		return message.PutRequestMessage{
			ID:   n.ID(),
			Body: n.Value(),
		}
	}
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

	for _, m := range []*textinput.Node{
		n.image,
		n.score,
		n.genres,
		n.studios,
		n.seasons,
		n.authors,
		n.illustrators,
	} {
		_, c = m.Update(msg)
		cmds = append(cmds, c)
	}

	_, c = n.save_button.Update(msg)
	cmds = append(cmds, c)

	switch msg := msg.(type) {
	case button.SubmitMessage:
		if msg.ID == n.save_button.ID() {
			cmds = append(cmds, n.do_save())
		}
	}

	return n, tea.Batch(cmds...)
}

func split(v string) []string {
	re := regexp.MustCompile(",\\s*")
	var res []string
	for _, v := range re.Split(v, -1) {
		if u := strings.Trim(v, " "); u != "" {
			res = append(res, u)
		}
	}
	return res
}

func sanitize(v string) string { return strings.Trim(v, " ") }

func (n *Node) Value() form.Value[source.S] {
	pb := &dpb.Source{
		Header: n.source.Header().PB(),
		NodeId: n.source.NodeID(),
	}
	for _, _t := range n.titles {
		if title := sanitize(_t.Title.Value().Value); title != "" {
			pb.Titles = append(pb.Titles, &dpb.Title{
				Title:        title,
				Localization: sanitize(_t.Localization.Value().Value),
			})
		}
	}
	pb.PreviewUrl = sanitize(n.image.Value().Value)
	score, err := strconv.Atoi(sanitize(n.score.Value().Value))
	if err != nil {
		score = 0
	}
	pb.Score = int64(score)
	pb.Genres = split(sanitize(n.genres.Value().Value))
	pb.Studios = split(sanitize(n.studios.Value().Value))
	pb.Seasons = split(sanitize(n.seasons.Value().Value))
	pb.Authors = split(sanitize(n.authors.Value().Value))
	pb.Illustrators = split(sanitize(n.illustrators.Value().Value))

	return form.Value[source.S]{
		Key:   key,
		Value: source.Make(pb),
	}
}

func (n *Node) View() string {
	parts := []string{
		lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(
			fmt.Sprintf(
				"%v%v%v",
				view.WithHeader("Type", view.R(n.source).Type(), view.RenderHeaderModeNone),
				lipgloss.NewStyle().Foreground(color_profile.ForegroundNegligible).Render(" > "),
				view.WithHeader("ID", view.R(n.source).ID(), view.RenderHeaderModeNone),
			),
		),
	}

	for _, _t := range n.titles {
		if !_t.Title.IsInvisible() {
			parts = append(
				parts,
				lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						lipgloss.NewStyle().Margin(0, 1, 0, 0).Render(
							_t.Title.View(),
						),
						_t.Localization.View(),
					),
				),
			)
		}
	}

	parts = append(parts, lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Margin(0, 1, 0, 0).Render(
				n.score.View(),
			),
			lipgloss.NewStyle().Foreground(color_profile.ForegroundNegligible).Render(
				"(0 - 100)",
			),
		),
	))

	for _, m := range []interface {
		View() string
		IsInvisible() bool
	}{
		n.genres,
		n.studios,
		n.seasons,
		n.authors,
		n.illustrators,
	} {
		if !m.IsInvisible() {
			parts = append(parts, lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(m.View()))
		}
	}

	parts = append(parts, n.save_button.View())

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
