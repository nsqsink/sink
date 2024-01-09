package http

import (
	"fmt"
	"net/url"
)

// isValidURL check if given url is valid url or not
func isValidURL(link string) (bool, error) {
	if _, err := url.ParseRequestURI(link); err != nil {
		return false, err
	}

	return true, nil
}

// isValidHTTPMethod validate supported http method
func isValidHTTPMethod(method string) (bool, error) {
	listMethod := map[string]bool{
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
	}

	available := listMethod[method]
	if !available {
		return false, fmt.Errorf("method %s not supported yet", method)
	}

	return true, nil
}
