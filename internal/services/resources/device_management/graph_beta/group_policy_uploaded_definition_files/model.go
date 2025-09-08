// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicyuploadeddefinitionfile?view=graph-rest-beta
package graphBetaGroupPolicyUploadedDefinitionFiles

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupPolicyUploadedDefinitionFileResourceModel struct {
	ID                               types.String   `tfsdk:"id"`
	DisplayName                      types.String   `tfsdk:"display_name"`
	Description                      types.String   `tfsdk:"description"`
	FileName                         types.String   `tfsdk:"file_name"`
	Content                          types.String   `tfsdk:"content"`
	DefaultLanguageCode              types.String   `tfsdk:"default_language_code"`
	LanguageCodes                    types.List     `tfsdk:"language_codes"`
	TargetPrefix                     types.String   `tfsdk:"target_prefix"`
	TargetNamespace                  types.String   `tfsdk:"target_namespace"`
	PolicyType                       types.String   `tfsdk:"policy_type"`
	Revision                         types.String   `tfsdk:"revision"`
	Status                           types.String   `tfsdk:"status"`
	UploadDateTime                   types.String   `tfsdk:"upload_date_time"`
	LastModifiedDateTime             types.String   `tfsdk:"last_modified_date_time"`
	GroupPolicyUploadedLanguageFiles types.Set      `tfsdk:"group_policy_uploaded_language_files"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
}

type GroupPolicyUploadedLanguageFileModel struct {
	FileName     types.String `tfsdk:"file_name"`
	LanguageCode types.String `tfsdk:"language_code"`
	Content      types.String `tfsdk:"content"`
}
