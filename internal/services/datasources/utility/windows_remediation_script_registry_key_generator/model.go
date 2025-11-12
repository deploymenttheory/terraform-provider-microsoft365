package utilityWindowsRemediationScriptRegistryKeyGenerator

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsRemediationScriptRegistryKeyGeneratorDataSourceModel describes the data source data model.
type WindowsRemediationScriptRegistryKeyGeneratorDataSourceModel struct {
	ID                types.String   `tfsdk:"id"`
	Context           types.String   `tfsdk:"context"`
	RegistryKeyPath   types.String   `tfsdk:"registry_key_path"`
	ValueName         types.String   `tfsdk:"value_name"`
	ValueType         types.String   `tfsdk:"value_type"`
	ValueData         types.String   `tfsdk:"value_data"`
	DetectionScript   types.String   `tfsdk:"detection_script"`
	RemediationScript types.String   `tfsdk:"remediation_script"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}
