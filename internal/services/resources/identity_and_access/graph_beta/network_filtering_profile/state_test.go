package graphBetaNetworkFilteringProfile

import (
	"context"
	"testing"
	"time"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func TestMapRemoteStateToTerraformMapsProfileFields(t *testing.T) {
	remote := models.NewFilteringProfile()
	id := "00000000-0000-0000-0000-000000000001"
	name := "test-profile"
	description := "Managed by Terraform"
	priority := int64(100)
	state := models.ENABLED_STATUS
	created := time.Date(2026, 7, 5, 1, 2, 3, 0, time.UTC)
	modified := time.Date(2026, 7, 5, 2, 3, 4, 0, time.UTC)
	version := "1.0.0"

	remote.SetId(&id)
	remote.SetName(&name)
	remote.SetDescription(&description)
	remote.SetPriority(&priority)
	remote.SetState(&state)
	remote.SetCreatedDateTime(&created)
	remote.SetLastModifiedDateTime(&modified)
	remote.SetVersion(&version)

	var data NetworkFilteringProfileResourceModel
	MapRemoteStateToTerraform(context.Background(), &data, remote)

	if data.ID.ValueString() != id {
		t.Fatalf("id = %q, expected %q", data.ID.ValueString(), id)
	}
	if data.Name.ValueString() != name {
		t.Fatalf("name = %q, expected %q", data.Name.ValueString(), name)
	}
	if data.Description.ValueString() != description {
		t.Fatalf("description = %q, expected %q", data.Description.ValueString(), description)
	}
	if data.Priority.ValueInt64() != priority {
		t.Fatalf("priority = %d, expected %d", data.Priority.ValueInt64(), priority)
	}
	if data.State.ValueString() != "enabled" {
		t.Fatalf("state = %q, expected enabled", data.State.ValueString())
	}
	if data.CreatedDateTime.ValueString() == "" {
		t.Fatal("created_date_time was empty")
	}
	if data.LastModifiedDateTime.ValueString() == "" {
		t.Fatal("last_modified_date_time was empty")
	}
	if data.Version.ValueString() != version {
		t.Fatalf("version = %q, expected %q", data.Version.ValueString(), version)
	}
}
