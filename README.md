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

### Use with Server
this example use ace framework 

`github.com/plimble/ace`

```
	a := ace.New()
	a.GET("/img/*imgname", func(c *ace.C) {
		fp.Get(c.Params.ByName("imgname"), c.Writer, c.Request)
	})

	a.Run(":3000")

```

### Example Request Image
get image with original size

```
	http://localhost:3000/img/test.jpg
```

get image with size width 100px , height 100px

```
	http://localhost:3000/img/w_100,h_100/test.jpg
```

get image with size and crop by scale ratio(c_s) image

```
	http://localhost:3000/img/c_s,w_300,h_200/test.jpg
```

get image with size and crop position exact(c_e) middle center(p_mc) from image

```
	http://localhost:3000/img/c_e,p_mc,w_500,h_300/test.jpg
```

###Params Table

| Param | Description                            |
|-------|----------------------------------------|
| c_e   | Crop mode exact                        |
| c_s   | Crop mode scale ratio                  |
| s_2   | use for retina                         |
| q_100 | Quality image maximum 100 (default 80) |
| p_tl  | Crop position top left                 |
| p_tc  | Crop position top center               |
| p_tr  | Crop position top right                |
| p_ml  | Crop position middle left              |
| p_mc  | Crop position middle center            |
| p_mr  | Crop position middle right             |
| p_bl  | Crop position bottom left              |
| p_bc  | Crop position bottom center            |
| p_br  | Crop position bottom right             |



