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

	requestBody, ok := body.(*filteringProfileRequestBody)
	if !ok {
		t.Fatalf("constructResource returned %T, expected filteringProfileRequestBody", body)
	}

	if got := requestBody.name; got == nil || *got != "test-profile" {
		t.Fatalf("name = %v, expected test-profile", got)
	}
	if got := requestBody.description; got == nil || *got != "Managed by Terraform" {
		t.Fatalf("description = %v, expected Managed by Terraform", got)
	}
	if got := requestBody.priority; got == nil || *got != 100 {
		t.Fatalf("priority = %v, expected 100", got)
	}
	if got := requestBody.state; got == nil || *got != "enabled" {
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
