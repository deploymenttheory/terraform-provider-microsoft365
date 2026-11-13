package graphBetaNetworkWebFilteringPolicyRule

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestNewWebFilteringPolicyRuleRequestInformationSerializesObservedPortalPayload(t *testing.T) {
	ctx := context.Background()
	urlsOrFqdns, diags := types.SetValueFrom(ctx, types.StringType, []string{"*.example.com", "www.MySite.com/a/*"})
	if diags.HasError() {
		t.Fatalf("failed to build urls_or_fqdns set: %s", diags.Errors()[0].Detail())
	}
	webCategories, diags := types.SetValueFrom(ctx, types.StringType, []string{"AlcoholAndTobacco"})
	if diags.HasError() {
		t.Fatalf("failed to build web_categories set: %s", diags.Errors()[0].Detail())
	}
	httpMethods, diags := types.SetValueFrom(ctx, types.StringType, []string{"get"})
	if diags.HasError() {
		t.Fatalf("failed to build http_methods set: %s", diags.Errors()[0].Detail())
	}
	sessionTypes, diags := types.SetValueFrom(ctx, types.StringType, []string{"user", "agent"})
	if diags.HasError() {
		t.Fatalf("failed to build session_types set: %s", diags.Errors()[0].Detail())
	}

	body, err := constructResource(ctx, &NetworkWebFilteringPolicyRuleResourceModel{
		Name:          types.StringValue("sample-rule-for-codex"),
		Description:   types.StringValue("sample rule for codex"),
		Action:        types.StringValue("allow"),
		Priority:      types.Int64Value(100),
		Status:        types.StringValue("enabled"),
		UrlsOrFqdns:   urlsOrFqdns,
		WebCategories: webCategories,
		HTTPMethods:   httpMethods,
		SessionTypes:  sessionTypes,
	})
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	requestInfo, err := newWebFilteringPolicyRuleRequestInformation(
		ctx,
		webFilteringPolicyRuleTestRequestAdapter{},
		abstractions.POST,
		"05bd8400-14ae-4eb0-b7d7-339cd312d8f2",
		"",
		body,
	)
	if err != nil {
		t.Fatalf("newWebFilteringPolicyRuleRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}

	if payload["@odata.type"] != webFilteringRuleODataType {
		t.Fatalf("@odata.type = %#v, expected %#v", payload["@odata.type"], webFilteringRuleODataType)
	}
	action := payload["action"].(map[string]any)
	if action["@odata.type"] != webFilteringActionAllowODataType {
		t.Fatalf("action @odata.type = %#v, expected %#v", action["@odata.type"], webFilteringActionAllowODataType)
	}
	settings := payload["settings"].(map[string]any)
	if settings["status"] != "enabled" {
		t.Fatalf("settings.status = %#v, expected enabled", settings["status"])
	}

	matchingConditions := payload["matchingConditions"].(map[string]any)
	destinations := matchingConditions["destinations"].(map[string]any)
	if destinations["httpRequestMethod"] != "get" {
		t.Fatalf("httpRequestMethod = %#v, expected get", destinations["httpRequestMethod"])
	}
	targets := destinations["targets"].([]any)
	if len(targets) != 2 {
		t.Fatalf("targets length = %d, expected 2", len(targets))
	}
	urlTarget := targets[0].(map[string]any)
	if urlTarget["@odata.type"] != webFilteringURLDestinationODataType {
		t.Fatalf("url target @odata.type = %#v, expected %#v", urlTarget["@odata.type"], webFilteringURLDestinationODataType)
	}
	if urlTarget["values"].([]any)[0] != "*.example.com" {
		t.Fatalf("url target value = %#v, expected *.example.com", urlTarget["values"].([]any)[0])
	}
	if urlTarget["values"].([]any)[1] != "www.MySite.com/a/*" {
		t.Fatalf("url target value = %#v, expected www.MySite.com/a/*", urlTarget["values"].([]any)[1])
	}
	categoryTarget := targets[1].(map[string]any)
	if categoryTarget["@odata.type"] != webFilteringWebCategoryDestinationODataType {
		t.Fatalf("category target @odata.type = %#v, expected %#v", categoryTarget["@odata.type"], webFilteringWebCategoryDestinationODataType)
	}
	if categoryTarget["values"].([]any)[0] != "AlcoholAndTobacco" {
		t.Fatalf("category value = %#v, expected AlcoholAndTobacco", categoryTarget["values"].([]any)[0])
	}
	sources := matchingConditions["sources"].(map[string]any)
	if sources["sessionType"] != "user, agent" {
		t.Fatalf("sessionType = %#v, expected user, agent", sources["sessionType"])
	}
}

func TestNewWebFilteringPolicyRuleRequestInformationSerializesObservedCategoryOnlyPayload(t *testing.T) {
	ctx := context.Background()
	webCategories, diags := types.SetValueFrom(ctx, types.StringType, []string{"AIAgents"})
	if diags.HasError() {
		t.Fatalf("failed to build web_categories set: %s", diags.Errors()[0].Detail())
	}
	sessionTypes, diags := types.SetValueFrom(ctx, types.StringType, []string{"user", "agent"})
	if diags.HasError() {
		t.Fatalf("failed to build session_types set: %s", diags.Errors()[0].Detail())
	}

	payload := serializedRulePayload(t, ctx, &NetworkWebFilteringPolicyRuleResourceModel{
		Name:          types.StringValue("sample"),
		Description:   types.StringValue(""),
		Action:        types.StringValue("block"),
		Priority:      types.Int64Value(1000),
		Status:        types.StringValue("enabled"),
		WebCategories: webCategories,
		SessionTypes:  sessionTypes,
	})

	action := payload["action"].(map[string]any)
	if action["@odata.type"] != webFilteringActionBlockODataType {
		t.Fatalf("action @odata.type = %#v, expected %#v", action["@odata.type"], webFilteringActionBlockODataType)
	}

	matchingConditions := payload["matchingConditions"].(map[string]any)
	destinations := matchingConditions["destinations"].(map[string]any)
	if destinations["httpRequestMethod"] != nil {
		t.Fatalf("httpRequestMethod = %#v, expected nil", destinations["httpRequestMethod"])
	}

	targets := destinations["targets"].([]any)
	if len(targets) != 1 {
		t.Fatalf("targets length = %d, expected 1", len(targets))
	}
	categoryTarget := targets[0].(map[string]any)
	if categoryTarget["@odata.type"] != webFilteringWebCategoryDestinationODataType {
		t.Fatalf("category target @odata.type = %#v, expected %#v", categoryTarget["@odata.type"], webFilteringWebCategoryDestinationODataType)
	}
	if categoryTarget["values"].([]any)[0] != "AIAgents" {
		t.Fatalf("category value = %#v, expected AIAgents", categoryTarget["values"].([]any)[0])
	}

	sources := matchingConditions["sources"].(map[string]any)
	if sources["sessionType"] != "user, agent" {
		t.Fatalf("sessionType = %#v, expected user, agent", sources["sessionType"])
	}
}

func TestConstructResourceRequiresAtLeastOneDestination(t *testing.T) {
	ctx := context.Background()

	_, err := constructResource(ctx, &NetworkWebFilteringPolicyRuleResourceModel{
		Name:        types.StringValue("missing-destination"),
		Description: types.StringValue(""),
		Action:      types.StringValue("block"),
		Priority:    types.Int64Value(100),
		Status:      types.StringValue("enabled"),
	})
	if err == nil {
		t.Fatal("constructResource returned nil error, expected destination validation error")
	}
	if !strings.Contains(err.Error(), "at least one destination") {
		t.Fatalf("constructResource error = %q, expected destination validation error", err.Error())
	}
}

func serializedRulePayload(t *testing.T, ctx context.Context, model *NetworkWebFilteringPolicyRuleResourceModel) map[string]any {
	t.Helper()

	body, err := constructResource(ctx, model)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	requestInfo, err := newWebFilteringPolicyRuleRequestInformation(
		ctx,
		webFilteringPolicyRuleTestRequestAdapter{},
		abstractions.POST,
		"05bd8400-14ae-4eb0-b7d7-339cd312d8f2",
		"",
		body,
	)
	if err != nil {
		t.Fatalf("newWebFilteringPolicyRuleRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}

	return payload
}

type webFilteringPolicyRuleTestRequestAdapter struct{}

func (webFilteringPolicyRuleTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}

func (webFilteringPolicyRuleTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (webFilteringPolicyRuleTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}

func (webFilteringPolicyRuleTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (webFilteringPolicyRuleTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (webFilteringPolicyRuleTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (webFilteringPolicyRuleTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}

func (webFilteringPolicyRuleTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}

func (webFilteringPolicyRuleTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {
}

func (webFilteringPolicyRuleTestRequestAdapter) SetBaseUrl(string) {}

func (webFilteringPolicyRuleTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}

func (webFilteringPolicyRuleTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

func TestConstructResourceSerializesCustomHeadersInOrder(t *testing.T) {
	ctx := context.Background()
	urlsOrFqdns, diags := types.SetValueFrom(ctx, types.StringType, []string{"example.com"})
	if diags.HasError() {
		t.Fatalf("failed to build urls_or_fqdns set: %s", diags.Errors()[0].Detail())
	}
	headers, diags := types.ListValueFrom(ctx, customHeaderObjectType(), []CustomHeaderModel{
		{Name: types.StringValue("test"), Value: types.StringValue("a1")},
		{Name: types.StringValue("test2"), Value: types.StringValue("a2")},
	})
	if diags.HasError() {
		t.Fatalf("failed to build custom_headers list: %s", diags.Errors()[0].Detail())
	}

	payload := serializedRulePayload(t, ctx, &NetworkWebFilteringPolicyRuleResourceModel{
		Name:          types.StringValue("headers"),
		Description:   types.StringValue(""),
		Action:        types.StringValue("allow"),
		Priority:      types.Int64Value(100),
		Status:        types.StringValue("enabled"),
		UrlsOrFqdns:   urlsOrFqdns,
		CustomHeaders: headers,
	})

	got, err := json.Marshal(payload["action"])
	if err != nil {
		t.Fatalf("failed to marshal action: %v", err)
	}
	want := `{"@odata.type":"#microsoft.graph.networkaccess.webFilteringActionAllow","headerSettings":{"modifications":[` +
		`{"@odata.type":"#microsoft.graph.networkaccess.headerModificationAdd","headerName":"test","headerValue":"a1"},` +
		`{"@odata.type":"#microsoft.graph.networkaccess.headerModificationAdd","headerName":"test2","headerValue":"a2"}]}}`
	if string(got) != want {
		t.Fatalf("action = %s, expected %s", got, want)
	}
}

func TestConstructResourceSendsNullHeaderSettingsWithoutCustomHeaders(t *testing.T) {
	ctx := context.Background()
	webCategories, diags := types.SetValueFrom(ctx, types.StringType, []string{"AlcoholAndTobacco"})
	if diags.HasError() {
		t.Fatalf("failed to build web_categories set: %s", diags.Errors()[0].Detail())
	}

	payload := serializedRulePayload(t, ctx, &NetworkWebFilteringPolicyRuleResourceModel{
		Name:          types.StringValue("no-headers"),
		Description:   types.StringValue(""),
		Action:        types.StringValue("allow"),
		Priority:      types.Int64Value(100),
		Status:        types.StringValue("enabled"),
		WebCategories: webCategories,
		CustomHeaders: types.ListNull(customHeaderObjectType()),
	})

	headerSettings, ok := payload["action"].(map[string]any)["headerSettings"]
	if !ok || headerSettings != nil {
		t.Fatalf("action = %#v, expected explicit null headerSettings", payload["action"])
	}
}

func TestMapRemoteStateToTerraformReadsCustomHeaders(t *testing.T) {
	ctx := context.Background()
	cases := map[string]struct {
		body string
		want []CustomHeaderModel
	}{
		"headers in Graph order": {
			body: `{"id":"rule","action":{"@odata.type":"#microsoft.graph.networkaccess.webFilteringActionAllow","headerSettings":{"modifications":[` +
				`{"@odata.type":"#microsoft.graph.networkaccess.headerModificationAdd","headerName":"zeta","headerValue":"1"},` +
				`{"@odata.type":"#microsoft.graph.networkaccess.headerModificationAdd","headerName":"alpha","headerValue":"2"}]}}}`,
			want: []CustomHeaderModel{
				{Name: types.StringValue("zeta"), Value: types.StringValue("1")},
				{Name: types.StringValue("alpha"), Value: types.StringValue("2")},
			},
		},
		"non-add modifications skipped": {
			body: `{"id":"rule","action":{"@odata.type":"#microsoft.graph.networkaccess.webFilteringActionAllow","headerSettings":{"modifications":[` +
				`{"@odata.type":"#microsoft.graph.networkaccess.headerModificationRemove","headerName":"x-removed"},` +
				`{"@odata.type":"#microsoft.graph.networkaccess.headerModificationAdd","headerName":"kept","headerValue":"1"}]}}}`,
			want: []CustomHeaderModel{
				{Name: types.StringValue("kept"), Value: types.StringValue("1")},
			},
		},
		"null headerSettings": {
			body: `{"id":"rule","action":{"@odata.type":"#microsoft.graph.networkaccess.webFilteringActionAllow","headerSettings":null}}`,
		},
		"empty modifications": {
			body: `{"id":"rule","action":{"@odata.type":"#microsoft.graph.networkaccess.webFilteringActionAllow","headerSettings":{"modifications":[]}}}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			parseNode, err := jsonserialization.NewJsonParseNode([]byte(tc.body))
			if err != nil {
				t.Fatalf("failed to build parse node: %v", err)
			}
			parsed, err := parseNode.GetObjectValue(createWebFilteringPolicyRuleResponseFromDiscriminatorValue)
			if err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			var data NetworkWebFilteringPolicyRuleResourceModel
			MapRemoteStateToTerraform(ctx, &data, parsed.(*webFilteringPolicyRuleResponse))

			if tc.want == nil {
				if !data.CustomHeaders.IsNull() {
					t.Fatalf("custom_headers = %s, expected null", data.CustomHeaders)
				}
				return
			}
			var got []CustomHeaderModel
			if diags := data.CustomHeaders.ElementsAs(ctx, &got, false); diags.HasError() {
				t.Fatalf("failed to read custom_headers: %s", diags.Errors()[0].Detail())
			}
			if len(got) != len(tc.want) {
				t.Fatalf("custom_headers = %#v, expected %#v", got, tc.want)
			}
			for i := range got {
				if !got[i].Name.Equal(tc.want[i].Name) || !got[i].Value.Equal(tc.want[i].Value) {
					t.Fatalf("custom_headers[%d] = %#v, expected %#v", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestUniqueCustomHeaderNamesValidator(t *testing.T) {
	ctx := context.Background()
	cases := map[string]struct {
		names   []string
		wantErr bool
	}{
		"unique":                 {names: []string{"test", "test2"}},
		"duplicate":              {names: []string{"test", "test"}, wantErr: true},
		"case-insensitive match": {names: []string{"X-Test", "x-test"}, wantErr: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			headers := make([]CustomHeaderModel, 0, len(tc.names))
			for _, n := range tc.names {
				headers = append(headers, CustomHeaderModel{Name: types.StringValue(n), Value: types.StringValue("v")})
			}
			list, diags := types.ListValueFrom(ctx, customHeaderObjectType(), headers)
			if diags.HasError() {
				t.Fatalf("failed to build list: %s", diags.Errors()[0].Detail())
			}

			resp := &validator.ListResponse{}
			uniqueCustomHeaderNamesValidator{}.ValidateList(ctx, validator.ListRequest{
				Path:        path.Root("custom_headers"),
				ConfigValue: list,
			}, resp)

			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Fatalf("HasError = %v, expected %v: %v", resp.Diagnostics.HasError(), tc.wantErr, resp.Diagnostics)
			}
		})
	}
}

func TestConstructResourceOmitsHeaderSettingsForBlock(t *testing.T) {
	ctx := context.Background()
	webCategories, diags := types.SetValueFrom(ctx, types.StringType, []string{"AlcoholAndTobacco"})
	if diags.HasError() {
		t.Fatalf("failed to build web_categories set: %s", diags.Errors()[0].Detail())
	}

	payload := serializedRulePayload(t, ctx, &NetworkWebFilteringPolicyRuleResourceModel{
		Name:          types.StringValue("block"),
		Description:   types.StringValue(""),
		Action:        types.StringValue("block"),
		Priority:      types.Int64Value(100),
		Status:        types.StringValue("enabled"),
		WebCategories: webCategories,
		CustomHeaders: types.ListNull(customHeaderObjectType()),
	})

	if _, ok := payload["action"].(map[string]any)["headerSettings"]; ok {
		t.Fatalf("action = %#v, expected no headerSettings", payload["action"])
	}
}

func TestCustomHeadersRequireAllowValidatorSkipsUnknown(t *testing.T) {
	ctx := context.Background()
	schemaResp := &resource.SchemaResponse{}
	(&NetworkWebFilteringPolicyRuleResource{}).Schema(ctx, resource.SchemaRequest{}, schemaResp)
	objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)

	attrs := map[string]tftypes.Value{}
	for name, typ := range objectType.AttributeTypes {
		attrs[name] = tftypes.NewValue(typ, nil)
	}
	attrs["action"] = tftypes.NewValue(tftypes.String, "block")

	cases := map[string]struct {
		headers tftypes.Value
		wantErr bool
	}{
		"unknown": {headers: tftypes.NewValue(objectType.AttributeTypes["custom_headers"], tftypes.UnknownValue)},
		"set": {
			headers: tftypes.NewValue(objectType.AttributeTypes["custom_headers"], []tftypes.Value{
				tftypes.NewValue(objectType.AttributeTypes["custom_headers"].(tftypes.List).ElementType, map[string]tftypes.Value{
					"name":  tftypes.NewValue(tftypes.String, "Test1"),
					"value": tftypes.NewValue(tftypes.String, "a1"),
				}),
			}),
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			attrs["custom_headers"] = tc.headers
			resp := &resource.ValidateConfigResponse{}
			customHeadersRequireAllowValidator{}.ValidateResource(ctx, resource.ValidateConfigRequest{
				Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, attrs)},
			}, resp)

			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Fatalf("HasError = %v, expected %v: %v", resp.Diagnostics.HasError(), tc.wantErr, resp.Diagnostics)
			}
		})
	}
}
