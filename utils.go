package ivy

import (
	"strconv"
	"strings"

	"github.com/h2non/bimg"
)

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i
}

func splitWidthHeightString(s string) (int, int) {
	splits := strings.Split(s, "x")
	if len(splits) != 2 {
		return 0, 0
	}

	w, err := strconv.Atoi(splits[0])
	if err != nil {
		w = 0
	}

	h, err := strconv.Atoi(splits[1])
	if err != nil {
		h = 0
	}

	return w, h
}

func getGravityFromString(s string) bimg.Gravity {
	switch s {
	case "n":
		return bimg.GravityNorth
	case "e":
		return bimg.GravityEast
	case "s":
		return bimg.GravitySouth
	case "w":
		return bimg.GravityWest
	case "c":
		return bimg.GravityCentre
	default:
		return -1
	}
}

func getContentType(s bimg.ImageType) string {
	switch s {
	case bimg.JPEG:
		return "image/jpg"
	case bimg.PNG:
		return "image/png"
	case bimg.GIF:
		return "image/gif"
	case bimg.WEBP:
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func getBoolFromString(s string) bool {
	if s == "true" || s == "1" {
		return true
	}

	return false
}

func stringToImageType(s string) bimg.ImageType {
	switch s {
	case "jpg":
		return bimg.JPEG
	case "png":
		return bimg.PNG
	case "gif":
		return bimg.GIF
	case "webp":
		return bimg.WEBP
	default:
		return bimg.UNKNOWN
	}
}
