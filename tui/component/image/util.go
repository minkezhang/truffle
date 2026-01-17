package image

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	_ "golang.org/x/image/webp"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/x/mosaic"
	"github.com/disintegration/imaging"
)

func height(width int) int {
	return int( // A4 ratio; in pixels
		math.Trunc(float64(width) * 1.414),
	)
}

func data(url string, cache_directory string, width int) (string, error) {
	buf := md5.Sum([]byte(url))
	fn := filepath.Join(
		cache_directory,
		fmt.Sprintf(
			"%s%s",
			hex.EncodeToString(buf[:]),
			filepath.Ext(url),
		),
	)
	var data []byte

	// Save data locally if not exists
	if _, err := os.Stat(fn); errors.Is(err, os.ErrNotExist) {
		slog.Debug(fmt.Sprintf("downloading remote asset: %s", url))
		res, err := http.Get(url)
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
		return "", fmt.Errorf("cannot decode %s image: %v", filepath.Ext(url), err)
	} else {
		return render(img, width), nil
	}
}

// paletted converts an image to be a paletted image.
//
// In a paletted image, each pixel must correspond to a limited palette provided
// as an encoding artifact, instead of storing (potentially N * M colors) per
// pixel.
//
// See https://groups.google.com/g/golang-nuts/c/28Kk1FfG5XE.
func paletted(img image.Image) *image.Paletted {
	if img, ok := img.(*image.Paletted); ok {
		return img
	}

	opts := gif.Options{
		NumColors: len(palette.Plan9),
		Drawer:    draw.FloydSteinberg,
	}

	res := image.NewPaletted(img.Bounds(), palette.Plan9[:opts.NumColors])

	if opts.Quantizer != nil {
		res.Palette = opts.Quantizer.Quantize(
			make(color.Palette, 0, opts.NumColors),
			img,
		)
	}

	opts.Drawer.Draw(res, img.Bounds(), img, image.ZP)

	return res
}

// render converts an image to a string representation of an image.
//
// TODO(minkezhang): Use Sixel support when it lands. See
// https://github.com/charmbracelet/bubbletea/issues/163.
func render(img image.Image, w int) string {
	img = imaging.Fit(img, w, height(w), imaging.Lanczos)

	// TODO(minkezhang): Width must be manually adjusted. See
	//
	// https://github.com/charmbracelet/x/issues/705.
	n := mosaic.New().Width(img.Bounds().Max.X * 2).Symbol(mosaic.Half)
	return n.Render(img)
}
