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
	// m.search.Cursor TODO
	return m
}

func (m *M) Init() tea.Cmd { return m.search.Focus() }

type QueryMsg string

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{m.Base.Update(msg)}

	switch msg := msg.(type) {
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
		switch msg.Type {
		case tea.KeyEnter:
			v := m.search.Value()
			cmds = append(cmds, func() tea.Msg { return QueryMsg(v) })
			m.search.SetValue("")
		}
	}

	var c tea.Cmd

	if m.Focus() {
		m.search, c = m.search.Update(msg)
		cmds = append(cmds, c)
	}

	m.search.Prompt = map[bool]string{
		true:  "⚲ ", // TODO
		false: "  ",
	}[m.Focus()]

	return m, tea.Batch(cmds...)

}

func (m *M) View() string {
	return m.RenderOrDie(lipgloss.JoinVertical(
		lipgloss.Left,
		zone.Mark(
			m.ID(),
			lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).Padding(0, 1).Render(
				m.search.View(),
			),
		),
	))
}
