// Package handwritingio provides a client library for interacting with the handwriting.io API
package handwritingio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Handwriting contains information about a handwriting style
type Handwriting struct {
	// ID is an unique identifier for the handwriting style
	ID string `json:"id"`
	// Title is human friendly name for the handwriting
	Title string `json:"title"`
	// Created is when the handwriting was created
	Created time.Time `json:"date_created"`
	// Modified is when then handwriting was last modified
	Modified time.Time `json:"date_modified"`
	// RatingNeatness is a rating of how "neat" a handwriting is relative to others
	RatingNeatness int `json:"rating_neatness"`
	// RatingEmbellishment is a rating of how "embellished" a handwriting is relative to others
	RatingEmbellishment int `json:"rating_embellishment"`
	// RatingCharacterWidth is a rating of how wide a handwriting is relative to others
	RatingCharacterWidth int `json:"rating_character_width"`
}

// HandwritingListParams contains the parameters for listing handwritings
type HandwritingListParams struct {
	// Offset is the number of handwritings to skip at the beginning of the list.  Useful for pagination.
	Offset int
	// Limit is the maximum number of handwritings to return.  Useful for pagination
	Limit int
	// OrderBy is the name of the field to use for sorting
	OrderBy string
	// OrderDir is the direction to sort.  Use "asc" for ascending, "desc" for descending.
	OrderDir string
}

// DefaultHandwritingListParams are the default values for HandwritingListParams
var DefaultHandwritingListParams = HandwritingListParams{
	Offset:   0,
	Limit:    200,
	OrderBy:  "id",
	OrderDir: "asc",
}

// Client is a client for making API calls
type Client struct {
	client *http.Client
	url    *url.URL
}

// NewClient constructs a Client from a URL
func NewClient(u *url.URL) *Client {
	client := http.DefaultClient
	c := Client{
		client: client,
		url:    u,
	}
	return &c
}

// ListHandwritings retreives a list of handwrings
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
