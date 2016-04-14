package handwritingio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// TokenError indicates a problem with API key or secret
type TokenError string

func (e TokenError) Error() string {
	return string(e)
}

// APIErrors indicates an error response returned by the API.
type APIErrors struct {
	StatusCode int        `json:"-"`
	Body       []byte     `json:"-"`
	Errors     []APIError `json:"errors"`
}

// APIError is a single error in the response from the API.
//
// There can be one or more per response.
type APIError struct {
	Error string `json:"error"`
	Field string `json:"field"`
}

func responseError(resp *http.Response) error {
	e := APIErrors{StatusCode: resp.StatusCode}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return err
	}

	e.Body = body

	err = json.Unmarshal(body, &e)
	if err != nil {
		return err
	}

	return e
}

// Error returns the message of the first APIError.
func (e APIErrors) Error() string {
	if e.Errors == nil || len(e.Errors) == 0 {
		return "unknown error"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error
	}
	return "multiple errors"
}
