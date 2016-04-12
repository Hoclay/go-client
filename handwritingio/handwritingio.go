// Package handwritingio provides a client library for interacting with the handwriting.io API
//
// Additional API documentation available at https://handwriting.io/docs/
package handwritingio

import (
	"encoding/json"
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
	RatingCursivity      int       `json:"rating_cursivity"`
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
	Key    string
	Secret string
}

// NewClient constructs a Client from a URL
func NewClient(u *url.URL) (*Client, error) {

	if u.User == nil {
		return nil, TokenError("token key and secret are required")
	}

	password, ok := u.User.Password()
	if !ok {
		return nil, TokenError("token secret is required")
	}

	c := Client{
		client: http.DefaultClient,
		url:    u,
		Key:    u.User.Username(),
		Secret: password,
	}
	return &c, nil
}

// ListHandwritings retrieves a list of handwritings
func (c *Client) ListHandwritings(params HandwritingListParams) (handwritings []Handwriting, err error) {
	values := url.Values{}
	values.Add("offset", strconv.Itoa(params.Offset))
	values.Add("limit", strconv.Itoa(params.Limit))
	values.Add("order_by", params.OrderBy)
	values.Add("order_dir", params.OrderDir)

	resp, err := c.get("/handwritings", values)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	err = json.Unmarshal(body, &handwritings)
	return
}

// GetHandwriting retrieves a single of handwriting
func (c *Client) GetHandwriting(id string) (handwriting Handwriting, err error) {
	resp, err := c.get("/handwritings/"+id, nil)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

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

	resp, err := c.get("/render/png", values)
	if err != nil {
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

	resp, err := c.get("/render/pdf", values)
	if err != nil {
		return
	}

	r = resp.Body
	return
}

func (c *Client) get(path string, values url.Values) (resp *http.Response, err error) {
	reqURL := c.url.Scheme + "://" + c.url.Host + path
	if values != nil && len(values) > 0 {
		reqURL += "?" + values.Encode()
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(c.Key, c.Secret)

	resp, err = c.client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = responseError(resp)
	}

	return
}
