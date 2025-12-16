package image

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
)

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
