package connector

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-vultr/pkg/client"
	"github.com/conductorone/baton-vultr/test"
)

var pageOptions = client.PageOptions{
	Next:     "",
	PageSize: 0,
}

// Tests that the client can fetch users based on the documented API below.
func TestVultrClient_GetUsers(t *testing.T) {
	body, err := test.ReadFile("usersMock.json")
	if err != nil {
		t.Fatalf("Error reading body: %s", err)
	}
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	testClient := test.NewTestClient(mockResponse, nil)

	ctx := context.Background()

	result, _, nextOptions, err := testClient.ListUsers(ctx, pageOptions)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	expectedCount := len(test.Users)
	if len(result) != expectedCount {
		t.Errorf("Expected count to be %d, got %d", expectedCount, len(result))
	}

	for index, user := range result {
		expectedUser := client.User{
			Id:         test.Users[index]["id"].(string),
			Name:       test.Users[index]["name"].(string),
			Email:      test.Users[index]["email"].(string),
			ApiEnabled: test.Users[index]["api_enabled"].(bool),
			ACLs:       test.Users[index]["acls"].([]string),
		}

		if !reflect.DeepEqual(user, expectedUser) {
			t.Errorf("Unexpected user: got %+v, want %+v", user, expectedUser)
		}
	}

	if nextOptions == nil {
		t.Fatal("Expected non-nil nextOptions")
	}
}

// mockUserByIDClient returns a *client.VultrClient whose GetUserByID call
// resolves to the single-user fixture (userSingleMock.json), regardless of
// the ID passed in.
func mockUserByIDClient(t *testing.T) *client.VultrClient {
	t.Helper()

	body, err := test.ReadFile("userSingleMock.json")
	if err != nil {
		t.Fatalf("Error reading body: %s", err)
	}
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	return test.NewTestClient(mockResponse, nil)
}

// TestUserBuilder_Grants_EmitsACLGrants_WhenSyncACLsEnabled verifies that
// userBuilder.Grants delegates to the acl grant-emission optimization and
// returns non-empty grants when the acl resource type is being synced.
func TestUserBuilder_Grants_EmitsACLGrants_WhenSyncACLsEnabled(t *testing.T) {
	u := &userBuilder{
		resourceType: userResourceType,
		client:       mockUserByIDClient(t),
		syncACLs:     true,
	}

	res := &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: userResourceType.Id,
			Resource:     "1e5ab26c3423",
		},
	}

	grants, _, err := u.Grants(context.Background(), res, rs.SyncOpAttrs{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(grants) != 2 {
		t.Fatalf("Expected 2 acl grants (manage_users, billing), got %d", len(grants))
	}
}

// TestUserBuilder_Grants_SkipsACLGrants_WhenSyncACLsDisabled verifies that
// userBuilder.Grants short-circuits to (nil, nil, nil) when the connector's
// sync filter excludes the acl resource type, so no grant referencing an
// unsynced resource type is emitted.
func TestUserBuilder_Grants_SkipsACLGrants_WhenSyncACLsDisabled(t *testing.T) {
	u := &userBuilder{
		resourceType: userResourceType,
		client:       mockUserByIDClient(t),
		syncACLs:     false,
	}

	res := &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: userResourceType.Id,
			Resource:     "1e5ab26c3423",
		},
	}

	grants, syncOpResults, err := u.Grants(context.Background(), res, rs.SyncOpAttrs{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if grants != nil {
		t.Fatalf("Expected nil grants when syncACLs is false, got %+v", grants)
	}
	if syncOpResults != nil {
		t.Fatalf("Expected nil SyncOpResults when syncACLs is false, got %+v", syncOpResults)
	}
}

func TestVultrClient_GetUsers_RequestDetails(t *testing.T) {
	var capturedRequest *http.Request
	mockTransport := &test.MockRoundTripper{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(`{
			"users": [],
			"meta": {
				"total": 0,
				"links": {
					"next": "",
					"prev": ""
				}
			}
		}`)),
			Header: make(http.Header),
		},
		Err: nil,
	}
	mockTransport.Response.Header.Set("Content-Type", "application/json")

	mockRoundTrip := func(req *http.Request) (*http.Response, error) {
		capturedRequest = req
		return mockTransport.Response, mockTransport.Err
	}

	mockTransport.SetRoundTrip(mockRoundTrip)

	httpClient := &http.Client{Transport: mockTransport}
	baseHttpClient := uhttp.NewBaseHttpClient(httpClient)

	testClient, err := client.NewClient("access-token-hash", baseHttpClient)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	_, _, nextOptions, err := testClient.ListUsers(ctx, pageOptions)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if capturedRequest == nil {
		t.Fatal("No request was captured")
	}

	expectedURL := "https://api.vultr.com/v2/users"
	actualURL := capturedRequest.URL.String()
	if !strings.HasPrefix(actualURL, expectedURL) {
		t.Errorf("Expected URL to start with %s, got %s", expectedURL, actualURL)
	}

	expectedHeaders := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": "Bearer access-token-hash",
	}

	for key, expectedValue := range expectedHeaders {
		if value := capturedRequest.Header.Get(key); value != expectedValue {
			t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, value)
		}
	}

	if nextOptions == nil {
		t.Fatal("Expected non-nil nextOptions")
	}
}
