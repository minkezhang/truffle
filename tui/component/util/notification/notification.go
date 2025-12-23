// TODO(minkezhang): Add timer.
// TODO(minkezhang): Add prefix.
package notification

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/util/input"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

const (
	ellipsis = "…"
)

type Type int

const (
	TypeNotice Type = iota
	TypeWarning
)

var (
	styles = map[Type]func(s lipgloss.Style) lipgloss.Style{
		TypeNotice: func(s lipgloss.Style) lipgloss.Style {
			return s.Background(lipgloss.Color("8"))
		},
		TypeWarning: func(s lipgloss.Style) lipgloss.Style {
			return s.Background(lipgloss.Color("3")).Foreground(lipgloss.Color("9")).Bold(true)
		},
	}

	prefix = map[Type]string{
		TypeNotice:  "NOTICE",
		TypeWarning: "WARNING",
	}
)

type O struct {
	model_ui.O
}

func New(o O) *M {
	return &M{
		Base: model_ui.New(o.O),
	}
}

type message struct {
	t Type
	v string
}

type M struct {
	*model_ui.Base

	messages []message
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if input.IsMouseJustPressed(msg) {
			if zone.Get(m.ID()).InBounds(msg) {
				m.messages = m.messages[:len(m.messages)-1]
			}
		}
	case model_ui.NoticeMsg:
		m.messages = append(m.messages, message{
			t: TypeNotice,
			v: string(msg),
		})
	case model_ui.WarningMsg:
		m.messages = append(m.messages, message{
			t: TypeWarning,
			v: error(msg).Error(),
		})
	}
	return m, nil
}

func (m *M) head() message {
	if len(m.messages) == 0 {
		return message{}
	}
	return m.messages[len(m.messages)-1]
}

func (m *M) body(n message) string {
	message := ""
	if n.v != "" {
		message = fmt.Sprintf("%s: %s", prefix[n.t], n.v)
	}
	return ansi.Truncate(message, m.Column().Content-2, ellipsis)
}

func (m *M) View() string {
	h := m.head()
	style := styles[h.t](m.Column().Style()).MarginTop(1)
	if s := styles[h.t](lipgloss.NewStyle()).Render(m.body(h)); lipgloss.Width(s) > 0 {
		return m.RenderOrDie(style.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				zone.Mark(
					m.ID(),
					styles[h.t](lipgloss.NewStyle()).Foreground(lipgloss.Color("15")).Bold(true).Render("X"),
				),
				styles[h.t](lipgloss.NewStyle()).Render(" "),
				s,
			),
		))
	}
	return m.RenderOrDie(style.Render(""))
}
