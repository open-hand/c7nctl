package gitee

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultBaseURL = "https://gitee.com/api/v5/"
	userAgent      = "go-gitee"

	headerRateRemaining = "X-RateLimit-Remaining"

	headerOTP = "X-GitHub-OTP"
)

type Client struct {
	clientMu sync.Mutex   // clientMu protects the client during calls that modify the CheckRedirect func.
	client   *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public GitHub API, but can be
	// set to a domain endpoint to use with GitHub Enterprise. BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// Base URL for uploading files.
	UploadURL *url.URL

	// User agent used when communicating with the GitHub API.
	UserAgent string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	Repositories *RepositoriesService
}

type service struct {
	client *Client
}

// NewClient returns a new GitHub API client. If a nil httpClient is
// provided, a new http.Client will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	//uploadURL, _ := url.Parse(uploadBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c

	c.Repositories = (*RepositoriesService)(&c.common)

	return c
}

// Response is a GitHub API response. This wraps the standard http.Response
// returned from GitHub and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response

	// These fields provide the page values for paginating through a set of
	// results. Any or all of these may be set to the zero value for
	// responses that are not part of a paginated set, or for which there
	// are no additional pages.
	//
	// These fields support what is called "offset pagination" and should
	// be used with the ListOptions struct.
	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int

	// Additionally, some APIs support "cursor pagination" instead of offset.
	// This means that a token points directly to the next record which
	// can lead to O(1) performance compared to O(n) performance provided
	// by offset pagination.
	//
	// For APIs that support cursor pagination (such as
	// TeamsService.ListIDPGroupsInOrganization), the following field
	// will be populated to point to the next page.
	//
	// To use this token, set ListCursorOptions.Page to this value before
	// calling the endpoint again.
	NextPageToken string
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// req.Header.Set("Accept", mediaTypeV3)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.populatePageValues()
	return response
}

// populatePageValues parses the HTTP Link response headers and populates the
// various pagination link values in the Response.
func (r *Response) populatePageValues() {
	if links, ok := r.Response.Header["Link"]; ok && len(links) > 0 {
		for _, link := range strings.Split(links[0], ",") {
			segments := strings.Split(strings.TrimSpace(link), ";")

			// link must at least have href and rel
			if len(segments) < 2 {
				continue
			}

			// ensure href is properly formatted
			if !strings.HasPrefix(segments[0], "<") || !strings.HasSuffix(segments[0], ">") {
				continue
			}

			// try to pull out page parameter
			url, err := url.Parse(segments[0][1 : len(segments[0])-1])
			if err != nil {
				continue
			}
			page := url.Query().Get("page")
			if page == "" {
				continue
			}

			for _, segment := range segments[1:] {
				switch strings.TrimSpace(segment) {
				case `rel="next"`:
					if r.NextPage, err = strconv.Atoi(page); err != nil {
						r.NextPageToken = page
					}
				case `rel="prev"`:
					r.PrevPage, _ = strconv.Atoi(page)
				case `rel="first"`:
					r.FirstPage, _ = strconv.Atoi(page)
				case `rel="last"`:
					r.LastPage, _ = strconv.Atoi(page)
				}

			}
		}
	}
}

/*
An Error reports more details on an individual error in an ErrorResponse.
These are the possible validation error codes:

    missing:
        resource does not exist
    missing_field:
        a required field on a resource has not been set
    invalid:
        the formatting of a field is invalid
    already_exists:
        another resource has the same valid as this field
    custom:
        some resources return this (e.g. github.User.CreateKey()), additional
        information is set in the Message field of the Error

GitHub error responses structure are often undocumented and inconsistent.
Sometimes error is just a simple string (Issue #540).
In such cases, Message represents an error message as a workaround.

GitHub API docs: https://docs.github.com/en/rest/reference/#client-errors
*/
type Error struct {
	Resource string `json:"resource"` // resource on which the error occurred
	Field    string `json:"field"`    // field on which the error occurred
	Code     string `json:"code"`     // validation error code
	Message  string `json:"message"`  // Message describing the error. Errors with Code == "custom" will always have this set.
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error caused by %v field on %v resource",
		e.Code, e.Field, e.Resource)
}

func (e *Error) UnmarshalJSON(data []byte) error {
	type aliasError Error // avoid infinite recursion by using type alias.
	if err := json.Unmarshal(data, (*aliasError)(e)); err != nil {
		return json.Unmarshal(data, &e.Message) // data can be json string.
	}
	return nil
}

/*
An ErrorResponse reports one or more errors caused by an API request.

GitHub API docs: https://docs.github.com/en/rest/reference/#client-errors
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
	Errors   []Error        `json:"errors"`  // more detail on individual errors
	// Block is only populated on certain types of errors such as code 451.
	// See https://developer.github.com/changes/2016-03-17-the-451-status-code-is-now-supported/
	// for more information.
	Block *struct {
		Reason    string     `json:"reason,omitempty"`
		CreatedAt *Timestamp `json:"created_at,omitempty"`
	} `json:"block,omitempty"`
	// Most errors will also include a documentation_url field pointing
	// to some content that might help you resolve the error, see
	// https://docs.github.com/en/rest/reference/#client-errors
	DocumentationURL string `json:"documentation_url,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message, r.Errors)
}

// TwoFactorAuthError occurs when using HTTP Basic Authentication for a user
// that has two-factor authentication enabled. The request can be reattempted
// by providing a one-time password in the request.
type TwoFactorAuthError ErrorResponse

func (r *TwoFactorAuthError) Error() string { return (*ErrorResponse)(r).Error() }

// AcceptedError occurs when GitHub returns 202 Accepted response with an
// empty body, which means a job was scheduled on the GitHub side to process
// the information needed and cache it.
// Technically, 202 Accepted is not a real error, it's just used to
// indicate that results are not ready yet, but should be available soon.
// The request can be repeated after some time.
type AcceptedError struct {
	// Raw contains the response body.
	Raw []byte
}

func (*AcceptedError) Error() string {
	return "job scheduled on GitHub side; try again later"
}

// AbuseRateLimitError occurs when GitHub returns 403 Forbidden response with the
// "documentation_url" field value equal to "https://docs.github.com/en/rest/reference/#abuse-rate-limits".
type AbuseRateLimitError struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message

	// RetryAfter is provided with some abuse rate limit errors. If present,
	// it is the amount of time that the client should wait before retrying.
	// Otherwise, the client should try again later (after an unspecified amount of time).
	RetryAfter *time.Duration
}

func (r *AbuseRateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message)
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it. If rate limit is exceeded and reset time is in the future,
// Do returns *RateLimitError immediately without making a network API call.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req = withContext(ctx, req)

	// rateLimitCategory := category(req.URL.Path)

	// If we've hit rate limit, don't make further requests before Reset time.
	/*
		if err := c.checkRateLimitBeforeDo(req, rateLimitCategory); err != nil {
			return &Response{
				Response: err.Response,
				Rate:     err.Rate,
			}, err
		}
	*/

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()
				return nil, e
			}
		}

		return nil, err
	}

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCP connection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			io.CopyN(ioutil.Discard, resp.Body, maxBodySlurpSize)
		}

		resp.Body.Close()
	}()

	response := newResponse(resp)

	// c.rateMu.Lock()
	// c.rateLimits[rateLimitCategory] = response.Rate
	// c.rateMu.Unlock()

	err = CheckResponse(resp)
	if err != nil {
		// Special case for AcceptedErrors. If an AcceptedError
		// has been encountered, the response's payload will be
		// added to the AcceptedError and returned.
		//
		// Issue #1022
		aerr, ok := err.(*AcceptedError)
		if ok {
			b, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				return response, readErr
			}

			aerr.Raw = b
			return response, aerr
		}

		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

func withContext(ctx context.Context, req *http.Request) *http.Request {
	return req.WithContext(ctx)
}

// sanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have response
// body, and a JSON response body that maps to ErrorResponse.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// *AcceptedError for 202 Accepted status codes,
// and *TwoFactorAuthError for two-factor authentication errors.
func CheckResponse(r *http.Response) error {
	if r.StatusCode == http.StatusAccepted {
		return &AcceptedError{}
	}
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	// Re-populate error response body because GitHub error responses are often
	// undocumented and inconsistent.
	// Issue #1136, #540.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	switch {
	case r.StatusCode == http.StatusUnauthorized && strings.HasPrefix(r.Header.Get(headerOTP), "required"):
		return (*TwoFactorAuthError)(errorResponse)
	case r.StatusCode == http.StatusForbidden && strings.HasSuffix(errorResponse.DocumentationURL, "#abuse-rate-limits"):
		abuseRateLimitError := &AbuseRateLimitError{
			Response: errorResponse.Response,
			Message:  errorResponse.Message,
		}
		if v := r.Header["Retry-After"]; len(v) > 0 {
			// According to GitHub support, the "Retry-After" header value will be
			// an integer which represents the number of seconds that one should
			// wait before resuming making requests.
			retryAfterSeconds, _ := strconv.ParseInt(v[0], 10, 64) // Error handling is noop.
			retryAfter := time.Duration(retryAfterSeconds) * time.Second
			abuseRateLimitError.RetryAfter = &retryAfter
		}
		return abuseRateLimitError
	default:
		return errorResponse
	}
}
