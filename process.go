package fileproxy

import (
	"bytes"
	"path"
)

func process(params *Params, filePath string, file *bytes.Buffer) (*bytes.Buffer, error) {
	if params.IsDefault {
		return file, nil
	}

	gm := NewGMBuilder()
	gm.Strip()

	if params.EnableResize {
		switch {
		case params.Width > 1 && params.Height > 1:
			gm.Resize(params.Width, params.Height, "!")
		case params.Width < 1 && params.Height > 1:
			gm.Resize(1, params.Height, "^")
		case params.Width > 1 && params.Height < 1:
			gm.Resize(params.Width, 1, "^")
		}
	}

	if params.EnableGravity {
		gm.Gravity(getGravity(params.Gravity))
	}

	if params.EnableCrop {
		gm.Crop(params.CropWidth, params.CropHeight)
	}

	if params.Quality != -1 {
		switch path.Ext(filePath) {
		case ".png":
			params.Quality = 100 - params.Quality
		}
		gm.Quality(params.Quality)
	}

	out := &bytes.Buffer{}
	err := gmProcess(file, out, gm)

	return out, err
}

func getGravity(cropPos string) string {
	switch cropPos {
	case "nw":
		return "NorthWest"
	case "n":
		return "North"
	case "ne":
		return "NorthEast"
	case "w":
		return "West"
	case "c":
		return "Center"
	case "e":
		return "East"
	case "sw":
		return "SouthWest"
	case "s":
		return "South"
	case "se":
		return "SouthEast"
	}

	return "NorthWest"
}
