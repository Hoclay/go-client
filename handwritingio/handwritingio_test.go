package handwritingio

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func ExampleClient_ListHandwritings() {

	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		fmt.Println("Make sure you have your environment variable HANDWRITINGIO_API_URL set correctly")
		fmt.Println(err)
		return
	}

	var params = DefaultHandwritingListParams
	params.Limit = 5
	c, err := NewClient(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	handwritings, err := c.ListHandwritings(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ListHandwritings returned %d Handwritings", len(handwritings))

	// Output:
	// ListHandwritings returned 5 Handwritings
}

func ExampleClient_GetHandwriting() {
	id := "2D5S46A80003"

	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		fmt.Println("Make sure you have your environment variable HANDWRITINGIO_API_URL set correctly")
		fmt.Println(err)
		return
	}

	c, err := NewClient(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	handwriting, err := c.GetHandwriting(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(handwriting.Title)

	// Output:
	// Perry
}

func TestClient_GetHandwriting(t *testing.T) {
	id := "2D5S46A80003"

	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		t.Error(err)
		return
	}

	c, err := NewClient(u)
	if err != nil {
		t.Error(err)
		return
	}

	handwriting, err := c.GetHandwriting(id)
	if err != nil {
		t.Error(err)
		return
	}

	if handwriting.Title != "Perry" {
		t.Fail()
	}

	if handwriting.ID != id {
		t.Fail()
	}

	// Ratings default to 1400, and go up or down relative to other handwritings
	// zero values would indicate deserialization problems
	if handwriting.RatingNeatness == 0 {
		t.Fail()
	}
	if handwriting.RatingCursivity == 0 {
		t.Fail()
	}
	if handwriting.RatingEmbellishment == 0 {
		t.Fail()
	}
	if handwriting.RatingCharacterWidth == 0 {
		t.Fail()
	}
}

func ExampleClient_RenderPNG() {
	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		fmt.Println("Make sure you have your environment variable HANDWRITINGIO_API_URL set correctly")
		fmt.Println(err)
		return
	}

	c, err := NewClient(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	var params = DefaultRenderParamsPNG
	// Perry
	params.HandwritingID = "2D5S46A80003"
	// https://groups.google.com/forum/#!topic/comp.os.plan9/VUUznNK2t4Q%5B151-175%5D
	params.Text = "object-oriented design is the roman numerals of computing.\n\n - Rob Pike"
	r, err := c.RenderPNG(params)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()

	// os.TempDir can be overridden with TMP_DIR environment variable
	filename := filepath.Join(os.TempDir(), "handwriting_io_render.png")
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("handwriting_io_render.png written to temporary directory")

	// Output:
	// handwriting_io_render.png written to temporary directory
}

func ExampleClient_RenderPDF() {
	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		fmt.Println("Make sure you have your environment variable HANDWRITINGIO_API_URL set correctly")
		fmt.Println(err)
		return
	}

	c, err := NewClient(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	var params = DefaultRenderParamsPDF
	// Perry
	params.HandwritingID = "2D5S46A80003"
	// https://groups.google.com/forum/#!topic/comp.os.plan9/VUUznNK2t4Q%5B151-175%5D
	params.Text = "object-oriented design is the roman numerals of computing.\n\n - Rob Pike"
	r, err := c.RenderPDF(params)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()

	// os.TempDir can be overridden with TMP_DIR environment variable
	filename := filepath.Join(os.TempDir(), "handwriting_io_render.pdf")
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("handwriting_io_render.pdf written to temporary directory")

	// Output:
	// handwriting_io_render.pdf written to temporary directory
}
