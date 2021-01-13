package resource

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/common/consts"
	"github.com/choerodon/c7nctl/pkg/utils"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"net/http"
	"net/url"
	"strings"
)

const (
	mediaTypeJson = "application/json; text/plain; charset=utf-8"
)

type Client struct {
	client  *http.Client
	BaseURL *url.URL

	UserAgent string

	Business     bool
	ResourcePath string
	Username     string
	Password     string
}

func NewClient(httpClient *http.Client, bUrl string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	if bUrl == "" {
		bUrl = consts.OpenSourceResourceURL
	}
	baseURL, _ := url.Parse(bUrl)

	c := &Client{client: httpClient, BaseURL: baseURL}
	return c
}

func (c *Client) Init() {
	if c.Business {
		c.BaseURL, _ = url.Parse(consts.BusinessResourcePath)
	}
}

func (c *Client) GetInstallDefinition(version string) (*InstallDefinition, error) {
	resource, err := c.GetResource(version, consts.ResourceInstallFile)
	if err != nil {
		return nil, err
	}
	rdJson, err := yaml.ToJSON([]byte(resource))
	if err != nil {
		return nil, err
	}
	i := &InstallDefinition{}
	// slaver 使用了 core_v1.ContainerPort, 必须先转 JSON
	if err = json.Unmarshal(rdJson, i); err != nil {
		log.Panic(err)
	}

	if i.Spec.Basic.DefaultAccessModes == nil {
		i.Spec.Basic.DefaultAccessModes = []v1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	}
	return i, nil
}

func (c *Client) GetHelmValueFile(version, releaseName string) (string, error) {
	rvurl := fmt.Sprintf("%s/%s.yaml", consts.DefaultHelmValuesPath, releaseName)
	return c.GetResource(version, rvurl)
}

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
	req.Header.Set("Accept", mediaTypeJson)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

func (c *Client) Login() (*Auth, error) {
	if c.Business {
		if c.Username == "" || c.Password == "" {
			return nil, std_errors.New("username and password cannot be blank, when installing the commercial version of the application")
		}
		u := fmt.Sprintf("auth?username=%v&password=%v", c.Username, c.Password)
		req, err := c.NewRequest("POST", u, nil)
		if err != nil {
			return nil, err
		}
		auth := new(Auth)
		err = c.Do(context.Background(), req, auth)
		return auth, err
	}
	return nil, nil
}

func (c *Client) GetResource(version, url string) (string, error) {
	// 获取本地资源
	if c.ResourcePath != "" {
		path := fmt.Sprintf("%s/%s", c.ResourcePath, url)
		resource, err := readLocalFile(path)
		if err != nil {
			return "", std_errors.WithMessage(err, fmt.Sprintf("Read local resource file %s failed", path))
		}
		return resource, nil
	}

	// 生成获取商业版或者开源版的 url : 商业版需要认证
	auth := new(Auth)
	fu := url
	auth, err := c.Login()
	if err != nil {
		return "", std_errors.WithMessage(err, "Authentication business resource failed: ")
	}
	if auth == nil {
		fu = fmt.Sprintf(consts.OpenSourceResourceBasePath, version, url)
	} else {
		fu = fmt.Sprintf(consts.BusinessResourceBasePath, version, url, *auth.Data.Token)
	}

	result := new(bytes.Buffer)
	freq, err := c.NewRequest("GET", fu, nil)
	if err != nil {
		return "", err
	}
	err = c.Do(context.Background(), freq, result)
	return result.String(), err
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	if ctx == nil {
		return std_errors.New("context must be non-nil")
	}
	req = withContext(ctx, req)

	resp, err := c.client.Do(req)

	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()
				return e
			}
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

		return err
	}

	err = checkResponse(resp)
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

	return err
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &map[string]interface{}{}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return std_errors.New(fmt.Sprintf("%+v", errorResponse))
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

func readLocalFile(path string) (resource string, err error) {
	if err, ok := utils.IsFileExist(path); ok {
		log.Debugf("Read Local file %s", path)
		data, err := ioutil.ReadFile(path)
		resource = string(data)
		if err != nil {
			return "", std_errors.WithMessage(err, fmt.Sprintf("Failed to Read %s", path))
		}
	} else if err != nil {
		log.Debugf("can't find file %s : %+v", path, err)
	}
	return resource, nil
}
