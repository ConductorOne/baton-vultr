package client

import (
	"context"
	"fmt"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/ratelimit"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
)

const (
	baseUrl             = "https://api.vultr.com/v2"
	getUsers            = "/users"
	getUserByID         = "/users/"
	getAccountPrincipal = "/account"
)

type VultrClient struct {
	wrapper     *uhttp.BaseHttpClient
	TokenSource oauth2.TokenSource
}

func New(ctx context.Context, bearerToken string) (*VultrClient, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	cli, err := uhttp.NewBaseHttpClientWithContext(context.Background(), httpClient)
	if err != nil {
		return nil, err
	}
	client := VultrClient{
		wrapper:     cli,
		TokenSource: getTokenSource(ctx, bearerToken),
	}
	return &client, nil
}

func NewClient(bearerToken string, httpClient ...*uhttp.BaseHttpClient) (*VultrClient, error) {
	tokenSource := getTokenSource(context.Background(), bearerToken)
	var wrapper = &uhttp.BaseHttpClient{}
	if len(httpClient) > 0 {
		wrapper = httpClient[0]
	}
	return &VultrClient{
		wrapper:     wrapper,
		TokenSource: tokenSource,
	}, nil
}

func (c *VultrClient) ListUsers(ctx context.Context, options PageOptions) ([]User, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res UserResponse
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseUrl, getUsers)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating UserResponse URL: %s", err))
		return nil, "", nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res, WithPageIndex(options.PageToken), WithPageLimit(options.PageSize))
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return nil, "", nil, err
	}

	result := res.Result
	if result != nil {
		var users []User
		for _, user := range result {
			users = append(users, user)
		}

		nextPageToken := ""
		if len(result) > 0 {
			totalRecords := len(res.Result)
			nextPageToken = getNextPageToken(options.PageToken, options.PageSize, totalRecords, len(result))
		}
		return users, nextPageToken, annotation, nil
	} else {
		return nil, "", annotation, nil
	}
}

func (c *VultrClient) ListAccountACLs(ctx context.Context) ([]string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res AccountResponse
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseUrl, getAccountPrincipal)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating Account URL: %s", err))
		return nil, nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting account data: %s", err))
		return nil, nil, err
	}

	if len(res.Account.ACLs) > 0 {
		return res.Account.ACLs, annotation, nil
	} else {
		return nil, annotation, nil
	}
}

func (c *VultrClient) GetUserByID(ctx context.Context, userID string) (User, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res UserSingleResponse
	var user User
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseUrl, getUserByID+userID)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating URL: %s", err))
		return user, nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res, WithQueryParam("userID", userID))
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resource: %s", err))
		return user, nil, err
	}

	if res.Result.Id != "" {
		user = res.Result
	}

	return user, annotation, nil
}

func (c *VultrClient) getResourcesFromAPI(
	ctx context.Context,
	urlAddress string,
	res any,
	reqOptions ...ReqOpt,
) (annotations.Annotations, error) {
	_, annotation, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res, reqOptions...)

	if err != nil {
		return nil, err
	}

	return annotation, nil
}

func (c *VultrClient) doRequest(
	ctx context.Context,
	method string,
	endpointUrl string,
	res interface{},
	reqOptions ...ReqOpt,
) (http.Header, annotations.Annotations, error) {
	var (
		resp *http.Response
		err  error
	)

	urlAddress, err := url.Parse(endpointUrl)

	if err != nil {
		return nil, nil, err
	}

	for _, o := range reqOptions {
		o(urlAddress)
	}

	authToken, err := c.TokenSource.Token()
	if err != nil {
		return nil, nil, err
	}

	req, err := c.wrapper.NewRequest(
		ctx,
		method,
		urlAddress,
		uhttp.WithContentTypeJSONHeader(),
		uhttp.WithAcceptJSONHeader(),
	)
	authToken.SetAuthHeader(req)

	if err != nil {
		return nil, nil, err
	}

	switch method {
	case http.MethodGet, http.MethodPut, http.MethodPost:
		var doOptions []uhttp.DoOption
		if res != nil {
			doOptions = append(doOptions, uhttp.WithResponse(&res))
		}
		resp, err = c.wrapper.Do(req, doOptions...)
		if resp != nil {
			defer resp.Body.Close()
		}
	case http.MethodDelete:
		resp, err = c.wrapper.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
	}

	if err != nil {
		return nil, nil, err
	}

	annotation := annotations.Annotations{}
	if resp != nil {
		if desc, err := ratelimit.ExtractRateLimitData(resp.StatusCode, &resp.Header); err == nil {
			annotation.WithRateLimiting(desc)
		} else {
			return nil, annotation, err
		}

		return resp.Header, annotation, nil
	}

	return nil, nil, err
}
