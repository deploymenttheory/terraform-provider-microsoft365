package graphBetaDeviceManagementGroupPolicyValueReference

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapDefinitionsToState converts a list of graph definitions to Terraform state
func mapDefinitionsToState(ctx context.Context, definitions []graphmodels.GroupPolicyDefinitionable, presentations map[string][]PresentationModel, diags *diag.Diagnostics) types.List {
	var definitionModels []DefinitionModel

	for _, def := range definitions {
		definitionID := convert.GraphToFrameworkString(def.GetId())
		displayName := convert.GraphToFrameworkString(def.GetDisplayName())
		classType := convert.GraphToFrameworkEnum(def.GetClassType())
		categoryPath := convert.GraphToFrameworkString(def.GetCategoryPath())
		explainText := convert.GraphToFrameworkString(def.GetExplainText())
		supportedOn := convert.GraphToFrameworkString(def.GetSupportedOn())
		policyType := convert.GraphToFrameworkEnum(def.GetPolicyType())

		// Get presentations for this definition (use the string value of definitionID)
		defIDString := definitionID.ValueString()
		presentationModels := presentations[defIDString]
		presentationsList, presentationDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: presentationAttrTypes(),
		}, presentationModels)
		diags.Append(presentationDiags...)

		definitionModels = append(definitionModels, DefinitionModel{
			Id:            definitionID,
			DisplayName:   displayName,
			ClassType:     classType,
			CategoryPath:  categoryPath,
			ExplainText:   explainText,
			SupportedOn:   supportedOn,
			PolicyType:    policyType,
			Presentations: presentationsList,
		})
	}

	definitionsList, listDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: definitionAttrTypes(),
	}, definitionModels)
	diags.Append(listDiags...)

	return definitionsList
}

// mapPresentationsToState converts graph presentations to presentation models
func mapPresentationsToState(presentations []graphmodels.GroupPolicyPresentationable) []PresentationModel {
	var presentationModels []PresentationModel

	for _, pres := range presentations {
		if pres == nil {
			continue
		}

		presentationID := convert.GraphToFrameworkString(pres.GetId())
		label := convert.GraphToFrameworkString(pres.GetLabel())
		presType := convert.GraphToFrameworkString(pres.GetOdataType())

		// Note: Required field is not available on the base GroupPolicyPresentationable interface
		// It would need type assertion to access specific presentation types
		// Default to false as the field is not reliably accessible
		required := types.BoolValue(false)

		presentationModels = append(presentationModels, PresentationModel{
			Id:       presentationID,
			Label:    label,
			Type:     presType,
			Required: required,
		})
	}

	return presentationModels
}

// definitionAttrTypes returns the attribute types for DefinitionModel
func definitionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"display_name":  types.StringType,
		"class_type":    types.StringType,
		"category_path": types.StringType,
		"explain_text":  types.StringType,
		"supported_on":  types.StringType,
		"policy_type":   types.StringType,
		"presentations": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: presentationAttrTypes(),
			},
		},
	}
}

// presentationAttrTypes returns the attribute types for PresentationModel
func presentationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       types.StringType,
		"label":    types.StringType,
		"type":     types.StringType,
		"required": types.BoolType,
	}
}
