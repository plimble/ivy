fileproxy [![godoc badge](http://godoc.org/github.com/plimble/fileproxy?status.png)](http://godoc.org/github.com/plimble/fileproxy)   [![gocover badge](http://gocover.io/_badge/github.com/plimble/fileproxy?t=1)](http://gocover.io/github.com/plimble/fileproxy) [![Build Status](https://api.travis-ci.org/plimble/fileproxy.svg?branch=master&t=1)](https://travis-ci.org/plimble/fileproxy) [![Go Report Card](http://goreportcard.com/badge/plimble/fileproxy?t=1)](http:/goreportcard.com/report/plimble/fileproxy)
=========

Assets & Image processing on the fly

### Installation
`go get -u github.com/plimble/ace`

### Documentation
 - [GoDoc](http://godoc.org/github.com/plimble/fileproxy)

### Sources

##### File System

```
	source := fileproxy.NewFileSystemSource("/path/to/asset")
	cache := fileproxy.NewFileSystemCache("/path/to/cache")

	config := &fileproxy.Config{
		IsDevelopment: false,
		HttpCache:     66000,
	}

	fp := fileproxy.New(source, cache, config)
```

##### AWS S3

```
	source := fileproxy.NewS3Source("accessKey", "secretKey")
	cache := fileproxy.NewFileSystemCache("/path/to/cache")

	config := &fileproxy.Config{
		IsDevelopment: false,
		HttpCache:     66000,
	}

	fp := fileproxy.New(source, cache, config)
```

### Cache

You can use file system for caching or set `nil` for CDN like Cloudfront or disable caching


### Server

You can run built-in server or implement in your server

For more config

```
./server -h
```

##### Ace Example

```
	a := ace.New()
	a.GET("/:bucket/:params/*path", func(c *ace.C) {
		fp.Get(
			c.Params.ByName("bucket"),
			c.Params.ByName("params"),
			c.Params.ByName("path"),
			c.Writer,
			c.Request,
		)
	})

	a.Run(":3000")
```

You can run as middleware also

### Example Request Image
Get image with original size or non image

```
	http://localhost:3000/bucket/_/test.jpg
```

```
	http://localhost:3000/bucket/0/test.jpg
```

Resize 100x100

```
	http://localhost:3000/bucket/r_100x100/test.jpg
```

Resize width 100px aspect ratio

```
	http://localhost:3000/bucket/r_100x0/test.jpg
```

Crop top left image 200x200

```
	http://localhost:3000/bucket/c_200x200/test.jpg
```

Crop with gravity NorthWest image 200x200

```
	http://localhost:3000/bucket/c_200x200,g_nw/test.jpg
```

Resize 400x400 then crop 200x200 and gravity center

```
	http://localhost:3000/bucket/r_400x400,c_200x200,g_c/test.jpg
```

Quality 100

```
	http://localhost:3000/bucket/q_100/test.jpg
```

###Params Table

| Param               | Description                            |
|---------------------|----------------------------------------|
| r_{width}x{height}  | Resize image, if 0 is aspect ratio     |
| c_{width}x{height}  | Crop image                             |
| g_{direction}       | Gravity image                          |
| q_{quality}         | Quality image maximum 100              |

###Gravity position

| Param | Description                            |
|-------|----------------------------------------|
| nw    | North West                             |
| n     | North                                  |
| ne    | North East                             |
| w     | West                                   |
| c     | Center                                 |
| e     | East                                   |
| sw    | South West                             |
| s     | South                                  |
| se    | South East                             |


###Contributing

If you'd like to help out with the project. You can put up a Pull Request.

