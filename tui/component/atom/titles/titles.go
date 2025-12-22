package title

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle/tui/util/titles"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

var (
	priority = map[string]int{
		"en": 0,
		"":   1,
		"ja": 2,
	}
)

type O struct {
	model_ui.O

	Titles []atom.T
}

type M struct {
	*model_ui.Base

	titles []atom.T
}

func Make(o O) M {
	m := M{
		titles: titles.Sort(o.Titles),
		Base:   model_ui.New(o.O),
	}
	return m
}

func (m M) Init() tea.Cmd                           { return nil }
func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m M) View() string {
	var parts []string
	titles := m.titles
	if len(m.titles) > 2 {
		titles = titles[:2]
	}

	style := m.Column().Style()

	if len(titles) == 0 {
		style = style.Bold(true).Foreground(lipgloss.Color("8"))
		parts = append(parts, style.Render("Unknown Title"))
	} else {
		for i, t := range titles {
			if i == 0 {
				style = style.Bold(true)
			} else {
				style = style.Foreground(lipgloss.Color("8"))
			}
			parts = append(parts, style.Render(t.Title))
		}
	}
	return m.RenderOrDie(strings.Join(parts, "\n"))
}
