package ivy

import (
	"bytes"
)

func process(params *params, filePath string, file *bytes.Buffer) (*bytes.Buffer, error) {
	if params.isDefault {
		return file, nil
	}

	gm := newGMBuilder()
	gm.Strip()

	if params.enableResize {
		switch {
		case params.width > 1 && params.height > 1:
			gm.Resize(params.width, params.height, "!")
		case params.width < 1 && params.height > 1:
			gm.Resize(1, params.height, "^")
		case params.width > 1 && params.height < 1:
			gm.Resize(params.width, 1, "^")
		}
	}

	if params.enableGravity {
		gm.Gravity(getGravity(params.gravity))
	}

	if params.enableCrop {
		gm.Crop(params.cropWidth, params.cropHeight)
	}

	if params.quality != -1 {
		gm.Quality(params.quality)
	}

	out := &bytes.Buffer{}
	err := gm.Process(file, out)

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
