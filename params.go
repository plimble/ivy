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

type params struct {
	width         int
	height        int
	cropWidth     int
	cropHeight    int
	gravity       string
	quality       int
	enableCrop    bool
	enableResize  bool
	enableGravity bool
	isDefault     bool
	str           string
}

func newParams() *params {
	return &params{
		width:         0,
		height:        0,
		cropWidth:     0,
		cropHeight:    0,
		gravity:       "",
		quality:       -1,
		enableCrop:    false,
		enableResize:  false,
		enableGravity: false,
		isDefault:     true,
		str:           "",
	}
}

func isEmptyParam(paramsStr string) bool {
	return paramsStr == "" || paramsStr == "_" || paramsStr == "0"
}

func parseParams(paramsStr string) (*params, error) {
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
	params.isDefault = false

	return params, nil
}

func setParams(params *params, key, value string) error {
	var err error

	switch key {
	case paramResize:
		if params.width, params.height, err = getParamDimentsion(key, value, 0); err != nil {
			return err
		}
		params.enableResize = true
	case paramCrop:
		if params.cropWidth, params.cropHeight, err = getParamDimentsion(key, value, 1); err != nil {
			return err
		}
		params.enableCrop = true
	case paramGravity:
		value = strings.ToLower(value)
		if !isValidGravity(value) {
			return fmt.Errorf("invalid value for %s", key)
		}
		params.gravity = value
		params.enableGravity = true
	case paramQuality:
		params.quality, err = strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("could not parse value for parameter: %s", key)
		}
		if params.quality < 0 || params.quality > 100 {
			return fmt.Errorf("value %s must be > 0 & <= 100: %s", value, key)
		}
	default:
		return fmt.Errorf("invalid parameter: %s_%s", key, value)
	}

	return nil
}

func (p *params) String() string {
	if p.str == "" {
		p.str = fmt.Sprintf("%d_%d_%d_%d_%d", p.width, p.height, p.cropWidth, p.cropHeight, p.quality)
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
