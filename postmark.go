// Package postmark encapsulates the Postmark API via Go
package postmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const postmarkURL = `https://api.postmarkapp.com`

// Client provides a connection to the Postmark API
type Client struct {
	// HTTPClient is &http.Client{} by default
	HTTPClient *http.Client
	// Server Token: Used for requests that require server level privileges. This token can be found on the Credentials tab under your Postmark server.
	ServerToken string
	// AccountToken: Used for requests that require account level privileges. This token is only accessible by the account owner, and can be found on the Account tab of your Postmark account.
	AccountToken string
	// BaseURL is the root API endpoint
	BaseURL string
}

const (
	accountToken = "account"
	serverToken  = "server"
)

// Options is an object to hold variable parameters to perform request.
type parameters struct {
	// Method is HTTP method type.
	Method string
	// Path is postfix for URI.
	Path string
	// Payload for the request.
	Payload interface{}
	// TokenType defines which token to use
	TokenType string
}

// NewClient builds a new Client pointer using the provided tokens, a default HTTPClient, and a default API base URL
// Accepts `Server Token`, and `Account Token` as arguments
// http://developer.postmarkapp.com/developer-api-overview.html#authentication
func NewClient(serverToken, accountToken string) *Client {
	return &Client{
		HTTPClient:   &http.Client{},
		ServerToken:  serverToken,
		AccountToken: accountToken,
		BaseURL:      postmarkURL,
	}
}

// doRequest performs the request to the Postmark API
func (client *Client) doRequest(ctx context.Context, opts parameters, dst interface{}) (err error) {
	url := fmt.Sprintf("%s/%s", client.BaseURL, opts.Path)

	var req *http.Request
	if req, err = http.NewRequestWithContext(
		ctx, opts.Method, url, nil,
	); err != nil {
		return err
	}

	if opts.Payload != nil {
		var payloadData []byte
		if payloadData, err = json.Marshal(opts.Payload); err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(payloadData))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewBuffer(payloadData)), nil
		}
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	switch opts.TokenType {
	case accountToken:
		req.Header.Add("X-Postmark-Account-Token", client.AccountToken)
	default:
		req.Header.Add("X-Postmark-Server-Token", client.ServerToken)
	}

	var res *http.Response
	if res, err = client.HTTPClient.Do(req); err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()
	var body []byte
	if body, err = io.ReadAll(res.Body); err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		// If the status code is not a success, attempt to unmarshall the body into the APIError struct.
		var apiErr APIError
		if err = json.Unmarshal(body, &apiErr); err != nil {
			return err
		}
		return apiErr
	}

	return json.Unmarshal(body, dst)
}

// APIError represents errors returned by Postmark
type APIError struct {
	// ErrorCode: see error codes here (https://postmarkapp.com/developer/api/overview#error-codes)
	ErrorCode int64 `json:"ErrorCode"`
	// Message contains error details
	Message string `json:"Message"`
}

// Error returns the error message details
func (res APIError) Error() string {
	return res.Message
}
