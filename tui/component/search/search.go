package search

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle/tui/component/button"
	"github.com/minkezhang/truffle/tui/component/checkbox_group"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/textinput"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	directory "github.com/minkezhang/truffle/tui/component/directory/focusable"
)

var (
	key_checkboxes = map[string]map[string]form.Key{
		"api": map[string]form.Key{
			"truffle": form.Key{"Truffle", epb.SourceAPI_SOURCE_API_TRUFFLE.String()},
			"mal":     form.Key{"MAL", epb.SourceAPI_SOURCE_API_MAL.String()},
		},
		"t": map[string]form.Key{
			"anime":       form.Key{"Anime", epb.SourceType_SOURCE_TYPE_SERIES_ANIME.String()},
			"anime_movie": form.Key{"Anime Movie", epb.SourceType_SOURCE_TYPE_MOVIE_ANIME.String()},
			"manga":       form.Key{"Manga", epb.SourceType_SOURCE_TYPE_BOOK_MANGA.String()},
			"light_novel": form.Key{"Light Novel", epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL.String()},
		},
	}
	key_options = map[string]form.Key{
		"nsfw": form.Key{"NSFW", "options-nsfw"},
	}

	key_submit_message = map[string]form.Key{
		"query":   form.Key{"Search", "search-textinput"},
		"apis":    form.Key{"APIs", "search-apis"},
		"types":   form.Key{"Media", "search-source-types"},
		"options": form.Key{"Options", "search-options"},
		"submit":  form.Key{"Search", "search-submit"},
	}
)

type Node struct {
	*focusable.Node

	column *column.C
	input  *textinput.Node

	apis          *checkbox_group.Node
	source_types  *checkbox_group.Node
	options       *checkbox_group.Node
	submit_button *button.Node
	is_expanded   bool
	clickable     *clickable.Node
}

type O struct {
	ParentID string
	Column   *column.C
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("search", o.ParentID, 1),
		column: o.Column,
	}
	n.input = textinput.New(textinput.O{
		Prefix:      "search-textinput",
		ParentID:    n.ID(),
		Width:       o.Column.Content() - 4,
		Placeholder: "Frieren t:manga t:anime api:mal nsfw:false, mal:manga/52991",
		Prompt:      "⚲ ",
		Value: form.Value[string]{
			Key: key_submit_message["query"],
		},
	})
	n.apis = checkbox_group.New(checkbox_group.O{
		Prefix:   "search-apis",
		ParentID: n.ID(),
		IsRadio:  false,
		Value: form.Value[[]form.Value[bool]]{
			Key: key_submit_message["apis"],
			Value: []form.Value[bool]{
				form.Value[bool]{key_checkboxes["api"]["truffle"], true},
				form.Value[bool]{key_checkboxes["api"]["mal"], true},
			},
		},
	})
	n.source_types = checkbox_group.New(checkbox_group.O{
		Prefix:   "search-source-types",
		ParentID: n.ID(),
		IsRadio:  false,
		Value: form.Value[[]form.Value[bool]]{
			Key: key_submit_message["types"],
			Value: []form.Value[bool]{
				form.Value[bool]{key_checkboxes["t"]["anime"], true},
				form.Value[bool]{key_checkboxes["t"]["anime_movie"], true},
				form.Value[bool]{key_checkboxes["t"]["manga"], true},
				form.Value[bool]{key_checkboxes["t"]["light_novel"], true},
			},
		},
	})
	n.options = checkbox_group.New(checkbox_group.O{
		Prefix:   "search-options",
		ParentID: n.ID(),
		IsRadio:  false,
		Value: form.Value[[]form.Value[bool]]{
			Key: key_submit_message["options"],
			Value: []form.Value[bool]{
				form.Value[bool]{key_options["nsfw"], false},
			},
		},
	})
	n.submit_button = button.New(button.O{
		ParentID: n.ID(),
		Key:      key_submit_message["submit"],
	})
	n.clickable = clickable.New(n.ID())
	return n
}

func (n *Node) Init() tea.Cmd {
	cmds := []tea.Cmd{
		n.input.Init(),
		n.clickable.Init(),
		n.apis.Init(),
		n.source_types.Init(),
		n.options.Init(),
		n.submit_button.Init(),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	}
	for _, c := range []directory.Node{
		n.apis,
		n.source_types,
		n.options,
		n.submit_button,
	} {
		cmds = append(cmds, c.SetIsInvisible(!n.is_expanded))
	}
	return tea.Sequence(cmds...)
}

type SubmitSearchMessage struct {
	ID          string
	Query       form.Value[string]
	APIs        form.Value[map[epb.SourceAPI]bool]
	SourceTypes form.Value[map[epb.SourceType]bool]
	Options     form.Value[[]option.O]
}

func (n *Node) submit() tea.Cmd {
	// Ignore blank queries.
	if n.input.Value().Value == "" {
		return nil
	}

	apis := map[epb.SourceAPI]bool{}
	types := map[epb.SourceType]bool{}
	_options := map[string]option.O{}

	for _, api := range n.apis.Value().Value {
		apis[epb.SourceAPI(epb.SourceAPI_value[api.Key.Key])] = api.Value
	}
	for _, t := range n.source_types.Value().Value {
		types[epb.SourceType(epb.SourceType_value[t.Key.Key])] = t.Value
	}
	for _, opt := range n.options.Value().Value {
		switch opt.Key.Key {
		case key_options["nsfw"].Key:
			_options[key_options["nsfw"].Key] = option.NSFW(opt.Value)
		}
	}

	// Set option overrides
	var invalid []string
	var query []string
	for _, token := range strings.Split(n.input.Value().Value, " ") {
		if head, tail, ok := strings.Cut(token, ":"); ok {
			switch head {
			case "-api":
				fallthrough
			case "api":
				if _, ok := key_checkboxes["api"]; ok {
					if k, ok := key_checkboxes["api"][tail]; ok {
						apis[epb.SourceAPI(epb.SourceAPI_value[k.Key])] = (head == "api")
					} else {
						invalid = append(invalid, token)
					}
				}
			case "-t":
				fallthrough
			case "t":
				if _, ok := key_checkboxes["t"]; ok {
					if k, ok := key_checkboxes["t"][tail]; ok {
						types[epb.SourceType(epb.SourceType_value[k.Key])] = (head == "t")
					} else {
						invalid = append(invalid, token)
					}
				}
			case "nsfw":
				if v, err := strconv.ParseBool(tail); err == nil {
					_options[key_options["nsfw"].Key] = option.NSFW(v)
				} else {
					invalid = append(invalid, token)
				}
			default:
				query = append(query, token)
			}
		} else {
			query = append(query, token)
		}
	}

	var options []option.O
	for _, opt := range _options {
		options = append(options, opt)
	}

	m := SubmitSearchMessage{
		ID:          n.ID(),
		Query:       form.Value[string]{key_submit_message["query"], strings.Join(query, " ")},
		APIs:        form.Value[map[epb.SourceAPI]bool]{key_submit_message["apis"], apis},
		SourceTypes: form.Value[map[epb.SourceType]bool]{key_submit_message["types"], types},
		Options:     form.Value[[]option.O]{key_submit_message["options"], options},
	}

	cmds := []tea.Cmd{
		func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelDebug,
				fmt.Sprintf("%v: submitting search query %v", n.ID(), m),
			)
		},
		func() tea.Msg {
			return m
		},
	}
	if len(invalid) > 0 {
		cmds = append(cmds, func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf("%v: invalid search tokens: %v", n.ID(), strings.Join(invalid, " ")),
			)
		})
	}
	return tea.Sequence(cmds...)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	for _, m := range []tea.Model{
		n.input,
		n.apis,
		n.source_types,
		n.options,
		n.submit_button,
		n.clickable,
	} {
		_, c = m.Update(msg)
		cmds = append(cmds, c)
	}

	switch msg := msg.(type) {
	case button.SubmitMessage:
		if msg.ID == n.submit_button.ID() {
			cmds = append(cmds, n.submit())
		}
	case tea.KeyMsg:
		if n.input.FocusState() == types.FocusStateActive {
			if msg.Type == tea.KeyEnter {
				cmds = append(cmds, n.submit())
			}
		}
		if n.FocusState() == types.FocusStateActive {
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				n.is_expanded = !n.is_expanded
				for _, c := range []directory.Node{
					n.apis,
					n.source_types,
					n.options,
					n.submit_button,
				} {
					cmds = append(cmds, c.SetIsInvisible(!n.is_expanded))
				}
			}
		}
	case clickable.Click:
		if msg.ID == n.clickable.ID() {
			n.is_expanded = !n.is_expanded
			cmds = append(cmds, func() tea.Msg {
				return types.FocusMessage{
					ID: n.ID(),
				}
			})
			for _, c := range []directory.Node{
				n.apis,
				n.source_types,
				n.options,
				n.submit_button,
			} {
				cmds = append(cmds, c.SetIsInvisible(!n.is_expanded))
			}
		}
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	parts := []string{
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			n.input.View(),
			lipgloss.NewStyle().PaddingLeft(1).Foreground(
				color_profile.UIForeground[n.FocusState()],
			).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(
				color_profile.UIForeground[n.input.FocusState()],
			).Render(
				zone.Mark(
					n.clickable.ID(),
					map[bool]string{
						false: "(+)",
						true:  "(-)",
					}[n.is_expanded],
				),
			),
		),
	}
	if n.is_expanded {
		parts = append(parts,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				n.apis.View(),
				lipgloss.NewStyle().MarginLeft(1).Render(
					n.source_types.View(),
				),
				lipgloss.NewStyle().MarginLeft(1).Render(
					n.options.View(),
				),
			),
			n.submit_button.View(),
		)
	}
	return n.column.RenderOrDie(lipgloss.JoinVertical(lipgloss.Left, parts...))
}
