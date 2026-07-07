package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	store "github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestNewPolicyLinkCreateRequestInformation(t *testing.T) {
	requestInfo, err := newPolicyLinkRequestInformation(
		context.Background(),
		policyLinkTestRequestAdapter{},
		abstractions.POST,
		"profile-id",
		"",
		&policyLinkRequestBody{ODataType: filteringPolicyLinkODataType, State: "enabled"},
	)
	if err != nil {
		t.Fatalf("newPolicyLinkRequestInformation returned error: %v", err)
	}

	if requestInfo.Method != abstractions.POST {
		t.Fatalf("method = %s, want POST", requestInfo.Method)
	}
	if requestInfo.UrlTemplate != filteringProfilePoliciesURLTemplate {
		t.Fatalf("UrlTemplate = %q, want %q", requestInfo.UrlTemplate, filteringProfilePoliciesURLTemplate)
	}
	if requestInfo.PathParameters["filteringProfileId"] != "profile-id" {
		t.Fatalf("filteringProfileId path parameter not set")
	}
	if _, ok := requestInfo.PathParameters["policyLinkId"]; ok {
		t.Fatalf("policyLinkId path parameter should not be set for create")
	}
}

func TestNewPolicyLinkUpdateRequestInformation(t *testing.T) {
	requestInfo, err := newPolicyLinkRequestInformation(
		context.Background(),
		policyLinkTestRequestAdapter{},
		abstractions.PATCH,
		"profile-id",
		"link-id",
		&policyLinkRequestBody{ODataType: filteringPolicyLinkODataType, State: "disabled"},
	)
	if err != nil {
		t.Fatalf("newPolicyLinkRequestInformation returned error: %v", err)
	}

	if requestInfo.Method != abstractions.PATCH {
		t.Fatalf("method = %s, want PATCH", requestInfo.Method)
	}
	if requestInfo.UrlTemplate != filteringProfilePolicyLinkURLTemplate {
		t.Fatalf("UrlTemplate = %q, want %q", requestInfo.UrlTemplate, filteringProfilePolicyLinkURLTemplate)
	}
	if requestInfo.PathParameters["filteringProfileId"] != "profile-id" {
		t.Fatalf("filteringProfileId path parameter not set")
	}
	if requestInfo.PathParameters["policyLinkId"] != "link-id" {
		t.Fatalf("policyLinkId path parameter not set")
	}
}

type policyLinkTestRequestAdapter struct{}

func (policyLinkTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}
func (policyLinkTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (policyLinkTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}
func (policyLinkTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (policyLinkTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (policyLinkTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (policyLinkTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}
func (policyLinkTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}
func (policyLinkTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}
func (policyLinkTestRequestAdapter) SetBaseUrl(string)                            {}
func (policyLinkTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}
func (policyLinkTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}
