package stripesearch

import "errors"

type ResponseError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e ResponseError) HasError() bool {
	return e.Message != "" || e.Type != ""
}

func (e ResponseError) Error() error {
	return errors.New(e.Message)
}
