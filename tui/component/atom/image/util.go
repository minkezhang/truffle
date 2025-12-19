package image

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"log/slog"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

func height(width int) int {
	return int( // A4 ratio; in pixels
		math.Trunc(float64(width) * 1.414),
	)
}

func data(m M) (string, error) {
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
