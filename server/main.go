package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/plimble/fileproxy"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func upload(services *services, c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	storage := c.Params.ByName("storage")
	path := c.Params.ByName("path")

	filename := path + "/" + header.Filename

	_, err = services.proxy.Save(storage, filename, file)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, "")
}

func load(services *services, c *gin.Context) {
	storage := c.Params.ByName("storage")
	path := c.Params.ByName("path")
	paramsStr := c.Params.ByName("params")

	ext := filepath.Ext(path)

	data, err := services.proxy.Load(storage, paramsStr, path, ext)
	if fileproxy.IsNotFound(err) {
		c.String(404, err.Error())
		return
	}
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Writer.Header().Set("Content-Type", mime.TypeByExtension(ext))
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	c.Writer.Header().Set("Cache-Control", "no-transform,public,max-age=14286411")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Last-Modified", "Tue, 09 Sep 1980 07:05:01 GMT")
	c.Writer.Header().Set("Vary", "Accept-Encoding")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(data)
}

func remove(services *services, c *gin.Context) {
	storage := c.Params.ByName("storage")
	path := c.Params.ByName("path")

	err := services.proxy.Delete(storage, path)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.String(200, "")
}

type appHandleFunc func(services *services, c *gin.Context)
type services struct {
	proxy *fileproxy.Proxy
}

func appHandle(services *services, handle appHandleFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handle(services, c)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//read env
	filePath := os.Getenv("CAYL_FILE_PATH")
	IP := os.Getenv("CAYL_IP")

	s := &services{
		proxy: fileproxy.New(fileproxy.Config{
			FileStoragePath: filePath,
		}),
	}

	g := gin.New()
	g.Use(gin.Recovery())

	g.POST("/upload/:storage/*path", appHandle(s, upload))
	g.GET("/:storage/:params/*path", appHandle(s, load))
	g.DELETE("/:storage/*path", appHandle(s, remove))

	if IP == "" {
		IP = ":3001"
	}
	g.Run(IP)
}
