package graphBetaSettingsCatalog

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// SettingsCatalogProfileResourceModel holds the configuration for a Settings Catalog profile.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
type SettingsCatalogProfileResourceModel struct {
	ID                   types.String                                                 `tfsdk:"id"`
	Name                 types.String                                                 `tfsdk:"name"`
	Description          types.String                                                 `tfsdk:"description"`
	Platforms            types.String                                                 `tfsdk:"platforms"`
	Technologies         []types.String                                               `tfsdk:"technologies"`
	RoleScopeTagIds      []types.String                                               `tfsdk:"role_scope_tag_ids"`
	SettingsCount        types.Int64                                                  `tfsdk:"settings_count"`
	IsAssigned           types.Bool                                                   `tfsdk:"is_assigned"`
	LastModifiedDateTime types.String                                                 `tfsdk:"last_modified_date_time"`
	CreatedDateTime      types.String                                                 `tfsdk:"created_date_time"`
	Settings             types.String                                                 `tfsdk:"settings"`
	Assignments          *sharedmodels.SettingsCatalogSettingsAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts             timeouts.Value                                               `tfsdk:"timeouts"`
}

// DeviceConfigV2GraphServiceModel is a struct that represents the JSON structure of settings catalog settings
// fors windows, linux, macOS, and iOS.
// This struct is used for both marshalling and unmarshalling the settings JSON string into a structured format.
// Keys are ordered alphabetically to ensure consistent ordering. This differs from the graph schema.
// This doesn't affect requests to the Graph API, but it ensures that the state is consistent and can be compared
// when using plan modifers.
var DeviceConfigV2GraphServiceModel struct {
	SettingsDetails []struct {
		ID              string `json:"id"`
		SettingInstance struct {
			// For choice setting collections
			ChoiceSettingCollectionValue []struct {
				Children []struct {
					ODataType                        string                                                                     `json:"@odata.type"`
					SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
					SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
					// For nested simple setting collection within choice setting collection
					SimpleSettingCollectionValue []struct {
						ODataType                     string                                                                     `json:"@odata.type"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						Value                         string                                                                     `json:"value"`
					} `json:"simpleSettingCollectionValue,omitempty"`
					// For nested simple settings within choice setting collection
					SimpleSettingValue *struct {
						ODataType                     string                                                                     `json:"@odata.type"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						Value                         interface{}                                                                `json:"value"`
					} `json:"simpleSettingValue,omitempty"`
				} `json:"children"`
				SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
				Value                         string                                                                     `json:"value"`
			} `json:"choiceSettingCollectionValue,omitempty"`

			// For choice settings
			ChoiceSettingValue *struct {
				Children []struct {
					// For nested choice settings within choice children
					ChoiceSettingValue *struct {
						Children []struct {
							ODataType           string `json:"@odata.type"`
							SettingDefinitionId string `json:"settingDefinitionId"`
						} `json:"children"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						Value                         string                                                                     `json:"value"`
					} `json:"choiceSettingValue,omitempty"`
					// For GroupSettingCollectionValue within Choice children
					GroupSettingCollectionValue []struct {
						Children []struct {
							ODataType                        string                                                                     `json:"@odata.type"`
							SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
							SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
							SimpleSettingValue               *struct {
								ODataType                     string                                                                     `json:"@odata.type"`
								SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								Value                         interface{}                                                                `json:"value"`
							} `json:"simpleSettingValue,omitempty"`
						} `json:"children"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
					} `json:"groupSettingCollectionValue,omitempty"`
					ODataType                        string                                                                     `json:"@odata.type"`
					SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
					SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
					// For SimpleSettingCollectionValue within Choice children
					SimpleSettingCollectionValue []struct {
						ODataType                     string                                                                     `json:"@odata.type"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						Value                         string                                                                     `json:"value"`
					} `json:"simpleSettingCollectionValue,omitempty"`
					// For simple settings within choice children
					SimpleSettingValue *struct {
						ODataType                     string                                                                     `json:"@odata.type"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						Value                         interface{}                                                                `json:"value"`
					} `json:"simpleSettingValue,omitempty"`
				} `json:"children"`
				SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
				Value                         string                                                                     `json:"value"`
			} `json:"choiceSettingValue,omitempty"`

			// For group setting collections (Level 1)
			GroupSettingCollectionValue []struct {
				Children []struct {
					// For nested choice settings within group setting collection (Level 2)
					ChoiceSettingValue *struct {
						Children []struct {
							ChoiceSettingValue *struct {
								Children []struct {
									ODataType           string `json:"@odata.type"`
									SettingDefinitionId string `json:"settingDefinitionId"`
								} `json:"children"`
								SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								Value                         string                                                                     `json:"value"`
							} `json:"choiceSettingValue,omitempty"`
							ODataType                        string                                                                     `json:"@odata.type"`
							SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
							SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
							SimpleSettingValue               *struct {
								ODataType                     string                                                                     `json:"@odata.type"`
								SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								Value                         interface{}                                                                `json:"value"`
							} `json:"simpleSettingValue,omitempty"`
						} `json:"children"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						Value                         string                                                                     `json:"value"`
					} `json:"choiceSettingValue,omitempty"`
					// For nested group setting collections within group setting collection (Level 2)
					GroupSettingCollectionValue []struct {
						Children []struct {
							// For nested choice settings within group setting collection within group setting collection (Level 3)
							ChoiceSettingValue *struct {
								Children []struct {
									ODataType           string `json:"@odata.type"`
									SettingDefinitionId string `json:"settingDefinitionId"`
									SimpleSettingValue  *struct {
										ODataType                     string                                                                     `json:"@odata.type"`
										SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
										Value                         interface{}                                                                `json:"value"`
										ValueState                    string                                                                     `json:"valueState,omitempty"`
									} `json:"simpleSettingValue,omitempty"`
								} `json:"children"`
								SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								Value                         string                                                                     `json:"value"`
							} `json:"choiceSettingValue,omitempty"`
							// For nested group setting collections within group setting collection within group setting collection (Level 3)
							GroupSettingCollectionValue []struct {
								Children []struct {
									// For nested choice settings within group setting collection within group setting collection within group setting collection (Level 4)
									ChoiceSettingValue *struct {
										Children []struct {
											ODataType                        string                                                                     `json:"@odata.type"`
											SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
											SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference,omitempty"`
											// For nested simple settings within choice settings within group setting collection within group setting collection within group setting collection (Level 5)
											SimpleSettingValue *struct {
												ODataType                     string                                                                     `json:"@odata.type"`
												SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
												Value                         interface{}                                                                `json:"value"`
											} `json:"simpleSettingValue,omitempty"`
										} `json:"children"`
										SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
										Value                         string                                                                     `json:"value"`
									} `json:"choiceSettingValue,omitempty"`
									ODataType                        string                                                                     `json:"@odata.type"`
									SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
									SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
									// For simple settings collection within group setting collection within group setting collection within group setting collection (Level 4)
									SimpleSettingCollectionValue []struct {
										ODataType                     string                                                                     `json:"@odata.type"`
										SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
										Value                         string                                                                     `json:"value"`
									} `json:"simpleSettingCollectionValue,omitempty"`
									// For simple settings within group setting collection within group setting collection within group setting collection (Level 4)
									SimpleSettingValue *struct {
										ODataType                     string                                                                     `json:"@odata.type"`
										SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
										Value                         interface{}                                                                `json:"value"`
									} `json:"simpleSettingValue,omitempty"`
								} `json:"children"`
								SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
							} `json:"groupSettingCollectionValue,omitempty"`
							ODataType                        string                                                                     `json:"@odata.type"`
							SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
							SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
							// For nested simple setting collections within group setting collection within group setting collection (Level 3)
							SimpleSettingCollectionValue []struct {
								ODataType                     string                                                                     `json:"@odata.type"`
								SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								Value                         string                                                                     `json:"value"`
							} `json:"simpleSettingCollectionValue,omitempty"`
							// For nested simple settings within group setting collection within group setting collection (Level 3)
							SimpleSettingValue *struct {
								ODataType                     string                                                                     `json:"@odata.type"`
								SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								Value                         interface{}                                                                `json:"value"`
								ValueState                    string                                                                     `json:"valueState,omitempty"`
							} `json:"simpleSettingValue,omitempty"`
						} `json:"children"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
					} `json:"groupSettingCollectionValue,omitempty"`
					ODataType                        string                                                                     `json:"@odata.type"`
					SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
					SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
					// For nested simple settings (string, integer, secret) within group setting collection  (Level 2)
					SimpleSettingValue *struct {
						ODataType                     string                                                                     `json:"@odata.type"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						Value                         interface{}                                                                `json:"value"`
						ValueState                    string                                                                     `json:"valueState,omitempty"`
					} `json:"simpleSettingValue,omitempty"`
					// For nested simple setting collections within group setting collection (Level 2)
					SimpleSettingCollectionValue []struct {
						ODataType                     string                                                                     `json:"@odata.type"`
						SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						Value                         string                                                                     `json:"value"`
					} `json:"simpleSettingCollectionValue,omitempty"`
				} `json:"children"`
				SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
			} `json:"groupSettingCollectionValue,omitempty"`

			ODataType                        string                                                                     `json:"@odata.type"`
			SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
			SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`

			// For simple collection settings
			SimpleSettingCollectionValue []struct {
				ODataType                     string                                                                     `json:"@odata.type"`
				SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
				Value                         string                                                                     `json:"value"`
			} `json:"simpleSettingCollectionValue,omitempty"`

			// For simple settings
			SimpleSettingValue *struct {
				ODataType                     string                                                                     `json:"@odata.type"`
				SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
				Value                         interface{}                                                                `json:"value"`
			} `json:"simpleSettingValue,omitempty"`
		} `json:"settingInstance"`
	} `json:"settingsDetails"`
}
