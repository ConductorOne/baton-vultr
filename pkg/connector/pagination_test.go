package connector

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-vultr/pkg/client"
	"github.com/conductorone/baton-vultr/test"
)

// testClientServing builds a client whose transport replies with body once. The response
// is constructed here rather than returned, so it never crosses a function boundary as a
// loose *http.Response.
func testClientServing(body string) *client.VultrClient {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	resp.Header.Set("Content-Type", "application/json")

	return test.NewTestClient(resp, nil)
}

// Vultr keeps meta on the last page and just empties links.next. That is the end of the
// sync, not a failure.
func TestVultrClient_ListUsersLastPage(t *testing.T) {
	testClient := testClientServing(`{"users":[{"id":"1"}],"meta":{"total":1,"links":{"next":"","prev":""}}}`)

	users, next, _, err := testClient.ListUsers(context.Background(), client.PageOptions{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(users))
	}
	if next != "" {
		t.Errorf("Expected an empty cursor on the last page, got %q", next)
	}
}

func TestVultrClient_ListUsersNextCursor(t *testing.T) {
	testClient := testClientServing(`{"users":[{"id":"1"}],"meta":{"total":2,"links":{"next":"bmV4dA","prev":""}}}`)

	_, next, _, err := testClient.ListUsers(context.Background(), client.PageOptions{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if next != "bmV4dA" {
		t.Errorf("Expected cursor %q, got %q", "bmV4dA", next)
	}
}

// The case the opt-in exists for. A 200 carrying users but no meta used to read a nil
// Meta straight through to meta.links.next.
func TestVultrClient_ListUsersMissingMeta(t *testing.T) {
	testClient := testClientServing(`{"users":[{"id":"1"}]}`)

	_, _, _, err := testClient.ListUsers(context.Background(), client.PageOptions{})
	if err == nil {
		t.Fatal("Expected an error when Vultr omits meta")
	}
	if !errors.Is(err, uhttp.ErrMissingPaginationData) {
		t.Fatalf("Expected ErrMissingPaginationData, got %v", err)
	}
}

// Single-resource reads carry no pagination data and must not start demanding it.
func TestVultrClient_GetUserByIDDoesNotRequirePagination(t *testing.T) {
	testClient := testClientServing(`{"user":{"id":"1","name":"root"}}`)

	user, _, err := testClient.GetUserByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.Id != "1" {
		t.Errorf("Expected user id %q, got %q", "1", user.Id)
	}
}

// Account reads go through the same generic path and are not paginated either.
func TestVultrClient_ListAccountACLsDoesNotRequirePagination(t *testing.T) {
	testClient := testClientServing(`{"account":{"name":"acme","acls":["manage_users"]}}`)

	acls, _, err := testClient.ListAccountACLs(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(acls) != 1 || acls[0] != "manage_users" {
		t.Errorf("Expected [manage_users], got %v", acls)
	}
}
