package graphBetaGroupPolicyCategories

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read refreshes the Terraform state with the latest data.
func (d *GroupPolicyCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupPolicyCategoryDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the setting name from the configuration
	settingName := data.SettingName.ValueString()
	if settingName == "" {
		resp.Diagnostics.AddError(
			"Missing Setting Name",
			"The setting_name attribute is required",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Searching for group policy setting: %s", settingName))

	// Step 1: Get all group policy categories with expanded definitions
	requestConfiguration := &devicemanagement.GroupPolicyCategoriesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.GroupPolicyCategoriesRequestBuilderGetQueryParameters{
			Expand: []string{"parent,definitions"},
			Select: []string{"id", "displayName", "isRoot", "ingestionSource"},
		},
	}

	categories, err := d.client.
		DeviceManagement().
		GroupPolicyCategories().
		Get(ctx, requestConfiguration)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Group Policy Categories",
			fmt.Sprintf("Could not read group policy categories: %v", err),
		)
		return
	}

	if categories == nil || categories.GetValue() == nil {
		resp.Diagnostics.AddError(
			"No Categories Found",
			"No group policy categories were returned from the API",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found %d categories", len(categories.GetValue())))

	// Step 2: Search through categories and their definitions for the matching setting name
	var foundCategory graphmodels.GroupPolicyCategoryable
	var foundDefinition graphmodels.GroupPolicyDefinitionable
	var definitionID string

	for _, category := range categories.GetValue() {
		if definitions := category.GetDefinitions(); definitions != nil {
			for _, definition := range definitions {
				if definition.GetDisplayName() != nil &&
					strings.EqualFold(*definition.GetDisplayName(), settingName) {
					foundCategory = category
					foundDefinition = definition
					definitionID = *definition.GetId()
					tflog.Debug(ctx, fmt.Sprintf("Found matching definition with ID: %s in category: %s",
						definitionID, *category.GetDisplayName()))
					break
				}
			}
			if foundDefinition != nil {
				break
			}
		}
	}

	if foundDefinition == nil {
		resp.Diagnostics.AddError(
			"Setting Not Found",
			fmt.Sprintf("Could not find group policy setting with name: %s", settingName),
		)
		return
	}

	// Step 3: Get the full definition details using the definition ID
	fullDefinition, err := d.client.
		DeviceManagement().
		GroupPolicyDefinitions().
		ByGroupPolicyDefinitionId(definitionID).
		Get(ctx, nil)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Group Policy Definition",
			fmt.Sprintf("Could not read group policy definition %s: %v", definitionID, err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved full definition details for: %s", *fullDefinition.GetDisplayName()))

	// Step 4: Get the presentations for the definition
	presentations, err := d.client.DeviceManagement().
		GroupPolicyDefinitions().
		ByGroupPolicyDefinitionId(definitionID).
		Presentations().
		Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Group Policy Presentations",
			fmt.Sprintf("Could not read presentations for definition %s: %v", definitionID, err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d presentations",
		func() int {
			if presentations != nil && presentations.GetValue() != nil {
				return len(presentations.GetValue())
			}
			return 0
		}()))

	// Map the data to our model
	data.Category = MapCategoryToDataSource(foundCategory)
	data.Definition = MapDefinitionToDataSource(fullDefinition)

	// Map presentations
	if presentations != nil && presentations.GetValue() != nil {
		data.Presentations = make([]GroupPolicyPresentationModel, len(presentations.GetValue()))
		for i, presentation := range presentations.GetValue() {
			data.Presentations[i] = MapPresentationToDataSource(presentation)
		}
	}

	// Set the ID for the datasource
	data.ID = types.StringValue(fmt.Sprintf("group_policy_category_setting_%s", settingName))

	tflog.Debug(ctx, fmt.Sprintf("Successfully mapped data for setting: %s", settingName))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
