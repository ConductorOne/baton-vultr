package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-vultr/pkg/client"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *client.VultrClient
}

func (u *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (u *userBuilder) List(ctx context.Context, _ *v2.ResourceId, opts rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	var resources []*v2.Resource
	pToken := &opts.PageToken

	bag, pageToken, err := getToken(pToken, userResourceType)
	if err != nil {
		return nil, nil, err
	}

	users, nextPageToken, _, err := u.client.ListUsers(ctx, client.PageOptions{
		Next:     pageToken,
		PageSize: pToken.Size,
	})
	if err != nil {
		return nil, nil, err
	}

	for _, user := range users {
		userCopy := user
		userResource, err := parseIntoUserResource(&userCopy)
		if err != nil {
			return nil, nil, err
		}
		resources = append(resources, userResource)
	}

	err = bag.Next(nextPageToken)
	if err != nil {
		return nil, nil, err
	}

	nextPageTokenStr, err := bag.Marshal()
	if err != nil {
		return nil, nil, err
	}

	return resources, &rs.SyncOpResults{NextPageToken: nextPageTokenStr}, nil
}

// Entitlements always returns an empty slice for users.
func (u *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

// The Grants function in the acls resource is performed in users for a better performance,
// since in this way for each user there is, the grants are directly assigned depending on which acls he has.
func (u *userBuilder) Grants(ctx context.Context, res *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	var grants []*v2.Grant
	var userID = res.Id.Resource

	user, _, err := u.client.GetUserByID(ctx, userID)

	if err != nil {
		return nil, nil, err
	}

	for _, userACL := range user.ACLs {
		if userACL != "" {
			aclResource := &v2.Resource{
				Id: &v2.ResourceId{
					ResourceType: aclResourceType.Id,
					Resource:     userACL,
				},
			}
			userCopy := user
			userResource, _ := parseIntoUserResource(&userCopy)
			userGrant := grant.NewGrant(aclResource, "assigned", userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("acl-grant:\n:%s:%s:%s", userACL, userID, "assigned"),
			}))
			grants = append(grants, userGrant)
		}
	}
	return grants, nil, nil
}

func parseIntoUserResource(user *client.User) (*v2.Resource, error) {
	var userStatus = v2.UserTrait_Status_STATUS_ENABLED

	profile := map[string]interface{}{
		"user_id":   user.Id,
		"user_name": user.Name,
		"email_id":  user.Email,
	}

	displayName := user.Name

	userTraits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
		rs.WithUserLogin(displayName),
		rs.WithEmail(user.Email, true),
	}

	ret, err := rs.NewUserResource(
		displayName,
		userResourceType,
		user.Id,
		userTraits,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newUserBuilder(c *client.VultrClient) *userBuilder {
	return &userBuilder{
		resourceType: userResourceType,
		client:       c,
	}
}
