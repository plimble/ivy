package ivy

import (
	"github.com/h2non/bimg"
	"github.com/kataras/iris"
	"github.com/plimble/errors"
)

func NewServer(config *Config) (*iris.Application, error) {
	app := iris.New()

	source := NewS3(config.SourceAwsId, config.SourceAwsSecret, config.SourceAwsS3Bucket, config.SourceAwsS3Region)

	p := NewProcessor(source)

	app.Get("/{filepath:path}", handler(p))

	returnappir, nil
}

func handler(p Processor) iris.Handler {
	return func(ctx iris.Context) {
		path := ctx.Params().Get("filepath")
		if path == "" {
			ctx.StatusCode(iris.StatusNotFound)
			return
		}

		opt := &bimg.Options{}
		opt.Width, opt.Height = splitWidthHeightString(ctx.FormValue("r"))
		opt.Crop = getBoolFromString(ctx.FormValue("c"))

		g := getGravityFromString(ctx.FormValue("g"))
		if g > -1 {
			opt.Gravity = g
		}

		f := getBoolFromString(ctx.FormValue("f"))
		if f {
			opt.Force = true
		} else {
			opt.Embed = true
		}

		opt.Quality = stringToInt(ctx.FormValue("q"))
		opt.Type = stringToImageType(ctx.FormValue("t"))

		img, imgType, err := p.Process(path, opt)
		if err != nil {
			status, err := errors.ErrorStatus(err)
			ctx.StatusCode(status)
			ctx.Text(err.Error())
			return
		}

		ctx.ContentType(getContentType(imgType))
		ctx.Header("Content-Length", len(img))
		ctx.ResponseWriter().Write(img)
	}
}
