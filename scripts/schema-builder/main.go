package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"
	// "github.com/hashicorp/terraform-plugin-framework/path" // No longer needed for dependency modifiers
)

// --- Constants for attribute types ---
const (
	AttrTypeString       = "schema.StringAttribute"
	AttrTypeBool         = "schema.BoolAttribute"
	AttrTypeInt64        = "schema.Int64Attribute"
	AttrTypeFloat64      = "schema.Float64Attribute"
	AttrTypeList         = "schema.ListAttribute"
	AttrTypeListNested   = "schema.ListNestedAttribute"
	AttrTypeSet          = "schema.SetAttribute"       // Note: Set requires specific element type handling
	AttrTypeSetNested    = "schema.SetNestedAttribute" // Note: Set requires specific element type handling
	AttrTypeMap          = "schema.MapAttribute"       // Note: Map requires specific element type handling
	AttrTypeMapNested    = "schema.MapNestedAttribute" // Note: Map requires specific element type handling
	AttrTypeSingleNested = "schema.SingleNestedAttribute"
	AttrTypeObject       = "schema.ObjectAttribute" // Generally used internally for nested structures
)

// --- Type definitions ---
// (Keep CommandLineOptions, MetadataFile, Metadata, PlatformGroup, Setting, SettingDefinition, Applicability, Option, DefaultValue, ValueDefinition, DependentOn, DependedOnBy as before)
type CommandLineOptions struct {
	inputPath    string
	outputPath   string
	resourceName string
	resourceType string
}
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
	Platforms interface{} `json:"platforms"`
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

type SchemaTemplateData struct {
	ResourceName         string
	ResourceType         string
	Description          string
	AttributeDefinitions []AttributeDefinition // Represents only TOP-LEVEL attributes in the schema
}
type AttributeDefinition struct {
	Name                         string
	Type                         string
	Description                  string // Internal use for building markdown
	Required                     bool   // Only relevant for NESTED attributes
	Optional                     bool   // Relevant for both top-level and nested
	Computed                     bool   // e.g., 'id'
	Sensitive                    bool
	ElementType                  string                         // For List/Set/Map types
	NestedAttributes             map[string]AttributeDefinition // Holds the schema for nested children
	Validators                   []string
	PlanModifiers                []string  // Mainly for 'id'
	DefaultValue                 string    // For documentation
	MarkdownDescription          string    // Raw markdown built internally
	FormattedMarkdownDescription string    // Formatted for Go template
	ODataInfo                    ODataInfo // Store original ID/Type
	SourceSetting                Setting   // Keep reference to original setting
}
type ODataInfo struct {
	Type string
	ID   string
}

// --- Global maps ---
var (
	allSettings     map[string]Setting  // Setting ID -> Setting
	settingChildren map[string][]string // Parent ID -> Child IDs (from DependedOnBy)
	settingIsChild  map[string]bool     // Child ID -> true (from DependentOn)
	// visited map is now local to recursion to handle diamond dependencies

	startsWithDigit = regexp.MustCompile(`^[0-9]`)
	goKeywords      = map[string]bool{ /* Keep keywords map */ }
)

// --- Main function ---
func main() {
	rand.Seed(time.Now().UnixNano())
	options := parseCommandLineOptions()
	metadata := readAndParseMetadata(options.inputPath)

	fmt.Println("\n=== Starting Schema Generation (Nested Structure) ===")
	schema := buildNestedSchema(metadata, options.resourceName, options.resourceType) // Use the nested builder

	fmt.Printf("Info: Generating Go schema file at: %s\n", options.outputPath)
	err := generateOutputFile(schema, options.outputPath)
	if err != nil {
		fmt.Printf("Error generating schema file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("=== Schema Generation Successful ===\nOutput written to: %s\n", options.outputPath)
}

// --- Core functions ---

// parseCommandLineOptions (Keep as is)
func parseCommandLineOptions() CommandLineOptions { /* ... implementation ... */
	inputPath := flag.String("input", "", "Path to the JSON metadata file or directory")
	outputPath := flag.String("output", "generated_schema.go", "Path to the output Go file")
	resourceName := flag.String("resource", "DeviceManagementConfigurationResource", "Resource name for the schema struct")
	resourceType := flag.String("type", "intune_device_management_configuration", "Terraform resource type identifier")
	flag.Parse()
	if *inputPath == "" {
		fmt.Println("Error: Input path (--input) is required.")
		flag.Usage()
		os.Exit(1)
	}
	if strings.TrimSpace(*outputPath) == "" {
		fmt.Println("Error: Output path (--output) cannot be empty.")
		os.Exit(1)
	}
	if !strings.HasSuffix(*outputPath, ".go") {
		fmt.Println("Warning: Output path does not end with .go.")
	}
	return CommandLineOptions{inputPath: *inputPath, outputPath: *outputPath, resourceName: *resourceName, resourceType: *resourceType}
}

// readAndParseMetadata (Keep as is)
func readAndParseMetadata(inputPath string) MetadataFile { /* ... implementation ... */
	var combinedMetadata MetadataFile
	combinedMetadata.PlatformDefinitions = make(map[string][]Setting)
	info, err := os.Stat(inputPath)
	if err != nil {
		fmt.Printf("Error accessing input path '%s': %v\n", inputPath, err)
		os.Exit(1)
	}
	var filesToProcess []string
	if info.IsDir() {
		fmt.Printf("Info: Input path '%s' is a directory, processing JSON files within...\n", inputPath)
		entries, err := os.ReadDir(inputPath)
		if err != nil {
			fmt.Printf("Error reading directory '%s': %v\n", inputPath, err)
			os.Exit(1)
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
				filesToProcess = append(filesToProcess, filepath.Join(inputPath, entry.Name()))
			}
		}
		if len(filesToProcess) == 0 {
			fmt.Printf("Error: No JSON files found in directory '%s'.\n", inputPath)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Info: Input path '%s' is a file.\n", inputPath)
		filesToProcess = append(filesToProcess, inputPath)
	}
	totalFilesProcessed := 0
	for _, file := range filesToProcess {
		fmt.Printf("--- Processing file: %s ---\n", file)
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file '%s': %v\n", file, err)
			continue
		}
		metadata, err := parseFlexible(data)
		if err != nil {
			fmt.Printf("Error parsing file '%s': %v\n", file, err)
			continue
		}
		if combinedMetadata.Metadata.Version == "" && metadata.Metadata.Version != "" {
			combinedMetadata.Metadata = metadata.Metadata
		}
		settingsAdded := 0
		for platform, settings := range metadata.PlatformDefinitions {
			if _, exists := combinedMetadata.PlatformDefinitions[platform]; exists {
				combinedMetadata.PlatformDefinitions[platform] = append(combinedMetadata.PlatformDefinitions[platform], settings...)
			} else {
				combinedMetadata.PlatformDefinitions[platform] = settings
			}
			settingsAdded += len(settings)
		}
		fmt.Printf("Info: Added %d settings from file '%s'.\n", settingsAdded, file)
		fmt.Printf("--- Finished processing file: %s ---\n", file)
		totalFilesProcessed++
	}
	if totalFilesProcessed == 0 {
		fmt.Println("Error: No files were successfully processed.")
		os.Exit(1)
	}
	if len(combinedMetadata.PlatformDefinitions) == 0 {
		fmt.Println("Error: No settings could be loaded.")
		os.Exit(1)
	}
	return combinedMetadata
}

// generateOutputFile (Keep as is)
func generateOutputFile(data SchemaTemplateData, outputPath string) error { /* ... implementation ... */
	dir := filepath.Dir(outputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory '%s': %w", dir, err)
		}
	}
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", outputPath, err)
	}
	defer file.Close()
	tmpl, err := template.New("schema").Funcs(template.FuncMap{
		"quote": func(s string) string { return fmt.Sprintf("%q", s) }, // Add quote function if needed in template directly
	}).Parse(schemaTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}
	return nil
}

// --- Schema Building Logic (Nested Approach) ---

// initializeAndPreprocess: Populates global maps including dependency info
func initializeAndPreprocess(metadata MetadataFile) {
	allSettings = make(map[string]Setting)
	settingChildren = make(map[string][]string) // parentID -> []childID
	settingIsChild = make(map[string]bool)      // childID -> true

	totalSettingsFound := 0
	uniqueSettingsAdded := 0
	for platform, settingsList := range metadata.PlatformDefinitions {
		fmt.Printf("Info: Reading %d settings from platform '%s'.\n", len(settingsList), platform)
		totalSettingsFound += len(settingsList)
		for i, setting := range settingsList {
			if setting.ID == "" {
				fmt.Printf("Warning: Skipping setting #%d from '%s' - missing ID.\n", i+1, platform)
				continue
			}
			if _, exists := allSettings[setting.ID]; !exists {
				allSettings[setting.ID] = setting
				uniqueSettingsAdded++
			} else {
				fmt.Printf("Debug: Duplicate ID '%s' ('%s'). Using first.\n", setting.ID, platform)
			}
		}
	}
	fmt.Printf("Info: Total settings found: %d.\n", totalSettingsFound)
	fmt.Printf("Info: Indexed %d unique settings.\n", len(allSettings))
	if len(allSettings) == 0 {
		fmt.Println("Error: No unique settings indexed.")
		os.Exit(1)
	}

	// Build parent->child and child->isChild maps
	parentChildLinks := 0
	childCount := 0
	for id, setting := range allSettings {
		// Build parent -> children map from DependedOnBy
		if len(setting.DependedOnBy) > 0 {
			childIDs := []string{}
			for _, dep := range setting.DependedOnBy {
				if childID := dep.DependedOnBy; childID != "" {
					if _, childExists := allSettings[childID]; childExists {
						childIDs = append(childIDs, childID)
						parentChildLinks++
					} else {
						fmt.Printf("Warning: Setting '%s' depended on by '%s', but child not found.\n", id, childID)
					}
				}
			}
			if len(childIDs) > 0 {
				settingChildren[id] = childIDs
			}
		}

		// Mark settings that declare a parent via DependentOn
		if len(setting.DependentOn) > 0 {
			isActuallyChild := false
			for _, dep := range setting.DependentOn {
				if parentID := dep.ParentSettingID; parentID != "" {
					if _, parentExists := allSettings[parentID]; parentExists {
						isActuallyChild = true // Mark as child if it declares any valid parent
					} else {
						fmt.Printf("Warning: Setting '%s' depends on parent '%s', but parent not found.\n", id, parentID)
					}
				}
			}
			if isActuallyChild {
				settingIsChild[id] = true
				childCount++
			}
		}
	}
	fmt.Printf("Info: Established %d parent->child links (from DependedOnBy).\n", parentChildLinks)
	fmt.Printf("Info: Identified %d settings that are children (based on DependentOn).\n", childCount)
}

// buildNestedSchema: Entry point for building the nested schema structure
func buildNestedSchema(metadata MetadataFile, resourceName, resourceType string) SchemaTemplateData {
	initializeAndPreprocess(metadata)

	templateData := SchemaTemplateData{
		ResourceName: resourceName,
		ResourceType: resourceType,
		Description:  fmt.Sprintf("Schema for managing %s settings using the Intune Settings Catalog API (Nested Structure).", resourceType),
		AttributeDefinitions: []AttributeDefinition{
			{
				Name:                         "id",
				Type:                         AttrTypeString,
				MarkdownDescription:          "The Terraform resource identifier.",
				Computed:                     true,
				FormattedMarkdownDescription: `"The Terraform resource identifier."`,
				PlanModifiers:                []string{"stringplanmodifier.UseStateForUnknown()"},
			},
		},
	}

	processedTopLevel := make(map[string]bool)
	nameCollisionCheck := make(map[string]bool) // Global check for all generated names

	fmt.Println("Info: Building nested schema structure...")
	for id := range allSettings {
		if !settingIsChild[id] { // Start recursion only for top-level settings
			if processedTopLevel[id] {
				continue
			}

			visited := make(map[string]bool) // Path-specific cycle detection

			attr, ok := buildAttributeRecursively(id, false, visited, nameCollisionCheck)
			if ok {
				// Top-level attributes in the resource schema are always optional
				attr.Optional = true
				attr.Required = false
				// Final description formatting for the top-level attribute
				attr.FormattedMarkdownDescription = formatDescriptionForGo(attr.MarkdownDescription)
				templateData.AttributeDefinitions = append(templateData.AttributeDefinitions, attr)
				processedTopLevel[id] = true
			} else {
				fmt.Printf("Error: Failed to build top-level attribute for setting ID %s.\n", id)
			}
		}
	}

	// Sort top-level attributes alphabetically (after 'id')
	sort.SliceStable(templateData.AttributeDefinitions[1:], func(i, j int) bool {
		return templateData.AttributeDefinitions[i+1].Name < templateData.AttributeDefinitions[j+1].Name
	})

	fmt.Printf("Info: Generated schema with %d top-level attributes (excluding 'id').\n", len(templateData.AttributeDefinitions)-1)
	if len(templateData.AttributeDefinitions) <= 1 {
		fmt.Println("Warning: No top-level settings found or processed. Check input data.")
	}

	return templateData
}

// buildAttributeRecursively: Creates an AttributeDefinition, handling nesting.
func buildAttributeRecursively(
	settingID string,
	isRequiredByParent bool,
	visited map[string]bool, // Path-specific cycle detection
	nameCollisionCheck map[string]bool, // Global name uniqueness check
) (AttributeDefinition, bool) {

	// Cycle detection
	if visited[settingID] {
		fmt.Printf("Error: Cycle detected: %s\n", settingID)
		return AttributeDefinition{}, false
	}
	visited[settingID] = true
	defer delete(visited, settingID)

	setting, exists := allSettings[settingID]
	if !exists {
		fmt.Printf("Error: Setting not found: %s\n", settingID)
		return AttributeDefinition{}, false
	}

	// Generate and resolve name
	baseAttrName := generateAttributeNameFromID(settingID)
	if baseAttrName == "" {
		fmt.Printf("Error: Failed name generation: %s\n", settingID)
		return AttributeDefinition{}, false
	}
	if baseAttrName == "id" {
		baseAttrName = "setting_id"
	}
	resolvedAttrName := resolveAttributeNameCollisionBasic(nameCollisionCheck, baseAttrName, settingID)
	if resolvedAttrName == "" {
		fmt.Printf("Error: Failed name collision resolution: %s\n", settingID)
		return AttributeDefinition{}, false
	}
	nameCollisionCheck[resolvedAttrName] = true

	attr := AttributeDefinition{
		Name:             resolvedAttrName,
		ODataInfo:        ODataInfo{Type: setting.ODataType, ID: setting.ID},
		SourceSetting:    setting,
		NestedAttributes: make(map[string]AttributeDefinition),
	}

	// Determine initial type
	setAttributeType(&attr, setting)
	if attr.Type == "" { // Check if skipped (e.g., Redirect type)
		return attr, true // Return true but with empty type, will be filtered out later
	}

	// --- Check for Schema Mismatch (Simple Type + Children) ---
	isSimpleType := attr.Type == AttrTypeString || attr.Type == AttrTypeBool || attr.Type == AttrTypeInt64 || attr.Type == AttrTypeFloat64 || attr.Type == AttrTypeList
	childIDs, hasChildren := settingChildren[settingID]
	isNestedType := attr.Type == AttrTypeListNested || attr.Type == AttrTypeSingleNested || attr.Type == AttrTypeSetNested || attr.Type == AttrTypeMapNested

	if isSimpleType && hasChildren {
		fmt.Printf("Warning: Setting %s (%s) identified as simple type %s but has children in DependedOnBy. Overriding type to SingleNestedAttribute due to children.\n", settingID, attr.Name, attr.Type)
		attr.Type = AttrTypeSingleNested // Prioritize structure over simple type if children exist
		isNestedType = true              // Update flag for child processing logic
		isSimpleType = false             // Update flag
		attr.ElementType = ""            // Clear element type if it was a List before override
	} else if !isNestedType && hasChildren {
		// This case remains an error - a non-simple, non-nested type shouldn't have children
		fmt.Printf("Error: Setting %s (%s) has unexpected type %s but has children. Cannot proceed.\n", settingID, attr.Name, attr.Type)
		return AttributeDefinition{}, false
	}
	// --- End Mismatch Check ---

	// Determine core properties (Validators, Default, Description)
	addValidators(&attr, setting)
	setDefaultValue(&attr, setting)
	enhanceDescription(&attr, setting)

	// Set Required/Optional based on parent
	attr.Required = isRequiredByParent
	attr.Optional = !isRequiredByParent

	// --- Handle Nested Children (if applicable) ---
	if isNestedType && hasChildren {
		fmt.Printf("Info: Processing children for nested setting %s (%s)\n", settingID, attr.Name)
		childRequiredMap := make(map[string]bool)
		for _, dep := range setting.DependedOnBy {
			if childDepID := dep.DependedOnBy; childDepID != "" {
				childRequiredMap[childDepID] = dep.Required
			}
		}

		nestedAttrMap := make(map[string]AttributeDefinition)
		for _, childID := range childIDs {
			childAttr, ok := buildAttributeRecursively(childID, childRequiredMap[childID], visited, nameCollisionCheck)
			if ok && childAttr.Type != "" { // Also check if child wasn't skipped
				if _, exists := nestedAttrMap[childAttr.Name]; exists {
					fmt.Printf("Error: Duplicate nested attribute name '%s' under parent '%s'. Skipping.\n", childAttr.Name, attr.Name)
				} else {
					childAttr.FormattedMarkdownDescription = formatDescriptionForGo(childAttr.MarkdownDescription)
					nestedAttrMap[childAttr.Name] = childAttr
				}
			} else if !ok { // Failure case
				fmt.Printf("Warning: Failed to build child attribute for ID %s (child of %s).\n", childID, settingID)
			} // Implicitly skips if childAttr.Type == ""
		}
		// Sort nested attributes alphabetically by name
		sortedNestedNames := make([]string, 0, len(nestedAttrMap))
		for name := range nestedAttrMap {
			sortedNestedNames = append(sortedNestedNames, name)
		}
		sort.Strings(sortedNestedNames)
		sortedNestedAttrMap := make(map[string]AttributeDefinition)
		for _, name := range sortedNestedNames {
			sortedNestedAttrMap[name] = nestedAttrMap[name]
		}
		attr.NestedAttributes = sortedNestedAttrMap // Assign sorted map

	} else if isNestedType && !hasChildren {
		fmt.Printf("Warning: Setting %s (%s) is type %s but has no children defined in DependedOnBy.\n", settingID, attr.Name, attr.Type)
	}
	// No error needed for !isNestedType && !hasChildren

	// Final FormattedMarkdownDescription set one level up
	return attr, true
}

// --- Helper Functions ---

// generateAttributeNameFromID (Keep as is)
func generateAttributeNameFromID(id string) string {
	if id == "" {
		return ""
	}
	name := strings.ReplaceAll(id, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	name = reg.ReplaceAllString(name, "")
	name = strings.Trim(name, "_")
	reg = regexp.MustCompile(`_+`)
	name = reg.ReplaceAllString(name, "_")
	if startsWithDigit.MatchString(name) {
		name = "setting_" + name
	}
	name = strings.ToLower(name)
	if goKeywords[name] {
		name = name + "_prop"
	}
	if name == "" {
		fmt.Printf("Warning: ID '%s' resulted in empty name.\n", id)
		return ""
	}
	return name
}

// resolveAttributeNameCollisionBasic (Keep as is)
func resolveAttributeNameCollisionBasic(existingNames map[string]bool, desiredName string, settingIdentifier string) string {
	if !existingNames[desiredName] {
		return desiredName
	}
	originalName := desiredName
	collisionCounter := 1
	for {
		newName := fmt.Sprintf("%s_%d", originalName, collisionCounter)
		if !existingNames[newName] {
			fmt.Printf("Warning: Name collision '%s' (from '%s'). Renamed to '%s'.\n", originalName, settingIdentifier, newName)
			return newName
		}
		collisionCounter++
		if collisionCounter > 100 {
			fmt.Printf("Error: Could not resolve collision for '%s' (from '%s').\n", originalName, settingIdentifier)
			return ""
		}
	}
}

// findOptionDisplayName (Keep as is)
func findOptionDisplayName(setting Setting, optionDetail string) string { /* ... implementation ... */
	for _, opt := range setting.Options {
		optIDStr := fmt.Sprintf("%v", opt.OptionID)
		if optIDStr != "<nil>" && optIDStr == optionDetail {
			return opt.DisplayName
		}
	}
	for _, opt := range setting.Options {
		optIDStr := fmt.Sprintf("%v", opt.OptionID)
		if optIDStr != "<nil>" && strings.HasSuffix(optIDStr, "_"+optionDetail) {
			return opt.DisplayName
		}
	}
	if setting.SettingDefinition.DefaultOptionID != "" && setting.SettingDefinition.DefaultOptionID == optionDetail {
		for _, opt := range setting.Options {
			optIDStr := fmt.Sprintf("%v", opt.OptionID)
			if optIDStr == setting.SettingDefinition.DefaultOptionID {
				return opt.DisplayName
			}
		}
	}
	return optionDetail
}

// formatDescriptionForGo (Keep as is)
func formatDescriptionForGo(markdownDesc string) string { /* ... implementation ... */
	goFormattedDesc := ""
	goLines := strings.Split(markdownDesc, "\\n\\n")
	isFirstLine := true
	for _, line := range goLines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		processedLine := strings.ReplaceAll(trimmedLine, "\"", "\\\"")
		processedLine = strings.ReplaceAll(processedLine, "`", "\\`")
		quotedLine := fmt.Sprintf("%q", processedLine)
		if isFirstLine {
			goFormattedDesc = quotedLine
			isFirstLine = false
		} else {
			goFormattedDesc = goFormattedDesc + " + \"\\n\\n\" + \n\t\t\t\t" + quotedLine
		}
	}
	if goFormattedDesc == "" {
		return `""`
	}
	return goFormattedDesc
}

// parseFlexible (Keep as is)
func parseFlexible(data []byte) (MetadataFile, error) { /* ... implementation ... */
	var result MetadataFile
	result.PlatformDefinitions = make(map[string][]Setting)
	err := json.Unmarshal(data, &result)
	if err == nil && (len(result.PlatformDefinitions) > 0 || result.Metadata.Version != "") {
		validPlatformDefs := false
		if len(result.PlatformDefinitions) > 0 {
			for _, settings := range result.PlatformDefinitions {
				if len(settings) > 0 && settings[0].ID != "" {
					validPlatformDefs = true
					break
				}
			}
		}
		if validPlatformDefs || result.Metadata.Version != "" {
			fmt.Println("Info: Parsed standard MetadataFile structure.")
			return result, nil
		} else {
			fmt.Println("Debug: Initial parse ok but no valid data. Trying alternatives.")
			result.PlatformDefinitions = make(map[string][]Setting)
		}
	}
	if err != nil {
		fmt.Printf("Debug: Initial parse failed: %v. Trying alternatives.\n", err)
	}
	var rawData map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawData); err != nil {
		return result, fmt.Errorf("error parsing JSON into raw map: %w", err)
	}
	if metadataRaw, ok := rawData["metadata"]; ok && result.Metadata.Version == "" {
		var metadata Metadata
		if err := json.Unmarshal(metadataRaw, &metadata); err == nil {
			result.Metadata = metadata
		} else {
			fmt.Printf("Warning: Found 'metadata' key but failed to parse: %v\n", err)
		}
	}
	if platformDefsRaw, ok := rawData["platformDefinitions"]; ok {
		var platformDefs map[string][]Setting
		if err := json.Unmarshal(platformDefsRaw, &platformDefs); err == nil {
			if result.PlatformDefinitions == nil {
				result.PlatformDefinitions = make(map[string][]Setting)
			}
			for key, val := range platformDefs {
				result.PlatformDefinitions[key] = val
			}
			fmt.Println("Info: Found 'platformDefinitions' object.")
		} else {
			fmt.Printf("Warning: Found 'platformDefinitions' key but failed parse: %v\n", err)
		}
	}
	if valueRaw, ok := rawData["value"]; ok {
		var settingsInValue []Setting
		if err := json.Unmarshal(valueRaw, &settingsInValue); err == nil && len(settingsInValue) > 0 {
			fmt.Println("Info: Found settings under 'value' key. Merging into 'defaultPlatform'.")
			if _, exists := result.PlatformDefinitions["defaultPlatform"]; exists {
				result.PlatformDefinitions["defaultPlatform"] = append(result.PlatformDefinitions["defaultPlatform"], settingsInValue...)
			} else {
				result.PlatformDefinitions["defaultPlatform"] = settingsInValue
			}
		} else if err != nil {
			fmt.Printf("Warning: Found 'value' key but failed parse: %v\n", err)
		}
	}
	var settingsArray []Setting
	if err := json.Unmarshal(data, &settingsArray); err == nil && len(settingsArray) > 0 {
		fmt.Println("Info: Found direct array of settings. Merging into 'defaultPlatform'.")
		if _, exists := result.PlatformDefinitions["defaultPlatform"]; exists {
			result.PlatformDefinitions["defaultPlatform"] = append(result.PlatformDefinitions["defaultPlatform"], settingsArray...)
		} else {
			result.PlatformDefinitions["defaultPlatform"] = settingsArray
		}
	}
	if len(result.PlatformDefinitions) == 0 {
		fmt.Println("Warning: Could not find recognizable settings data structure.")
	}
	return result, nil
}

// setAttributeType (Keep as is - primarily used by buildAttributeRecursively)
// setAttributeType determines the schema type based on ODataType and ValueDefinition
func setAttributeType(attr *AttributeDefinition, setting Setting) {
	attr.Type = AttrTypeString // Default assumption
	attr.Sensitive = false
	attr.ElementType = "" // Reset element type

	switch setting.ODataType {

	// --- Skip Redirect Type ---
	case "#microsoft.graph.deviceManagementConfigurationRedirectSettingDefinition":
		fmt.Printf("Info: Skipping RedirectSettingDefinition: %s (%s) - Configure this via its dedicated resource type.\n", setting.ID, attr.Name)
		attr.Type = "" // Mark as invalid/to be skipped
		return

	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition":
		if len(setting.Options) == 2 { // Check for boolean-like choices
			isBoolLike := false
			optionNames := make(map[string]bool)
			for _, opt := range setting.Options {
				optionNames[strings.ToLower(opt.DisplayName)] = true
			}
			// Add more robust checks if needed (e.g., based on values, defaultOptionId suffix)
			if (optionNames["true"] && optionNames["false"]) || (optionNames["enabled"] && optionNames["disabled"]) || (optionNames["yes"] && optionNames["no"]) || (optionNames["allow"] && optionNames["block"]) {
				isBoolLike = true
			}
			if isBoolLike {
				attr.Type = AttrTypeBool
				return // It's a boolean
			}
		}
		attr.Type = AttrTypeString // Default for non-boolean choice

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingDefinition":
		if setting.ValueDefinition != nil {
			valDefType := fmt.Sprintf("%v", setting.ValueDefinition.ODataType)
			if strings.Contains(valDefType, "String") {
				attr.Type = AttrTypeString
				if setting.ValueDefinition.IsSecret {
					attr.Sensitive = true
				}
			} else if strings.Contains(valDefType, "Integer") || strings.Contains(valDefType, "Int64") {
				attr.Type = AttrTypeInt64
			} else if strings.Contains(valDefType, "Boolean") || strings.Contains(valDefType, "Bool") {
				attr.Type = AttrTypeBool
			} else if strings.Contains(valDefType, "Double") || strings.Contains(valDefType, "Float") {
				attr.Type = AttrTypeFloat64
			} else {
				fmt.Printf("Warning: Unknown ValueDef Type '%s' for simple setting %s (%s), defaulting to String.\n", valDefType, setting.ID, attr.Name)
				attr.Type = AttrTypeString
			}
		} else {
			fmt.Printf("Warning: Missing ValueDef for simple setting %s (%s), defaulting to String.\n", setting.ID, attr.Name)
			attr.Type = AttrTypeString
		}

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionDefinition":
		attr.Type = AttrTypeList
		attr.ElementType = "StringType" // Default, determine based on ValueDefinition
		if setting.ValueDefinition != nil {
			valDefType := fmt.Sprintf("%v", setting.ValueDefinition.ODataType)
			if strings.Contains(valDefType, "Integer") || strings.Contains(valDefType, "Int64") {
				attr.ElementType = "Int64Type"
			} else if strings.Contains(valDefType, "Boolean") || strings.Contains(valDefType, "Bool") {
				attr.ElementType = "BoolType"
			} else if strings.Contains(valDefType, "Double") || strings.Contains(valDefType, "Float") {
				attr.ElementType = "Float64Type"
			} else if strings.Contains(valDefType, "String") {
				attr.ElementType = "StringType"
				if setting.ValueDefinition.IsSecret {
					attr.Sensitive = true
				} // Mark list sensitive if elements are
			} else {
				fmt.Printf("Warning: Unknown ValueDef Type '%s' for simple collection %s (%s), defaulting to StringType.\n", valDefType, setting.ID, attr.Name)
			}
		} else {
			fmt.Printf("Warning: Missing ValueDef for simple collection %s (%s), defaulting to StringType.\n", setting.ID, attr.Name)
		}

	// --- Handle Choice Collection ---
	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionDefinition":
		attr.Type = AttrTypeList
		attr.ElementType = "StringType" // Represent list of selected choices as strings
		fmt.Printf("Info: %s (%s) -> ListAttribute of StringType (Choice Collection)\n", setting.ID, attr.Name)

	case "#microsoft.graph.deviceManagementConfigurationSettingGroupCollectionDefinition", "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstanceTemplate":
		attr.Type = AttrTypeListNested // Collections of groups are lists of nested objects

	case "#microsoft.graph.deviceManagementConfigurationGroupSettingInstanceTemplate", "#microsoft.graph.deviceManagementConfigurationSettingGroupDefinition":
		// Individual group instances are single nested objects
		// Check if it *also* has children defined - if so, it acts nested.
		// If not, it might be a simple grouping label, potentially skip or handle differently?
		// For now, assume it implies a nested structure if the type indicates a group.
		_, hasChildren := settingChildren[setting.ID]
		if hasChildren {
			attr.Type = AttrTypeSingleNested
		} else {
			// If it's a group definition but has NO children according to dependedOnBy,
			// it might be a category header or an unused definition. Let's treat it as potentially skippable
			// or maybe a simple string if it has options/valueDef? For now, lean towards SingleNested if it's a group type.
			attr.Type = AttrTypeSingleNested
			fmt.Printf("Warning: Setting %s (%s) is Group type %s but has no children in dependedOnBy. Treating as SingleNestedObject.\n", setting.ID, attr.Name, setting.ODataType)
		}

	case "#microsoft.graph.deviceManagementConfigurationSecretSettingDefinition":
		attr.Type = AttrTypeString
		attr.Sensitive = true

	default:
		fmt.Printf("Warning: Unhandled ODataType '%s' for %s (%s), defaulting to String.\n", setting.ODataType, setting.ID, attr.Name)
		attr.Type = AttrTypeString
	}
}

// addValidators (Keep as is - primarily used by buildAttributeRecursively)
func addValidators(attr *AttributeDefinition, setting Setting) { /* ... implementation ... */
	if attr.Type == AttrTypeString && setting.ODataType == "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition" && len(setting.Options) > 0 {
		isBoolLike := false
		if len(setting.Options) == 2 {
			optionNames := make(map[string]bool)
			for _, opt := range setting.Options {
				optionNames[strings.ToLower(opt.DisplayName)] = true
			}
			if (optionNames["true"] && optionNames["false"]) || (optionNames["enabled"] && optionNames["disabled"]) {
				isBoolLike = true
			}
		}
		if !isBoolLike {
			var options []string
			for _, opt := range setting.Options {
				options = append(options, fmt.Sprintf("%q", opt.DisplayName))
			}
			if len(options) > 0 {
				attr.Validators = append(attr.Validators, fmt.Sprintf("stringvalidator.OneOf(%s)", strings.Join(options, ", ")))
			}
		}
	}
	if setting.ValueDefinition != nil {
		if attr.Type == AttrTypeInt64 {
			if setting.ValueDefinition.MinimumValue != 0 {
				attr.Validators = append(attr.Validators, fmt.Sprintf("int64validator.AtLeast(%d)", setting.ValueDefinition.MinimumValue))
			}
			if setting.ValueDefinition.MaximumValue != 0 {
				attr.Validators = append(attr.Validators, fmt.Sprintf("int64validator.AtMost(%d)", setting.ValueDefinition.MaximumValue))
			}
		} else if attr.Type == AttrTypeFloat64 {
			if setting.ValueDefinition.MinimumValue != 0 {
				attr.Validators = append(attr.Validators, fmt.Sprintf("float64validator.AtLeast(%f)", float64(setting.ValueDefinition.MinimumValue)))
			}
			if setting.ValueDefinition.MaximumValue != 0 {
				attr.Validators = append(attr.Validators, fmt.Sprintf("float64validator.AtMost(%f)", float64(setting.ValueDefinition.MaximumValue)))
			}
		}
	}
	if attr.Type == AttrTypeString && setting.ValueDefinition != nil {
		if setting.ValueDefinition.MinimumLength > 0 {
			attr.Validators = append(attr.Validators, fmt.Sprintf("stringvalidator.LengthAtLeast(%d)", setting.ValueDefinition.MinimumLength))
		}
		if setting.ValueDefinition.MaximumLength > 0 {
			attr.Validators = append(attr.Validators, fmt.Sprintf("stringvalidator.LengthAtMost(%d)", setting.ValueDefinition.MaximumLength))
		}
	}
}

// setDefaultValue (Keep as is - primarily used by buildAttributeRecursively)
func setDefaultValue(attr *AttributeDefinition, setting Setting) { /* ... implementation ... */
	if setting.DefaultValue == nil || setting.DefaultValue.Value == nil {
		attr.DefaultValue = ""
		return
	}
	var defaultValStr string
	switch attr.Type {
	case AttrTypeBool:
		if boolVal, ok := setting.DefaultValue.Value.(bool); ok {
			defaultValStr = fmt.Sprintf("%t", boolVal)
		} else if strVal, ok := setting.DefaultValue.Value.(string); ok {
			lowerVal := strings.ToLower(strVal)
			if lowerVal == "true" || lowerVal == "enabled" || lowerVal == "1" {
				defaultValStr = "true"
			} else if lowerVal == "false" || lowerVal == "disabled" || lowerVal == "0" {
				defaultValStr = "false"
			}
		}
	case AttrTypeInt64:
		if numVal, ok := setting.DefaultValue.Value.(float64); ok {
			defaultValStr = fmt.Sprintf("%d", int64(numVal))
		} else if intVal, ok := setting.DefaultValue.Value.(int); ok {
			defaultValStr = fmt.Sprintf("%d", int64(intVal))
		} else if int64Val, ok := setting.DefaultValue.Value.(int64); ok {
			defaultValStr = fmt.Sprintf("%d", int64Val)
		} else if strVal, ok := setting.DefaultValue.Value.(string); ok {
			defaultValStr = strVal
		}
	case AttrTypeFloat64:
		if numVal, ok := setting.DefaultValue.Value.(float64); ok {
			defaultValStr = fmt.Sprintf("%g", numVal)
		} else if strVal, ok := setting.DefaultValue.Value.(string); ok {
			defaultValStr = strVal
		}
	case AttrTypeString:
		val := fmt.Sprintf("%v", setting.DefaultValue.Value)
		defaultValStr = val
		if setting.ODataType == "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition" && setting.SettingDefinition.DefaultOptionID != "" {
			for _, opt := range setting.Options {
				if fmt.Sprintf("%v", opt.OptionID) == setting.SettingDefinition.DefaultOptionID {
					defaultValStr = opt.DisplayName
					break
				}
			}
		}
		if defaultValStr == "<nil>" {
			defaultValStr = ""
		}
		defaultValStr = strings.ReplaceAll(defaultValStr, "`", "\\`")
		defaultValStr = strings.ReplaceAll(defaultValStr, "\"", "\\\"")
		if defaultValStr == "" {
			defaultValStr = `""`
		}
	case AttrTypeList, AttrTypeListNested, AttrTypeSet, AttrTypeSetNested, AttrTypeMap, AttrTypeMapNested, AttrTypeSingleNested:
		jsonBytes, err := json.Marshal(setting.DefaultValue.Value)
		if err == nil {
			defaultValStr = string(jsonBytes)
		} else {
			defaultValStr = fmt.Sprintf("%v", setting.DefaultValue.Value)
		}
		defaultValStr = strings.ReplaceAll(defaultValStr, "`", "\\`")
		defaultValStr = strings.ReplaceAll(defaultValStr, "\"", "\\\"")
		if defaultValStr == "[]" || defaultValStr == "{}" {
			defaultValStr = ""
		}
	default:
		defaultValStr = fmt.Sprintf("%v", setting.DefaultValue.Value)
		defaultValStr = strings.ReplaceAll(defaultValStr, "`", "\\`")
		defaultValStr = strings.ReplaceAll(defaultValStr, "\"", "\\\"")
		if defaultValStr == "<nil>" {
			defaultValStr = ""
		}
	}
	meaningful := defaultValStr != "" && defaultValStr != `""` && defaultValStr != "<nil>" && !(attr.Type == AttrTypeInt64 && defaultValStr == "0") && !(attr.Type == AttrTypeFloat64 && defaultValStr == "0") && !(attr.Type == AttrTypeBool && defaultValStr == "false")
	if meaningful {
		attr.DefaultValue = defaultValStr
	} else {
		attr.DefaultValue = ""
	}
}

// enhanceDescription (Keep as is - builds raw markdown)
func enhanceDescription(attr *AttributeDefinition, setting Setting) { /* ... implementation ... */
	var parts []string
	firstPart := ""
	cleanDesc := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(setting.Description, "\n", " "), "\r", ""))
	cleanDispName := strings.TrimSpace(setting.DisplayName)
	if cleanDispName != "" && cleanDispName != setting.Name {
		firstPart = cleanDispName
		if cleanDesc != "" {
			firstPart += " - " + cleanDesc
		}
	} else if cleanDesc != "" {
		firstPart = cleanDesc
	} else {
		firstPart = setting.Name
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
	if setting.ValueDefinition != nil {
		var constraints []string
		if setting.ValueDefinition.MinimumLength > 0 {
			constraints = append(constraints, fmt.Sprintf("min length: %d", setting.ValueDefinition.MinimumLength))
		}
		if setting.ValueDefinition.MaximumLength > 0 {
			constraints = append(constraints, fmt.Sprintf("max length: %d", setting.ValueDefinition.MaximumLength))
		}
		if attr.Type == AttrTypeInt64 || attr.Type == AttrTypeFloat64 {
			if setting.ValueDefinition.MinimumValue != 0 {
				constraints = append(constraints, fmt.Sprintf("min value: %d", setting.ValueDefinition.MinimumValue))
			}
			if setting.ValueDefinition.MaximumValue != 0 {
				constraints = append(constraints, fmt.Sprintf("max value: %d", setting.ValueDefinition.MaximumValue))
			}
		}
		if len(constraints) > 0 {
			parts = append(parts, fmt.Sprintf("Constraints: %s", strings.Join(constraints, ", ")))
		}
	}
	if attr.DefaultValue != "" {
		if strings.ContainsAny(attr.DefaultValue, "[{") {
			parts = append(parts, fmt.Sprintf("Default value: %s", attr.DefaultValue))
		} else {
			parts = append(parts, fmt.Sprintf("Default value: `%s`", attr.DefaultValue))
		}
	}
	attr.MarkdownDescription = strings.Join(parts, "\\n\\n")
	attr.FormattedMarkdownDescription = "" // Reset, formatted when finalized
}

// --- Template definition (Keep as is from previous response) ---
const schemaTemplate = `package resource

	import (
		"context"
	
		"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
		"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
		"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
		"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
		"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
		// "github.com/hashicorp/terraform-plugin-framework/path" // No longer needed for dependency modifiers
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

		// planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers" // No longer needed for dependency
	)

	// {{.ResourceName}} represents the resource implementation.
	// type {{.ResourceName}} struct { /* ... */ }

	// Generated with Intune Schema Generator - Nested Structure - DO NOT EDIT MANUALLY

	{{ define "attributeDefinition" -}}
	{{ .Type }}{
		{{- if .FormattedMarkdownDescription }}
		MarkdownDescription: {{ .FormattedMarkdownDescription }},
		{{- else }}
		MarkdownDescription: "",
		{{- end }}
		{{- if .Required }}
		Required: true,
		{{- else if .Optional }}
		Optional: true,
		{{- else if .Computed }}
		Computed: true,
		{{- end }}
		{{- if .Sensitive }}
		Sensitive: true,
		{{- end }}

		{{- if or (eq .Type "schema.ListNestedAttribute") (eq .Type "schema.SetNestedAttribute") (eq .Type "schema.MapNestedAttribute") (eq .Type "schema.SingleNestedAttribute") }}
			{{- if .NestedAttributes }}
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
			{{- range $nestedName, $nestedAttr := .NestedAttributes }}
				"{{ $nestedName }}": {{ template "attributeDefinition" $nestedAttr }},
			{{- end }}
			},
		},
			{{- else }}
		// NestedObject: schema.NestedAttributeObject{ Attributes: map[string]schema.Attribute{} }, // No children generated or defined
			{{- end }}
		{{- else if and (or (eq .Type "schema.ListAttribute") (eq .Type "schema.SetAttribute") (eq .Type "schema.MapAttribute")) .ElementType }}
		ElementType: types.{{ .ElementType }},
		{{- end }}

		{{- if .Validators }}
			{{- if eq .Type "schema.StringAttribute" }}
		Validators: []validator.String{ {{- range .Validators }} {{ . }}, {{- end }} },
			{{- else if eq .Type "schema.BoolAttribute" }}
		Validators: []validator.Bool{ {{- range .Validators }} {{ . }}, {{- end }} },
			{{- else if eq .Type "schema.Int64Attribute" }}
		Validators: []validator.Int64{ {{- range .Validators }} {{ . }}, {{- end }} },
			{{- else if eq .Type "schema.Float64Attribute" }}
		Validators: []validator.Float64{ {{- range .Validators }} {{ . }}, {{- end }} },
			{{- else if or (eq .Type "schema.ListAttribute") (eq .Type "schema.ListNestedAttribute") }}
		Validators: []validator.List{ {{- range .Validators }} {{ . }}, {{- end }} },
			{{- else if or (eq .Type "schema.SetAttribute") (eq .Type "schema.SetNestedAttribute") }}
		Validators: []validator.Set{ {{- range .Validators }} {{ . }}, {{- end }} },
			{{- else if or (eq .Type "schema.MapAttribute") (eq .Type "schema.MapNestedAttribute") }}
		Validators: []validator.Map{ {{- range .Validators }} {{ . }}, {{- end }} },
			{{- else if eq .Type "schema.SingleNestedAttribute" }}
		Validators: []validator.Object{ {{- range .Validators }} {{ . }}, {{- end }} },
			{{- end }}
		{{- end }}

		{{- /* PlanModifiers only needed for specific cases like 'id' */}}
		{{- if .PlanModifiers }}
			{{- if eq .Type "schema.StringAttribute" }}
		PlanModifiers: []planmodifier.String{ {{- range .PlanModifiers }} {{ . }}, {{- end }} },
			{{- /* Add other types if needed, but dependency modifiers are gone */}}
			{{- end }}
		{{- end }}

		{{- if .DefaultValue }}
		// DefaultValue: /* Requires type-specific default implementation */ // Documented default: {{.DefaultValue}}
		{{- end }}
	}
	{{- end }}

	// Schema defines the schema for the {{.ResourceName}} resource.
	func (r *{{.ResourceName}}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
		resp.Schema = schema.Schema{
			MarkdownDescription: {{ quote .Description }},
			Attributes: map[string]schema.Attribute{
			{{- /* Iterate only over top-level attributes defined in templateData */}}
			{{- range .AttributeDefinitions }}
				"{{ .Name }}": {{ template "attributeDefinition" . }},
			{{- end }}
			},
		}
	}

	// Provider Configure, Create, Read, Update, Delete, ImportState methods are required.
	`
