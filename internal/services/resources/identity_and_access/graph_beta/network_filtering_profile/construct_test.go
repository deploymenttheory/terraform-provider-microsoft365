package graphBetaNetworkFilteringProfile

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConstructResourceMapsWritableFields(t *testing.T) {
	body, err := constructResource(context.Background(), &NetworkFilteringProfileResourceModel{
		Name:        types.StringValue("test-profile"),
		Description: types.StringValue("Managed by Terraform"),
		Priority:    types.Int64Value(100),
		State:       types.StringValue("enabled"),
	})
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	if got := body.GetName(); got == nil || *got != "test-profile" {
		t.Fatalf("name = %v, expected test-profile", got)
	}
	if got := body.GetDescription(); got == nil || *got != "Managed by Terraform" {
		t.Fatalf("description = %v, expected Managed by Terraform", got)
	}
	if got := body.GetPriority(); got == nil || *got != 100 {
		t.Fatalf("priority = %v, expected 100", got)
	}
	if got := body.GetState(); got == nil || got.String() != "enabled" {
		t.Fatalf("state = %v, expected enabled", got)
	}
}

func TestConstructResourceRejectsInvalidState(t *testing.T) {
	_, err := constructResource(context.Background(), &NetworkFilteringProfileResourceModel{
		Name:  types.StringValue("test-profile"),
		State: types.StringValue("invalid"),
	})
	if err == nil {
		t.Fatal("constructResource returned nil error for invalid state")
	}
}
