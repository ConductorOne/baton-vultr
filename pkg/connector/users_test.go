package connector

import (
	"context"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-vultr/pkg/client"
	"github.com/conductorone/baton-vultr/test"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var pageOptions = client.PageOptions{
	PageSize:  5,
	PageToken: "",
}

// Tests that the client can fetch users based on the documented API below.
func TestVultrClient_GetUsers(t *testing.T) {
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(test.ReadFile("usersMock.json"))),
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

		if !reflect.DeepEqual(user.Id, expectedUser.Id) ||
			!reflect.DeepEqual(user.Name, expectedUser.Name) ||
			!reflect.DeepEqual(user.Email, expectedUser.Email) ||
			!reflect.DeepEqual(user.ApiEnabled, expectedUser.ApiEnabled) ||
			!reflect.DeepEqual(user.ACLs, expectedUser.ACLs) {
			t.Errorf("Unexpected user: got %+v, want %+v", user, expectedUser)
		}
	}

	if nextOptions == nil {
		t.Fatal("Expected non-nil nextOptions")
	}
}

func TestVultrClient_GetUsers_RequestDetails(t *testing.T) {
	var capturedRequest *http.Request
	mockTransport := &test.MockRoundTripper{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"users": []}`)),
			Header:     make(http.Header),
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
	if capturedRequest.URL.String() != expectedURL && !strings.Contains(capturedRequest.URL.String(), "limit=5") {
		t.Errorf("Expected URL %s, got %s", expectedURL, capturedRequest.URL.String())
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
