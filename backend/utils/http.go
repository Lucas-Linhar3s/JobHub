package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

// HTTPClient is a shared client for http requests
type HTTPClient struct {
	// Clt Internal HTTP client
	Clt *http.Client
	// BaseURL base address for this service
	BaseURL *url.URL
	// Name defines a human readable name for this service
	Name string
}

// HTTPRequest is an http request that
// can be performed by an HTTPClient
type HTTPRequest struct {
	Inner *http.Request
	Tag   string
}

// NewHTTPClient creates a new HTTPClient
func NewHTTPClient(serviceBaseURL string) *HTTPClient {
	u, _ := url.Parse(serviceBaseURL)

	return &HTTPClient{
		BaseURL: u,
		Clt: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (hc *HTTPClient) buildRequest(tag string, method string, uri string, headers map[string][]string, body io.Reader, cookies ...*http.Cookie) (*HTTPRequest, error) {
	hc.BaseURL.Path = path.Join(hc.BaseURL.Path, uri)

	inner, err := http.NewRequest(method, hc.BaseURL.String(), body)
	if err != nil {
		return nil, err
	}
	inner.Header = headers

	for _, cookie := range cookies {
		if cookie != nil {
			inner.AddCookie(cookie)
		}
	}

	return &HTTPRequest{Inner: inner, Tag: tag}, err
}

// BuildGetRequest creates a new GET request
// to be sent using a client
func (hc *HTTPClient) BuildGetRequest(tag string, uri string, headers map[string][]string, cookies ...*http.Cookie) (*HTTPRequest, error) {
	return hc.buildRequest(tag, http.MethodGet, uri, headers, nil, cookies...)
}

// BuildPatchRequest creates a new PATCH request
// with a payload body to be sent using a client
func (hc *HTTPClient) BuildPatchRequest(tag string, uri string, headers map[string][]string, body []byte, cookies ...*http.Cookie) (*HTTPRequest, error) {
	bodyReader := io.NopCloser(bytes.NewReader(body))
	return hc.buildRequest(tag, http.MethodPatch, uri, headers, bodyReader, cookies...)
}

// BuildPostRequest creates a new POST request
// with a payload body to be sent using a client
func (hc *HTTPClient) BuildPostRequest(tag string, uri string, headers map[string][]string, body []byte, cookies ...*http.Cookie) (*HTTPRequest, error) {
	bodyReader := io.NopCloser(bytes.NewReader(body))
	return hc.buildRequest(tag, http.MethodPost, uri, headers, bodyReader, cookies...)
}

// BuildPutRequest creates a new PUT request
// with a payload body to be sent using a client
func (hc *HTTPClient) BuildPutRequest(tag string, uri string, headers map[string][]string, body []byte, cookies ...*http.Cookie) (*HTTPRequest, error) {
	bodyReader := io.NopCloser(bytes.NewReader(body))
	return hc.buildRequest(tag, http.MethodPut, uri, headers, bodyReader, cookies...)
}

// BuildDeleteRequest creates a new DELETE request
// for deleting a remote resource
func (hc *HTTPClient) BuildDeleteRequest(tag string, uri string, headers map[string][]string, cookies ...*http.Cookie) (*HTTPRequest, error) {
	return hc.buildRequest(tag, http.MethodDelete, uri, headers, nil, cookies...)
}

// BuildDeleteRequestWithBody creates a new DELETE request
// for deleting a remote resource with Body
func (hc *HTTPClient) BuildDeleteRequestWithBody(tag string, uri string, headers map[string][]string, body []byte, cookies ...*http.Cookie) (*HTTPRequest, error) {
	bodyReader := io.NopCloser(bytes.NewReader(body))
	return hc.buildRequest(tag, http.MethodDelete, uri, headers, bodyReader, cookies...)
}

func RedirectRequest(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// PerformRequest performs an HTTP request
// through a client with its environment
func (hc *HTTPClient) PerformRequest(req *HTTPRequest, data interface{}, errData error) (err error) {
	var resp *http.Response
	resp, err = hc.Clt.Do(req.Inner)
	if err != nil {
		return err
	}

	defer func() {
		deferErr := resp.Body.Close()
		if deferErr != nil {
			log.Println(err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusNoContent:
	case http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusPartialContent:
		err = json.NewDecoder(resp.Body).Decode(data)
	default:
		if err = json.NewDecoder(resp.Body).Decode(errData); err == nil {
			err = errData
		}
	}

	if err != nil {
		return err
	}

	return
}
