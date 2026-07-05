package graphBetaNetworkFilteringProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	s "github.com/microsoft/kiota-abstractions-go/serialization"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *NetworkFilteringProfileResourceModel) (s.Parsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := &filteringProfileRequestBody{}

	convert.FrameworkToGraphString(data.Name, func(value *string) {
		requestBody.name = value
	})

	convert.FrameworkToGraphString(data.Description, func(value *string) {
		requestBody.description = value
	})

	convert.FrameworkToGraphInt64(data.Priority, func(value *int64) {
		requestBody.priority = value
	})

	if !data.State.IsNull() && !data.State.IsUnknown() {
		parsed, err := models.ParseStatus(data.State.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid filtering profile state: %s", err)
		}
		if parsed == nil {
			return nil, fmt.Errorf("invalid filtering profile state: %s", data.State.ValueString())
		}
		state := data.State.ValueString()
		requestBody.state = &state
	}

	if requestBody.priority == nil {
		return nil, fmt.Errorf("filtering profile priority is required")
	}

	if requestBody.state == nil {
		return nil, fmt.Errorf("filtering profile state is required")
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

type filteringProfileRequestBody struct {
	name        *string
	description *string
	priority    *int64
	state       *string
}

func (b *filteringProfileRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("name", b.name); err != nil {
		return err
	}
	if err := writer.WriteStringValue("description", b.description); err != nil {
		return err
	}
	if err := writer.WriteInt64Value("priority", b.priority); err != nil {
		return err
	}
	if err := writer.WriteStringValue("state", b.state); err != nil {
		return err
	}

	return nil
}

func (b *filteringProfileRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}
