# Handwriting.io Client

[![GoDoc](https://godoc.org/github.com/handwritingio/go-client/handwritingio?status.svg)](http://godoc.org/github.com/handwritingio/go-client/handwritingio)

## Installation

    go get github.com/handwritingio/go-client/handwritingio

## Basic Example

Set up the client, render an image, and write it to a file:

```go
package main

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/handwritingio/go-client/handwritingio"
)

func main() {
	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		fmt.Println("Make sure you have your environment variable HANDWRITINGIO_API_URL set correctly")
		fmt.Println(err)
		return
	}

	c, err := handwritingio.NewClient(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	var params = handwritingio.DefaultRenderParamsPNG
	params.HandwritingID = "31SB3NWR00E0" // found in our catalog or by listing handwritings
	params.Text = "Handwriting with Go!"
	params.Height = "auto"
	r, err := c.RenderPNG(params)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()

	f, err := os.Create("handwriting.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		fmt.Println(err)
	}

	return
}
```

If all goes well, this should create an image similar to the following:

![image](https://s3.amazonaws.com/hwio-cdn-production/go-client/handwriting.png)

## Advanced Examples

Overlaying handwriting on a background image:
```go
package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/url"
	"os"

	"github.com/handwritingio/go-client/handwritingio"
)

func main() {

	// Original https://golang.org/doc/gopher/fiveyears.jpg
	// By Renee French, used under Creative Commons license
	// More gophers at https://blog.golang.org/gopher
	f, err := os.Open("fiveyears.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// draw.Draw doesn't work with JPEG image as destination
	// So we're creating a new RGBA to hold the result
	b := src.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)

	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		fmt.Println("Make sure you have your environment variable HANDWRITINGIO_API_URL set correctly")
		log.Fatal(err)
		return
	}

	c, err := handwritingio.NewClient(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	var params = handwritingio.DefaultRenderParamsPNG
	params.HandwritingID = "31SB3NWR00E0" // found in our catalog or by listing handwritings
	params.Text = "Handwriting with Go!"
	params.Height = "auto"
	r, err := c.RenderPNG(params)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()

	src, _, err = image.Decode(r)
	if err != nil {
		log.Fatal(err)
	}

	sr := src.Bounds().Add(image.Pt(476, 10))
	draw.Draw(m, sr, src, image.ZP, draw.Over)

	f, err = os.Create("handwriting_overlay.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = png.Encode(f, m)
	if err != nil {
		fmt.Println(err)
	}

	return
}
```

It should create an image like this : [handwriting_overlay.png](https://s3.amazonaws.com/hwio-cdn-production/go-client/handwriting_overlay.png)

## Reference

See the [API Documentation](https://www.handwriting.io/docs) for details on all endpoints and parameters. For the most part, the Client passes parameters through to the API directly.

The endpoints map to client methods as follows:

- [GET /handwritings](https://handwriting.io/docs/#get-handwritings) -> `client.ListHandwritings(params)`
- [GET /handwritings/{id}](https://handwriting.io/docs/#get-handwritings--id-) -> `client.GetHandwriting(handwriting_id)`
- [GET /render/png](https://handwriting.io/docs/#get-render-png) -> `client.RenderPNG(params)`
- [GET /render/pdf](https://handwriting.io/docs/#get-render-pdf) -> `client.RenderPDF(params)`

## Version Numbers

Version numbers for this package work slightly differently than standard
[semantic versioning](http://semver.org/). For this package, the `major`
version number will match the Handwriting.io API version number, and the
`minor` version will be  incremented for any breaking changes to this package.
The `patch` version will be incremented for bug fixes and changes that add
functionality only.

Tagged releases are compatible with [gopkg.in](http://labix.org/gopkg.in) versioning.

```
go get gopkg.in/handwritingio/go-client.v1/handwritingio
```

## Issues

Please open an issue on [Github](https://github.com/handwritingio/go-client/issues)
or [contact us](https://handwriting.io/contact) directly for help with any
problems you find.
