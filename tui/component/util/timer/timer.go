package timer

import (
	"time"

	"github.com/charmbracelet/bubbletea"
)

const (
	granularity = 100 * time.Millisecond
)

type TimeoutMsg struct {
	ID string
}

type tick struct {
	id       string
	n        int
	enqueued time.Time
}

type start struct {
	id    string
	n     int
	start time.Time
}

type reset start

type pause tick

type stop struct {
	id string
	n  int
}

type M struct {
	id string
	n  int // Number of times this is started
	d  time.Duration

	running bool
	start   time.Time
	elapsed time.Duration

	now   func() time.Time
	sleep func(d time.Duration)
}

func New(d time.Duration) *M {
	return &M{
		id:    "foo",
		d:     d,
		now:   time.Now,
		sleep: time.Sleep,
	}
}

func (m *M) ID() string    { return m.id }
func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case reset:
		if m.id == msg.id && m.n == msg.n {
			m.running = true
			m.n += 1
			m.start = msg.start
			m.elapsed = 0
			enqueued := m.now()
			cmds = append(cmds, func() tea.Msg {
				return tick{
					id:       m.id,
					enqueued: enqueued,
					n:        m.n,
				}
			})
		}
	case start:
		if !m.running && m.id == msg.id && m.n == msg.n {
			m.running = true
			m.n += 1
			m.start = msg.start
			m.elapsed = 0
			enqueued := m.now()
			cmds = append(cmds, func() tea.Msg {
				return tick{
					id:       m.id,
					enqueued: enqueued,
					n:        m.n,
				}
			})
		}
	case stop:
		if m.id == msg.id && m.n == msg.n {
			m.running = false
			m.elapsed = 0
		}
	case pause:
		if m.id == msg.id && m.n == msg.n {
			m.running = false
			m.elapsed += m.now().Sub(msg.enqueued)
			if m.elapsed >= m.d && m.running {
				cmds = append(cmds, func() tea.Msg {
					return TimeoutMsg{
						ID: m.id,
					}
				})
			}
		}
	case tick:
		if m.running && m.id == msg.id && m.n == msg.n {
			m.elapsed += m.now().Sub(msg.enqueued)
			if m.elapsed >= m.d {
				if m.running {
					m.running = false
					cmds = append(cmds, func() tea.Msg {
						return TimeoutMsg{
							ID: m.id,
						}
					})
				}
			} else {
				enqueued := m.now()
				cmds = append(cmds, func() tea.Msg {
					m.sleep(granularity)
					return tick{
						id:       m.id,
						enqueued: enqueued,
						n:        msg.n,
					}
				})
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *M) View() string { return (m.d - m.elapsed).Truncate(granularity).String() }

func (m *M) Pause() tea.Cmd {
	enqueued := m.now()
	return func() tea.Msg {
		return pause{
			id:       m.id,
			enqueued: enqueued,
			n:        m.n,
		}
	}
}

func (m *M) Stop() tea.Cmd {
	return func() tea.Msg {
		return stop{
			id: m.id,
			n:  m.n,
		}
	}
}

func (m *M) Start() tea.Cmd {
	enqueued := m.now()
	return func() tea.Msg {
		return start{
			id:    m.id,
			start: enqueued,
			n:     m.n,
		}
	}
}

func (m *M) Reset() tea.Cmd {
	enqueued := m.now()
	return func() tea.Msg {
		return reset{
			id:    m.id,
			start: enqueued,
			n:     m.n,
		}
	}
}
