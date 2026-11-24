package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		// e.g. "/microsoft.graph.iosScepCertificateProfile/rootCertificate", addressed below the profile ID.
		route := "/" + strings.TrimPrefix(typ, "#") + "/" + name
		for relatedID, binding := range after {
			if _, exists := before[relatedID]; exists {
				continue
			}
			reference := map[string]string{"@odata.id": binding}
			if collection {
				postConfig := customrequests.PostRequestConfig{
					APIVersion: customrequests.GraphAPIBeta,
					Endpoint:   r.ResourcePath + "/" + id + route + "/$ref",
					RawBody:    reference,
				}
				err = customrequests.PostRequestNoContent(ctx, r.client.GetAdapter(), postConfig)
			} else {
				putConfig := customrequests.PutRequestConfig{
					APIVersion:        customrequests.GraphAPIBeta,
					Endpoint:          r.ResourcePath,
					ResourceIDPattern: customrequests.ResourceIDPatternSegment,
					ResourceID:        id,
					EndpointSuffix:    route + "/$ref",
					RawBody:           reference,
				}
				err = customrequests.PutRequestByResourceId(ctx, r.client.GetAdapter(), putConfig)
			}
			if err != nil {
				return fmt.Errorf("set %s reference: %w", name, err)
			}
		}
		for relatedID := range before {
			if _, exists := after[relatedID]; exists || (!collection && len(after) != 0) {
				continue
			}
			// relatedID is a parsed GUID (see relationshipIDs), so it is safe to write into the URL.
			suffix := route + "/$ref"
			if collection {
				suffix = route + "/" + relatedID + "/$ref"
			}
			deleteConfig := customrequests.DeleteRequestConfig{
				APIVersion:        customrequests.GraphAPIBeta,
				Endpoint:          r.ResourcePath,
				ResourceIDPattern: customrequests.ResourceIDPatternSegment,
				ResourceID:        id,
				EndpointSuffix:    suffix,
			}
			if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
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
