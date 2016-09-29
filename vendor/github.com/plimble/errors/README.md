errors [![godoc badge](http://godoc.org/github.com/plimble/errors?status.png)](http://godoc.org/github.com/plimble/errors)
========


## Installation

```go
$ go get -u github.com/plimble/errors
```

## Usage

```go

func main(){

  err1 = errors.New("error message")
  err2 = errors.Newf("%s", "error message")
  err2 = errors.NewCode(404,"error message")
  err2 = errors.NewCodef(404, "%s", "error message")
  err3 = errors.InternalError("error message")
}


```