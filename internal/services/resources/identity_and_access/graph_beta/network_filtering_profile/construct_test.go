package graphBetaNetworkFilteringProfile

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConstructResourceSerializesObservedGraphPayload(t *testing.T) {
	body, err := constructResource(context.Background(), &NetworkFilteringProfileResourceModel{
		Name:        types.StringValue("test-profile"),
		Description: types.StringValue("Managed by Terraform"),
		Priority:    types.Int64Value(100),
		State:       types.StringValue("enabled"),
	})
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	content, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	want := `{"name":"test-profile","description":"Managed by Terraform","priority":100,"state":"enabled"}`
	if string(content) != want {
		t.Fatalf("request body = %s, expected %s", content, want)
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
