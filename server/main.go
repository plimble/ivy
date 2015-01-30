package main

import (
	"github.com/bradhe/stopwatch"
	"github.com/codegangsta/cli"
	"github.com/plimble/ace"
	"github.com/plimble/fileproxy"
	"log"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "fileproxy"
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
		var source fileproxy.Source
		switch c.String("source") {
		case "s3":
			source = fileproxy.NewS3Source(c.String("s3key"), c.String("s3secret"))
		default:
			source = fileproxy.NewFileSystemSource(c.String("sourceroot"))
		}

		var cache fileproxy.Cache
		switch c.String("cache") {
		case "file":
			cache = fileproxy.NewFileSystemCache("cacheroot")
		default:
			cache = nil
		}

		fp := fileproxy.New(source, cache, &fileproxy.Config{
			HttpCache:     int64(c.Int("httpcache")),
			IsDevelopment: false,
		})

		a := ace.New()

		a.GET("/:bucket/:params/*path", func(c *ace.C) {
			bucket := c.Params.ByName("bucket")
			path := c.Params.ByName("path")
			params := c.Params.ByName("params")
			start := stopwatch.Start()
			fp.Get(bucket, params, path, c.Writer, c.Request)
			watch := stopwatch.Stop(start)
			log.Printf("[fileproxy] %d %s %s %s %v ms", c.Writer.Status(), bucket, params, path, watch.Milliseconds())
		})

		url := c.String("url")
		log.Println("[fileproxy] Start server on " + url)
		a.Run(url)
	}

	app.Run(os.Args)
}
