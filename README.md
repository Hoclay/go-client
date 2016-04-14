# Handwriting.io Client

[![GoDoc](https://godoc.org/github.com/handwritingio/go-client/handwritingio?status.svg)](http://godoc.org/github.com/handwritingio/go-client/handwritingio)

## Installation

    go get github.com/handwritingio/go-client/handwritingio

## Basic Example

Set up the client, render an image, and write it to a file:

```golang
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

## Issues

Please open an issue on [Github](https://github.com/handwritingio/go-client/issues)
or [contact us](https://handwriting.io/contact) directly for help with any
problems you find.
