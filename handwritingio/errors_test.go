package handwritingio

import (
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestAuthError(t *testing.T) {
	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		t.Error(err)
		return
	}

	c, err := NewClientURL(u)
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

func TestMultipleErrors(t *testing.T) {
	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		t.Error(err)
		return
	}

	c, err := NewClientURL(u)
	if err != nil {
		t.Error(err)
		return
	}

	params := DefaultRenderParamsPNG
	params.HandwritingID = "2D5S46A80003" // Perry
	params.Width = "5gophers"
	params.Height = "80%"
	params.HandwritingSize = "-5px"
	_, err = c.RenderPNG(params)

	if es, ok := err.(APIErrors); ok {
		fields := []string{}
		for _, e := range es.Errors {
			fields = append(fields, e.Field)
		}
		sort.Strings(fields)
		expected := "handwriting_size height width"
		actual := strings.Join(fields, " ")
		if expected != actual {
			t.Logf("expected fields: %#v", expected)
			t.Logf("actual fields: %#v", actual)
			t.Fail()
		}
	} else {
		t.Error("error returned was not APIErrors")
	}

	if err.Error() != "multiple errors" {
		t.Fail()
	}

}
