package graphBetaUsersUser

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
)

func TestUnitListResourceUser_01_NewResource(t *testing.T) {
	resource := NewUserListResource()
	if resource == nil {
		t.Fatal("Expected resource to be created, got nil")
	}

	// Verify it implements the list.ListResource interface
	var _ list.ListResource = resource
}

func TestUnitListResourceUser_02_Metadata(t *testing.T) {
	resource := NewUserListResource()
	listResource, ok := resource.(*UserListResource)
	if !ok {
		t.Fatal("Expected resource to be *UserListResource")
	}

	if listResource.ResourcePath != "/users" {
		t.Errorf("Expected ResourcePath to be '/users', got %s", listResource.ResourcePath)
	}

	if len(listResource.ReadPermissions) != 2 {
		t.Errorf("Expected 2 read permissions, got %d", len(listResource.ReadPermissions))
	}

	expectedPermissions := []string{
		"User.Read.All",
		"Directory.Read.All",
	}

	for i, expected := range expectedPermissions {
		if i >= len(listResource.ReadPermissions) {
			t.Errorf("Missing expected read permission: %s", expected)
			continue
		}
		if listResource.ReadPermissions[i] != expected {
			t.Errorf("Expected read permission '%s', got '%s'", expected, listResource.ReadPermissions[i])
		}
	}
}
