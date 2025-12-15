package score

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"
)

type O struct {
	Score int64
}

type M struct {
	score    int64
	progress progress.Model
}

func New(o O) *M {
	return &M{
		score: o.Score,
		progress: progress.New(
			progress.WithFillCharacters('★', '☆'),
			progress.WithWidth(10),
		),
	}
}

func (m *M) Init() tea.Cmd                           { return nil }
func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *M) View() string                            { return m.progress.ViewAs(float64(m.score) / 100) }
