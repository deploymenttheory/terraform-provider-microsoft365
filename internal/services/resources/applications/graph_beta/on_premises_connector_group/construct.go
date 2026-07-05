package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
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
