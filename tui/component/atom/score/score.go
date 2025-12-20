package score

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type O struct {
	model_ui.O

	Score int64
}

type M struct {
	model_ui.Base

	score    int64
	progress progress.Model
}

func Make(o O) M {
	return M{
		Base:  model_ui.Make(o.O),
		score: o.Score,
		progress: progress.New(
			progress.WithFillCharacters('★', '☆'),
			progress.WithWidth(10),
		),
	}
}

func (m M) Init() tea.Cmd                           { return nil }
func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m M) View() string { return m.RenderOrDie(m.progress.ViewAs(float64(m.score) / 100)) }
