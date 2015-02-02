fileproxy
=========

file and image on the fly proxy

### Config fileproxy
```
    fsource := fileproxy.NewFileSystemSource("sourcefolder")
	csource := fileproxy.NewFileSystemCache("cachefolder")

	fconfig := &fileproxy.Config{
		IsDevelopment: false,
		HttpCache:     66000,
	}

	fp := fileproxy.New(
		fsource,
		csource,
		fconfig,
	)
```

### Stand alone server

```
./server -h
```

### Use with Server
this example use ace framework

`github.com/plimble/ace`

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

### Example Request Image
get image with original size

```
	http://localhost:3000/bucket/_/test.jpg
```

```
	http://localhost:3000/bucket/0/test.jpg
```

resize 100x100

```
	http://localhost:3000/bucket/r_100x100/test.jpg
```

resize width 100px aspect ratio

```
	http://localhost:3000/bucket/r_100x0/test.jpg
```

crop top left image 200x200

```
	http://localhost:3000/bucket/c_200x200/test.jpg
```

crop with gravity NorthWest image 200x200

```
	http://localhost:3000/bucket/c_200x200,g_nw/test.jpg
```

resize 400x400 then crop 200x200 and gravity center

```
	http://localhost:3000/bucket/r_400x400,c_200x200,g_c/test.jpg
```

quality 100

```
	http://localhost:3000/bucket/q_100/test.jpg
```

###Params Table

| Param | Description                            |
|-------|----------------------------------------|
| r_{width}x{height}   | Resize image, if 0 is aspect ratio                        |
| c_{width}x{height}   | Crop image                  |
| g_{direction}   | Gravity image                         |
| q_{quality} | Quality image maximum 100 |

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

