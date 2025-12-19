package card

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/node"
	"github.com/minkezhang/truffle/tui/util/titles"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

var (
	_ list.Item        = &I{}
	_ list.DefaultItem = &I{}
)

type I struct {
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

	list list.Model
}

func Make(o O) M {
	items := []list.Item{}
	for _, n := range o.Nodes {
		items = append(items, &I{
			node: n,
		})
	}
	l := list.New(
		items,
		list.NewDefaultDelegate(),
		o.Column.Content,
		50,
	)
	l.DisableQuitKeybindings()
	return M{
		Base: model_ui.New(o.O),
		list: l,
	}
}

func (m M) Init() tea.Cmd { return nil }
func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var c tea.Cmd

	m.list, c = m.list.Update(msg)

	return m, c
}

func (m M) View() string { return m.RenderOrDie(m.list.View()) }
