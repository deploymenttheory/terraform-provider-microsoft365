package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

func (r *DeviceConfigurationTemplatesJsonResource) createRelationships(
	ctx context.Context,
	id string,
	desired graphmodels.DeviceConfigurationable,
) error {
	typ := *desired.GetOdataType()
	hasBindings := false
	for name := range deviceConfigurationRelationships[typ] {
		if _, exists := desired.GetAdditionalData()[name+"@odata.bind"]; exists {
			hasBindings = true
			break
		}
	}
	if !hasBindings {
		return nil
	}
	// Some types accept embedded bindings without persisting them. Read the actual
	// references and establish missing ones explicitly before completing creation.
	actual := map[string]any{"@odata.type": typ}
	if err := r.readRelationships(ctx, id, types.StringNull(), actual); err != nil {
		return err
	}
	previous := graphmodels.NewDeviceConfiguration()
	previous.SetAdditionalData(actual)
	return r.updateRelationships(ctx, id, previous, desired)
}

func (r *DeviceConfigurationTemplatesJsonResource) updateRelationships(
	ctx context.Context,
	id string,
	prior, desired graphmodels.DeviceConfigurationable,
) error {
	typ := *desired.GetOdataType()
	for name, collection := range deviceConfigurationRelationships[typ] {
		key := name + "@odata.bind"
		before, err := relationshipIDs(prior.GetAdditionalData()[key], collection)
		if err != nil {
			return fmt.Errorf("read previous %s: %w", name, err)
		}
		after, err := relationshipIDs(desired.GetAdditionalData()[key], collection)
		if err != nil {
			return fmt.Errorf("read configured %s: %w", name, err)
		}
		route := r.ResourcePath + "/{id}/" + strings.TrimPrefix(typ, "#") + "/" + name
		params := map[string]string{"id": id}
		for relatedID, binding := range after {
			if _, exists := before[relatedID]; exists {
				continue
			}
			method := abstractions.PUT
			if collection {
				method = abstractions.POST
			}
			if err := customrequests.JSONRequest(
				ctx,
				r.client.GetAdapter(),
				method,
				route+"/$ref",
				params,
				map[string]string{"@odata.id": binding},
				nil,
			); err != nil {
				return fmt.Errorf("set %s reference: %w", name, err)
			}
		}
		for relatedID := range before {
			if _, exists := after[relatedID]; exists || (!collection && len(after) != 0) {
				continue
			}
			deleteRoute := route + "/$ref"
			if collection {
				deleteRoute = route + "/{relatedId}/$ref"
				params["relatedId"] = relatedID
			}
			if err := customrequests.JSONRequest(
				ctx,
				r.client.GetAdapter(),
				abstractions.DELETE,
				deleteRoute,
				params,
				nil,
				nil,
			); err != nil {
				return fmt.Errorf("remove %s reference: %w", name, err)
			}
		}
	}
	return nil
}

func relationshipIDs(value any, collection bool) (map[string]string, error) {
	result := make(map[string]string)
	if value == nil {
		return result, nil
	}
	values := []any{value}
	if collection {
		if strings, ok := value.([]string); ok {
			values = make([]any, len(strings))
			for i, value := range strings {
				values[i] = value
			}
			value = values
		}
		var ok bool
		values, ok = value.([]any)
		if !ok {
			return nil, fmt.Errorf(
				"%w: collection binding must be an array of URLs",
				errInvalidProfileResponse,
			)
		}
	}
	for _, value := range values {
		binding, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf(
				"%w: binding must be a device configuration URL",
				errInvalidProfileResponse,
			)
		}
		var id string
		if _, suffix, found := strings.Cut(
			binding,
			"/deviceManagement/deviceConfigurations('",
		); found &&
			strings.HasSuffix(suffix, "')") {
			id = strings.TrimSuffix(suffix, "')")
		} else if _, suffix, found := strings.Cut(
			binding,
			"/deviceManagement/deviceConfigurations/",
		); found {
			id = suffix
		}
		parsed, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: binding must contain a device configuration ID",
				errInvalidProfileResponse,
			)
		}
		result[parsed.String()] = binding
	}
	return result, nil
}
