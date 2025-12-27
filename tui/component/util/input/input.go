package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/util/input"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type M struct {
	*model_ui.Base

	input textinput.Model

	prompts map[bool]string
}

type O struct {
	model_ui.O

	Placeholder     string
	PromptUnfocused string
	PromptFocused   string
}

func New(o O) *M {
	t := textinput.New()
	t.Width = 50 // TODO
	t.Placeholder = o.Placeholder
	t.Prompt = o.PromptUnfocused
	m := &M{
		Base:  model_ui.New(o.O.WithNTabs(1)),
		input: t,
		prompts: map[bool]string{
			false: o.PromptUnfocused,
			true:  o.PromptFocused,
		},
	}
	return m
}

type SubmitMsg struct {
	ID string
	V  string
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{
		m.Base.Update(msg),
	}

	switch msg := msg.(type) {
	case model_ui.FocusMsg:
		if m.ID() == msg.ID {
			cmds = append(cmds, m.input.Focus())
		} else {
			m.input.Blur()
		}
	case model_ui.EOFMsg:
		if m.ID() == msg.ID {
			m.input.Blur()
		}
	case tea.MouseMsg:
		if input.IsMouseJustPressed(msg) && zone.Get(m.ID()).InBounds(msg) {
			cmds = append(cmds, model_ui.ToCommand(
				model_ui.FocusMsg{
					ID: m.ID(),
				},
			))
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
			m.input, c = m.input.Update(msg)
			cmds = append(cmds, c)
		}

		switch msg.Type {
		case tea.KeyEnter:
			v := m.Value()
			m.input.SetValue("")
			cmds = append(
				cmds,
				tea.Sequence(
					model_ui.ToCommand(SubmitMsg{
						ID: m.ID(),
						V:  v,
					}),
				),
			)
		}
	}

	m.input.Prompt = m.prompts[m.Focus()]

	return m, tea.Batch(cmds...)
}

func (m *M) Value() string { return m.input.Value() }

func (m *M) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).Padding(0, 1)
	if !m.Focus() {
		style = style.BorderForeground(lipgloss.Color("8"))
	}
	return m.RenderOrDie(
		zone.Mark(
			m.ID(),
			style.Render(m.input.View()),
		),
	)
}
