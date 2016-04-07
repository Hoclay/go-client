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

type Handwriting struct {
	ID                   string    `json:"id"`
	Title                string    `json:"title"`
	Created              time.Time `json:"date_created"`
	Modified             time.Time `json:"date_modified"`
	RatingNeatness       int       `json:"rating_neatness"`
	RatingEmbellishment  int       `json:"rating_embellishment"`
	RatingCharacterWidth int       `json:"rating_character_width"`
}

type Client struct {
	client *http.Client
	url    *url.URL
}

func NewClient(u *url.URL) *Client {
	client := http.DefaultClient
	c := Client{
		client: client,
		url:    u,
	}
	return &c
}

func (c *Client) ListHandwritings(offset, limit int, order_by, order_direction string) (handwritings []Handwriting, err error) {
	values := url.Values{}
	values.Add("offset", strconv.Itoa(offset))
	values.Add("limit", strconv.Itoa(limit))
	values.Add("order_by", order_by)
	values.Add("order_dir", order_direction)
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
