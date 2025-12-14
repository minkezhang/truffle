// Package image loads a bubbletea model which returns a sixel-formatted
// representation of a URL.
//
// TOOD(minkezhang): Implement graphics via Kitty when available
// https://github.com/charmbracelet/bubbletea/issues/163.
package image

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/webp"
)

type O struct {
	URL            string
	Width          int
	Height         int
	CacheDirectory string
}

type M struct {
	url       string
	filepath  string
	directory string

	width  int
	height int
	cache  string
}

func New(o O) *M {
	return &M{
		url:       o.URL,
		width:     o.Width,
		height:    o.Height,
		directory: o.CacheDirectory,
	}
}

type message struct {
	s string
	e error
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

		// write the whole body at once
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

	var err error
	var img image.Image

	switch t := filepath.Ext(m.url); t {
	case ".webp":
		img, err = webp.Decode(bytes.NewReader(data))
		if err != nil {
			return "", fmt.Errorf("cannot decode .webp image: %v", err)
		}
	default:
		img, _, err = image.Decode(bytes.NewReader(data))
		if err != nil {
			return "", fmt.Errorf("cannot decode %s image: %v", t, err)
		}
	}
	return m.render(img), nil
}

// resize converts an image to a string representation of an image.
//
// This relies on the fact that each character in the terminal has a 2 x 1
// aspect ratio. The top half of the character can represent a pixel at some y
// coordinate using a tinted ▀ character, and the y + 1 pixel can be represented
// by the character background color (using lipgloss styles).
//
// From https://github.com/knipferrc/fm.
func (m *M) render(img image.Image) string {
	img = imaging.Fit(img, m.width, m.height, imaging.Lanczos)

	buf := strings.Builder{}

	for y := 0; y < img.Bounds().Max.Y; y += 2 {
		for x := 0; x < img.Bounds().Max.X; x++ {
			c1, _ := colorful.MakeColor(img.At(x, y))
			c2, _ := colorful.MakeColor(img.At(x, y+1))

			ct := lipgloss.Color(c1.Hex()) // Top pixel
			cb := lipgloss.Color(c2.Hex())

			style := lipgloss.NewStyle().Foreground(ct).Background(cb)
			buf.WriteString(style.Render("▀"))
		}

		buf.WriteString("\n")
	}

	return buf.String()
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
	return nil, nil
}

func (m *M) View() string { return m.cache }
