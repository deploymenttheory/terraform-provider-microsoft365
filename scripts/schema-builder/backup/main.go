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
)

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
	MinimumValue          int         `json:"minimumValue,omitempty"`
	MaximumValue          int         `json:"maximumValue,omitempty"`
}

type DependentOn struct {
	ParentSettingID string `json:"parentSettingId"`
	DependentOn     string `json:"dependentOn"`
}

type DependedOnBy struct {
	Required     bool   `json:"required"`
	DependedOnBy string `json:"dependedOnBy"`
}

// SchemaTemplateData holds structured data for the template
type SchemaTemplateData struct {
	ResourceName         string
	ResourceType         string
	Description          string
	AttributeDefinitions []AttributeDefinition
}

type AttributeDefinition struct {
	Name                         string
	Type                         string
	Description                  string
	Required                     bool
	Optional                     bool
	Computed                     bool
	Sensitive                    bool
	ElementType                  string
	NestedType                   string
	Validators                   []string
	PlanModifiers                []string
	DefaultValue                 string
	MarkdownDescription          string
	FormattedMarkdownDescription string
	ODataInfo                    ODataInfo
}

type ODataInfo struct {
	Type string
	ID   string
}

// Constants for attribute types
const (
	AttrTypeString       = "schema.StringAttribute"
	AttrTypeBool         = "schema.BoolAttribute"
	AttrTypeInt64        = "schema.Int64Attribute"
	AttrTypeFloat64      = "schema.Float64Attribute"
	AttrTypeList         = "schema.ListAttribute"
	AttrTypeListNested   = "schema.ListNestedAttribute"
	AttrTypeSet          = "schema.SetAttribute"
	AttrTypeSetNested    = "schema.SetNestedAttribute"
	AttrTypeMap          = "schema.MapAttribute"
	AttrTypeMapNested    = "schema.MapNestedAttribute"
	AttrTypeSingleNested = "schema.SingleNestedAttribute"
	AttrTypeObject       = "schema.ObjectAttribute"
)

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
		// If there are no platformDefinitions, check if it's a direct array of settings
		var settings []Setting
		if err := json.Unmarshal(data, &settings); err == nil && len(settings) > 0 {
			result.PlatformDefinitions = map[string][]Setting{
				"default": settings,
			}
		}
	}

	return result, nil
}

// cleanAttributeName creates a clean, valid Terraform attribute name
func cleanAttributeName(name string) string {
	// Replace dots with underscores
	result := strings.ReplaceAll(name, ".", "_")

	// Replace dashes with underscores
	result = strings.ReplaceAll(result, "-", "_")

	// Replace multiple consecutive underscores with a single underscore
	multipleUnderscores := regexp.MustCompile(`__+`)
	result = multipleUnderscores.ReplaceAllString(result, "_")

	// Trim leading and trailing underscores
	result = strings.Trim(result, "_")

	// If we end up with an empty string, return a generic name
	if result == "" {
		return "attribute"
	}

	return result
}

// convertToAttributeName converts CamelCase to snake_case
func convertToAttributeName(name string) string {
	// Convert CamelCase to snake_case
	result := ""
	for i, char := range name {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result += "_"
		}
		result += strings.ToLower(string(char))
	}
	return result
}

// convertSettingToAttribute converts a setting from the JSON metadata to a Terraform schema attribute
func convertSettingToAttribute(setting Setting) AttributeDefinition {
	// Convert the setting to a schema attribute
	attr := AttributeDefinition{
		Name:                cleanAttributeName(setting.ID),
		Description:         setting.Description,
		MarkdownDescription: setting.Description,
		ODataInfo: ODataInfo{
			Type: setting.ODataType,
			ID:   setting.ID,
		},
	}

	// Handle missing or empty ID
	if attr.Name == "" {
		// Use Name as a fallback
		if setting.Name != "" {
			attr.Name = cleanAttributeName(setting.Name)
		} else {
			// Generate a unique placeholder name
			attr.Name = fmt.Sprintf("setting_%d", rand.Intn(10000))
		}
		fmt.Printf("Warning: Setting has no ID, using '%s' instead\n", attr.Name)
	}

	// Enhance description with display name, keywords, and OData type info
	enhanceDescription(&attr, setting)

	// Determine attribute type based on setting type and value definition
	setAttributeType(&attr, setting)

	// Set default to optional
	attr.Optional = true

	// Add validators if needed
	addValidators(&attr, setting)

	// Set default value if available
	setDefaultValue(&attr, setting)

	return attr
}

// enhanceDescription adds display name, keywords, and OData type information to the description
// using newlines instead of parentheses
func enhanceDescription(attr *AttributeDefinition, setting Setting) {
	// Start with parts that will be joined with newlines
	var parts []string

	// Add the display name and original description as the first part
	firstPart := ""
	if setting.DisplayName != "" && setting.DisplayName != setting.Name {
		if attr.Description != "" {
			firstPart = fmt.Sprintf("%s - %s", setting.DisplayName, attr.Description)
		} else {
			firstPart = setting.DisplayName
		}
	} else {
		firstPart = attr.Description
	}

	// Only add non-empty first parts
	if firstPart != "" {
		parts = append(parts, firstPart)
	}

	// Add keywords if available
	if len(setting.Keywords) > 0 {
		keywordsStr := strings.Join(setting.Keywords, ", ")
		parts = append(parts, fmt.Sprintf("Keywords: %s", keywordsStr))
	}

	// Add OData type information
	if setting.ODataType != "" {
		parts = append(parts, fmt.Sprintf("OData Type: %s", setting.ODataType))
	}

	// Add setting ID
	if setting.ID != "" {
		parts = append(parts, fmt.Sprintf("Settings Catalog ID: %s", setting.ID))
	}

	// Join all parts with newlines for markdown format
	description := strings.Join(parts, "\n")

	// For the Go string literal format in the schema file, we need to format it
	// with string concatenation for multiline strings
	goFormattedDesc := ""
	lines := strings.Split(description, "\n")
	for i, line := range lines {
		if i == 0 {
			goFormattedDesc = fmt.Sprintf("%q", line)
		} else {
			goFormattedDesc = fmt.Sprintf("%s +\n\t\t\t\t%q", goFormattedDesc, line)
		}
	}

	// Update both description fields
	attr.Description = description
	attr.MarkdownDescription = description

	// Store the Go formatted version for template rendering
	attr.FormattedMarkdownDescription = goFormattedDesc
}

// setAttributeType determines and sets the proper attribute type based on the setting
func setAttributeType(attr *AttributeDefinition, setting Setting) {
	switch setting.ODataType {
	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition":
		// For boolean choices like true/false, use BoolAttribute when appropriate
		if len(setting.Options) == 2 {
			boolValues := map[string]bool{
				"true": true, "True": true, "Enabled": true, "enabled": true,
				"false": true, "False": true, "Disabled": true, "disabled": true,
			}

			// Check if this is a boolean choice by examining the options
			allBooleanOptions := true
			for _, opt := range setting.Options {
				if !boolValues[opt.DisplayName] {
					allBooleanOptions = false
					break
				}
			}

			// Additional check: if ID or name contains "enabled" or similar, likely boolean
			nameHintsBool := strings.Contains(strings.ToLower(setting.ID), "enabled") ||
				strings.Contains(strings.ToLower(setting.ID), "disabled") ||
				strings.Contains(strings.ToLower(setting.Name), "enabled") ||
				strings.Contains(strings.ToLower(setting.Name), "disabled")

			// If options look boolean OR name hints it's boolean
			if allBooleanOptions || nameHintsBool {
				attr.Type = AttrTypeBool
				return
			}
		}

		// Default to string for choice settings
		attr.Type = AttrTypeString

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingDefinition":
		if setting.ValueDefinition != nil {
			odataType := fmt.Sprintf("%v", setting.ValueDefinition.ODataType)

			if strings.Contains(odataType, "String") {
				attr.Type = AttrTypeString
				if setting.ValueDefinition.IsSecret {
					attr.Sensitive = true
				}
			} else if strings.Contains(odataType, "Integer") || strings.Contains(odataType, "Int") {
				attr.Type = AttrTypeInt64
			} else if strings.Contains(odataType, "Bool") {
				attr.Type = AttrTypeBool
			} else if strings.Contains(odataType, "Float") {
				attr.Type = AttrTypeFloat64
			} else {
				attr.Type = AttrTypeString // Default to string
				fmt.Printf("Warning: Unknown value definition type %s for setting %s, defaulting to StringAttribute\n",
					odataType, attr.Name)
			}
		} else {
			attr.Type = AttrTypeString // Default to string
			fmt.Printf("Warning: No value definition for setting %s, defaulting to StringAttribute\n", attr.Name)
		}

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionDefinition":
		// Handle collection definitions as list attributes
		attr.Type = AttrTypeList
		attr.ElementType = "StringType" // Default to string collections

		// If value definition exists and contains type information, refine the element type
		if setting.ValueDefinition != nil {
			odataType := fmt.Sprintf("%v", setting.ValueDefinition.ODataType)

			// Determine element type from OData type
			if strings.Contains(odataType, "Int") {
				attr.ElementType = "Int64Type"
			} else if strings.Contains(odataType, "Bool") {
				attr.ElementType = "BoolType"
			} else if strings.Contains(odataType, "Float") {
				attr.ElementType = "Float64Type"
			} else if strings.Contains(odataType, "Secret") {
				attr.ElementType = "StringType"
				attr.Sensitive = true
			}

			// Add descriptive information
			if setting.ValueDefinition.MaximumLength > 0 || setting.ValueDefinition.MinimumLength > 0 {
				lengthInfo := ""
				if setting.ValueDefinition.MinimumLength > 0 {
					lengthInfo += fmt.Sprintf(" (min length: %d", setting.ValueDefinition.MinimumLength)
					if setting.ValueDefinition.MaximumLength > 0 {
						lengthInfo += fmt.Sprintf(", max length: %d)", setting.ValueDefinition.MaximumLength)
					} else {
						lengthInfo += ")"
					}
				} else if setting.ValueDefinition.MaximumLength > 0 {
					lengthInfo += fmt.Sprintf(" (max length: %d)", setting.ValueDefinition.MaximumLength)
				}

				if attr.Description != "" {
					attr.Description += lengthInfo
				}
				if attr.MarkdownDescription != "" {
					attr.MarkdownDescription += lengthInfo
				}
			}
		}

	case "#microsoft.graph.deviceManagementConfigurationSettingGroupCollectionDefinition":
		attr.Type = AttrTypeListNested
		attr.NestedType = "ObjectType"

	case "#microsoft.graph.deviceManagementConfigurationSecretSettingDefinition":
		attr.Type = AttrTypeString
		attr.Sensitive = true

	default:
		attr.Type = AttrTypeString // Default to string
		fmt.Printf("Warning: Unknown setting type %s for setting %s, defaulting to StringAttribute\n",
			setting.ODataType, attr.Name)
	}
}

// addValidators adds appropriate validators for the attribute
func addValidators(attr *AttributeDefinition, setting Setting) {
	// For choice settings, add validators for the options
	if setting.ODataType == "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition" && len(setting.Options) > 0 {
		// Handle different validator types based on attribute type
		if attr.Type == AttrTypeBool {
			// If we've determined it's a boolean, we don't need validators
			return
		} else if attr.Type == AttrTypeString {
			// For string attributes, use OneOf validator
			options := make([]string, 0, len(setting.Options))
			for _, opt := range setting.Options {
				if opt.DisplayName != "" {
					// Use the original DisplayName as the validator option
					options = append(options, fmt.Sprintf("%q", opt.DisplayName))
				}
			}

			if len(options) > 0 {
				attr.Validators = []string{
					fmt.Sprintf("stringvalidator.OneOf(%s)", strings.Join(options, ", ")),
				}
			}
		}
	}

	// Add range validators for numeric types
	if (attr.Type == AttrTypeInt64 || attr.Type == AttrTypeFloat64) && setting.ValueDefinition != nil {
		// For integer attributes with min/max constraints
		if setting.ValueDefinition.MinimumValue != 0 || setting.ValueDefinition.MaximumValue != 0 {
			var validators []string

			if setting.ValueDefinition.MinimumValue != 0 {
				if attr.Type == AttrTypeInt64 {
					validators = append(validators, fmt.Sprintf("int64validator.AtLeast(%d)", setting.ValueDefinition.MinimumValue))
				} else {
					validators = append(validators, fmt.Sprintf("float64validator.AtLeast(%d)", setting.ValueDefinition.MinimumValue))
				}
			}

			if setting.ValueDefinition.MaximumValue != 0 {
				if attr.Type == AttrTypeInt64 {
					validators = append(validators, fmt.Sprintf("int64validator.AtMost(%d)", setting.ValueDefinition.MaximumValue))
				} else {
					validators = append(validators, fmt.Sprintf("float64validator.AtMost(%d)", setting.ValueDefinition.MaximumValue))
				}
			}

			if len(validators) > 0 {
				attr.Validators = validators
			}
		}
	}
}

// setDefaultValue sets the default value for the attribute if one is specified
func setDefaultValue(attr *AttributeDefinition, setting Setting) {
	if setting.DefaultValue == nil || setting.DefaultValue.Value == nil {
		return
	}

	var defaultVal string

	// Handle the default value based on the attribute type
	switch attr.Type {
	case AttrTypeBool:
		// Extract boolean default value
		boolVal, isBool := setting.DefaultValue.Value.(bool)
		if isBool {
			defaultVal = fmt.Sprintf("%v", boolVal)
		} else if strVal, isStr := setting.DefaultValue.Value.(string); isStr {
			// Handle string representations of booleans
			lowerVal := strings.ToLower(strVal)
			if lowerVal == "true" || lowerVal == "enabled" {
				defaultVal = "true"
			} else if lowerVal == "false" || lowerVal == "disabled" {
				defaultVal = "false"
			}
		}

	case AttrTypeInt64, AttrTypeFloat64:
		// Extract numeric default value
		if numVal, isNum := setting.DefaultValue.Value.(float64); isNum {
			defaultVal = fmt.Sprintf("%v", numVal)
		} else if strVal, isStr := setting.DefaultValue.Value.(string); isStr {
			// Try to parse string as number
			defaultVal = strVal
		}

	case AttrTypeString:
		// Extract string default value
		if strVal, isStr := setting.DefaultValue.Value.(string); isStr {
			defaultVal = fmt.Sprintf("%q", strVal)
		} else {
			// For other types, use fmt.Sprintf and wrap in quotes
			defaultVal = fmt.Sprintf("%q", fmt.Sprintf("%v", setting.DefaultValue.Value))
		}

	default:
		// For complex types or unhandled types, we skip default values
		return
	}

	// Only set non-empty, non-zero default values
	if defaultVal != "" && defaultVal != "\"\"" && defaultVal != "0" && defaultVal != "false" && defaultVal != "<nil>" {
		attr.DefaultValue = defaultVal
	}
}

// main entrypoint for the schema generator
// This tool parses Microsoft Graph API metadata files and generates Terraform Provider Framework schema code
// Usage: go run schema-generator.go --input=metadata.json --output=schema.go --resource=ResourceName --type=resourceType
func main() {
	// Define command-line flags
	inputPath := flag.String("input", "", "Path to the JSON metadata file")
	outputPath := flag.String("output", "schema.go", "Path to the output Go file")
	resourceName := flag.String("resource", "DeviceManagementConfigurationResource", "Resource name for the schema")
	resourceType := flag.String("type", "deviceManagementConfiguration", "Resource type identifier")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Error: Input file path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Read and parse the JSON file
	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// First try parsing as a raw map to explore the structure
	var rawData map[string]interface{}
	err = json.Unmarshal(data, &rawData)
	if err != nil {
		fmt.Printf("Error parsing JSON (initial exploration): %v\n", err)
		os.Exit(1)
	}

	// Print top-level keys to understand structure
	fmt.Println("Top-level JSON keys:")
	for k := range rawData {
		fmt.Printf("- %s\n", k)
	}

	// Now try to unmarshal into our expected structure
	var metadata MetadataFile
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		fmt.Printf("Error parsing JSON into structured format: %v\n", err)
		fmt.Println("Attempting to parse with a more flexible approach...")

		// If standard parsing fails, try a more flexible approach
		metadata, err = parseFlexible(data)
		if err != nil {
			fmt.Printf("Flexible parsing also failed: %v\n", err)
			os.Exit(1)
		}
	}

	// Generate schema definitions
	templateData := SchemaTemplateData{
		ResourceName: *resourceName,
		ResourceType: *resourceType,
		Description:  fmt.Sprintf("Manages %s settings in Microsoft Intune.", *resourceType),
		AttributeDefinitions: []AttributeDefinition{
			{
				Name:                "id",
				Type:                AttrTypeString,
				MarkdownDescription: "Unique identifier for this resource.",
				Computed:            true,
				PlanModifiers: []string{
					"stringplanmodifier.UseStateForUnknown()",
				},
			},
		},
	}

	// Process all platform definitions
	if len(metadata.PlatformDefinitions) == 0 {
		fmt.Println("Warning: No platform definitions found in the metadata")
	}

	settingCount := 0
	for platform, settings := range metadata.PlatformDefinitions {
		fmt.Printf("Processing platform: %s with %d settings\n", platform, len(settings))

		// Add settings to the schema
		for _, setting := range settings {
			// Get the attribute definition
			attrDef := convertSettingToAttribute(setting)
			templateData.AttributeDefinitions = append(templateData.AttributeDefinitions, attrDef)
			settingCount++
		}
	}

	if settingCount == 0 {
		fmt.Println("Warning: No settings were processed! The output schema will only contain the ID attribute.")
	} else {
		fmt.Printf("Successfully processed %d settings across all platforms\n", settingCount)
	}

	// Generate the schema file
	err = generateSchemaFile(templateData, *outputPath)
	if err != nil {
		fmt.Printf("Error generating schema file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated schema at: %s\n", *outputPath)
}

// Modified template rendering function to use the formatted description
func generateSchemaFile(data SchemaTemplateData, outputPath string) error {
	// Create the output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Define the template
	tmplText := `package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
)

// Generated with Microsoft Graph API Schema Generator
// This file was automatically generated - manual modifications may be overwritten.

func (r *{{.ResourceName}}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
			Description: "{{.Description}}",
			Attributes: map[string]schema.Attribute{
					{{- range .AttributeDefinitions}}
					"{{.Name}}": {{.Type}}{
							{{- if .FormattedMarkdownDescription}}
							MarkdownDescription: {{.FormattedMarkdownDescription}},
							{{- end}}
							{{- if .Required}}
							Required: true,
							{{- else if .Optional}}
							Optional: true,
							{{- else if .Computed}}
							Computed: true,
							{{- end}}
							{{- if .Sensitive}}
							Sensitive: true,
							{{- end}}
							
							{{- if eq .Type "schema.ListNestedAttribute" "schema.SetNestedAttribute" "schema.MapNestedAttribute" "schema.SingleNestedAttribute"}}
							NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
											// Nested attributes would be defined here
											// This is a placeholder that should be manually populated
											"name": schema.StringAttribute{
													Required: true,
													MarkdownDescription: "Name of the nested object",
											},
											"value": schema.StringAttribute{
													Required: true,
													MarkdownDescription: "Value of the nested object",
											},
									},
							},
							{{- else if and (eq .Type "schema.ListAttribute" "schema.SetAttribute" "schema.MapAttribute") .ElementType}}
							ElementType: types.{{.ElementType}},
							{{- end}}
							
							{{- if and (eq .Type "schema.StringAttribute") .Validators}}
							Validators: []validator.String{
									{{- range .Validators}}
									{{.}},
									{{- end}}
							},
							{{- else if and (eq .Type "schema.BoolAttribute") .Validators}}
							Validators: []validator.Bool{
									{{- range .Validators}}
									{{.}},
									{{- end}}
							},
							{{- else if and (eq .Type "schema.Int64Attribute") .Validators}}
							Validators: []validator.Int64{
									{{- range .Validators}}
									{{.}},
									{{- end}}
							},
							{{- else if and (eq .Type "schema.Float64Attribute") .Validators}}
							Validators: []validator.Float64{
									{{- range .Validators}}
									{{.}},
									{{- end}}
							},
							{{- end}}
							
							{{- if and (eq .Type "schema.StringAttribute") .PlanModifiers}}
							PlanModifiers: []planmodifier.String{
									{{- range .PlanModifiers}}
									{{.}},
									{{- end}}
							},
							{{- else if and (eq .Type "schema.BoolAttribute") .PlanModifiers}}
							PlanModifiers: []planmodifier.Bool{
									{{- range .PlanModifiers}}
									{{.}},
									{{- end}}
							},
							{{- else if and (eq .Type "schema.Int64Attribute") .PlanModifiers}}
							PlanModifiers: []planmodifier.Int64{
									{{- range .PlanModifiers}}
									{{.}},
									{{- end}}
							},
							{{- else if and (eq .Type "schema.Float64Attribute") .PlanModifiers}}
							PlanModifiers: []planmodifier.Float64{
									{{- range .PlanModifiers}}
									{{.}},
									{{- end}}
							},
							{{- end}}
							
							{{- if and (eq .Type "schema.StringAttribute") .DefaultValue}}
							DefaultValue: {{.DefaultValue}},
							{{- else if and (eq .Type "schema.BoolAttribute") .DefaultValue}}
							DefaultValue: {{.DefaultValue}},
							{{- else if and (eq .Type "schema.Int64Attribute") .DefaultValue}}
							DefaultValue: {{.DefaultValue}},
							{{- else if and (eq .Type "schema.Float64Attribute") .DefaultValue}}
							DefaultValue: {{.DefaultValue}},
							{{- end}}
					},
					{{- end}}
			},
	}
}
`

	// Parse and execute the template
	tmpl, err := template.New("schema").Parse(tmplText)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, data)
}
