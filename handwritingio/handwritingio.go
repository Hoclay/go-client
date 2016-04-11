// Package handwritingio provides a client library for interacting with the handwriting.io API
//
// Additional API documentation available at https://handwriting.io/docs/
package handwritingio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Handwriting contains information about a handwriting style
type Handwriting struct {
	ID                   string    `json:"id"`
	Title                string    `json:"title"`
	Created              time.Time `json:"date_created"`
	Modified             time.Time `json:"date_modified"`
	RatingNeatness       int       `json:"rating_neatness"`
	RatingEmbellishment  int       `json:"rating_embellishment"`
	RatingCharacterWidth int       `json:"rating_character_width"`
}

// HandwritingListParams contains the parameters for listing handwritings
type HandwritingListParams struct {
	Offset   int
	Limit    int
	OrderBy  string
	OrderDir string
}

// DefaultHandwritingListParams are the default values for HandwritingListParams
var DefaultHandwritingListParams = HandwritingListParams{
	Offset:   0,
	Limit:    200,
	OrderBy:  "id",
	OrderDir: "asc",
}

// RenderParamsPNG contains the parameters for rendering a PNG image
type RenderParamsPNG struct {
	HandwritingID       string
	Text                string
	HandwritingSize     string
	HandwritingColor    string
	Width               string
	Height              string
	LineSpacing         float64
	LineSpacingVariance float64
	WordSpacingVariance float64
	RandomSeed          int64
}

// DefaultRenderParamsPNG are the default values for RenderParamsPNG
var DefaultRenderParamsPNG = RenderParamsPNG{
	HandwritingID:       "",
	Text:                "",
	HandwritingSize:     "20px",
	HandwritingColor:    "#000000",
	Width:               "504px",
	Height:              "360px",
	LineSpacing:         1.5,
	LineSpacingVariance: 0,
	WordSpacingVariance: 0,
	RandomSeed:          -1,
}

// RenderParamsPDF contains the parameters for rendering a PDF image
type RenderParamsPDF struct {
	HandwritingID       string
	Text                string
	HandwritingSize     string
	HandwritingColor    string
	Width               string
	Height              string
	LineSpacing         float64
	LineSpacingVariance float64
	WordSpacingVariance float64
	RandomSeed          int64
}

// DefaultRenderParamsPDF are the default values for RenderParamsPDF
var DefaultRenderParamsPDF = RenderParamsPDF{
	HandwritingID:       "",
	Text:                "",
	HandwritingSize:     "20pt",
	HandwritingColor:    "(0, 0, 0, 1)",
	Width:               "7in",
	Height:              "5in",
	LineSpacing:         1.5,
	LineSpacingVariance: 0,
	WordSpacingVariance: 0,
	RandomSeed:          -1,
}

// Client is a client for making API calls
type Client struct {
	client *http.Client
	url    *url.URL
}

// NewClient constructs a Client from a URL
func NewClient(u *url.URL) *Client {
	// FIXME move url validation up into the constructor for better error messaging
	client := http.DefaultClient
	c := Client{
		client: client,
		url:    u,
	}
	return &c
}

// ListHandwritings retrieves a list of handwritings
func (c *Client) ListHandwritings(params HandwritingListParams) (handwritings []Handwriting, err error) {
	values := url.Values{}
	values.Add("offset", strconv.Itoa(params.Offset))
	values.Add("limit", strconv.Itoa(params.Limit))
	values.Add("order_by", params.OrderBy)
	values.Add("order_dir", params.OrderDir)
	reqURL := c.url.Scheme + "://" + c.url.Host + "/handwritings?" + values.Encode()
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return
	}

	if c.url.User == nil {
		err = errors.New("token key and secret are required")
		return
	}

	password, ok := c.url.User.Password()
	if !ok {
		err = errors.New("token secret is required")
		return
	}
	req.SetBasicAuth(c.url.User.Username(), password)

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// FIXME
		err = errors.New("NOT IMPLEMENTED")
		return
	}

	err = json.Unmarshal(body, &handwritings)
	return
}

// GetHandwriting retrieves a single of handwriting
func (c *Client) GetHandwriting(id string) (handwriting Handwriting, err error) {
	reqURL := c.url.Scheme + "://" + c.url.Host + "/handwritings/" + id

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return
	}

	if c.url.User == nil {
		err = errors.New("token key and secret are required")
		return
	}

	password, ok := c.url.User.Password()
	if !ok {
		err = errors.New("token secret is required")
		return
	}
	req.SetBasicAuth(c.url.User.Username(), password)

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// FIXME
		err = errors.New("NOT IMPLEMENTED")
		return
	}

	err = json.Unmarshal(body, &handwriting)
	return
}

// RenderPNG calls the API to produce a PNG image
func (c *Client) RenderPNG(params RenderParamsPNG) (r io.ReadCloser, err error) {
	values := url.Values{}
	values.Add("handwriting_id", params.HandwritingID)
	values.Add("text", params.Text)
	values.Add("handwriting_size", params.HandwritingSize)
	values.Add("handwriting_color", params.HandwritingColor)
	values.Add("width", params.Width)
	values.Add("height", params.Height)
	values.Add("line_spacing", strconv.FormatFloat(params.LineSpacing, 'f', -1, 64))
	values.Add("line_spacing_variance", strconv.FormatFloat(params.LineSpacingVariance, 'f', -1, 64))
	values.Add("word_spacing_variance", strconv.FormatFloat(params.WordSpacingVariance, 'f', -1, 64))
	values.Add("random_seed", strconv.FormatInt(params.RandomSeed, 10))

	reqURL := c.url.Scheme + "://" + c.url.Host + "/render/png?" + values.Encode()
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return
	}

	if c.url.User == nil {
		err = errors.New("token key and secret are required")
		return
	}

	password, ok := c.url.User.Password()
	if !ok {
		err = errors.New("token secret is required")
		return
	}
	req.SetBasicAuth(c.url.User.Username(), password)

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println(err)
		bs, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bs))
		err = errors.New("NOT IMPLEMENTED")
		return
	}

	r = resp.Body
	return
}

// RenderPDF calls the API to produce a PDF image
func (c *Client) RenderPDF(params RenderParamsPDF) (r io.ReadCloser, err error) {
	values := url.Values{}
	values.Add("handwriting_id", params.HandwritingID)
	values.Add("text", params.Text)
	values.Add("handwriting_size", params.HandwritingSize)
	values.Add("handwriting_color", params.HandwritingColor)
	values.Add("width", params.Width)
	values.Add("height", params.Height)
	values.Add("line_spacing", strconv.FormatFloat(params.LineSpacing, 'f', -1, 64))
	values.Add("line_spacing_variance", strconv.FormatFloat(params.LineSpacingVariance, 'f', -1, 64))
	values.Add("word_spacing_variance", strconv.FormatFloat(params.WordSpacingVariance, 'f', -1, 64))
	values.Add("random_seed", strconv.FormatInt(params.RandomSeed, 10))

	reqURL := c.url.Scheme + "://" + c.url.Host + "/render/pdf?" + values.Encode()
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return
	}

	if c.url.User == nil {
		err = errors.New("token key and secret are required")
		return
	}

	password, ok := c.url.User.Password()
	if !ok {
		err = errors.New("token secret is required")
		return
	}
	req.SetBasicAuth(c.url.User.Username(), password)

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println(err)
		bs, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bs))
		err = errors.New("NOT IMPLEMENTED")
		return
	}

	r = resp.Body
	return
}
