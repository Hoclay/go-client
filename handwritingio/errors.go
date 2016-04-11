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

// Error indicates an error response returned by the API
type Error struct {
	StatusCode int
	Body       []byte
	Message    string
}

type jsonErrors struct {
	Errors []jsonError `json:"errors"`
}

type jsonError struct {
	Message string `json:"error"`
}

func responseError(resp *http.Response) Error {
	e := Error{StatusCode: resp.StatusCode}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		e.Message = err.Error()
		return e
	}

	e.Body = body

	jes := jsonErrors{}
	err = json.Unmarshal(body, &jes)
	if err != nil {
		e.Message = err.Error()
		return e
	}

	if len(jes.Errors) == 0 {
		e.Message = "unknown error"
		return e
	}
	je := jes.Errors[0]
	e.Message = je.Message
	if e.Message == "" {
		e.Message = "unknown error"
	}

	return e
}

func (e Error) Error() string {
	return e.Message
}
