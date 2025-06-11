package configurationPolicyTemplateBuilders

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// GroupChildConfig represents a child setting within a group
type GroupChildConfig struct {
	SettingDefinitionId string
	SettingType         string // "simple_string", "simple_integer", "choice", "simple_collection", "choice_collection", "group_collection"
	Value               string
	Values              []string          // For collections
	ChoiceValues        []string          // For choice collections
	ChildConfig         *GroupChildConfig // For nested settings
	InstanceTemplateId  string
	ValueTemplateId     string
	UseTemplateDefault  bool // For template references
}

// ChoiceChildConfig represents a child setting within a choice setting (more specific than GroupChildConfig)
type ChoiceChildConfig struct {
	SettingDefinitionId string
	SettingType         string // "choice", "simple_string", "simple_integer", "simple_collection", "choice_collection"
	Value               string
	Values              []string // For simple collections
	ChoiceValues        []string // For choice collections
	InstanceTemplateId  string
	ValueTemplateId     string
	UseTemplateDefault  bool
}

// ========================================================================================
// MAIN DISPATCHER
// ========================================================================================

// SettingsDispatcher handles the conversion from HCL flat structure to Graph API settings
type SettingsDispatcher struct {
	ctx      context.Context
	registry *SettingsCatalogTemplate
}

// NewSettingsDispatcher creates a new dispatcher instance
func NewSettingsDispatcher(ctx context.Context) *SettingsDispatcher {
	return &SettingsDispatcher{
		ctx:      ctx,
		registry: NewSettingsCatalogTemplate(),
	}
}

// DispatchSettings converts HCL flat structure to Graph API settings
func (d *SettingsDispatcher) DispatchSettings(hclInput HCLSettingsInput) ([]models.DeviceManagementConfigurationSettingable, error) {
	tflog.Debug(d.ctx, "Starting settings dispatch from HCL", map[string]interface{}{
		"inputSettingsCount": len(hclInput),
	})

	// Step 1: Group settings by parent-child relationships
	groupedSettings := d.groupSettingsByRelationships(hclInput)

	// Step 2: Process each group
	var settings []models.DeviceManagementConfigurationSettingable
	for parentID, group := range groupedSettings {
		tflog.Debug(d.ctx, "Processing setting group", map[string]interface{}{
			"parentSettingId": parentID,
			"childrenCount":   len(group.Children),
		})

		setting, err := d.processSettingGroup(group)
		if err != nil {
			return nil, fmt.Errorf("failed to process setting group %s: %w", parentID, err)
		}

		settings = append(settings, setting)
	}

	// Step 3: Process standalone settings (those without parent-child relationships)
	standaloneSettings := d.findStandaloneSettings(hclInput, groupedSettings)
	for settingID, value := range standaloneSettings {
		tflog.Debug(d.ctx, "Processing standalone setting", map[string]interface{}{
			"settingId": settingID,
			"value":     value,
		})

		setting, err := d.processStandaloneSetting(settingID, value)
		if err != nil {
			return nil, fmt.Errorf("failed to process standalone setting %s: %w", settingID, err)
		}

		settings = append(settings, setting)
	}

	tflog.Debug(d.ctx, "Successfully dispatched all settings", map[string]interface{}{
		"totalSettings": len(settings),
	})

	return settings, nil
}

// ========================================================================================
// GROUPING AND ORGANIZATION
// ========================================================================================

// SettingGroup represents a parent setting with its children
type SettingGroup struct {
	Parent   SettingItem
	Children []SettingItem
}

// SettingItem represents a single setting item
type SettingItem struct {
	SettingID  string
	Value      string
	Definition SettingDefinition
}

// groupSettingsByRelationships organizes settings into parent-child groups
func (d *SettingsDispatcher) groupSettingsByRelationships(hclInput HCLSettingsInput) map[string]*SettingGroup {
	groups := make(map[string]*SettingGroup)

	// First pass: identify all parent settings
	for settingID, value := range hclInput {
		definition, exists := d.registry.definitions[settingID]
		if !exists {
			tflog.Warn(d.ctx, "Unknown setting definition", map[string]interface{}{
				"settingId": settingID,
			})
			continue
		}

		// If this setting has children, it's a parent
		if len(definition.ChildSettings) > 0 {
			groups[settingID] = &SettingGroup{
				Parent: SettingItem{
					SettingID:  settingID,
					Value:      value,
					Definition: definition,
				},
				Children: []SettingItem{},
			}
		}
	}

	// Second pass: assign children to their parents
	for settingID, value := range hclInput {
		definition, exists := d.registry.definitions[settingID]
		if !exists {
			continue
		}

		// If this setting has a parent, add it as a child
		if definition.ParentSetting != "" {
			if group, exists := groups[definition.ParentSetting]; exists {
				group.Children = append(group.Children, SettingItem{
					SettingID:  settingID,
					Value:      value,
					Definition: definition,
				})
			}
		}
	}

	return groups
}

// findStandaloneSettings finds settings that are not part of parent-child relationships
func (d *SettingsDispatcher) findStandaloneSettings(hclInput HCLSettingsInput, groups map[string]*SettingGroup) map[string]string {
	standalone := make(map[string]string)

	for settingID, value := range hclInput {
		definition, exists := d.registry.definitions[settingID]
		if !exists {
			continue
		}

		// Skip if it's a parent (already in groups)
		if len(definition.ChildSettings) > 0 {
			continue
		}

		// Skip if it's a child (already assigned to a parent)
		if definition.ParentSetting != "" {
			continue
		}

		// This is a standalone setting
		standalone[settingID] = value
	}

	return standalone
}

// ========================================================================================
// SETTING PROCESSORS
// ========================================================================================

// processSettingGroup processes a parent setting with its children
func (d *SettingsDispatcher) processSettingGroup(group *SettingGroup) (models.DeviceManagementConfigurationSettingable, error) {
	parent := group.Parent

	tflog.Debug(d.ctx, "Processing setting group", map[string]interface{}{
		"parentSettingId": parent.SettingID,
		"parentType":      parent.Definition.SettingType,
		"childrenCount":   len(group.Children),
	})

	switch parent.Definition.SettingType {
	case "choice_with_child":
		return d.processChoiceWithChild(parent, group.Children)
	default:
		return nil, fmt.Errorf("unsupported parent setting type: %s", parent.Definition.SettingType)
	}
}

// processChoiceWithChild processes a choice setting that has child settings
func (d *SettingsDispatcher) processChoiceWithChild(parent SettingItem, children []SettingItem) (models.DeviceManagementConfigurationSettingable, error) {
	if len(children) == 0 {
		return nil, fmt.Errorf("choice with child must have at least one child")
	}

	// For now, handle single child case (like LAPS backup directory + password age)
	if len(children) == 1 {
		child := children[0]

		switch child.Definition.ValueType {
		case "string":
			return ConstructChoiceWithStringSetting(
				d.ctx,
				parent.SettingID,
				parent.Value,
				child.SettingID,
				child.Value,
				parent.Definition.InstanceTemplateID,
				parent.Definition.ValueTemplateID,
				child.Definition.ValueTemplateID,
			)

		case "integer":
			return ConstructChoiceWithIntegerSetting(
				d.ctx,
				parent.SettingID,
				parent.Value,
				child.SettingID,
				child.Value,
				parent.Definition.InstanceTemplateID,
				parent.Definition.ValueTemplateID,
				child.Definition.ValueTemplateID,
			)

		case "choice":
			return ConstructChoiceWithBooleanSetting(
				d.ctx,
				parent.SettingID,
				parent.Value,
				child.SettingID,
				child.Value,
				parent.Definition.InstanceTemplateID,
				parent.Definition.ValueTemplateID,
				child.Definition.ValueTemplateID,
			)

		default:
			return nil, fmt.Errorf("unsupported child value type: %s", child.Definition.ValueType)
		}
	}

	// Multiple children case (Windows ASR pattern)
	var childConfigs []ChoiceChildConfig
	for _, child := range children {
		childConfig := ChoiceChildConfig{
			SettingDefinitionId: child.SettingID,
			Value:               child.Value,
			InstanceTemplateId:  child.Definition.InstanceTemplateID,
			ValueTemplateId:     child.Definition.ValueTemplateID,
			UseTemplateDefault:  child.Definition.UseTemplateDefault,
		}

		switch child.Definition.ValueType {
		case "string":
			childConfig.SettingType = "simple_string"
		case "integer":
			childConfig.SettingType = "simple_integer"
		case "choice":
			childConfig.SettingType = "choice"
		default:
			return nil, fmt.Errorf("unsupported child value type: %s", child.Definition.ValueType)
		}

		childConfigs = append(childConfigs, childConfig)
	}

	return ConstructChoiceWithMultipleChoiceChildren(
		d.ctx,
		parent.SettingID,
		parent.Value,
		childConfigs,
		parent.Definition.InstanceTemplateID,
		parent.Definition.ValueTemplateID,
		parent.Definition.UseTemplateDefault,
	)
}

// processStandaloneSetting processes a setting that has no parent-child relationships
func (d *SettingsDispatcher) processStandaloneSetting(settingID string, value string) (models.DeviceManagementConfigurationSettingable, error) {
	definition, exists := d.registry.definitions[settingID]
	if !exists {
		return nil, fmt.Errorf("unknown setting definition: %s", settingID)
	}

	tflog.Debug(d.ctx, "Processing standalone setting", map[string]interface{}{
		"settingId":   settingID,
		"settingType": definition.SettingType,
		"valueType":   definition.ValueType,
	})

	switch definition.SettingType {
	case "simple_choice":
		return ConstructSimpleChoiceSettingWithTemplate(
			d.ctx,
			settingID,
			value,
			definition.InstanceTemplateID,
			definition.ValueTemplateID,
			definition.UseTemplateDefault,
		)

	case "simple_string":
		return ConstructSimpleStringSetting(
			d.ctx,
			settingID,
			value,
			definition.InstanceTemplateID,
			definition.ValueTemplateID,
		)

	case "simple_integer":
		return ConstructSimpleIntegerSetting(
			d.ctx,
			settingID,
			value,
			definition.InstanceTemplateID,
			definition.ValueTemplateID,
		)

	case "simple_boolean":
		return ConstructSimpleBooleanChoiceSetting(
			d.ctx,
			settingID,
			value,
			definition.InstanceTemplateID,
			definition.ValueTemplateID,
		)

	case "simple_secret":
		return ConstructSimpleSecretSetting(
			d.ctx,
			settingID,
			value,
			definition.InstanceTemplateID,
			definition.ValueTemplateID,
		)

	case "string_collection":
		// Parse comma-separated values
		values := d.parseCollectionValue(value)
		return ConstructStringCollectionSetting(
			d.ctx,
			settingID,
			values,
			definition.InstanceTemplateID,
		)

	case "integer_collection":
		// Parse comma-separated values
		values := d.parseCollectionValue(value)
		return ConstructIntegerCollectionSetting(
			d.ctx,
			settingID,
			values,
			definition.InstanceTemplateID,
		)

	case "choice_collection":
		// Parse comma-separated values
		values := d.parseCollectionValue(value)
		return ConstructChoiceCollectionSetting(
			d.ctx,
			settingID,
			values,
			definition.InstanceTemplateID,
		)

	default:
		return nil, fmt.Errorf("unsupported setting type: %s", definition.SettingType)
	}
}

// ========================================================================================
// HELPER FUNCTIONS
// ========================================================================================

// parseCollectionValue parses a comma-separated string into a slice
func (d *SettingsDispatcher) parseCollectionValue(value string) []string {
	if value == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// RegisterCustomSetting allows external registration of custom setting definitions
func (d *SettingsDispatcher) RegisterCustomSetting(settingID string, definition SettingDefinition) {
	d.registry.definitions[settingID] = definition
}

// GetSupportedSettings returns a list of all supported setting IDs
func (d *SettingsDispatcher) GetSupportedSettings() []string {
	var settings []string
	for settingID := range d.registry.definitions {
		settings = append(settings, settingID)
	}
	return settings
}
