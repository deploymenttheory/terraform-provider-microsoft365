package graphBetaMacosDeviceConfigurationTemplates

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_macos_device_configuration_templates"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &MacosDeviceConfigurationTemplatesResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &MacosDeviceConfigurationTemplatesResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &MacosDeviceConfigurationTemplatesResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &MacosDeviceConfigurationTemplatesResource{}
)

func NewMacosDeviceConfigurationTemplatesResource() resource.Resource {
	return &MacosDeviceConfigurationTemplatesResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceConfigurations",
	}
}

type MacosDeviceConfigurationTemplatesResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *MacosDeviceConfigurationTemplatesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *MacosDeviceConfigurationTemplatesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *MacosDeviceConfigurationTemplatesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *MacosDeviceConfigurationTemplatesResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages macOS configuration templates in Microsoft Intune. " +
			"This resource creates device configurations for macOS devices including custom configuration profiles, " +
			"preference files, trusted certificates, and certificate profiles (SCEP/PKCS).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the macOS configuration template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the macOS configuration template.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description for the macOS configuration template.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"custom_configuration": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The custom configuration template allows IT admins to assign settings that aren't built into Intune yet. For macOS devices, you can import a .mobileconfig file that you created using Profile Manager or a different tool.",
				Attributes: map[string]schema.Attribute{
					"deployment_channel": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Select the channel you want to use to deploy your configuration profile. If the channel doesn’t match what’s listed for the payload in Apple documentation, deployment could fail. The selected channel cannot be changed once the profile has been created. Possible values are: deviceChannel, userChannel.",
						Validators: []validator.String{
							stringvalidator.OneOf("deviceChannel", "userChannel"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"payload_file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The profile name displayed to users.",
					},
					"payload": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The macOS configuration payload (.mobileconfig / .plist) file content.",
					},
					"payload_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The name of the payload configuration.",
					},
				},
			},
			"preference_file": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configure a preference file that uses the standard property list (.plist) format to define preferences for apps and the device.",
				Attributes: map[string]schema.Attribute{
					"file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The file name of the preference file (.plist file).",
					},
					"configuration_xml": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The base64-encoded XML configuration content (.plist file content).",
					},
					"bundle_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The bundle ID (Preference domain name) of the application this preference file applies to. Typically in the format com.company.appname.",
					},
				},
			},
			"trusted_certificate": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Trusted root certificate configuration for macOS devices.",
				Attributes: map[string]schema.Attribute{
					"deployment_channel": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The deployment channel for the certificate. Possible values are: deviceChannel, userChannel.",
						Validators: []validator.String{
							stringvalidator.OneOf("deviceChannel", "userChannel"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"cert_file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The file name of the certificate file (.cer file).",
					},
					"trusted_root_certificate": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The base64-encoded trusted root certificate content. This should be a filebase64() encoded string. e.g filebase64(\"my-root-cert.cer\")",
					},
				},
			},
			"scep_certificate": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "SCEP certificate profile configuration for macOS devices.",
				Attributes: map[string]schema.Attribute{
					"deployment_channel": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The deployment channel for the certificate. Possible values are: deviceChannel, userChannel.",
						Validators: []validator.String{
							stringvalidator.OneOf("deviceChannel", "userChannel"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"renewal_threshold_percentage": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "The certificate renewal threshold percentage (1-99).",
						Validators: []validator.Int32{
							int32validator.Between(1, 99),
						},
					},
					"certificate_store": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The certificate store location. Possible values are: user, machine.",
						Validators: []validator.String{
							stringvalidator.OneOf("user", "machine"),
						},
					},
					"certificate_validity_period_scale": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The certificate validity period scale. Possible values are: days, months, years.",
						Validators: []validator.String{
							stringvalidator.OneOf("days", "months", "years"),
						},
					},
					"certificate_validity_period_value": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "The certificate validity period value.",
						Validators: []validator.Int32{
							int32validator.AtLeast(1),
						},
					},
					"subject_name_format": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Defaults to custom.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"commonName",
								"commonNameIncludingEmail",
								"commonNameAsEmail",
								"custom",
								"commonNameAsIMEI",
								"commonNameAsSerialNumber",
								"commonNameAsAadDeviceId",
								"commonNameAsIntuneDeviceId",
								"commonNameAsDurableDeviceId",
								"commonNameAsOnPremisesSamAccountName",
							),
						},
					},
					"subject_name_format_string": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "Select how Intune automatically creates the subject name in the certificate request. " +
							"If the certificate is for a user, you can also include the user's email address in the subject name. " +
							"Please review subject name documentation 'https://learn.microsoft.com/en-us/intune/intune-service/protect/certificates-profile-scep'" +
							"on how to best use the Subject name format field." +
							"Custom. Example: CN={{AAD_Device_ID}},O={{Organization}}",
					},
					"root_certificate_odata_bind": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "Reference to the pre existing trusted root certificate configuration for the odata bind." +
							"Valid format is \"https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('00000000-0000-0000-0000-000000000000')\"." +
							"Or you can supply just the ID of the certificate configuration. e.g. '00000000-0000-0000-0000-000000000000'",
					},
					"key_size": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The key size in bits for the certificate.2048 is the recommended minimum key length.Possible values are: size1024, size2048, size4096.",
						Validators: []validator.String{
							stringvalidator.OneOf("size1024", "size2048", "size4096"),
						},
					},
					"key_usage": schema.SetAttribute{
						ElementType:         types.StringType,
						Required:            true,
						MarkdownDescription: "Key usage options for the certificate. Possible values are: keyEncipherment, digitalSignature.",
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(stringvalidator.OneOf("keyEncipherment", "digitalSignature")),
						},
					},
					"custom_subject_alternative_names": schema.SetNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Custom Subject Alternative Names for the certificate.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"san_type": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The SAN type. Possible values are: emailAddress, userPrincipalName, customAzureADAttribute, domainNameService, universalResourceIdentifier.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"emailAddress", "userPrincipalName", "customAzureADAttribute",
											"domainNameService", "universalResourceIdentifier",
										),
									},
								},
								"name": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The SAN value/name.",
								},
							},
						},
					},
					"extended_key_usages": schema.SetNestedAttribute{
						Required:            true,
						MarkdownDescription: "Extended key usage settings for the certificate.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The extended key usage name.",
								},
								"object_identifier": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The extended key usage object identifier (OID).",
								},
							},
						},
					},
					"scep_server_urls": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "SCEP server URL(s) for certificate enrollment.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"allow_all_apps_access": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether to allow all applications to access the certificate.",
					},
				},
			},
			"pkcs_certificate": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "PKCS certificate profile configuration for macOS devices.",
				Attributes: map[string]schema.Attribute{
					"deployment_channel": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The deployment channel for the certificate. Possible values are: deviceChannel, userChannel.",
						Validators: []validator.String{
							stringvalidator.OneOf("deviceChannel", "userChannel"),
						},
						Default: stringdefault.StaticString("deviceChannel"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"renewal_threshold_percentage": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The certificate renewal threshold percentage (1-99).",
						Validators: []validator.Int32{
							int32validator.Between(1, 99),
						},
					},
					"certificate_store": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The certificate store location. Possible values are: user, machine.",
						Validators: []validator.String{
							stringvalidator.OneOf("user", "machine"),
						},
					},
					"certificate_validity_period_scale": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The certificate validity period scale. Possible values are: days, months, years.",
						Validators: []validator.String{
							stringvalidator.OneOf("days", "months", "years"),
						},
					},
					"certificate_validity_period_value": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The certificate validity period value.",
						Validators: []validator.Int32{
							int32validator.AtLeast(1),
						},
					},
					"subject_name_format": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Defaults to custom.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"commonName",
								"commonNameIncludingEmail",
								"commonNameAsEmail",
								"custom",
								"commonNameAsIMEI",
								"commonNameAsSerialNumber",
								"commonNameAsAadDeviceId",
								"commonNameAsIntuneDeviceId",
								"commonNameAsDurableDeviceId",
								"commonNameAsOnPremisesSamAccountName",
							),
						},
					},
					"subject_name_format_string": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom format to use with SubjectNameFormat = Custom. Example: CN={{AAD_Device_ID}},O={{Organization}}",
					},
					"certification_authority": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The certification authority for PKCS certificates.",
					},
					"certification_authority_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The certification authority name for PKCS certificates.",
					},
					"certificate_template_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The certificate template name for PKCS certificates.",
					},
					"custom_subject_alternative_names": schema.SetNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Custom Subject Alternative Names for the certificate.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"san_type": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The SAN type. Possible values are: emailAddress, userPrincipalName, customAzureADAttribute, domainNameService, universalResourceIdentifier.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"none", "emailAddress", "userPrincipalName", "customAzureADAttribute", "domainNameService", "universalResourceIdentifier",
										),
									},
								},
								"name": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The SAN value/name.",
								},
							},
						},
					},
					"allow_all_apps_access": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether to allow all applications to access the certificate.",
					},
				},
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
