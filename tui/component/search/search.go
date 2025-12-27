package search

import (
	"fmt" // DEBUG

	"github.com/charmbracelet/bubbletea"

	input_ui "github.com/minkezhang/truffle/tui/component/util/input"
	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type O struct {
	model_ui.O
}

type M struct {
	*model_ui.Base

	input tea.Model
}

func New(o O) *M {
	return &M{
		Base: model_ui.New(o.O),
		input: input_ui.New(input_ui.O{
			O:               o.O.WithNTabs(1),
			Placeholder:     "The Apothecary Diaries",
			PromptUnfocused: "  ",
			PromptFocused:   "⚲ ",
		}),
	}
}

type QueryMsg string

func (m *M) Init() tea.Cmd {
	return tea.Sequence( // DEBUG
		m.input.Init(),
		model_ui.ToCommand(model_ui.ToNoticeMsg(fmt.Sprintf("Search Bar ID: %v", m.ID()))),
	)
}

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case model_ui.BlurMsg:
		if m.input.(*input_ui.M).ID() == msg.ID {
			cmds = append(cmds,
				model_ui.ToCommand(model_ui.BlurMsg{
					BaseMsg: model_ui.BaseMsg{
						m.ID(),
					},
					IsEnd: msg.IsEnd,
				}),
				model_ui.ToCommand(model_ui.ToNoticeMsg(
					fmt.Sprintf("search caught input blur and blurring self search: self = %v", m.ID()),
				)),
			)
		}
	case model_ui.FocusMsg:
		if m.ID() == msg.ID {
			cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
				BaseMsg: model_ui.BaseMsg{
					m.input.(*input_ui.M).ID(),
				},
				Index: 0,
			}))
		}
	case input_ui.SubmitMsg:
		if m.input.(*input_ui.M).ID() == msg.ID {
			cmds = append(
				cmds,
				model_ui.ToCommand(QueryMsg(msg.V)),
			)
		}
	}

	var c tea.Cmd
	m.input, c = m.input.Update(msg)
	cmds = append(cmds, c)

	return m, tea.Batch(cmds...)

}

func (m *M) View() string { return m.RenderOrDie(m.input.View()) }
