package test

import (
	"net/http"
	"os"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-vultr/pkg/client"
)

var (
	Users = []map[string]interface{}{
		{
			"id":          "1e5ab26c3423",
			"name":        "User_Root",
			"email":       "userroot@gmail.com",
			"api_enabled": true,
			"acls":        []string{"manage_users"},
		},
		{
			"id":          "776da6a156565",
			"name":        "deployServerUser",
			"email":       "deployServerUser@gmail.com",
			"api_enabled": true,
			"acls": []string{
				"subscriptions",
				"subscriptions_view",
				"provisioning",
			},
		},
		{
			"id":          "e97d4084dfd6",
			"name":        "upgradeServersUser",
			"email":       "upgradeServerUser@gmail.com",
			"api_enabled": true,
			"acls": []string{
				"subscriptions",
				"subscriptions_view",
				"upgrade",
			},
		},
		{
			"id":          "b6fd6483-343",
			"name":        "receiveAUO/TOSUser",
			"email":       "receiveAUO/TOSUser@gmail.com",
			"api_enabled": true,
			"acls": []string{
				"support",
				"abuse",
			},
		},
		{
			"id":          "257b3c9f-bddbc8a54dd330",
			"name":        "withoutRootPermissionUser",
			"email":       "withoutRootPermissionUser@gmail.com",
			"api_enabled": true,
			"acls": []string{
				"subscriptions",
				"subscriptions_view",
				"provisioning",
				"upgrade",
				"billing",
				"support",
				"dns",
				"firewall",
				"objstore",
				"loadbalancer",
				"vke",
				"abuse",
				"alerts",
			},
		},
	}
)

// Custom RoundTripper for testing.
type TestRoundTripper struct {
	response *http.Response
	err      error
}

type MockRoundTripper struct {
	Response  *http.Response
	Err       error
	roundTrip func(*http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

func (m *MockRoundTripper) SetRoundTrip(roundTrip func(*http.Request) (*http.Response, error)) {
	m.roundTrip = roundTrip
}

func (t *TestRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return t.response, t.err
}

// Helper function to create a test client with custom transport.
func NewTestClient(response *http.Response, err error) *client.VultrClient {
	transport := &TestRoundTripper{response: response, err: err}
	httpClient := &http.Client{Transport: transport}
	baseHttpClient := uhttp.NewBaseHttpClient(httpClient)

	bearerToken := ""

	newClientT, _ := client.NewClient(bearerToken, "", baseHttpClient)

	return newClientT
}

func ReadFile(fileName string) (string, error) {
	data, err := os.ReadFile("../../test/mockResponses/" + fileName)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
