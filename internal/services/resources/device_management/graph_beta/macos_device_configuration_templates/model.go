package graphBetaMacosDeviceConfigurationTemplates

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MacosDeviceConfigurationTemplatesResourceModel describes the resource data model.
type MacosDeviceConfigurationTemplatesResourceModel struct {
	ID              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	RoleScopeTagIds types.Set    `tfsdk:"role_scope_tag_ids"`
	// Nested configuration blocks (mutually exclusive)
	CustomConfiguration types.Object   `tfsdk:"custom_configuration"`
	PreferenceFile      types.Object   `tfsdk:"preference_file"`
	TrustedCertificate  types.Object   `tfsdk:"trusted_certificate"`
	ScepCertificate     types.Object   `tfsdk:"scep_certificate"`
	PkcsCertificate     types.Object   `tfsdk:"pkcs_certificate"`
	Assignments         types.Set      `tfsdk:"assignments"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

// CustomConfigurationResourceModel describes macOSCustomConfiguration.
type CustomConfigurationResourceModel struct {
	DeploymentChannel types.String `tfsdk:"deployment_channel"`
	PayloadFileName   types.String `tfsdk:"payload_file_name"`
	Payload           types.String `tfsdk:"payload"`
	PayloadName       types.String `tfsdk:"payload_name"`
}

// PreferenceFileResourceModel describes macOSCustomAppConfiguration.
type PreferenceFileResourceModel struct {
	FileName         types.String `tfsdk:"file_name"`
	ConfigurationXml types.String `tfsdk:"configuration_xml"`
	BundleId         types.String `tfsdk:"bundle_id"`
}

// TrustedCertificateResourceModel describes macOSTrustedRootCertificate.
type TrustedCertificateResourceModel struct {
	DeploymentChannel      types.String `tfsdk:"deployment_channel"`
	CertFileName           types.String `tfsdk:"cert_file_name"`
	TrustedRootCertificate types.String `tfsdk:"trusted_root_certificate"`
}

// ScepCertificateResourceModel describes macOSScepCertificateProfile.
type ScepCertificateResourceModel struct {
	DeploymentChannel              types.String `tfsdk:"deployment_channel"`
	RenewalThresholdPercentage     types.Int32  `tfsdk:"renewal_threshold_percentage"`
	CertificateStore               types.String `tfsdk:"certificate_store"`
	CertificateValidityPeriodScale types.String `tfsdk:"certificate_validity_period_scale"`
	CertificateValidityPeriodValue types.Int32  `tfsdk:"certificate_validity_period_value"`
	SubjectNameFormat              types.String `tfsdk:"subject_name_format"`
	SubjectNameFormatString        types.String `tfsdk:"subject_name_format_string"`
	RootCertificateOdataBind       types.String `tfsdk:"root_certificate_odata_bind"`
	KeySize                        types.String `tfsdk:"key_size"`
	KeyUsage                       types.Set    `tfsdk:"key_usage"`
	CustomSubjectAlternativeNames  types.Set    `tfsdk:"custom_subject_alternative_names"`
	ExtendedKeyUsages              types.Set    `tfsdk:"extended_key_usages"`
	ScepServerUrls                 types.Set    `tfsdk:"scep_server_urls"`
	AllowAllAppsAccess             types.Bool   `tfsdk:"allow_all_apps_access"`
}

// PkcsCertificateResourceModel describes macOSPkcsCertificateProfile.
type PkcsCertificateResourceModel struct {
	DeploymentChannel              types.String `tfsdk:"deployment_channel"`
	RenewalThresholdPercentage     types.Int32  `tfsdk:"renewal_threshold_percentage"`
	CertificateStore               types.String `tfsdk:"certificate_store"`
	CertificateValidityPeriodScale types.String `tfsdk:"certificate_validity_period_scale"`
	CertificateValidityPeriodValue types.Int32  `tfsdk:"certificate_validity_period_value"`
	SubjectNameFormat              types.String `tfsdk:"subject_name_format"`
	SubjectNameFormatString        types.String `tfsdk:"subject_name_format_string"`
	CertificationAuthority         types.String `tfsdk:"certification_authority"`
	CertificationAuthorityName     types.String `tfsdk:"certification_authority_name"`
	CertificateTemplateName        types.String `tfsdk:"certificate_template_name"`
	CustomSubjectAlternativeNames  types.Set    `tfsdk:"custom_subject_alternative_names"`
	AllowAllAppsAccess             types.Bool   `tfsdk:"allow_all_apps_access"`
}

// CustomSubjectAlternativeNameResourceModel describes custom SAN entries.
type CustomSubjectAlternativeNameResourceModel struct {
	SanType types.String `tfsdk:"san_type"`
	Name    types.String `tfsdk:"name"`
}

// ExtendedKeyUsageResourceModel describes extended key usage entries.
type ExtendedKeyUsageResourceModel struct {
	Name             types.String `tfsdk:"name"`
	ObjectIdentifier types.String `tfsdk:"object_identifier"`
}
