package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time" // Added for seeding random number generator
)

// --- JSON Structure Types ---

// JSON structure types matching the input format
type MetadataFile struct {
	Metadata            Metadata             `json:"metadata"`
	PlatformDefinitions map[string][]Setting `json:"platformDefinitions"`
}

type Metadata struct {
	PlatformGroups     map[string]PlatformGroup `json:"platformGroups"`
	ExportTimestamp    string                   `json:"exportTimestamp"`
	TotalSettingsCount int                      `json:"totalSettingsCount"`
	Version            string                   `json:"version"`
}

type PlatformGroup struct {
	Platforms interface{} `json:"platforms"` // Can be either []string or string
	Count     int         `json:"count"`
}

type Setting struct {
	ODataType         string            `json:"oDataType"`
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	DisplayName       string            `json:"displayName"`
	Description       string            `json:"description"`
	Keywords          []string          `json:"keywords"`
	Visibility        string            `json:"visibility"`
	AccessTypes       string            `json:"accessTypes"`
	SettingDefinition SettingDefinition `json:"settingDefinition"`
	Applicability     Applicability     `json:"applicability"`
	Options           []Option          `json:"options,omitempty"`
	DefaultValue      *DefaultValue     `json:"defaultValue,omitempty"`
	ValueDefinition   *ValueDefinition  `json:"valueDefinition,omitempty"`
	DependentOn       []DependentOn     `json:"dependentOn,omitempty"`
	DependedOnBy      []DependedOnBy    `json:"dependedOnBy,omitempty"`
	InfoUrls          []string          `json:"infoUrls"`
}

type SettingDefinition struct {
	RootDefinitionID string `json:"rootDefinitionId"`
	CategoryID       string `json:"categoryId"`
	SettingUsage     string `json:"settingUsage"`
	BaseURI          string `json:"baseUri"`
	OffsetURI        string `json:"offsetUri"`
	Version          string `json:"version"`
	DefaultOptionID  string `json:"defaultOptionId,omitempty"`
}

type Applicability struct {
	Technologies string      `json:"technologies"`
	Platform     string      `json:"platform"`
	DeviceMode   string      `json:"deviceMode"`
	Description  interface{} `json:"description"`
}

type Option struct {
	DisplayName string      `json:"displayName"`
	Description interface{} `json:"description"`
	Value       interface{} `json:"value"`
	OptionID    interface{} `json:"optionId"`
}

type DefaultValue struct {
	ODataType                     string      `json:"@odata.type,omitempty"`
	SettingValueTemplateReference interface{} `json:"settingValueTemplateReference"`
	Value                         interface{} `json:"value"`
	ValueState                    string      `json:"valueState,omitempty"`
}

type ValueDefinition struct {
	ODataType             string      `json:"@odata.type,omitempty"`
	Format                string      `json:"format"`
	FileTypes             []string    `json:"fileTypes"`
	MinimumLength         int         `json:"minimumLength"`
	MaximumLength         int         `json:"maximumLength"`
	InputValidationSchema interface{} `json:"inputValidationSchema"`
	IsSecret              bool        `json:"isSecret,omitempty"`
	MinimumValue          int         `json:"minimumValue,omitempty"` // Assuming int for simplicity, might be float
	MaximumValue          int         `json:"maximumValue,omitempty"` // Assuming int for simplicity, might be float
}

type DependentOn struct {
	ParentSettingID string `json:"parentSettingId"`
	DependentOn     string `json:"dependentOn"` // Often the same as ParentSettingID
}

type DependedOnBy struct {
	Required     bool   `json:"required"`
	DependedOnBy string `json:"dependedOnBy"` // The ID of the child setting
}

// --- Schema Generation Structures ---

// SchemaTemplateData holds structured data for the template
type SchemaTemplateData struct {
	ResourceName         string
	ResourceType         string
	Description          string
	AttributeDefinitions []AttributeDefinition // Top-level attributes
}

// AttributeDefinition holds processed data for a single Terraform attribute
type AttributeDefinition struct {
	Name                         string
	Type                         string
	Description                  string // Raw description from setting
	Required                     bool   // Calculated based on parent dependency
	Optional                     bool   // Default unless Required or Computed
	Computed                     bool
	Sensitive                    bool
	ElementType                  string                         // For List/Set/Map of simple types
	NestedAttributes             map[string]AttributeDefinition // For nested object types!
	Validators                   []string
	PlanModifiers                []string
	DefaultValue                 string // Formatted default value string (for comments, actual default needs DefaultProvider)
	MarkdownDescription          string // Enhanced description for docs
	FormattedMarkdownDescription string // Enhanced description formatted for Go template string literal
	ODataInfo                    ODataInfo
	SourceSetting                Setting // Keep original setting for reference if needed
}

type ODataInfo struct {
	Type string
	ID   string
}

// --- Constants ---

// Constants for attribute types
const (
	AttrTypeString       = "schema.StringAttribute"
	AttrTypeBool         = "schema.BoolAttribute"
	AttrTypeInt64        = "schema.Int64Attribute"
	AttrTypeFloat64      = "schema.Float64Attribute"
	AttrTypeList         = "schema.ListAttribute"
	AttrTypeListNested   = "schema.ListNestedAttribute"
	AttrTypeSet          = "schema.SetAttribute"          // Note: Not currently generated, might need ODataType hints
	AttrTypeSetNested    = "schema.SetNestedAttribute"    // Note: Not currently generated
	AttrTypeMap          = "schema.MapAttribute"          // Note: Not currently generated
	AttrTypeMapNested    = "schema.MapNestedAttribute"    // Note: Not currently generated
	AttrTypeSingleNested = "schema.SingleNestedAttribute" // Note: Not currently generated, needs distinction from ListNested
	AttrTypeObject       = "schema.ObjectAttribute"       // Generally covered by *Nested types
)

// --- Global Processing Maps ---
var (
	allSettings        map[string]Setting  // All unique settings indexed by ID
	settingChildren    map[string][]string // Map parent ID -> list of child IDs that depend on it
	settingIsProcessed map[string]bool     // Track settings processed globally (added to schema structure)
	settingIsChild     map[string]bool     // Track settings that appear in any 'dependentOn' list
)

// --- Utility Functions ---

// parseFlexible attempts to parse the JSON data in a more flexible way
func parseFlexible(data []byte) (MetadataFile, error) {
	var result MetadataFile

	// First parse into a generic map
	var rawData map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawData); err != nil {
		return result, fmt.Errorf("error in flexible parsing step 1: %w", err)
	}

	// Parse metadata separately
	if metadataRaw, ok := rawData["metadata"]; ok {
		var metadata Metadata
		if err := json.Unmarshal(metadataRaw, &metadata); err != nil {
			fmt.Printf("Warning: Could not parse metadata: %v\n", err)
			// Continue without metadata
		} else {
			result.Metadata = metadata
		}
	}

	// Parse platformDefinitions separately
	if platformDefsRaw, ok := rawData["platformDefinitions"]; ok {
		var platformDefs map[string]json.RawMessage
		if err := json.Unmarshal(platformDefsRaw, &platformDefs); err != nil {
			return result, fmt.Errorf("error parsing platformDefinitions: %w", err)
		}

		// Initialize the map in the result
		result.PlatformDefinitions = make(map[string][]Setting)

		// For each platform, parse the settings
		for platform, settingsRaw := range platformDefs {
			var settings []Setting
			if err := json.Unmarshal(settingsRaw, &settings); err != nil {
				fmt.Printf("Warning: Could not parse settings for platform %s: %v\n", platform, err)
				continue
			}
			result.PlatformDefinitions[platform] = settings
		}
	} else {
		// Check if the root itself is an array of settings (less common structure)
		var settings []Setting
		if err := json.Unmarshal(data, &settings); err == nil && len(settings) > 0 {
			fmt.Println("Info: Input JSON appears to be a direct array of settings.")
			result.PlatformDefinitions = map[string][]Setting{
				"defaultPlatform": settings, // Assign a default platform key
			}
		} else {
			// Check if the root *contains* an array named 'value' (common in Graph API list responses)
			if valueRaw, ok := rawData["value"]; ok {
				var settingsInValue []Setting
				if err := json.Unmarshal(valueRaw, &settingsInValue); err == nil && len(settingsInValue) > 0 {
					fmt.Println("Info: Input JSON appears to be Graph API response with settings under 'value' key.")
					result.PlatformDefinitions = map[string][]Setting{
						"defaultPlatform": settingsInValue,
					}
				}
			}
		}
	}

	if len(result.PlatformDefinitions) == 0 {
		fmt.Println("Warning: No platformDefinitions found or recognized settings array in the input JSON.")
	}

	return result, nil
}

// cleanAttributeName prepares a string segment for use as part of a Go identifier/HCL name
func cleanAttributeNameSegment(segment string) string {
	// Replace disallowed chars with underscores
	result := strings.ReplaceAll(segment, ".", "_")
	result = strings.ReplaceAll(result, "-", "_")
	// Remove other potentially problematic chars (adjust regex as needed)
	nonAlphaNum := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	result = nonAlphaNum.ReplaceAllString(result, "")

	// Clean up underscores
	multipleUnderscores := regexp.MustCompile(`__+`)
	result = multipleUnderscores.ReplaceAllString(result, "_")
	result = strings.Trim(result, "_")

	return result
}

// convertToAttributeName converts CamelCase to snake_case
func convertToAttributeName(name string) string {
	if name == "" {
		return ""
	}
	// Simple CamelCase to snake_case conversion
	var result strings.Builder
	result.WriteRune(rune(strings.ToLower(string(name[0]))[0])) // First char lower
	for i := 1; i < len(name); i++ {
		char := rune(name[i])
		if char >= 'A' && char <= 'Z' {
			// Check if previous char was also uppercase (e.g., URL -> url, not u_r_l)
			// Or if next char is lowercase (e.g. APNName -> apn_name)
			prevChar := rune(name[i-1])
			isNextLower := (i+1 < len(name) && rune(name[i+1]) >= 'a' && rune(name[i+1]) <= 'z')
			if (prevChar >= 'a' && prevChar <= 'z') || isNextLower {
				result.WriteRune('_')
			}
			result.WriteRune(rune(strings.ToLower(string(char))[0]))
		} else {
			result.WriteRune(char) // Keep lowercase, digits, etc.
		}
	}
	return result.String()
}

// Regular expressions for name generation
var (
	nameSplitter    = regexp.MustCompile(`[._]`)
	startsWithDigit = regexp.MustCompile(`^[0-9]`)
	goKeywords      = map[string]bool{
		"break": true, "default": true, "func": true, "interface": true, "select": true,
		"case": true, "defer": true, "go": true, "map": true, "struct": true,
		"chan": true, "else": true, "goto": true, "package": true, "switch": true,
		"const": true, "fallthrough": true, "if": true, "range": true, "type": true,
		"continue": true, "for": true, "import": true, "return": true, "var": true,
		// Add common terraform/framework keywords if needed
		"resource": true, "schema": true, "provider": true, "data": true, "context": true,
	}
)

// generateNestedAttributeName derives a suitable HCL attribute/block name from a setting's ID and Name.
func generateNestedAttributeName(id string, name string) string {
	var baseName string

	if id != "" {
		// 1. Split the ID into segments
		segments := nameSplitter.Split(id, -1)

		// 2. Find the last relevant, non-empty segment
		for i := len(segments) - 1; i >= 0; i-- {
			segment := segments[i]
			// Ignore empty segments, placeholders like "[{0}]", or overly generic terms like "item" if possible
			if segment != "" && !strings.Contains(segment, "{") && strings.ToLower(segment) != "item" {
				baseName = segment
				break
			}
		}
	}

	// Fallback or if ID didn't yield a good name
	if baseName == "" {
		if name != "" {
			// Use the setting's 'name' field if ID was insufficient
			baseName = name
		} else if id != "" {
			// Last resort for ID: use the full cleaned ID
			baseName = cleanAttributeNameSegment(id) // Use segment cleaner
		} else {
			// No ID and no Name - generate placeholder
			return fmt.Sprintf("setting_%d", rand.Intn(100000))
		}
	}

	// 3. Convert base name (potentially CamelCase) to snake_case
	snakeCaseName := convertToAttributeName(baseName)

	// 4. Clean the snake_case name
	cleanedName := cleanAttributeNameSegment(snakeCaseName)

	// 5. Handle leading digits
	if startsWithDigit.MatchString(cleanedName) {
		cleanedName = "attr_" + cleanedName // Prefix to make it valid
	}

	// 6. Handle Go Keywords / common framework names
	// Check against lowercase version as keywords are lowercase
	lowerCleanedName := strings.ToLower(cleanedName)
	if goKeywords[lowerCleanedName] || lowerCleanedName == "id" { // Avoid collision with TF 'id'
		cleanedName = cleanedName + "_setting" // Suffix to avoid collision
	}

	// Ensure it's not empty after all cleaning
	if cleanedName == "" {
		return fmt.Sprintf("setting_%d", rand.Intn(100000))
	}

	return strings.ToLower(cleanedName) // Ensure final output is lowercase
}

// enhanceDescription adds details to the markdown description and formats it for Go.
func enhanceDescription(attr *AttributeDefinition, setting Setting) {
	var parts []string

	// Use DisplayName if different from Name, otherwise start with Description
	firstPart := setting.Description
	if setting.DisplayName != "" && setting.DisplayName != setting.Name {
		if firstPart != "" {
			firstPart = fmt.Sprintf("%s - %s", setting.DisplayName, firstPart)
		} else {
			firstPart = setting.DisplayName
		}
	}
	if firstPart != "" {
		parts = append(parts, firstPart)
	}

	if len(setting.Keywords) > 0 {
		parts = append(parts, fmt.Sprintf("Keywords: %s", strings.Join(setting.Keywords, ", ")))
	}
	if setting.ODataType != "" {
		parts = append(parts, fmt.Sprintf("OData Type: `%s`", setting.ODataType))
	}
	if setting.ID != "" {
		parts = append(parts, fmt.Sprintf("Settings Catalog ID: `%s`", setting.ID))
	}
	if setting.Applicability.Platform != "" {
		parts = append(parts, fmt.Sprintf("Platform: %s", setting.Applicability.Platform))
	}
	// Add min/max length/value if present
	if setting.ValueDefinition != nil {
		var constraints []string
		if setting.ValueDefinition.MinimumLength > 0 {
			constraints = append(constraints, fmt.Sprintf("min length: %d", setting.ValueDefinition.MinimumLength))
		}
		if setting.ValueDefinition.MaximumLength > 0 {
			constraints = append(constraints, fmt.Sprintf("max length: %d", setting.ValueDefinition.MaximumLength))
		}
		if setting.ValueDefinition.MinimumValue != 0 {
			constraints = append(constraints, fmt.Sprintf("min value: %d", setting.ValueDefinition.MinimumValue))
		}
		if setting.ValueDefinition.MaximumValue != 0 {
			constraints = append(constraints, fmt.Sprintf("max value: %d", setting.ValueDefinition.MaximumValue))
		}
		if len(constraints) > 0 {
			parts = append(parts, fmt.Sprintf("Constraints: %s", strings.Join(constraints, ", ")))
		}
	}

	// Join parts for Markdown
	markdownDesc := strings.Join(parts, "\\n\\n") // Use escaped newlines for TF docs
	attr.MarkdownDescription = markdownDesc       // Store raw markdown description

	// Format for Go multi-line string literal
	goFormattedDesc := ""
	goLines := strings.Split(markdownDesc, "\\n\\n") // Split by the explicit marker
	for i, line := range goLines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue // Skip empty lines
		}
		quotedLine := fmt.Sprintf("%q", trimmedLine) // Quote each part
		if i == 0 {
			goFormattedDesc = quotedLine
		} else {
			// Add concatenation operator and newline/tab for readability
			goFormattedDesc = goFormattedDesc + " + \"\\n\\n\" + \n\t\t\t\t" + quotedLine
		}
	}
	attr.FormattedMarkdownDescription = goFormattedDesc
}

// setAttributeType determines the Terraform schema type.
func setAttributeType(attr *AttributeDefinition, setting Setting) {
	attr.Type = AttrTypeString // Default to String

	switch setting.ODataType {
	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition":
		// Check for boolean-like choices
		if len(setting.Options) == 2 {
			isBoolLike := false
			optionNames := make(map[string]bool)
			for _, opt := range setting.Options {
				optionNames[strings.ToLower(opt.DisplayName)] = true
			}
			if (optionNames["true"] && optionNames["false"]) || (optionNames["enabled"] && optionNames["disabled"]) {
				isBoolLike = true
			}
			// Also consider if default option ID hints at boolean
			if !isBoolLike && setting.SettingDefinition.DefaultOptionID != "" {
				defaultOptSuffix := strings.ToLower(setting.SettingDefinition.DefaultOptionID)
				if strings.HasSuffix(defaultOptSuffix, "_true") || strings.HasSuffix(defaultOptSuffix, "_false") ||
					strings.HasSuffix(defaultOptSuffix, "_enabled") || strings.HasSuffix(defaultOptSuffix, "_disabled") {
					isBoolLike = true
				}
			}

			if isBoolLike {
				attr.Type = AttrTypeBool
				return // Found boolean type
			}
		}
		// Otherwise, it's a string choice
		attr.Type = AttrTypeString

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingDefinition":
		if setting.ValueDefinition != nil {
			odataType := fmt.Sprintf("%v", setting.ValueDefinition.ODataType) // Use ODataType from ValueDefinition
			if strings.Contains(odataType, "String") {
				attr.Type = AttrTypeString
				if setting.ValueDefinition.IsSecret {
					attr.Sensitive = true
				}
			} else if strings.Contains(odataType, "Integer") || strings.Contains(odataType, "Int64") {
				attr.Type = AttrTypeInt64
			} else if strings.Contains(odataType, "Boolean") || strings.Contains(odataType, "Bool") {
				attr.Type = AttrTypeBool
			} else if strings.Contains(odataType, "Double") || strings.Contains(odataType, "Float") {
				attr.Type = AttrTypeFloat64
			} else {
				fmt.Printf("Warning: Unknown ValueDefinition ODataType '%s' for simple setting %s, defaulting to String.\n", odataType, setting.ID)
				attr.Type = AttrTypeString
			}
		} else {
			fmt.Printf("Warning: Missing ValueDefinition for simple setting %s, defaulting to String.\n", setting.ID)
			attr.Type = AttrTypeString // Default if no value definition
		}

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionDefinition":
		attr.Type = AttrTypeList
		attr.ElementType = "StringType" // Default element type
		if setting.ValueDefinition != nil {
			odataType := fmt.Sprintf("%v", setting.ValueDefinition.ODataType)
			if strings.Contains(odataType, "Integer") || strings.Contains(odataType, "Int64") {
				attr.ElementType = "Int64Type"
			} else if strings.Contains(odataType, "Boolean") || strings.Contains(odataType, "Bool") {
				attr.ElementType = "BoolType"
			} else if strings.Contains(odataType, "Double") || strings.Contains(odataType, "Float") {
				attr.ElementType = "Float64Type"
			} else if strings.Contains(odataType, "String") {
				attr.ElementType = "StringType"
				if setting.ValueDefinition.IsSecret {
					attr.Sensitive = true // Mark list as sensitive if elements are secrets
				}
			} else {
				fmt.Printf("Warning: Unknown ValueDefinition ODataType '%s' for collection setting %s, defaulting to StringType elements.\n", odataType, setting.ID)
			}
		} else {
			fmt.Printf("Warning: Missing ValueDefinition for collection setting %s, defaulting to StringType elements.\n", setting.ID)
		}

	case "#microsoft.graph.deviceManagementConfigurationSettingGroupCollectionDefinition",
		"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstanceTemplate": // Handle instance template as well
		// This indicates a list of nested objects.
		attr.Type = AttrTypeListNested
		// Nested attributes will be populated recursively.

	case "#microsoft.graph.deviceManagementConfigurationGroupSettingInstanceTemplate",
		"#microsoft.graph.deviceManagementConfigurationSettingGroupDefinition": // Treat group definition similar to instance
		// This indicates a single nested object.
		// Need a way to differentiate Single vs List reliably. Maybe check if ID contains "[{0}]" or similar patterns?
		// For now, defaulting Group definitions/instances to SingleNested. Review if ListNested is more appropriate.
		attr.Type = AttrTypeSingleNested // Assuming single complex object
		fmt.Printf("Info: Setting %s (%s) mapped to SingleNestedAttribute. Verify if ListNestedAttribute is more appropriate.\n", setting.ID, setting.ODataType)

	case "#microsoft.graph.deviceManagementConfigurationSecretSettingDefinition":
		attr.Type = AttrTypeString
		attr.Sensitive = true

	default:
		fmt.Printf("Warning: Unknown ODataType '%s' for setting %s, defaulting to String.\n", setting.ODataType, setting.ID)
		attr.Type = AttrTypeString
	}
}

// addValidators adds schema validators based on setting definition.
func addValidators(attr *AttributeDefinition, setting Setting) {
	// String Choice Validator
	if attr.Type == AttrTypeString && setting.ODataType == "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition" && len(setting.Options) > 0 {
		var options []string
		for _, opt := range setting.Options {
			// Use DisplayName for the validator value
			options = append(options, fmt.Sprintf("%q", opt.DisplayName))
		}
		if len(options) > 0 {
			attr.Validators = append(attr.Validators, fmt.Sprintf("stringvalidator.OneOf(%s)", strings.Join(options, ", ")))
		}
	}

	// Numeric Range Validators
	if setting.ValueDefinition != nil {
		if attr.Type == AttrTypeInt64 {
			if setting.ValueDefinition.MinimumValue != 0 {
				attr.Validators = append(attr.Validators, fmt.Sprintf("int64validator.AtLeast(%d)", setting.ValueDefinition.MinimumValue))
			}
			if setting.ValueDefinition.MaximumValue != 0 {
				attr.Validators = append(attr.Validators, fmt.Sprintf("int64validator.AtMost(%d)", setting.ValueDefinition.MaximumValue))
			}
		} else if attr.Type == AttrTypeFloat64 {
			// Assuming Min/Max values are integers in JSON, cast if needed for float
			if setting.ValueDefinition.MinimumValue != 0 {
				attr.Validators = append(attr.Validators, fmt.Sprintf("float64validator.AtLeast(%f)", float64(setting.ValueDefinition.MinimumValue)))
			}
			if setting.ValueDefinition.MaximumValue != 0 {
				attr.Validators = append(attr.Validators, fmt.Sprintf("float64validator.AtMost(%f)", float64(setting.ValueDefinition.MaximumValue)))
			}
		}
	}

	// String Length Validators
	if attr.Type == AttrTypeString && setting.ValueDefinition != nil {
		if setting.ValueDefinition.MinimumLength > 0 {
			attr.Validators = append(attr.Validators, fmt.Sprintf("stringvalidator.LengthAtLeast(%d)", setting.ValueDefinition.MinimumLength))
		}
		if setting.ValueDefinition.MaximumLength > 0 {
			attr.Validators = append(attr.Validators, fmt.Sprintf("stringvalidator.LengthAtMost(%d)", setting.ValueDefinition.MaximumLength))
		}
	}

	// List Length Validators (Example - not directly in metadata, but could be added if needed)
	// if attr.Type == AttrTypeList || attr.Type == AttrTypeListNested {
	//  attr.Validators = append(attr.Validators, "listvalidator.SizeAtLeast(1)") // e.g., require at least one item
	// }
}

// setDefaultValue formats the default value for inclusion in comments (actual default requires provider impl).
func setDefaultValue(attr *AttributeDefinition, setting Setting) {
	if setting.DefaultValue == nil || setting.DefaultValue.Value == nil {
		return
	}

	// Attempt to format based on expected type
	var defaultValStr string
	switch attr.Type {
	case AttrTypeBool:
		if boolVal, ok := setting.DefaultValue.Value.(bool); ok {
			defaultValStr = fmt.Sprintf("%t", boolVal)
		} else if strVal, ok := setting.DefaultValue.Value.(string); ok { // Handle string representation like "True"
			lowerVal := strings.ToLower(strVal)
			if lowerVal == "true" || lowerVal == "enabled" {
				defaultValStr = "true"
			} else if lowerVal == "false" || lowerVal == "disabled" {
				defaultValStr = "false"
			}
		}
	case AttrTypeInt64:
		if numVal, ok := setting.DefaultValue.Value.(float64); ok { // JSON numbers are often float64
			defaultValStr = fmt.Sprintf("%d", int64(numVal))
		} else if strVal, ok := setting.DefaultValue.Value.(string); ok {
			defaultValStr = strVal // Keep as string, might need parsing
		}
	case AttrTypeFloat64:
		if numVal, ok := setting.DefaultValue.Value.(float64); ok {
			defaultValStr = fmt.Sprintf("%f", numVal)
		} else if strVal, ok := setting.DefaultValue.Value.(string); ok {
			defaultValStr = strVal // Keep as string, might need parsing
		}
	case AttrTypeString:
		if strVal, ok := setting.DefaultValue.Value.(string); ok {
			// Quote the string value properly for Go syntax
			defaultValStr = fmt.Sprintf("%q", strVal)
		} else {
			// Handle non-string default for string attribute (less common)
			defaultValStr = fmt.Sprintf("%q", fmt.Sprintf("%v", setting.DefaultValue.Value))
		}
	default:
		// Don't set default string for complex types or unhandled cases
		return
	}

	// Only store potentially non-empty/non-zero values
	if defaultValStr != "" && defaultValStr != "\"\"" && defaultValStr != "0" && defaultValStr != "false" && defaultValStr != "<nil>" {
		attr.DefaultValue = defaultValStr
	}
}

// --- Core Processing Logic ---

// Initializes global maps and preprocesses settings
func initializeAndPreprocess(metadata MetadataFile) {
	allSettings = make(map[string]Setting)
	settingChildren = make(map[string][]string)
	settingIsProcessed = make(map[string]bool)
	settingIsChild = make(map[string]bool)

	totalSettingsFound := 0
	// 1. Index all unique settings
	for platform, settingsList := range metadata.PlatformDefinitions {
		fmt.Printf("Info: Reading %d settings from platform '%s'.\n", len(settingsList), platform)
		totalSettingsFound += len(settingsList)
		for _, setting := range settingsList {
			if setting.ID == "" {
				fmt.Printf("Warning: Skipping setting with empty ID (Name: '%s', Platform: '%s').\n", setting.Name, platform)
				continue
			}
			if _, exists := allSettings[setting.ID]; !exists {
				allSettings[setting.ID] = setting
			} else {
				// Could potentially merge applicability or warn about overrides here
				// fmt.Printf("Info: Setting ID '%s' already indexed, potentially defined on multiple platforms.\n", setting.ID)
			}
		}
	}
	fmt.Printf("Info: Total settings found across all platforms: %d.\n", totalSettingsFound)
	fmt.Printf("Info: Indexed %d unique settings by ID.\n", len(allSettings))

	// 2. Build Children Map and Identify Child Status
	parentChildLinks := 0
	for id, setting := range allSettings {
		// Build parent -> children map from dependedOnBy
		// A setting ID listed in 'dependedOnBy' of 'parentID' means the setting ID is a child of 'parentID'
		if len(setting.DependedOnBy) > 0 {
			childIDs := []string{}
			for _, dep := range setting.DependedOnBy {
				if childID := dep.DependedOnBy; childID != "" {
					// Check if child actually exists in our indexed settings
					if _, childExists := allSettings[childID]; childExists {
						childIDs = append(childIDs, childID)
						// Mark the child setting as being a child
						settingIsChild[childID] = true
						parentChildLinks++
					} else {
						fmt.Printf("Warning: Setting '%s' depends on '%s', but child setting not found in indexed list.\n", id, childID)
					}
				}
			}
			if len(childIDs) > 0 {
				settingChildren[id] = childIDs
			}
		}

		// Alternative check: Mark settings that explicitly declare a parent via 'dependentOn'
		// This might be redundant if 'dependedOnBy' is comprehensive, but good for cross-check.
		for _, dep := range setting.DependentOn {
			if parentID := dep.ParentSettingID; parentID != "" {
				if _, parentExists := allSettings[parentID]; parentExists {
					if !settingIsChild[id] { // Only mark if not already marked by dependedOnBy
						// settingIsChild[id] = true // Optionally enable this for broader child detection
						// fmt.Printf("Info: Setting '%s' marked as child via dependentOn pointing to '%s'.\n", id, parentID)
					}
				} else {
					fmt.Printf("Warning: Setting '%s' depends on parent '%s', but parent setting not found.\n", id, parentID)
				}
			}
		}

	}
	fmt.Printf("Info: Established %d parent->child links.\n", parentChildLinks)
	childCount := 0
	for _, isChild := range settingIsChild {
		if isChild {
			childCount++
		}
	}
	fmt.Printf("Info: Identified %d settings that are children of other settings.\n", childCount)
}

// buildAttributeRecursively creates an AttributeDefinition, handling nesting.
// parentRequiredInfo maps CHILD ID -> isRequiredByParent.
// visited map is used for cycle detection within a single recursive path.
func buildAttributeRecursively(settingID string, parentRequiredInfo map[string]bool, visited map[string]bool) (AttributeDefinition, bool) {

	// --- Cycle Detection ---
	if visited[settingID] {
		fmt.Printf("Error: Cycle detected involving setting ID %s. Aborting build for this path.\n", settingID)
		return AttributeDefinition{}, false // Indicate failure/skip
	}
	visited[settingID] = true
	defer delete(visited, settingID) // Remove from path when returning up the stack

	// --- Get Setting ---
	setting, exists := allSettings[settingID]
	if !exists {
		// This should ideally not happen if called correctly from main/parent
		fmt.Printf("Error: Setting ID %s requested for build but not found in indexed settings.\n", settingID)
		return AttributeDefinition{}, false
	}

	// --- Basic Attribute Definition ---
	attr := AttributeDefinition{
		// Name will be set by the caller based on context (top-level vs nested)
		ODataInfo: ODataInfo{
			Type: setting.ODataType,
			ID:   setting.ID,
		},
		NestedAttributes: make(map[string]AttributeDefinition), // Initialize map
		SourceSetting:    setting,                              // Store original setting
	}

	// Determine base type (String, Bool, Int, List, ListNested, etc.)
	setAttributeType(&attr, setting) // Sets attr.Type, attr.ElementType, attr.Sensitive

	// Enhance description *after* type determination might influence it
	enhanceDescription(&attr, setting) // Sets attr.MarkdownDescription, attr.FormattedMarkdownDescription

	// Add validators
	addValidators(&attr, setting) // Sets attr.Validators

	// Set default value string (for comments)
	setDefaultValue(&attr, setting) // Sets attr.DefaultValue

	// --- Determine Required/Optional/Computed ---
	attr.Optional = true // Default to Optional
	attr.Required = false
	attr.Computed = false

	// Check if this setting is required by its parent (passed down)
	if required, specified := parentRequiredInfo[settingID]; specified && required {
		// Only set Required=true if explicitly marked as required by parent
		attr.Required = true
		attr.Optional = false
	}

	// --- Handle Nesting (if applicable) ---
	childIDs, hasChildren := settingChildren[settingID]

	// Check if the type itself implies it SHOULD have children (even if none listed/found)
	isGroupType := attr.Type == AttrTypeListNested || attr.Type == AttrTypeSingleNested ||
		attr.Type == AttrTypeSetNested || attr.Type == AttrTypeMapNested

	if hasChildren || isGroupType {
		if !isGroupType {
			// If children were found via 'dependedOnBy' but type wasn't set to nested, update it.
			// Default to ListNested, might need refinement.
			fmt.Printf("Warning: Setting %s has children but type was %s. Forcing to ListNestedAttribute.\n", settingID, attr.Type)
			attr.Type = AttrTypeListNested
		}

		// Build map to pass requirement info down to direct children based on 'DependedOnBy'
		childRequiredInfo := make(map[string]bool)
		for _, dep := range setting.DependedOnBy {
			if childDepID := dep.DependedOnBy; childDepID != "" {
				childRequiredInfo[childDepID] = dep.Required
			}
		}

		// Recursively process children
		processedChildrenCount := 0
		for _, childID := range childIDs {
			// Check if child exists before recursing (should always exist due to preprocessing check)
			if _, childExists := allSettings[childID]; !childExists {
				continue
			}

			// Do not re-process if already processed globally in another branch (avoids redundant work, assumes tree structure)
			// If DAGs are possible, this might need adjustment or cycle detection becomes even more critical.
			if settingIsProcessed[childID] {
				// fmt.Printf("Debug: Child setting %s already processed globally, skipping recursive call from parent %s.\n", childID, settingID)
				continue
			}

			// fmt.Printf("Debug: Recursing into child %s from parent %s\n", childID, settingID)
			childAttr, ok := buildAttributeRecursively(childID, childRequiredInfo, visited)
			if ok {
				// Generate the HCL name for the child attribute *within this parent*
				childAttrName := generateNestedAttributeName(childAttr.SourceSetting.ID, childAttr.SourceSetting.Name)
				childAttr.Name = childAttrName // Set the name used in HCL/map key

				// Check for name collision within the *same* parent
				if _, exists := attr.NestedAttributes[childAttrName]; exists {
					fmt.Printf("Error: Nested attribute name collision for '%s' under parent '%s' (from child ID '%s'). Skipping duplicate.\n", childAttrName, settingID, childID)
				} else {
					attr.NestedAttributes[childAttrName] = childAttr
					processedChildrenCount++
					settingIsProcessed[childID] = true // Mark child as processed globally *after* successful recursive build
				}
			} else {
				fmt.Printf("Warning: Failed to build child attribute for ID %s (child of %s).\n", childID, settingID)
			}
		}
		// fmt.Printf("Debug: Processed %d children for parent %s\n", processedChildrenCount, settingID)

		// If it's a group type but no children were found/processed, log a warning.
		if isGroupType && processedChildrenCount == 0 && len(childIDs) == 0 {
			fmt.Printf("Warning: Setting %s is type %s but no children were found or processed.\n", settingID, attr.Type)
		}
	}

	// Mark current setting as processed *after* handling children.
	// Note: This happens even if recursive calls failed, to prevent re-processing from other roots.
	// Only top-level call should mark roots processed. This line might be redundant if roots handle it.
	// settingIsProcessed[settingID] = true // Let the caller (main loop) mark roots as processed.

	return attr, true
}

// --- Main Execution ---

func main() {
	// Seed random number generator (for fallback names)
	rand.Seed(time.Now().UnixNano())

	// Define command-line flags
	inputPath := flag.String("input", "", "Path to the JSON metadata file or directory")
	outputPath := flag.String("output", "generated_schema.go", "Path to the output Go file")
	resourceName := flag.String("resource", "DeviceManagementConfigurationResource", "Resource name for the schema struct (e.g., IosGeneralPolicyResource)")
	resourceType := flag.String("type", "intune_device_management_configuration", "Terraform resource type identifier (e.g., intune_ios_general_policy)")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Error: Input path (--input) is required.")
		flag.Usage()
		os.Exit(1)
	}

	// --- Read Input ---
	var combinedMetadata MetadataFile
	combinedMetadata.PlatformDefinitions = make(map[string][]Setting)

	info, err := os.Stat(*inputPath)
	if err != nil {
		fmt.Printf("Error accessing input path '%s': %v\n", *inputPath, err)
		os.Exit(1)
	}

	var filesToProcess []string
	if info.IsDir() {
		fmt.Printf("Info: Input path '%s' is a directory, processing JSON files within...\n", *inputPath)
		entries, err := os.ReadDir(*inputPath)
		if err != nil {
			fmt.Printf("Error reading directory '%s': %v\n", *inputPath, err)
			os.Exit(1)
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
				filesToProcess = append(filesToProcess, filepath.Join(*inputPath, entry.Name()))
			}
		}
		if len(filesToProcess) == 0 {
			fmt.Printf("Error: No JSON files found in directory '%s'.\n", *inputPath)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Info: Input path '%s' is a file.\n", *inputPath)
		filesToProcess = append(filesToProcess, *inputPath)
	}

	// --- Parse and Combine ---
	for _, file := range filesToProcess {
		fmt.Printf("--- Processing file: %s ---\n", file)
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file '%s': %v\n", file, err)
			continue // Skip this file
		}

		metadata, err := parseFlexible(data)
		if err != nil {
			fmt.Printf("Error parsing file '%s': %v\n", file, err)
			continue // Skip this file
		}

		// Merge metadata (simple merge, last file wins for top-level fields)
		if metadata.Metadata.Version != "" {
			combinedMetadata.Metadata = metadata.Metadata
		}
		// Merge platform definitions
		for platform, settings := range metadata.PlatformDefinitions {
			if _, exists := combinedMetadata.PlatformDefinitions[platform]; exists {
				fmt.Printf("Info: Merging settings for platform '%s' from file '%s'.\n", platform, file)
				combinedMetadata.PlatformDefinitions[platform] = append(combinedMetadata.PlatformDefinitions[platform], settings...)
			} else {
				combinedMetadata.PlatformDefinitions[platform] = settings
			}
		}
		fmt.Printf("--- Finished processing file: %s ---\n", file)
	}

	if len(combinedMetadata.PlatformDefinitions) == 0 {
		fmt.Println("Error: No settings could be loaded from the provided input path(s).")
		os.Exit(1)
	}

	fmt.Println("\n=== Starting Schema Generation ===")

	// Initialize global maps and preprocess all combined settings
	initializeAndPreprocess(combinedMetadata)

	// --- Build Schema ---
	templateData := SchemaTemplateData{
		ResourceName: *resourceName,
		ResourceType: *resourceType,
		Description:  fmt.Sprintf("Manages %s settings using the Intune Settings Catalog.", *resourceType), // More specific description
		AttributeDefinitions: []AttributeDefinition{
			// Add standard 'id' attribute first
			{
				Name:                         "id",
				Type:                         AttrTypeString,
				MarkdownDescription:          "The Terraform resource identifier.",
				Computed:                     true,
				Optional:                     false,
				Required:                     false,
				PlanModifiers:                []string{"stringplanmodifier.UseStateForUnknown()"},
				FormattedMarkdownDescription: `"The Terraform resource identifier."`, // Simple quote
			},
			// Example: Add other common top-level fields if your resource needs them
			// These are NOT generated from settings data but are part of the resource definition.
			/*
				{
					Name:                         "display_name",
					Type:                         AttrTypeString,
					MarkdownDescription:          "The display name of the policy.",
					Required:                     true,
					FormattedMarkdownDescription: `"The display name of the policy."`,
				},
				{
					Name:                         "description",
					Type:                         AttrTypeString,
					MarkdownDescription:          "The description of the policy.",
					Optional:                     true,
					FormattedMarkdownDescription: `"The description of the policy."`,
				},
			*/
		},
	}

	// Identify and process root settings (those not children of other settings)
	processedRootsCount := 0
	fmt.Println("Info: Identifying and processing root settings...")
	for id := range allSettings {
		if !settingIsChild[id] {
			// Check if already processed (e.g., if a setting somehow appeared as both root and child,
			// or if processing one root implicitly processed another due to shared non-root nodes - less likely with tree assumption)
			if settingIsProcessed[id] {
				// fmt.Printf("Debug: Skipping potential root %s as it was already processed.\n", id)
				continue
			}

			// fmt.Printf("Debug: Processing root setting: %s\n", id)
			// Empty map for requirement info as roots have no parent in this context
			// New visited map for each root path's cycle detection
			rootAttr, ok := buildAttributeRecursively(id, make(map[string]bool), make(map[string]bool))
			if ok {
				// Set the top-level HCL attribute name for the root
				rootAttrName := generateNestedAttributeName(rootAttr.SourceSetting.ID, rootAttr.SourceSetting.Name)
				rootAttr.Name = rootAttrName // Set the final HCL name

				// Check for top-level name collision (e.g., with 'id', 'display_name')
				collision := false
				for _, existingAttr := range templateData.AttributeDefinitions {
					if existingAttr.Name == rootAttrName {
						fmt.Printf("Error: Top-level attribute name collision for '%s' (from setting ID '%s'). Skipping this root setting.\n", rootAttrName, id)
						collision = true
						break
					}
				}
				if !collision {
					templateData.AttributeDefinitions = append(templateData.AttributeDefinitions, rootAttr)
					processedRootsCount++
					settingIsProcessed[id] = true // Mark root as processed globally
				}

			} else {
				fmt.Printf("Warning: Failed to build root attribute for ID %s.\n", id)
				// Mark as processed anyway to avoid retrying if it failed due to a cycle etc.
				settingIsProcessed[id] = true
			}
		}
	}

	fmt.Printf("Info: Processed %d root settings.\n", processedRootsCount)
	totalProcessedCount := 0
	for _, processed := range settingIsProcessed {
		if processed {
			totalProcessedCount++
		}
	}
	fmt.Printf("Info: Total settings incorporated into schema structure (including children): %d / %d unique settings.\n", totalProcessedCount, len(allSettings))
	unprocessedCount := len(allSettings) - totalProcessedCount
	if unprocessedCount > 0 {
		fmt.Printf("Warning: %d settings were not incorporated into the schema. They might be orphans or part of failed/skipped build paths.\n", unprocessedCount)
	}

	if len(templateData.AttributeDefinitions) <= 1 { // Only the manual 'id' is present
		fmt.Println("Error: No settings were successfully processed into the schema! Check input data and logs.")
		os.Exit(1)
	}

	// --- Generate Schema File ---
	fmt.Printf("Info: Generating Go schema file at: %s\n", *outputPath)
	err = generateSchemaFile(templateData, *outputPath)
	if err != nil {
		fmt.Printf("Error generating schema file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("=== Schema Generation Successful ===\nOutput written to: %s\n", *outputPath)
}

// --- Template Generation ---

// generateSchemaFile renders the Go schema file using templates, now with recursion
func generateSchemaFile(data SchemaTemplateData, outputPath string) error {
	// Create the output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory '%s': %w", dir, err)
		}
	}

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", outputPath, err)
	}
	defer file.Close()

	// Define the main template and the recursive part
	// Note: Added imports for List/Object validators/planmodifiers
	// Note: Commented out DefaultValue rendering as it needs DefaultProvider impl.
	tmplText := `package resource // Assuming 'resource' package, adjust if needed

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// {{.ResourceName}} represents the resource implementation.
// Ensure this struct exists elsewhere in your package with Create/Read/Update/Delete methods.
// type {{.ResourceName}} struct {
//	 client *GraphClient // Example client field
// }

// Ensure the implementation satisfies the expected interfaces.
// var _ resource.Resource = &{{.ResourceName}}{}
// var _ resource.ResourceWithImportState = &{{.ResourceName}}{} // If import is supported

// Generated with Intune Schema Generator
// 

// Defines the schema definition for the {{.ResourceName}} resource.
// Schema definition generated by 
{{ define "attributeDefinition" -}}
{{ .Type }}{
	{{- if .FormattedMarkdownDescription }}
	MarkdownDescription: {{ .FormattedMarkdownDescription }}, // 
	{{- end }}
	{{- if .Required }}
	Required: true, // 
	{{- else if .Optional }}
	Optional: true, // 
	{{- else if .Computed }}
	Computed: true, // 
	{{- end }}
	{{- if .Sensitive }}
	Sensitive: true, // 
	{{- end }}

	{{- /* Handling Nested Objects */}}
	{{- if and (or (eq .Type "schema.ListNestedAttribute") (eq .Type "schema.SetNestedAttribute") (eq .Type "schema.MapNestedAttribute") (eq .Type "schema.SingleNestedAttribute")) .NestedAttributes }}
	NestedObject: schema.NestedAttributeObject{ // 
		Attributes: map[string]schema.Attribute{
		{{- range $nestedName, $nestedAttr := .NestedAttributes }}
			"{{ $nestedName }}": {{ template "attributeDefinition" $nestedAttr }}, // 
		{{- end }}
		},
	},
	{{- /* Handling Simple Lists/Sets/Maps */}}
	{{- else if and (or (eq .Type "schema.ListAttribute") (eq .Type "schema.SetAttribute") (eq .Type "schema.MapAttribute")) .ElementType }}
	ElementType: types.{{ .ElementType }}, // 
	{{- end }}

	{{- /* Validators based on Type */}}
	{{- if and (eq .Type "schema.StringAttribute") .Validators }}
	Validators: []validator.String{ // 
		{{- range .Validators }}
		{{ . }},
		{{- end }}
	},
	{{- else if and (eq .Type "schema.BoolAttribute") .Validators }}
	Validators: []validator.Bool{ // 
		{{- range .Validators }}
		{{ . }},
		{{- end }}
	},
	{{- else if and (eq .Type "schema.Int64Attribute") .Validators }}
	Validators: []validator.Int64{ // 
		{{- range .Validators }}
		{{ . }},
		{{- end }}
	},
	{{- else if and (eq .Type "schema.Float64Attribute") .Validators }}
	Validators: []validator.Float64{ // 
		{{- range .Validators }}
		{{ . }},
		{{- end }}
	},
    {{- else if and (or (eq .Type "schema.ListAttribute") (eq .Type "schema.ListNestedAttribute")) .Validators }}
	Validators: []validator.List{ // 
	    {{- range .Validators }}
		{{ . }},
		{{- end }}
	},
    {{- else if and (eq .Type "schema.SingleNestedAttribute") .Validators }}
	Validators: []validator.Object{ // 
	    {{- range .Validators }}
		{{ . }},
		{{- end }}
	},
	{{- /* Add other validator types (Set, Map) if needed */}}
	{{- end }}

	{{- /* PlanModifiers based on Type */}}
	{{- if .PlanModifiers }}
		{{- if eq .Type "schema.StringAttribute" }}
	PlanModifiers: []planmodifier.String{ // 
		{{- range .PlanModifiers }}
		{{ . }},
		{{- end }}
	},
		{{- else if eq .Type "schema.BoolAttribute" }}
	PlanModifiers: []planmodifier.Bool{ // 
		{{- range .PlanModifiers }}
		{{ . }},
		{{- end }}
	},
		{{- else if eq .Type "schema.Int64Attribute" }}
	PlanModifiers: []planmodifier.Int64{ // 
		{{- range .PlanModifiers }}
		{{ . }},
		{{- end }}
	},
		{{- else if eq .Type "schema.Float64Attribute" }}
	PlanModifiers: []planmodifier.Float64{ // 
		{{- range .PlanModifiers }}
		{{ . }},
		{{- end }}
	},
    	{{- else if or (eq .Type "schema.ListAttribute") (eq .Type "schema.ListNestedAttribute") }}
	PlanModifiers: []planmodifier.List{ // 
	    {{- range .PlanModifiers }}
		{{ . }},
		{{- end }}
	},
    	{{- else if or (eq .Type "schema.SingleNestedAttribute") (eq .Type "schema.MapNestedAttribute") (eq .Type "schema.SetNestedAttribute") }} // Assuming Object plan modifier applies to these
	PlanModifiers: []planmodifier.Object{ // 
	    {{- range .PlanModifiers }}
		{{ . }},
		{{- end }}
	},
		{{- /* Add other planmodifier types (Set, Map) if needed */}}
		{{- end }}
	{{- end }} {{- /* End PlanModifiers block */}}


	{{- /* DefaultValue - Requires DefaultProvider implementation, commenting out direct value */}}
	{{- if .DefaultValue }}
	// DefaultValue: {{.DefaultValue}} // : Default value requires provider implementation (schema*_default)
	{{- end }}
}
{{- end }}


// Schema defines the schema for the {{.ResourceName}} resource.
func (r *{{.ResourceName}}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "{{.Description}}", // 
		Attributes: map[string]schema.Attribute{
		{{- range .AttributeDefinitions }}
			"{{ .Name }}": {{ template "attributeDefinition" . }}, // 
		{{- end }}
		},
	}
}

// Configure prepares the resource implementation client. (Example implementation)
// func (r *{{.ResourceName}}) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}
// 	client, ok := req.ProviderData.(*GraphClient) // Adjust type assertion as needed
// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Unexpected Resource Configure Type",
// 			fmt.Sprintf("Expected *GraphClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
// 		)
// 		return
// 	}
// 	r.client = client
// }

// Create creates the resource and sets the initial Terraform state. (Placeholder)
// func (r *{{.ResourceName}}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
//  // Implementation required
//	resp.Diagnostics.AddError("Create not implemented", "Resource creation logic is missing.")
// }

// Read refreshes the Terraform state with the latest data. (Placeholder)
// func (r *{{.ResourceName}}) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
//  // Implementation required
//	resp.Diagnostics.AddError("Read not implemented", "Resource read logic is missing.")
// }

// Update updates the resource and sets the updated Terraform state on success. (Placeholder)
// func (r *{{.ResourceName}}) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
//  // Implementation required
//	resp.Diagnostics.AddError("Update not implemented", "Resource update logic is missing.")
// }

// Delete deletes the resource and removes the Terraform state on success. (Placeholder)
// func (r *{{.ResourceName}}) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
//  // Implementation required
//	resp.Diagnostics.AddError("Delete not implemented", "Resource deletion logic is missing.")
// }

// ImportState imports the resource into Terraform state. (Placeholder)
// func (r *{{.ResourceName}}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
//	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }

`

	// Parse the template (including the named definition)
	tmpl, err := template.New("schema").Parse(tmplText)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Execute the template
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
