package notification

import (
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
	timer_ui "github.com/minkezhang/truffle/tui/component/util/timer"
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
)

type O struct {
	model_ui.O
}

func New(o O) *M {
	return &M{
		Base:  model_ui.New(o.O),
		timer: timer_ui.New(5 * time.Second),
	}
}

type message struct {
	t Type
	v string
}

type M struct {
	*model_ui.Base

	id int

	messages []message
	timer    tea.Model
}

const (
	timeout = 5 * time.Second
)

func (m *M) Init() tea.Cmd { return m.timer.Init() }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var c tea.Cmd

	m.timer, c = m.timer.Update(msg)
	cmds = append(cmds, c)

	switch msg := msg.(type) {
	case timer_ui.TimeoutMsg:
		if m.timer.(*timer_ui.M).ID() == msg.ID {
			m.messages = m.messages[:len(m.messages)-1]
		}
		if len(m.messages) > 0 {
			cmds = append(cmds, m.timer.(*timer_ui.M).Reset())
		}
	case model_ui.NoticeMsg:
		cmds = append(cmds, m.timer.(*timer_ui.M).Reset())
		m.messages = append(m.messages, message{
			t: TypeNotice,
			v: string(msg),
		})
	case model_ui.WarningMsg:
		cmds = append(cmds, m.timer.(*timer_ui.M).Reset())
		m.messages = append(m.messages, message{
			t: TypeWarning,
			v: error(msg).Error(),
		})
	}

	return m, tea.Batch(cmds...)
}

func (m *M) head() message {
	if len(m.messages) == 0 {
		return message{}
	}
	return m.messages[len(m.messages)-1]
}

func (m *M) body(n message) string {
	return ansi.Truncate(n.v, m.Column().Content, ellipsis)
}

func (m *M) View() string {
	if len(m.messages) == 0 {
		return ""
	}
	h := m.head()
	style := styles[h.t](m.Column().Style()).MarginBottom(1)
	return m.RenderOrDie(style.Render(m.body(h)))
}
