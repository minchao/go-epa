package epa

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	libraryVersion   = "0.0.1"
	defaultBaseURL   = "https://opendata.epa.gov.tw/"
	defaultUserAgent = "go-epa/" + libraryVersion
)

// A Client manages communication with the EPA Open data API.
type Client struct {
	client *http.Client
	token  string

	// Base URL for API requests. Defaults to the public GitHub API, but can be
	// set to a domain endpoint to use with GitHub Enterprise. BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the EPA API.
	UserAgent string
}

// NewClient returns a new EPA Open data API client.
func NewClient(token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		client:    httpClient,
		token:     token,
		UserAgent: defaultUserAgent,
		BaseURL:   baseURL,
	}
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Do sends an API request and returns the API response.
//
// The provided ctx must be non-nil. If it is canceled or time out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	q := req.URL.Query()
	q.Set("token", c.token)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if u, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(u).String()
				return nil, e
			}
		}

		return nil, err
	}
	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return resp, nil
}

// sanitizeURL redacts the PWD parameter from the URL which may be exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("token")) > 0 {
		params.Set("token", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

// ErrorResponse reports error caused by an API request.
type ErrorResponse struct {
	Response *http.Response
	Message  string `json:"Message"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

// CheckResponse checks the API response for errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}
