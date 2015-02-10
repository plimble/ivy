package main

import (
	"github.com/bradhe/stopwatch"
	"github.com/codegangsta/cli"
	"github.com/plimble/ace"
	"github.com/plimble/ivy"
	"log"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "Ivy"
	app.Usage = "make an explosive entrance"
	app.Author = "Witoo Harianto"
	app.Email = "witooh@icloud.com"
	app.Version = "1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url,u",
			Value: ":4900",
			Usage: "server port",
		},
		cli.IntFlag{
			Name:  "httpcache,t",
			Value: 0,
			Usage: "if cache enable this is http cache in second",
		},
		cli.StringFlag{
			Name:  "cache,c",
			Value: "",
			Usage: "enable cache, specific  eg, file",
		},
		cli.StringFlag{
			Name:  "source,s",
			Value: "file",
			Usage: "source of image eg, file, s3",
		},
		cli.StringFlag{
			Name:   "s3key",
			Value:  "",
			Usage:  "if source is s3, AWS S3 access key",
			EnvVar: "AWS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "s3secret",
			Value:  "",
			Usage:  "if source is s3, AWS S3 secret key",
			EnvVar: "AWS_SECRET_KEY",
		},
		cli.StringFlag{
			Name:  "sourceroot",
			Value: "",
			Usage: "if source is file, specific root path of image",
		},

		cli.StringFlag{
			Name:  "cacheroot",
			Value: "",
			Usage: "if cache is file, specific root path of cache",
		},
	}
	app.Action = func(c *cli.Context) {
		var source ivy.Source
		switch c.String("source") {
		case "s3":
			source = ivy.NewS3Source(c.String("s3key"), c.String("s3secret"))
		default:
			source = ivy.NewFileSystemSource(c.String("sourceroot"))
		}

		var cache ivy.Cache
		switch c.String("cache") {
		case "file":
			cache = ivy.NewFileSystemCache("cacheroot")
		default:
			cache = nil
		}

		iv := ivy.New(source, cache, ivy.NewGMProcessor(), &ivy.Config{
			HTTPCache:     int64(c.Int("httpcache")),
			IsDevelopment: false,
		})

		a := ace.New()

		a.GET("/:bucket/:params/*path", func(c *ace.C) {
			start := stopwatch.Start()
			iv.Get(c.Params.ByName("bucket"), c.Params.ByName("params"), c.Params.ByName("path"), c.Writer, c.Request)
			watch := stopwatch.Stop(start)
			log.Printf("[Ivy] %d %s %vms", c.Writer.Status(), c.Request.URL.String(), watch.Milliseconds())
		})

		url := c.String("url")
		log.Println("[Ivy] Start server on " + url)
		a.Run(url)
	}

	app.Run(os.Args)
}
