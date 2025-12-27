package card

import (
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/db/node"
	"github.com/minkezhang/truffle/tui/util/grid"
	"github.com/minkezhang/truffle/tui/util/input"
	"github.com/minkezhang/truffle/tui/util/titles"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

var (
	_ list.Item        = &I{}
	_ list.DefaultItem = &I{}
)

type D struct {
	width int
	list.DefaultDelegate
}

func (d D) Render(w io.Writer, m list.Model, index int, item list.Item) {
	style := lipgloss.NewStyle().Width(d.width)

	var buf strings.Builder
	d.DefaultDelegate.Render(&buf, m, index, item)
	w.Write(
		[]byte(
			zone.Mark(
				item.(*I).key,
				style.Render(buf.String()),
			),
		),
	)
}

type I struct {
	key  string
	node *node.N

	title string
	apis  string
}

// FilterValue complies with the list.Item interface.
//
// We will not support filtering items.
func (i *I) FilterValue() string { return "" }

func (i *I) Title() string {
	if i.title == "" {
		i.title = titles.Title(i.node.Virtual().Titles())
	}
	return i.title
}

func (i *I) Description() string {
	f := map[bool]string{
		true:  "⚑",
		false: "⚐",
	}
	return fmt.Sprintf("%s %s", f[i.node.IsQueued()], i.APIs())

}

func (i *I) APIs() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	var parts []string

	for _, a := range i.node.Atoms() {
		switch api := a.APIType(); api {
		case epb.API_API_VIRTUAL:
			parts = append(parts, style.Render("TRUFFLE"))
		default:
			parts = append(
				parts,
				style.Render(fmt.Sprintf("%s/%s", strings.ReplaceAll(api.String(), "API_", ""), a.APIID())),
			)
		}
	}

	if len(parts) <= 3 {
		return strings.Join(parts, ", ")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("Multiple APIs")
}

type O struct {
	model_ui.O

	Nodes []*node.N
}

type M struct {
	*model_ui.Base

	order []string // { list.GlobalIndex: key }
	list  list.Model
}

func Make(o O) M {
	items := []list.Item{}
	order := []string{}
	for _, n := range o.Nodes {
		k := zone.NewPrefix()
		order = append(order, k)

		items = append(items, &I{
			key:  k,
			node: n,
		})
	}

	delegate := list.NewDefaultDelegate()
	w := o.Column.Content - grid.GetFrame(delegate.Styles.SelectedTitle)
	d := D{
		DefaultDelegate: delegate,
		width:           w,
	}

	l := list.New(
		items,
		d,
		w,
		17*d.Height(),
	)
	l.DisableQuitKeybindings()
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)

	return M{
		Base:  model_ui.New(o.O.WithNTabs(len(o.Nodes))),
		list:  l,
		order: order,
	}
}

type SelectMsg *node.N

func (m M) Init() tea.Cmd {
	return model_ui.ToCommand(model_ui.ToNoticeMsg(fmt.Sprintf("results ID: %v", m.ID())))
} // DEBUG
func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{
		m.Base.Update(msg),
	}

	switch msg := msg.(type) {
	case tea.MouseMsg:
		if input.IsMouseJustPressed(msg) {
			for i, k := range m.order {
				if zone.Get(k).InBounds(msg) {
					m.list.Select(i)
					v := m.list.SelectedItem().(*I)
					cmds = append(
						cmds,
						model_ui.ToCommand(SelectMsg(v.node)),
						model_ui.ToCommand(model_ui.FocusMsg{
							BaseMsg: model_ui.BaseMsg{
								ID: m.ID(),
							},
							Index: i,
						}),
					)
				}
			}
		}
	case model_ui.FocusMsg:
		if m.ID() == msg.ID {
			m.list.Select(msg.Index) // TODO(minkezhang): Investigate why this is being selected over and over.
		}
	case tea.KeyMsg:
		if m.Focus() {
			switch msg.Type {
			case tea.KeyEsc:
				cmds = append(cmds, model_ui.ToCommand(model_ui.BlurMsg{
					BaseMsg: model_ui.BaseMsg{ID: m.ID()},
				}))
			case tea.KeyEnter:
				v := m.list.SelectedItem().(*I)
				cmds = append(
					cmds,
					model_ui.ToCommand(SelectMsg(v.node)),
					model_ui.ToCommand(model_ui.BlurMsg{
						BaseMsg: model_ui.BaseMsg{ID: m.ID()},
					}),
				)
			}
		}
	}

	// Sync list tab index with Truffle state.
	if m.Focus() {
		var c tea.Cmd
		m.list, c = m.list.Update(msg)
		cmds = append(cmds, c)
		slog.Debug(fmt.Sprintf("global index: %d", m.list.GlobalIndex()))
		m.SetIndex(m.list.GlobalIndex())
	}

	return m, tea.Batch(cmds...)
}

func (m M) View() string {
	if m.Focus() {
		return m.RenderOrDie(m.Column().Style().Render(m.list.View()))
	}
	return ""
}
