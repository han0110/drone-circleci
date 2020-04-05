package circleci

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("prefix", "circleci-client")

// Client extends resty.Client with CircleCI API.
type Client struct {
	client     *resty.Client
	apiVersion APIVersion
	apiToken   string
}

// NewClient creates a client with default api version (v2).
func NewClient() *Client {
	return NewVersionedClient(DefaultAPIVersion)
}

// NewVersionedClient creates a client with api version assigned.
func NewVersionedClient(apiVersion APIVersion) *Client {
	return (&Client{
		client: resty.New().SetHeaders(map[string]string{
			"Accept":     "application/json",
			"User-Agent": "drone-circleci/0.1.0 (https://github.com/han0110/drone-circleci)",
		}),
	}).SetAPIVersion(apiVersion)
}

// C clones current client with api version and token.
func (c *Client) C() *Client {
	return NewVersionedClient(c.apiVersion).SetAPIToken(c.apiToken)
}

// SetAPIVersion sets api version of client.
func (c *Client) SetAPIVersion(apiVersion APIVersion) *Client {
	if !apiVersion.IsSupported() {
		panic(errors.Errorf("apiVersion %s not supported by circleci.Client yet", apiVersion))
	}
	c.apiVersion = apiVersion
	c.client.SetHostURL(fmt.Sprintf("%s/%s", APIEndpoint, c.apiVersion))
	return c
}

// SetAPIToken sets api token of client.
func (c *Client) SetAPIToken(apiToken string) *Client {
	c.apiToken = apiToken
	c.client.SetHeader(HeaderAPIToken, apiToken)
	return c
}

// Authenticate uses api token to get authenticated.
func (c *Client) Authenticate(apiToken string) error {
	c.SetAPIToken(apiToken)
	user, err := c.GetMyself(context.TODO())
	if err != nil {
		return err
	}
	logger.WithFields(logrus.Fields{
		"user": user.Name,
	}).Debug("authenticated successfully")
	return nil
}

type option func(*resty.Request)

func withContext(ctx context.Context) option {
	return func(req *resty.Request) {
		req.SetContext(ctx)
	}
}

func withQueryParams(params map[string]string) option {
	return func(req *resty.Request) {
		req.SetQueryParams(params)
	}
}

func withPathParams(params map[string]string) option {
	return func(req *resty.Request) {
		req.SetPathParams(params)
	}
}

func (c *Client) get(result interface{}, path string, options ...option) error {
	var errRes ErrorResponse
	req := c.client.R().
		SetResult(result).
		SetError(&errRes)
	for _, o := range options {
		o(req)
	}
	resp, err := req.Get(path)
	if err != nil {
		return errors.Wrapf(err, "failed to send GET %s", resp.Request.URL)
	}
	if !errRes.IsEmpty() {
		return errors.Wrapf(
			errRes, "sent GET %s with error response %d",
			resp.Request.URL, resp.StatusCode(),
		)
	}
	return nil
}

type listIterator struct {
	client        *Client
	path          string
	nextPageToken string
	extraOptions  []option
	err           error
}

func newListIterator(c *Client, path string) listIterator {
	return listIterator{client: c, path: path}
}

type listResponse struct {
	NextPageToken string          `json:"next_page_token"`
	Items         json.RawMessage `json:"items"`
}

func (i *listIterator) Error() error {
	return i.err
}

// Next get data of next page
func (i *listIterator) Next(ctx context.Context, result interface{}) bool {
	if i.err != nil || i.nextPageToken == NullPageToken {
		return false
	}
	var res listResponse
	options := []option{withContext(ctx)}
	if i.nextPageToken != "" {
		options = append(options, withQueryParams(map[string]string{"page-token": i.nextPageToken}))
	} else {
		options = append(options, i.extraOptions...)
	}
	if i.err = i.client.get(&res, i.path, options...); i.err != nil {
		return false
	}
	i.err = errors.Wrap(json.Unmarshal(res.Items, result), "failed to unmarshal items into result")
	if i.err != nil {
		return false
	}
	if i.nextPageToken = res.NextPageToken; i.nextPageToken == "" {
		i.nextPageToken = NullPageToken
	}
	return true
}

// All get data of all the pages
func (i *listIterator) All(ctx context.Context, result interface{}) error {
	if err := i.Error(); err != nil {
		return err
	}
	// Reset result value
	resultVal := reflect.ValueOf(result).Elem()
	resultVal.Set(reflect.MakeSlice(resultVal.Type(), 0, 0))
	// Iterate all pages
	pagedPtr := reflect.New(resultVal.Type())
	pagedPtr.Elem().Set(reflect.MakeSlice(resultVal.Type(), 0, 0))
	for i.Next(ctx, pagedPtr.Interface()) {
		if err := i.Error(); err != nil {
			return err
		}
		resultVal.Set(reflect.AppendSlice(resultVal, pagedPtr.Elem()))
	}
	// Reset nextPageToken
	i.nextPageToken = ""
	return nil
}
