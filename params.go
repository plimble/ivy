package ivy

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	paramResize  = "r"
	paramCrop    = "c"
	paramGravity = "g"
	paramQuality = "q"

	cropNorthWest = "nw"
	cropNorth     = "n"
	cropNorthEast = "ne"
	cropWest      = "w"
	cropCenter    = "c"
	cropEast      = "e"
	cropSouthWest = "sw"
	cropSouth     = "s"
	cropSouthEast = "se"
)

//Params for set up image
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

//NewParams create default params
func NewParams() *Params {
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

//ParseParams parse params string into params struct
func ParseParams(paramsStr string) (*Params, error) {
	params := NewParams()
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
	case paramResize:
		if params.Width, params.Height, err = getParamDimentsion(key, value, 0); err != nil {
			return err
		}
		params.EnableResize = true
	case paramCrop:
		if params.CropWidth, params.CropHeight, err = getParamDimentsion(key, value, 1); err != nil {
			return err
		}
		params.EnableCrop = true
	case paramGravity:
		value = strings.ToLower(value)
		if !isValidGravity(value) {
			return fmt.Errorf("invalid value for %s", key)
		}
		params.Gravity = value
		params.EnableGravity = true
	case paramQuality:
		params.Quality, err = strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("could not parse value for parameter: %s", key)
		}
		if params.Quality < 0 || params.Quality > 100 {
			return fmt.Errorf("value %s must be > 0 & <= 100: %s", value, key)
		}
	default:
		return fmt.Errorf("invalid parameter: %s_%s", key, value)
	}

	return nil
}

func (p *Params) String() string {
	if p.str == "" {
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
		return 0, 0, fmt.Errorf("value %s must be > 0: %s", value, key)
	}

	height, err := strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse value for parameter: %s", key)
	}
	if height < min {
		return 0, 0, fmt.Errorf("value %s must be > 0: %s", value, key)
	}

	return width, height, nil
}

func isValidGravity(str string) bool {
	return str == cropNorthWest || str == cropNorth || str == cropNorthEast || str == cropWest || str == cropCenter || str == cropEast || str == cropSouthWest || str == cropSouth || str == cropSouthEast
}
