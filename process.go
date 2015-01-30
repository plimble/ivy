package fileproxy

import (
	"bytes"
	"fmt"
	"github.com/disintegration/gift"
	"image"
	"image/gif"
	"path"

	"image/jpeg"
	"image/png"
	"io"
)

const (
	JPG  = ".jpg"
	JPEG = ".jpeg"
	GIF  = ".gif"
	PNG  = ".png"
)

func process(params *Params, filePath string, file *bytes.Buffer) (*bytes.Buffer, error) {
	ext := path.Ext(filePath)
	img, err := decode(ext, file)
	if err != nil {
		return nil, err
	}

	var dst image.Image

	if params.raw {
		return file, nil
	}

	if params.CropMode != "" {
		switch params.CropMode {
		case CropModeExact:
			dst = crop(img, params.Width, params.Height, params.CropPos)
		case CropModeScale:
			imgBound := img.Bounds()
			ratioX := float64(imgBound.Dx()) / float64(imgBound.Dy())
			ratioY := float64(imgBound.Dy()) / float64(imgBound.Dx())
			//wide
			if ratioX > ratioY {
				dst = resize(img, 0, params.Height)
			} else {
				dst = resize(img, params.Width, 0)
			}

			dst = crop(dst, params.Width, params.Height, params.CropPos)
		}
	} else {
		if params.Width > 1 || params.Height > 1 {
			dst = resize(img, params.Width, params.Height)
		} else {
			dst = img
		}
	}

	buffer := &bytes.Buffer{}
	err = encode(buffer, ext, dst, params)
	return buffer, err
}

func resize(img image.Image, width, height int) image.Image {
	imgBound := img.Bounds()

	g := gift.New()

	g.Add(gift.Resize(width, height, gift.LanczosResampling))

	dst := image.NewRGBA(g.Bounds(imgBound))
	g.Draw(dst, img)

	return dst
}

func crop(img image.Image, width, height int, pos string) image.Image {
	imgBound := img.Bounds()
	if width < 1 {
		width = imgBound.Dx()
	}
	if height < 1 {
		height = imgBound.Dy()
	}

	x0, y0, x1, y1 := getCropPos(pos, width, height, imgBound.Dx(), imgBound.Dy())

	g := gift.New()
	g.Add(gift.Crop(image.Rect(x0, y0, x1, y1)))

	dst := image.NewRGBA(g.Bounds(imgBound))
	g.Draw(dst, img)

	return dst
}

func getCropPos(pos string, width, height, bx, by int) (int, int, int, int) {
	switch pos {
	case CropTopLeft:
		return 0, 0, width, height
	case CropTopCenter:
		x0 := getTopLeft(bx, width)
		return x0, 0, x0 + width, height
	case CropTopRight:
		return bx - width, 0, bx, height
	case CropMiddleLeft:
		y0 := getTopLeft(by, height)
		return 0, y0, width, y0 + height
	case CropMiddleCenter:
		x0 := getTopLeft(bx, width)
		y0 := getTopLeft(by, height)
		return x0, y0, x0 + width, y0 + height
	case CropMiddleRight:
		y0 := getTopLeft(by, height)
		return bx - width, y0, bx, y0 + height
	case CropBottomLeft:
		return 0, by, width, by - height
	case CropBottomCenter:
		x0 := getTopLeft(bx, width)
		return x0, by - height, x0 + width, by
	case CropBottomRight:
		return bx - width, by - height, bx, by
	}

	return 0, 0, bx, by
}

func getTopLeft(bound, dis int) int {
	return (bound - dis) / 2
}

func encode(bf io.Writer, ext string, img image.Image, params *Params) error {
	switch ext {
	case JPG, JPEG:
		return jpeg.Encode(bf, img, &jpeg.Options{Quality: params.Quality})
	case PNG:
		return png.Encode(bf, img)
	case GIF:
		return gif.Encode(bf, img, nil)
	default:
		return fmt.Errorf(`unsupported image format: "%s"`, ext)
	}
}

func decode(ext string, r io.Reader) (image.Image, error) {
	switch ext {
	case JPG, JPEG:
		return jpeg.Decode(r)
	case PNG:
		return png.Decode(r)
	case GIF:
		return gif.Decode(r)
	default:
		return nil, fmt.Errorf(`unsupported image format: "%s"`, ext)
	}
}
