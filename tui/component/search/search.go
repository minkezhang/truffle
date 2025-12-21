package search

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type O struct {
	model_ui.O
}

type M struct {
	*model_ui.Base

	search textinput.Model
}

func New(o O) *M {
	m := &M{
		Base:   model_ui.New(o.O.WithNTabs(1)),
		search: textinput.New(),
	}
	m.search.Width = 50 // TODO
	m.search.Placeholder = "The Apothecary Diaries"
	m.search.Prompt = "  "
	return m
}

func (m *M) Init() tea.Cmd { return nil }

type QueryMsg string

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{m.Base.Update(msg)}

	switch msg := msg.(type) {
	case model_ui.FocusMsg:
		if m.ID() == msg.ID {
			cmds = append(cmds, m.search.Focus())
		} else {
			m.search.Blur()
		}
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress {
			if zone.Get(m.ID()).InBounds(msg) {
				cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
					BaseMsg: model_ui.BaseMsg{
						ID: m.ID(),
					},
					Index: m.Index(),
				}))
			}
		}
	case tea.KeyMsg:
		if !m.Focus() {
			break
		}

		// Mouse input is being passed into search input as KeyMsg; it
		// is unclear how or why this is happening.
		//
		// Such a KeyMsg is of the form
		//
		//   { Type: KeyRunes, Alt: true, Runes: []rune{'['} }
		//
		// or
		//
		//   { Type: KeyRunes, Alt: false, Runes: []rune{...} }
		if (len(msg.Runes) <= 1 && !msg.Alt) || msg.Paste {
			var c tea.Cmd
			m.search, c = m.search.Update(msg)
			cmds = append(cmds, c)
		}

		switch msg.Type {
		case tea.KeyEnter:
			v := m.search.Value()
			m.search.SetValue("")
			cmds = append(cmds,
				model_ui.ToCommand(QueryMsg(v)),
				model_ui.ToCommand(model_ui.BlurMsg{
					BaseMsg: model_ui.BaseMsg{ID: m.ID()},
				},
				))
		}
	}

	m.search.Prompt = map[bool]string{
		true:  "⚲ ",
		false: "  ",
	}[m.Focus()]

	return m, tea.Batch(cmds...)

}

func (m *M) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).Padding(0, 1)
	if !m.Focus() {
		style = style.BorderForeground(lipgloss.Color("8"))
	}
	return m.RenderOrDie(lipgloss.JoinVertical(
		lipgloss.Left,
		zone.Mark(
			m.ID(),
			style.Render(m.search.View()),
		),
	))
}
