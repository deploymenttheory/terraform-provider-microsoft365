package graphBetaGroupPolicyCategories

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapCategoryToDataSource maps a Group Policy Category to a model
func MapCategoryToDataSource(data graphmodels.GroupPolicyCategoryable) *GroupPolicyCategoryModel {
	model := &GroupPolicyCategoryModel{
		ID:              convert.GraphToFrameworkString(data.GetId()),
		DisplayName:     convert.GraphToFrameworkString(data.GetDisplayName()),
		IsRoot:          convert.GraphToFrameworkBool(data.GetIsRoot()),
		IngestionSource: convert.GraphToFrameworkEnum(data.GetIngestionSource()),
	}

	// Map parent if it exists
	if parent := data.GetParent(); parent != nil {
		model.Parent = &GroupPolicyCategoryParentModel{
			ID:          convert.GraphToFrameworkString(parent.GetId()),
			DisplayName: convert.GraphToFrameworkString(parent.GetDisplayName()),
			IsRoot:      convert.GraphToFrameworkBool(parent.GetIsRoot()),
		}
	}

	return model
}

// MapDefinitionToDataSource maps a Group Policy Definition to a model
func MapDefinitionToDataSource(data graphmodels.GroupPolicyDefinitionable) *GroupPolicyDefinitionModel {
	model := &GroupPolicyDefinitionModel{
		ID:                    convert.GraphToFrameworkString(data.GetId()),
		DisplayName:           convert.GraphToFrameworkString(data.GetDisplayName()),
		CategoryPath:          convert.GraphToFrameworkString(data.GetCategoryPath()),
		ClassType:             convert.GraphToFrameworkEnum(data.GetClassType()),
		PolicyType:            convert.GraphToFrameworkEnum(data.GetPolicyType()),
		Version:               convert.GraphToFrameworkString(data.GetVersion()),
		HasRelatedDefinitions: convert.GraphToFrameworkBool(data.GetHasRelatedDefinitions()),
		ExplainText:           convert.GraphToFrameworkString(data.GetExplainText()),
		SupportedOn:           convert.GraphToFrameworkString(data.GetSupportedOn()),
		GroupPolicyCategoryID: convert.GraphToFrameworkUUID(data.GetGroupPolicyCategoryId()),
		MinDeviceCSPVersion:   convert.GraphToFrameworkString(data.GetMinDeviceCspVersion()),
		MinUserCSPVersion:     convert.GraphToFrameworkString(data.GetMinUserCspVersion()),
		LastModifiedDateTime:  convert.GraphToFrameworkTime(data.GetLastModifiedDateTime()),
	}

	return model
}

// MapPresentationToDataSource maps a Group Policy Presentation to a model
func MapPresentationToDataSource(data graphmodels.GroupPolicyPresentationable) GroupPolicyPresentationModel {
	model := GroupPolicyPresentationModel{
		ID:                   convert.GraphToFrameworkString(data.GetId()),
		ODataType:            convert.GraphToFrameworkString(data.GetOdataType()),
		Label:                convert.GraphToFrameworkString(data.GetLabel()),
		Required:             types.BoolNull(), // Will be set by specific presentation type handlers if available
		LastModifiedDateTime: convert.GraphToFrameworkTime(data.GetLastModifiedDateTime()),
	}

	// Handle different presentation types using type assertion based on OData type
	if data.GetOdataType() != nil {
		switch *data.GetOdataType() {
		case "#microsoft.graph.groupPolicyPresentationDropdownList":
			mapDropdownListPresentation(data, &model)
		case "#microsoft.graph.groupPolicyPresentationTextBox":
			mapTextBoxPresentation(data, &model)
		case "#microsoft.graph.groupPolicyPresentationCheckBox":
			mapCheckBoxPresentation(data, &model)
		case "#microsoft.graph.groupPolicyPresentationDecimalTextBox":
			mapDecimalTextBoxPresentation(data, &model)
		case "#microsoft.graph.groupPolicyPresentationLongDecimalTextBox":
			mapLongDecimalTextBoxPresentation(data, &model)
		case "#microsoft.graph.groupPolicyPresentationMultiTextBox":
			mapMultiTextBoxPresentation(data, &model)
		case "#microsoft.graph.groupPolicyPresentationComboBox":
			mapComboBoxPresentation(data, &model)
		case "#microsoft.graph.groupPolicyPresentationListBox":
			mapListBoxPresentation(data, &model)
		}
	}

	return model
}

// mapDropdownListPresentation handles dropdown list specific properties
func mapDropdownListPresentation(data graphmodels.GroupPolicyPresentationable, model *GroupPolicyPresentationModel) {
	// Cast to the specific dropdown type based on odata.type
	if dropdownData, ok := data.(*graphmodels.GroupPolicyPresentationDropdownList); ok {
		// Handle required field
		model.Required = convert.GraphToFrameworkBool(dropdownData.GetRequired())

		// Handle default item
		if defaultItem := dropdownData.GetDefaultItem(); defaultItem != nil {
			model.DefaultItem = &GroupPolicyPresentationItemModel{
				DisplayName: convert.GraphToFrameworkString(defaultItem.GetDisplayName()),
				Value:       convert.GraphToFrameworkString(defaultItem.GetValue()),
			}
		}

		// Handle items list
		if items := dropdownData.GetItems(); items != nil {
			for _, item := range items {
				model.Items = append(model.Items, GroupPolicyPresentationItemModel{
					DisplayName: convert.GraphToFrameworkString(item.GetDisplayName()),
					Value:       convert.GraphToFrameworkString(item.GetValue()),
				})
			}
		}
	}
}

// mapTextBoxPresentation handles text box specific properties
func mapTextBoxPresentation(data graphmodels.GroupPolicyPresentationable, model *GroupPolicyPresentationModel) {
	// Cast to the specific text box type based on odata.type
	if textBoxData, ok := data.(*graphmodels.GroupPolicyPresentationTextBox); ok {
		model.DefaultValue = convert.GraphToFrameworkString(textBoxData.GetDefaultValue())
		model.MaxLength = convert.GraphToFrameworkInt64(textBoxData.GetMaxLength())
		model.Required = convert.GraphToFrameworkBool(textBoxData.GetRequired())
	}
}

// mapCheckBoxPresentation handles checkbox specific properties
func mapCheckBoxPresentation(data graphmodels.GroupPolicyPresentationable, model *GroupPolicyPresentationModel) {
	// Cast to the specific checkbox type based on odata.type
	if checkBoxData, ok := data.(*graphmodels.GroupPolicyPresentationCheckBox); ok {
		model.DefaultChecked = convert.GraphToFrameworkBool(checkBoxData.GetDefaultChecked())
	}
}

// mapDecimalTextBoxPresentation handles decimal text box specific properties
func mapDecimalTextBoxPresentation(data graphmodels.GroupPolicyPresentationable, model *GroupPolicyPresentationModel) {
	// Cast to the specific decimal text box type based on odata.type
	if decimalTextBoxData, ok := data.(*graphmodels.GroupPolicyPresentationDecimalTextBox); ok {
		model.DefaultDecimalValue = convert.GraphToFrameworkInt64(decimalTextBoxData.GetDefaultValue())
		model.MinValue = convert.GraphToFrameworkInt64(decimalTextBoxData.GetMinValue())
		model.MaxValue = convert.GraphToFrameworkInt64(decimalTextBoxData.GetMaxValue())
		model.Spin = convert.GraphToFrameworkBool(decimalTextBoxData.GetSpin())
		model.SpinStep = convert.GraphToFrameworkInt64(decimalTextBoxData.GetSpinStep())
		model.Required = convert.GraphToFrameworkBool(decimalTextBoxData.GetRequired())
	}
}

// mapLongDecimalTextBoxPresentation handles long decimal text box specific properties
func mapLongDecimalTextBoxPresentation(data graphmodels.GroupPolicyPresentationable, model *GroupPolicyPresentationModel) {
	// Cast to the specific long decimal text box type based on odata.type
	if longDecimalTextBoxData, ok := data.(*graphmodels.GroupPolicyPresentationLongDecimalTextBox); ok {
		model.DefaultDecimalValue = convert.GraphToFrameworkInt64(longDecimalTextBoxData.GetDefaultValue())
		model.MinValue = convert.GraphToFrameworkInt64(longDecimalTextBoxData.GetMinValue())
		model.MaxValue = convert.GraphToFrameworkInt64(longDecimalTextBoxData.GetMaxValue())
		model.Spin = convert.GraphToFrameworkBool(longDecimalTextBoxData.GetSpin())
		model.SpinStep = convert.GraphToFrameworkInt64(longDecimalTextBoxData.GetSpinStep())
		model.Required = convert.GraphToFrameworkBool(longDecimalTextBoxData.GetRequired())
	}
}

// mapMultiTextBoxPresentation handles multi text box specific properties
func mapMultiTextBoxPresentation(data graphmodels.GroupPolicyPresentationable, model *GroupPolicyPresentationModel) {
	// Cast to the specific multi text box type based on odata.type
	if multiTextBoxData, ok := data.(*graphmodels.GroupPolicyPresentationMultiTextBox); ok {
		model.MaxLength = convert.GraphToFrameworkInt64(multiTextBoxData.GetMaxLength())
		model.MaxValue = convert.GraphToFrameworkInt64(multiTextBoxData.GetMaxStrings()) // Store MaxStrings in MaxValue
		model.Required = convert.GraphToFrameworkBool(multiTextBoxData.GetRequired())
	}
}

// mapComboBoxPresentation handles combo box specific properties
func mapComboBoxPresentation(data graphmodels.GroupPolicyPresentationable, model *GroupPolicyPresentationModel) {
	// Cast to the specific combo box type based on odata.type
	if comboBoxData, ok := data.(*graphmodels.GroupPolicyPresentationComboBox); ok {
		model.DefaultValue = convert.GraphToFrameworkString(comboBoxData.GetDefaultValue())
		model.MaxLength = convert.GraphToFrameworkInt64(comboBoxData.GetMaxLength())
		model.Required = convert.GraphToFrameworkBool(comboBoxData.GetRequired())
	}
}

// mapListBoxPresentation handles list box specific properties
// REF: https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicypresentationlistbox-get?view=graph-rest-beta
func mapListBoxPresentation(data graphmodels.GroupPolicyPresentationable, model *GroupPolicyPresentationModel) {
	if listBoxData, ok := data.(*graphmodels.GroupPolicyPresentationListBox); ok {
		model.ExplicitValue = convert.GraphToFrameworkBool(listBoxData.GetExplicitValue())
		model.ValuePrefix = convert.GraphToFrameworkString(listBoxData.GetValuePrefix())
		// List box presentations don't have a 'required' field like other presentation types
	}
}
