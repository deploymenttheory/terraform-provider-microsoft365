// REF: https://learn.microsoft.com/en-us/graph/api/resources/synchronization-synchronizationschema?view=graph-rest-beta
package graphBetaServicePrincipalSynchronizationSchemaJson

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServicePrincipalSynchronizationSchemaJsonResourceModel represents the schema for the ServicePrincipalSynchronizationSchemaJson resource
type ServicePrincipalSynchronizationSchemaJsonResourceModel struct {
	ID                       types.String    `tfsdk:"id"`
	ServicePrincipalId       types.String    `tfsdk:"service_principal_id"`
	JobId                    types.String    `tfsdk:"job_id"`
	SynchronizationRulesJson SchemaJSONValue `tfsdk:"synchronization_rules_json"`
	DirectoriesJson          SchemaJSONValue `tfsdk:"directories_json"`
	Version                  types.String    `tfsdk:"version"`
	Timeouts                 timeouts.Value  `tfsdk:"timeouts"`
}

// synchronizationSchemaResponse is the part of the synchronizationSchema returned by Graph that the resource maps.
// The collections are kept as raw JSON so properties the SDK model does not know are not lost.
type synchronizationSchemaResponse struct {
	ID                   *string         `json:"id"`
	Version              *string         `json:"version"`
	SynchronizationRules json.RawMessage `json:"synchronizationRules"`
	Directories          json.RawMessage `json:"directories"`
}
