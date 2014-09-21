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

var proxy *fileproxy.Proxy

func upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	storage := c.Params.ByName("storage")
	path := c.Params.ByName("path")

	filename := path + "/" + header.Filename

	_, err = proxy.Save(storage, filename, file)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, "")
}

func load(c *gin.Context) {
	storage := c.Params.ByName("storage")
	path := c.Params.ByName("path")
	paramsStr := c.Params.ByName("params")

	ext := filepath.Ext(path)

	data, err := proxy.Load(storage, paramsStr, path, ext)
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

func remove(c *gin.Context) {
	storage := c.Params.ByName("storage")
	path := c.Params.ByName("path")

	err := proxy.Delete(storage, path)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.String(200, "")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//read env
	filePath := os.Getenv("FILEPROXY_FILE_PATH")
	IP := os.Getenv("FILEPROXY_IP")

	proxy = fileproxy.New(fileproxy.Config{
		FileStoragePath: filePath,
	})

	g := gin.New()
	g.Use(gin.Recovery())

	g.POST("/upload/:storage/*path", upload)
	g.GET("/:storage/:params/*path", load)
	g.DELETE("/:storage/*path", remove)

	if IP == "" {
		IP = ":3001"
	}
	g.Run(IP)
}
