package ivy

import (
	"github.com/h2non/bimg"
	"github.com/kataras/iris"
	"github.com/plimble/errors"
)

func NewServer(config *Config) (*iris.Framework, error) {
	ir := iris.New(
		iris.OptionDisableBanner(true),
	)

	source := NewS3(config.SourceAwsId, config.SourceAwsSecret, config.SourceAwsS3Bucket, config.SourceAwsS3Region)

	p := NewProcessor(source)

	ir.Get("/*path", handler(p))

	return ir, nil
}

func handler(p Processor) iris.HandlerFunc {
	return func(ctx *iris.Context) {
		path := ctx.Param("path")
		if path == "" {
			ctx.SetStatusCode(404)
			return
		}

		opt := &bimg.Options{}
		opt.Width, opt.Height = splitWidthHeightString(ctx.FormValueString("r"))
		opt.Crop = getBoolFromString(ctx.FormValueString("c"))

		g := getGravityFromString(ctx.FormValueString("g"))
		if g > -1 {
			opt.Gravity = g
		}

		f := getBoolFromString(ctx.FormValueString("f"))
		if f {
			opt.Force = true
		} else {
			opt.Embed = true
		}

		opt.Quality = stringToInt(ctx.FormValueString("q"))
		opt.Type = stringToImageType(ctx.FormValueString("t"))

		img, imgType, err := p.Process(path, opt)
		if err != nil {
			status, err := errors.ErrorStatus(err)
			ctx.Text(status, err.Error())
			return
		}

		ctx.Response.Header.SetContentType(getContentType(imgType))
		ctx.Response.Header.SetContentLength(len(img))
		ctx.Response.AppendBody(img)
		ctx.SetStatusCode(200)
	}
}
