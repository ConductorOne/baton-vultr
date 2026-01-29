package connector

import (
	"context"
	"fmt"
	"os"
	"testing"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-vultr/pkg/client"
	"github.com/stretchr/testify/assert"
)

var (
	ctx              = context.Background()
	apiToken         = os.Getenv("BEARERTOKEN")
	parentResourceID = &v2.ResourceId{}
	syncOpAttrs      = resource.SyncOpAttrs{PageToken: ""}
)

func initClient(t *testing.T) *client.VultrClient {
	if apiToken == "" {
		message :=
			fmt.Sprintf("Any of the required params not found. Api token: %s", apiToken)
		t.Skip(message)
	}

	c, err := client.New(ctx, apiToken)

	if err != nil {
		t.Errorf("ERROR: Failed to create client: %v", err)
	}
	return c
}

func TestUserBuilderList(t *testing.T) {
	c := initClient(t)

	u := newUserBuilder(c)
	res, _, err := u.List(ctx, parentResourceID, syncOpAttrs)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	message := fmt.Sprintf("Amount of users obtained: %d", len(res))
	t.Log(message)
}

func TestACLsBuilderList(t *testing.T) {
	c := initClient(t)

	acl := newACLbuilder(c)

	res, _, err := acl.List(ctx, parentResourceID, syncOpAttrs)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	message := fmt.Sprintf("Amount of ACLs obtained: %d", len(res))
	t.Log(message)
}
