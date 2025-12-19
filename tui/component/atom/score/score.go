package score

import (
	"fmt"
	"reflect"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type O struct {
	model_ui.O

	Score int64
}

type M struct {
	*model_ui.Base

	score    int64
	progress progress.Model
}

func New(o O) *M {
	return &M{
		Base:  model_ui.New(o.O),
		score: o.Score,
		progress: progress.New(
			progress.WithFillCharacters('★', '☆'),
			progress.WithWidth(10),
		),
	}
}

type updateScoreMsg struct {
	model_ui.BaseMsg
	payload int64
}

func UpdateScoreAsync(m tea.Model, v int64) tea.Cmd {
	if m, ok := m.(*M); ok {
		return func() tea.Msg {
			return updateScoreMsg{
				BaseMsg: model_ui.BaseMsg{
					ID: m.ID(),
				},
				payload: v,
			}
		}
	}
	return model_ui.ErrorCmd(fmt.Errorf("incorrect model type: %v", reflect.TypeOf(m)))
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateScoreMsg:
		if m.ID() == msg.ID {
			m.score = msg.payload
		}
	}
	return m, nil
}

func (m *M) View() string { return m.RenderOrDie(m.progress.ViewAs(float64(m.score) / 100)) }
