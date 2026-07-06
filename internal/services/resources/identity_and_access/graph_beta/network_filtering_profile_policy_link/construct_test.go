package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

func TestConstructCreateResourceFilteringPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyID:   types.StringValue("b0a6c790-c677-42a2-ab4e-6e0657322b5c"),
		PolicyType: types.StringValue(policyTypeFiltering),
		State:      types.StringValue("enabled"),
		Priority:   types.Int64Value(100),
	}

	body, err := constructCreateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructCreateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type":  filteringPolicyLinkODataType,
		"priority":     float64(100),
		"state":        "enabled",
		"loggingState": "enabled",
		"policy": map[string]any{
			"id":          "b0a6c790-c677-42a2-ab4e-6e0657322b5c",
			"@odata.type": filteringPolicyODataType,
		},
	}

	assertJSONMapEqual(t, expected, actual)
}

func TestConstructCreateResourceWebFilteringPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyID:   types.StringValue("9f217cd0-8192-4b1a-a673-4ead774897a4"),
		PolicyType: types.StringValue(policyTypeWebFiltering),
		State:      types.StringValue("enabled"),
	}

	body, err := constructCreateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructCreateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": webFilteringPolicyLinkODataType,
		"state":       "enabled",
		"policy": map[string]any{
			"id":          "9f217cd0-8192-4b1a-a673-4ead774897a4",
			"@odata.type": webFilteringPolicyODataType,
		},
	}

	assertJSONMapEqual(t, expected, actual)
}

func TestConstructCreateResourcePromptPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyID:   types.StringValue("e0759300-3367-4118-985b-e6a79c3d404f"),
		PolicyType: types.StringValue(policyTypePrompt),
		State:      types.StringValue("enabled"),
	}

	body, err := constructCreateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructCreateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": promptPolicyLinkODataType,
		"state":       "enabled",
		"policy": map[string]any{
			"id":          "e0759300-3367-4118-985b-e6a79c3d404f",
			"@odata.type": promptPolicyODataType,
		},
	}

	assertJSONMapEqual(t, expected, actual)
}

func TestConstructCreateResourceContentPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyID:   types.StringValue("44444444-4444-4444-4444-444444444444"),
		PolicyType: types.StringValue(policyTypeContent),
		State:      types.StringValue("enabled"),
	}

	body, err := constructCreateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructCreateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": filePolicyLinkODataType,
		"state":       "enabled",
		"policy": map[string]any{
			"id":          "44444444-4444-4444-4444-444444444444",
			"@odata.type": filePolicyODataType,
		},
	}

	assertJSONMapEqual(t, expected, actual)
}

func TestConstructCreateResourceNetskopeDlpPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyID:   types.StringValue("55555555-5555-5555-5555-555555555555"),
		PolicyType: types.StringValue(policyTypeNetskopeDlp),
		State:      types.StringValue("enabled"),
	}

	body, err := constructCreateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructCreateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": securityProviderPolicyLinkODataType,
		"state":       "enabled",
		"policy": map[string]any{
			"id":          "55555555-5555-5555-5555-555555555555",
			"@odata.type": securityProviderPolicyODataType,
		},
	}

	assertJSONMapEqual(t, expected, actual)
}

func TestConstructUpdateResourceFilteringPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyType: types.StringValue(policyTypeFiltering),
		State:      types.StringValue("disabled"),
	}

	body, err := constructUpdateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructUpdateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": filteringPolicyLinkODataType,
		"state":       "disabled",
	}

	assertJSONMapEqual(t, expected, actual)
}

func TestConstructUpdateResourceWebFilteringPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyType: types.StringValue(policyTypeWebFiltering),
		State:      types.StringValue("disabled"),
	}

	body, err := constructUpdateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructUpdateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": webFilteringPolicyLinkODataType,
		"state":       "disabled",
	}

	assertJSONMapEqual(t, expected, actual)
}

func TestConstructUpdateResourcePromptPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyType: types.StringValue(policyTypePrompt),
		State:      types.StringValue("disabled"),
	}

	body, err := constructUpdateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructUpdateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": promptPolicyLinkODataType,
		"state":       "disabled",
	}

	assertJSONMapEqual(t, expected, actual)
}

func TestConstructUpdateResourceContentPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyType: types.StringValue(policyTypeContent),
		State:      types.StringValue("disabled"),
	}

	body, err := constructUpdateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructUpdateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": filePolicyLinkODataType,
		"state":       "disabled",
	}

	assertJSONMapEqual(t, expected, actual)
}

func TestConstructUpdateResourceNetskopeDlpPolicy(t *testing.T) {
	model := &NetworkFilteringProfilePolicyLinkResourceModel{
		PolicyType: types.StringValue(policyTypeNetskopeDlp),
		State:      types.StringValue("disabled"),
	}

	body, err := constructUpdateResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructUpdateResource returned error: %v", err)
	}

	actual := serializeParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": securityProviderPolicyLinkODataType,
		"state":       "disabled",
	}

	assertJSONMapEqual(t, expected, actual)
}

func serializeParsableForTest(t *testing.T, body s.Parsable) map[string]any {
	t.Helper()

	requestInfo := abstractions.NewRequestInformation()
	if err := requestInfo.SetContentFromParsable(context.Background(), policyLinkTestRequestAdapter{}, "application/json", body); err != nil {
		t.Fatalf("SetContentFromParsable returned error: %v", err)
	}

	var actual map[string]any
	if err := json.Unmarshal(requestInfo.Content, &actual); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}
	return actual
}

func assertJSONMapEqual(t *testing.T, expected, actual map[string]any) {
	t.Helper()

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("marshal expected: %v", err)
	}
	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("marshal actual: %v", err)
	}
	if string(expectedJSON) != string(actualJSON) {
		t.Fatalf("unexpected JSON\nexpected: %s\nactual:   %s", expectedJSON, actualJSON)
	}
}
