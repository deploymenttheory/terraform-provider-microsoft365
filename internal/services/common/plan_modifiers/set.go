// planmodifiers/set.go
package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// defaultValueSet is a Set plan modifier that sets a default value when the config is null or empty

type SetModifier interface {
	planmodifier.Set
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}

type setModifier struct {
	description         string
	markdownDescription string
}

func (m setModifier) Description(ctx context.Context) string {
	return m.description
}

func (m setModifier) MarkdownDescription(ctx context.Context) string {
	return m.markdownDescription
}

type useStateForUnknownSet struct {
	setModifier
}

func (m useStateForUnknownSet) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	resp.PlanValue = req.StateValue
}

func UseStateForUnknownSet() SetModifier {
	return useStateForUnknownSet{
		setModifier: setModifier{
			description:         "Use state value if unknown",
			markdownDescription: "Use state value if unknown",
		},
	}
}

// defaultValueSet is a Set plan modifier that sets a default value when the config is null or empty
type defaultValueSet struct {
	setModifier
	defaultValue types.Set
}

// PlanModifySet sets the plan value to the default set if the config is null or empty.
func (m defaultValueSet) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.PlanValue.IsNull() && len(req.PlanValue.Elements()) > 0 {
		return
	}
	resp.PlanValue = m.defaultValue
}

// DefaultSetValue returns a SetModifier that sets the default value to the specified set.
func DefaultSetValue(defaultValue []attr.Value) planmodifier.Set {
	return defaultValueSet{
		setModifier: setModifier{
			description:         fmt.Sprintf("Default value set to %v", defaultValue),
			markdownDescription: fmt.Sprintf("Default value set to `%v`", defaultValue),
		},
		defaultValue: types.SetValueMust(types.StringType, defaultValue),
	}
}

// DefaultSetEmptyValue returns a SetModifier that sets the default value to an empty set.
func DefaultSetEmptyValue() planmodifier.Set {
	emptySet, _ := types.SetValue(types.StringType, []attr.Value{})
	return defaultValueSet{
		setModifier: setModifier{
			description:         "Default value set to empty set",
			markdownDescription: "Default value set to empty set",
		},
		defaultValue: emptySet,
	}
}

// RequiresOtherAttributeEnabledSet returns a plan modifier that ensures a set attribute
// can only be used when another specified attribute is enabled (set to true).
func RequiresOtherAttributeEnabledSet(dependencyPath path.Path) planmodifier.Set {
	return &requiresOtherAttributeEnabledSetModifier{
		dependencyPath: dependencyPath,
	}
}

type requiresOtherAttributeEnabledSetModifier struct {
	dependencyPath path.Path
}

func (m *requiresOtherAttributeEnabledSetModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when %s is enabled", m.dependencyPath)
}

func (m *requiresOtherAttributeEnabledSetModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when `%s` is enabled", m.dependencyPath)
}

func (m *requiresOtherAttributeEnabledSetModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.PlanValue.IsNull() {
		return
	}

	var dependencyValue types.Bool
	diags := req.Plan.GetAttribute(ctx, m.dependencyPath, &dependencyValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !dependencyValue.IsNull() && !dependencyValue.IsUnknown() && !dependencyValue.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid attribute usage",
			fmt.Sprintf("This attribute can only be used when %s is enabled (true)", m.dependencyPath),
		)
	}
}

// UseStateForUnknownOrNullSet returns a plan modifier that copies a known prior state
// Set value into the planned value if the planned value is null or unknown.
// This is useful for fields that are populated during creation but may not be
// explicitly set in configuration.
func UseStateForUnknownOrNullSet() planmodifier.Set {
	return useStateForUnknownOrNullSetModifier{}
}

// useStateForUnknownOrNullSetModifier implements the modifier
type useStateForUnknownOrNullSetModifier struct{}

// Description returns a plain text description of the modifier's behavior.
func (m useStateForUnknownOrNullSetModifier) Description(ctx context.Context) string {
	return "If the Set is unknown or null after plan creation, use the value from the state."
}

// MarkdownDescription returns a markdown formatted description of the modifier's behavior.
func (m useStateForUnknownOrNullSetModifier) MarkdownDescription(ctx context.Context) string {
	return "If the Set is unknown or null after plan creation, use the value from the state."
}

// PlanModifySet implements the plan modification logic.
func (m useStateForUnknownOrNullSetModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	tflog.Debug(ctx, "Using state value instead of null/unknown plan value for set", map[string]any{
		"path": req.Path.String(),
	})

	resp.PlanValue = req.StateValue
}

//------------------------------------------------------------------------------

// AllowSystemGeneratedSetValues returns a plan modifier that permits the system to add
// extra values to a set beyond what was specified in the configuration. This is useful
// for attributes where the service automatically adds system-managed values.
//
// The modifier ensures that:
// - All values from the configuration are present in the state
// - Additional values added by the system are permitted and preserved
// - State is updated to include both configured and system-generated values
//
// Example use case: Service principal tags where Microsoft automatically adds system tags
// like "WindowsAzureActiveDirectoryIntegratedApp" alongside user-specified tags.
func AllowSystemGeneratedSetValues() planmodifier.Set {
	return allowSystemGeneratedSetValuesModifier{}
}

type allowSystemGeneratedSetValuesModifier struct{}

func (m allowSystemGeneratedSetValuesModifier) Description(ctx context.Context) string {
	return "Allows the system to add additional values to the set beyond what was configured"
}

func (m allowSystemGeneratedSetValuesModifier) MarkdownDescription(ctx context.Context) string {
	return "Allows the system to add additional values to the set beyond what was configured"
}

func (m allowSystemGeneratedSetValuesModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// If config is null or unknown, don't modify anything
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// If plan value already equals state value, don't modify (no-op refresh)
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.PlanValue.Equal(req.StateValue) {
		return
	}

	// During creation (state is null), mark the plan as unknown to allow system-generated values
	// This tells Terraform to accept whatever the API returns
	if req.StateValue.IsNull() {
		tflog.Debug(ctx, "State is null (creation); marking plan as unknown to allow system-generated values", map[string]any{
			"path": req.Path.String(),
		})
		resp.PlanValue = types.SetUnknown(types.StringType)
		return
	}

	// During updates (state exists), check if config values are preserved in state
	if req.StateValue.IsUnknown() {
		return
	}

	// Get elements from config and state
	configElements := req.ConfigValue.Elements()
	stateElements := req.StateValue.Elements()

	// Build a map of config values for quick lookup
	configMap := make(map[string]bool)
	for _, elem := range configElements {
		if strVal, ok := elem.(types.String); ok {
			configMap[strVal.ValueString()] = true
		}
	}

	// Check if all config values are present in state
	allConfigValuesPresent := true
	for key := range configMap {
		found := false
		for _, stateElem := range stateElements {
			if strVal, ok := stateElem.(types.String); ok && strVal.ValueString() == key {
				found = true
				break
			}
		}
		if !found {
			allConfigValuesPresent = false
			break
		}
	}

	// If all config values are present in state, use the state value (which may include system-generated values)
	if allConfigValuesPresent {
		tflog.Debug(ctx, "All configured values present in state; using state value which includes system-generated values", map[string]any{
			"path":         req.Path.String(),
			"config_count": len(configElements),
			"state_count":  len(stateElements),
		})
		resp.PlanValue = req.StateValue
	}
}
