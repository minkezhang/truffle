// Package image loads a bubbletea model which returns a sixel-formatted
// representation of a URL.
//
// TOOD(minkezhang): Implement graphics via Kitty when available
// https://github.com/charmbracelet/bubbletea/issues/163.
//
// TODO(minkezhang): Support animated gifs.
package image

import (
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/mosaic"
	"github.com/disintegration/imaging"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type O struct {
	model_ui.O

	URL            string
	CacheDirectory string
}

type M struct {
	*model_ui.Base

	url       string
	filepath  string
	directory string

	cache string
}

func New(o O) *M {
	return &M{
		Base:      model_ui.New(o.O),
		url:       o.URL,
		directory: o.CacheDirectory,
	}
}

// render converts an image to a string representation of an image.
//
// TODO(minkezhang): Use Sixel support when it lands. See
// https://github.com/charmbracelet/bubbletea/issues/163.
func (m *M) render(img image.Image) string {
	img = imaging.Fit(img, m.Column().Content, height(m.Column().Content), imaging.Lanczos)

	// TODO(minkezhang): Width must be manually adjusted. See
	//
	// https://github.com/charmbracelet/x/issues/705.
	n := mosaic.New().Width(img.Bounds().Max.X * 2).Symbol(mosaic.Half)
	return n.Render(img)
}

type updateImageMsg struct {
	id      string
	payload string
}

func (m *M) Init() tea.Cmd {
	return func() tea.Msg {
		if m.url == "" {
			return nil
		}
		if s, err := data(*m); err == nil {
			return updateImageMsg{
				id:      m.ID(),
				payload: s,
			}
		} else {
			return model_ui.ToErrorMsg(err)
		}
	}
}

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateImageMsg:
		if m.ID() == msg.id {
			m.cache = string(msg.payload)
		}
	}
	return m, nil
}

func (m *M) View() string {
	return m.RenderOrDie(
		lipgloss.Place(
			m.Column().Content,
			// Two vertical pixels per character may leave a pixel
			// unaccounted for.
			height(m.Column().Content)/2+1,
			lipgloss.Top,
			lipgloss.Center,
			m.cache,
		),
	)
}
