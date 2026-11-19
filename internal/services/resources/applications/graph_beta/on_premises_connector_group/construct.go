package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
)

func constructCreateResource(ctx context.Context, data *OnPremisesConnectorGroupResourceModel) (s.Parsable, error) {
	return constructResource(ctx, data.Name, data.Region, true)
}

func constructUpdateResource(ctx context.Context, plan, state *OnPremisesConnectorGroupResourceModel) (s.Parsable, error) {
	includeRegion := shouldSendRegion(plan.Region, state.Region)
	return constructResource(ctx, plan.Name, plan.Region, includeRegion)
}

func constructResource(ctx context.Context, name types.String, region types.String, includeRegion bool) (s.Parsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	// Microsoft Learn documents connectorGroup create/update at:
	// https://learn.microsoft.com/en-us/graph/api/connectorgroup-post?view=graph-rest-beta
	// https://learn.microsoft.com/en-us/graph/api/connectorgroup-update?view=graph-rest-beta
	//
	// Direct API verification on 2026-07-05 showed that connectorGroupType and
	// isDefault are system-managed values. Graph accepted those fields in create
	// payloads but ignored them, and PATCH rejected them as read-only. Keep this
	// request body intentionally narrow so Terraform only sends writable fields.
	requestBody := &connectorGroupRequestBody{}

	if !name.IsNull() && !name.IsUnknown() {
		value := name.ValueString()
		requestBody.name = &value
	}

	if includeRegion && !region.IsNull() && !region.IsUnknown() {
		value := region.ValueString()
		requestBody.region = &value
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

func shouldSendRegion(plan, state types.String) bool {
	// Region is writable according to Learn, but only while the connector group
	// has no assigned connectors or applications:
	// https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta
	//
	// PATCH requests verified on 2026-07-05 return 204 with no response body.
	// Omitting an unchanged region avoids unnecessary Graph validation of a
	// field that can become non-writable after assignments are added.
	if plan.IsNull() || plan.IsUnknown() {
		return false
	}
	if state.IsNull() || state.IsUnknown() {
		return true
	}

	return plan.ValueString() != state.ValueString()
}

type connectorGroupRequestBody struct {
	name   *string
	region *string
}

func (b *connectorGroupRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("name", b.name); err != nil {
		return err
	}
	if b.region != nil {
		if err := writer.WriteStringValue("region", b.region); err != nil {
			return err
		}
	}

	return nil
}

func (b *connectorGroupRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

const (
	defaultGraphBetaBaseURL = "https://graph.microsoft.com/beta"

	// Microsoft Learn documents connectorGroup CRUD at:
	// https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta
	// https://learn.microsoft.com/en-us/graph/api/connectorgroup-post?view=graph-rest-beta
	//
	// The current msgraph-beta-sdk-go package exposes generated request builders
	// for this path. A dedicated Kiota request helper is still used because the
	// generated ConnectorGroup model parses region through an enum generated from
	// Microsoft Graph beta OData CSDL metadata, while direct API verification on
	// 2026-07-05 returned "region": "japan", which is not in that enum. This
	// same verification also showed create returning HTTP 200 instead of the
	// Learn-documented 201, read-only create body values such as isDefault being
	// ignored, and PATCH returning 204 with no response body. This helper keeps
	// the normal Kiota adapter pipeline but pairs it with the custom
	// request/response types in this package so actual API values are preserved.
	connectorGroupsURLTemplate = "{+baseurl}/onPremisesPublishingProfiles/applicationProxy/connectorGroups"
	connectorGroupURLTemplate  = connectorGroupsURLTemplate + "/{connectorGroupId}"
)

func newConnectorGroupRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	connectorGroupID string,
	requestBody s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl": graphBetaBaseURL(adapter),
	}

	urlTemplate := connectorGroupsURLTemplate
	if connectorGroupID != "" {
		urlTemplate = connectorGroupURLTemplate
		pathParameters["connectorGroupId"] = connectorGroupID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set connector group request content: %w", err)
		}
	}

	return requestInfo, nil
}

func graphBetaBaseURL(adapter abstractions.RequestAdapter) string {
	if adapter == nil || adapter.GetBaseUrl() == "" {
		return defaultGraphBetaBaseURL
	}

	return strings.TrimRight(adapter.GetBaseUrl(), "/")
}
