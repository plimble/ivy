package fileproxy

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ParamWidth    = "w"
	ParamHeight   = "h"
	ParamCropPos  = "p"
	ParamScale    = "s"
	ParamCropMode = "c"
	ParamQuality  = "q"

	CropModeExact = "e"
	CropModeScale = "s"

	CropTopLeft      = "tl"
	CropTopCenter    = "tc"
	CropTopRight     = "tr"
	CropMiddleLeft   = "ml"
	CropMiddleCenter = "mc"
	CropMiddleRight  = "mr"
	CropBottomLeft   = "bl"
	CropBottomCenter = "bc"
	CropBottomRight  = "br"

	DefaultScale    = 1
	DefaultQuality  = 100
	DefaultCropMode = ""
	DefaultCropPos  = ""
)

type Params struct {
	Width     int
	Height    int
	CropMode  string
	CropPos   string
	Scale     int
	Quality   int
	IsDefault bool
}

func parseParams(paramsStr string) (*Params, error) {
	params := &Params{0, 0, DefaultCropMode, DefaultCropPos, DefaultScale, DefaultQuality, true}
	if paramsStr == "" || paramsStr == "_" {
		return params, nil
	}

	parts := strings.Split(paramsStr, ",")
	for _, part := range parts {
		if len(part) < 3 {
			return nil, fmt.Errorf("invalid parameter: %s", part)
		}
		if string(part[1]) != "_" {
			return nil, fmt.Errorf("invalid parameter: %s", part)
		}
		key := string(part[0])
		value := string(part[2:])

		switch key {
		case ParamWidth, ParamHeight:
			value, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse value for parameter: %s", key)
			}
			if value < 1 {
				return nil, fmt.Errorf("value %d must be > 0: %s", value, key)
			}
			if key == ParamWidth {
				params.Width = value
			} else {
				params.Height = value
			}
			params.IsDefault = false
		case ParamCropMode:
			value = strings.ToLower(value)
			if !isValidCropMode(value) {
				return nil, fmt.Errorf("invalid value for %s", key)
			}
			params.CropMode = value
			params.IsDefault = false
		case ParamCropPos:
			value = strings.ToLower(value)
			if !isValidCropPos(value) {
				return nil, fmt.Errorf("invalid value for %s", key)
			}
			params.CropPos = value
			params.IsDefault = false
		case ParamScale:
			value, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse value for parameter: %s", key)
			}
			if value < 1 {
				return nil, fmt.Errorf("value %d must be > 0: %s", value, key)
			}
			params.Scale = value
			params.IsDefault = false
		case ParamQuality:
			value, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse value for parameter: %s", key)
			}
			if value < 1 {
				return nil, fmt.Errorf("value %d must be > 0: %s", value, key)
			}
			params.Quality = value
			params.IsDefault = false
		default:
			return nil, fmt.Errorf("invalid parameter: %s", part)
		}
	}

	return params, nil
}

func (p *Params) String() string {
	return fmt.Sprintf("%d_%d_%s_%s_%d_%d", p.Width, p.Height, p.CropMode, p.CropPos, p.Scale, p.Quality)
}

func isValidCropMode(str string) bool {
	return str == CropModeExact || str == CropModeScale
}

func isValidCropPos(str string) bool {
	return str == CropTopLeft || str == CropTopCenter || str == CropTopRight || str == CropMiddleLeft || str == CropMiddleCenter || str == CropMiddleRight || str == CropBottomLeft || str == CropBottomCenter || str == CropBottomRight
}
