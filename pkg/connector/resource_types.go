package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// The user resource type is for all user objects from the database.
var userResourceType = &v2.ResourceType{
	Id:          "user",
	DisplayName: "User",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_USER},
}

// ACLResourceTypeID is the resource type ID for ACLs. Exposed so the
// acl grant-emission gate in userBuilder can be checked against
// cli.ConnectorOpts.WillSyncResourceType.
const ACLResourceTypeID = "acl"

var aclResourceType = &v2.ResourceType{
	Id:          ACLResourceTypeID,
	DisplayName: "ACL",
}
