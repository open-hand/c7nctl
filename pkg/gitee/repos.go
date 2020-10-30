package gitee

import (
	"context"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"reflect"
)

type RepositoriesService service

type Repository struct {
}

// RepositoryTag represents a repository tag.
type RepositoryTag struct {
	Name    *string `json:"name,omitempty"`
	Message *string `json:"message,omitempty"`
	Commit  *Commit `json:"commit,omitempty"`
}

// ListOptions specifies the optional parameters to various List methods that
// support offset pagination.
type ListOptions struct {
	AccessToken string `url:"access_token,omitempty"`
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// ListTags lists tags for the specified repository.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-repository-tags
func (s *RepositoriesService) ListTags(ctx context.Context, owner string, repo string, opts *ListOptions) ([]*RepositoryTag, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/tags", owner, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tags []*RepositoryTag
	resp, err := s.client.Do(ctx, req, &tags)
	if err != nil {
		return nil, resp, err
	}

	return tags, resp, nil
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
