package graphBetaConditionalAccessPolicy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
)

func TestUnitListResourceConditionalAccessPolicy_01_NewResource(t *testing.T) {
	resource := NewConditionalAccessPolicyListResource()
	if resource == nil {
		t.Fatal("Expected resource to be created, got nil")
	}

	// Verify it implements the list.ListResource interface
	var _ list.ListResource = resource
}

func TestUnitListResourceConditionalAccessPolicy_02_Metadata(t *testing.T) {
	resource := NewConditionalAccessPolicyListResource()
	listResource, ok := resource.(*ConditionalAccessPolicyListResource)
	if !ok {
		t.Fatal("Expected resource to be *ConditionalAccessPolicyListResource")
	}

	if listResource.ResourcePath != "/identity/conditionalAccess/policies" {
		t.Errorf("Expected ResourcePath to be '/identity/conditionalAccess/policies', got %s", listResource.ResourcePath)
	}

	if len(listResource.ReadPermissions) != 2 {
		t.Errorf("Expected 2 read permissions, got %d", len(listResource.ReadPermissions))
	}

	expectedPermissions := []string{
		"Policy.Read.All",
		"Policy.Read.ConditionalAccess",
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
