package connector

import (
	"context"
	"fmt"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	resourceType "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-vultr/pkg/client"
)

type aclBuilder struct {
	resourceType *v2.ResourceType
	client       *client.VultrClient
}

func (a *aclBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (a *aclBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return aclResourceType
}

func (a *aclBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource
	ACLs, _, err := a.client.ListAccountACLs(ctx)
	if err != nil {
		fmt.Println("Error fetching ACLs:", err)
		return nil, "", nil, err
	}
	for _, acl := range ACLs {
		roleResource, err := parseIntoACLResource(acl)
		if err != nil {
			return nil, "", nil, err
		}

		resources = append(resources, roleResource)
	}
	return resources, "", nil, nil
}

func (a *aclBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var entitlements []*v2.Entitlement
	if resource.DisplayName == "" {
		return nil, "", nil, fmt.Errorf("DisplayName is empty for resource: %v", resource)
	}

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("ACL permission: %s", resource.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("Assigned ACL: %s", resource.DisplayName)),
	}

	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, "assigned", assigmentOptions...))

	return entitlements, "", nil, nil
}

func parseIntoACLResource(aclName string) (*v2.Resource, error) {
	ret, err := resourceType.NewResource(
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
