package fileproxy

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ParamResize  = "r"
	ParamCrop    = "c"
	ParamGravity = "g"
	ParamQuality = "q"

	CropNorthWest = "nw"
	CropNorth     = "n"
	CropNorthEast = "ne"
	CropWest      = "w"
	CropCenter    = "c"
	CropEast      = "e"
	CropSouthWest = "sw"
	CropSouth     = "s"
	CropSouthEast = "se"
)

type Params struct {
	Width         int
	Height        int
	CropWidth     int
	CropHeight    int
	Gravity       string
	Quality       int
	EnableCrop    bool
	EnableResize  bool
	EnableGravity bool
	IsDefault     bool
	str           string
}

func newParams() *Params {
	return &Params{
		Width:         0,
		Height:        0,
		CropWidth:     0,
		CropHeight:    0,
		Gravity:       "",
		Quality:       -1,
		EnableCrop:    false,
		EnableResize:  false,
		EnableGravity: false,
		IsDefault:     true,
		str:           "",
	}
}

func parseParams(paramsStr string) (*Params, error) {
	params := newParams()
	if paramsStr == "" || paramsStr == "_" || paramsStr == "0" {
		return params, nil
	}

	var err error

	parts := strings.Split(paramsStr, ",")
	for _, part := range parts {
		if len(part) < 3 {
			return nil, fmt.Errorf("invalid parameter: %s", part)
		}
		if part[1] != 95 {
			return nil, fmt.Errorf("invalid parameter: %s", part)
		}
		key := string(part[0])
		value := string(part[2:])

		switch key {
		case ParamResize:
			params.Width, params.Height, err = getParamDimentsion(key, value, 0)
			if err != nil {
				return nil, err
			}
			params.EnableResize = true
		case ParamCrop:
			params.CropWidth, params.CropHeight, err = getParamDimentsion(key, value, 1)
			if err != nil {
				return nil, err
			}
			params.EnableCrop = true
		case ParamGravity:
			value = strings.ToLower(value)
			if !isValidGravity(value) {
				return nil, fmt.Errorf("invalid value for %s", key)
			}
			params.Gravity = value
			params.EnableGravity = true
		// case ParamScale:
		// 	vint, err = strconv.Atoi(value)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("could not parse value for parameter: %s", key)
		// 	}
		// 	if vint < 1 {
		// 		return nil, fmt.Errorf("value %d must be > 0: %s", value, key)
		// 	}
		// 	params.Scale = vint
		case ParamQuality:
			params.Quality, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse value for parameter: %s", key)
			}
			if params.Quality < 0 || params.Quality > 100 {
				return nil, fmt.Errorf("value %d must be > 0 & <= 100: %s", value, key)
			}
		default:
			return nil, fmt.Errorf("invalid parameter: %s", part)
		}
	}
	params.IsDefault = false

	return params, nil
}

func (p *Params) String() string {
	if p.str != "" {
		p.str = fmt.Sprintf("%d_%d_%d_%d_%d", p.Width, p.Height, p.CropWidth, p.CropHeight, p.Quality)
	}
	return p.str
}

func getParamDimentsion(key, value string, min int) (int, int, error) {
	values := strings.Split(value, "x")
	if len(values) != 2 {
		return 0, 0, fmt.Errorf("invalid value for %s", key)
	}

	width, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse value for parameter: %s", key)
	}
	if width < min {
		return 0, 0, fmt.Errorf("value %d must be > 0: %s", value, key)
	}

	height, err := strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse value for parameter: %s", key)
	}
	if height < min {
		return 0, 0, fmt.Errorf("value %d must be > 0: %s", value, key)
	}

	return width, height, nil
}

func isValidGravity(str string) bool {
	return str == CropNorthWest || str == CropNorth || str == CropNorthEast || str == CropWest || str == CropCenter || str == CropEast || str == CropSouthWest || str == CropSouth || str == CropSouthEast
}
