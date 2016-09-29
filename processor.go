package ivy

import (
	"github.com/h2non/bimg"
	"github.com/plimble/errors"
)

type Processor interface {
	Process(path string, opt *bimg.Options) ([]byte, bimg.ImageType, error)
}

type processor struct {
	source Source
}

func NewProcessor(source Source) Processor {
	return &processor{source}
}

func (p *processor) Process(path string, opt *bimg.Options) ([]byte, bimg.ImageType, error) {
	data, err := p.source.Get(path)
	if err != nil {
		return nil, 0, err
	}

	result, err := bimg.NewImage(data).Process(*opt)
	imgType := bimg.DetermineImageType(result)

	return result, imgType, errors.InternalServerErrorErr(err)
}
