// Package image loads a bubbletea model which returns a sixel-formatted
// representation of a URL.
//
// TOOD(minkezhang): Implement graphics via Kitty when available
// https://github.com/charmbracelet/bubbletea/issues/163.
//
// TODO(minkezhang): Support animated gifs.
package image

import (
	"math"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

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

type message struct {
	s string
	e error
}

func height(width int) int {
	return int( // A4 ratio; in pixels
		math.Trunc(float64(width) * 1.414),
	)
}

func (m *M) data() (string, error) {
	buf := md5.Sum([]byte(m.url))
	fn := filepath.Join(
		m.directory,
		fmt.Sprintf(
			"%s%s",
			hex.EncodeToString(buf[:]),
			filepath.Ext(m.url),
		),
	)
	var data []byte

	// Save data locally if not exists
	if _, err := os.Stat(fn); errors.Is(err, os.ErrNotExist) {
		slog.Debug(fmt.Sprintf("downloading remote asset: %s", m.url))
		res, err := http.Get(m.url)
		if err != nil || res.StatusCode != 200 {
			return "", fmt.Errorf("cannot download image: %v", err)
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("cannot read image: %v", err)
		}

		err = os.WriteFile(fn, data, 0644)
		if err != nil {
			return "", fmt.Errorf("cannot write to file %s: %v", fn, err)
		}
	}

	if data == nil {
		var err error
		data, err = os.ReadFile(fn)
		if err != nil {
			return "", fmt.Errorf("cannot read file %s: %v", fn, err)
		}
	}

	if img, _, err := image.Decode(bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("cannot decode %s image: %v", filepath.Ext(m.url), err)
	} else {
		return m.render(img), nil
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

func (m *M) Init() tea.Cmd {
	return func() tea.Msg {
		s, err := m.data()
		return message{
			s: s,
			e: err,
		}
	}
}

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case message:
		if msg.e != nil {
			slog.Error(fmt.Sprintf("image.go: %v", msg.e))
		} else {
			m.cache = msg.s
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
