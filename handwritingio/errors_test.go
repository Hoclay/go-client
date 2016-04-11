package handwritingio

import (
	"net/url"
	"os"
	"testing"
)

func TestAuthError(t *testing.T) {
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

	c.Secret = "blatantlywrong"
	_, err = c.ListHandwritings(DefaultHandwritingListParams)
	if err.Error() != "unauthorized" {
		t.Error(err)
	}

}

func TestValidation(t *testing.T) {
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

	params := DefaultRenderParamsPNG
	params.HandwritingID = "2D5S46A80003" // Perry
	params.Width = "5gophers"
	_, err = c.RenderPNG(params)
	if err.Error() != `width invalid unit: "gophers"` {
		t.Error(err)
	}

}
