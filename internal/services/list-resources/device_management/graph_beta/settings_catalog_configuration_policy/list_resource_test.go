package graphBetaSettingsCatalogConfigurationPolicy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
)

func TestNewSettingsCatalogListResource(t *testing.T) {
	resource := NewSettingsCatalogListResource()
	if resource == nil {
		t.Fatal("Expected resource to be created, got nil")
	}

	// Verify it implements the list.ListResource interface
	var _ list.ListResource = resource
}

func TestSettingsCatalogListResource_Metadata(t *testing.T) {
	resource := NewSettingsCatalogListResource()
	listResource, ok := resource.(*SettingsCatalogListResource)
	if !ok {
		t.Fatal("Expected resource to be *SettingsCatalogListResource")
	}

	if listResource.ResourcePath != "/deviceManagement/configurationPolicies" {
		t.Errorf("Expected ResourcePath to be '/deviceManagement/configurationPolicies', got %s", listResource.ResourcePath)
	}

	if len(listResource.ReadPermissions) != 1 {
		t.Errorf("Expected 1 read permission, got %d", len(listResource.ReadPermissions))
	}

	if listResource.ReadPermissions[0] != "DeviceManagementConfiguration.Read.All" {
		t.Errorf("Expected read permission 'DeviceManagementConfiguration.Read.All', got %s", listResource.ReadPermissions[0])
	}
}
