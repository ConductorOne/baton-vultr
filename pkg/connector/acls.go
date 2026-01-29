package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-vultr/pkg/client"
)

type aclBuilder struct {
	resourceType *v2.ResourceType
	client       *client.VultrClient
}

func (a *aclBuilder) Grants(_ context.Context, _ *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Grant, *resource.SyncOpResults, error) {
	return nil, nil, nil
}

func (a *aclBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return aclResourceType
}

func (a *aclBuilder) List(ctx context.Context, _ *v2.ResourceId, _ resource.SyncOpAttrs) ([]*v2.Resource, *resource.SyncOpResults, error) {
	var resources []*v2.Resource
	ACLs, _, err := a.client.ListAccountACLs(ctx)
	if err != nil {
		return nil, nil, err
	}
	for _, acl := range ACLs {
		aclResource, err := parseIntoACLResource(acl)
		if err != nil {
			return nil, nil, err
		}

		resources = append(resources, aclResource)
	}
	return resources, nil, nil
}

func (a *aclBuilder) Entitlements(_ context.Context, res *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Entitlement, *resource.SyncOpResults, error) {
	var entitlements []*v2.Entitlement

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("ACL permission: %s", res.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("Assigned ACL: %s", res.DisplayName)),
	}

	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(res, "assigned", assigmentOptions...))

	return entitlements, nil, nil
}

func parseIntoACLResource(aclName string) (*v2.Resource, error) {
	ret, err := resource.NewResource(
		aclName,
		aclResourceType,
		aclName,
	)
	if err != nil {
		return nil, err
	}
	ret.DisplayName = aclName

	return ret, nil
}

func newACLbuilder(c *client.VultrClient) *aclBuilder {
	return &aclBuilder{
		resourceType: aclResourceType,
		client:       c,
	}
}
