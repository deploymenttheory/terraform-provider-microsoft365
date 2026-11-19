package graphBetaNetworkFilteringProfile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource maps Terraform schema values to the minimal Graph payload.
//
// The generated beta SDK FilteringProfile serializer emits read-only/null fields
// such as createdDateTime, lastModifiedDateTime, and version. The
// /networkAccess/filteringProfiles endpoint rejects those fields on create with
// "400 Invalid request parameters", so this resource intentionally uses a small
// custom Kiota Parsable that serializes only the fields observed as writable.
//
// Microsoft Learn documents filteringProfile and the supported read/update
// surface here:
// https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringprofile?view=graph-rest-beta
// https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-update?view=graph-rest-beta
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

const (
	// Microsoft Learn currently documents the filteringProfile resource plus
	// list/get/update operations:
	// https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringprofile?view=graph-rest-beta
	// https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-list?view=graph-rest-beta
	// https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-get?view=graph-rest-beta
	// https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-update?view=graph-rest-beta
	//
	// The beta SDK also exposes collection POST and item DELETE builders for the
	// same /networkAccess/filteringProfiles paths. Live Graph probing confirmed
	// POST works when the payload includes name, description, priority, and state.
	filteringProfilesURLTemplate    = "{+baseurl}/networkAccess/filteringProfiles"
	filteringProfileItemURLTemplate = filteringProfilesURLTemplate + "/{filteringProfileId}"
)

func newFilteringProfileRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	filteringProfileID string,
	requestBody s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl": adapter.GetBaseUrl(),
	}

	urlTemplate := filteringProfilesURLTemplate
	if filteringProfileID != "" {
		urlTemplate = filteringProfileItemURLTemplate
		pathParameters["filteringProfileId"] = filteringProfileID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set filtering profile request content: %w", err)
		}
	}

	return requestInfo, nil
}
