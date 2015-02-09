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

func isEmptyParam(paramsStr string) bool {
	return paramsStr == "" || paramsStr == "_" || paramsStr == "0"
}

func parseParams(paramsStr string) (*Params, error) {
	params := newParams()
	if isEmptyParam(paramsStr) {
		return params, nil
	}

	parts := strings.Split(paramsStr, ",")
	for _, part := range parts {
		if len(part) < 3 || part[1] != 95 {
			return nil, fmt.Errorf("invalid parameter: %s", part)
		}

		key := part[:1]
		value := part[2:]

		if err := setParams(params, key, value); err != nil {
			return nil, err
		}
	}
	params.IsDefault = false

	return params, nil
}

func setParams(params *Params, key, value string) error {
	var err error

	switch key {
	case ParamResize:
		if params.Width, params.Height, err = getParamDimentsion(key, value, 0); err != nil {
			return err
		}
		params.EnableResize = true
	case ParamCrop:
		if params.CropWidth, params.CropHeight, err = getParamDimentsion(key, value, 1); err != nil {
			return err
		}
		params.EnableCrop = true
	case ParamGravity:
		value = strings.ToLower(value)
		if !isValidGravity(value) {
			return fmt.Errorf("invalid value for %s", key)
		}
		params.Gravity = value
		params.EnableGravity = true
	case ParamQuality:
		params.Quality, err = strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("could not parse value for parameter: %s", key)
		}
		if params.Quality < 0 || params.Quality > 100 {
			return fmt.Errorf("value %d must be > 0 & <= 100: %s", value, key)
		}
	default:
		return fmt.Errorf("invalid parameter: %s_%s", key, value)
	}

	return nil
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
