package graphBetaMacosDeviceConfigurationTemplates

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Attribute type helpers for ObjectNull calls

func CustomConfigurationType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"deployment_channel": types.StringType,
			"payload_file_name":  types.StringType,
			"payload":            types.StringType,
			"payload_name":       types.StringType,
		},
	}
}

func PreferenceFileType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"file_name":         types.StringType,
			"configuration_xml": types.StringType,
			"bundle_id":         types.StringType,
		},
	}
}

func TrustedCertificateType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"deployment_channel":       types.StringType,
			"cert_file_name":           types.StringType,
			"trusted_root_certificate": types.StringType,
		},
	}
}

func ScepCertificateType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"deployment_channel":                types.StringType,
			"renewal_threshold_percentage":      types.Int32Type,
			"certificate_store":                 types.StringType,
			"certificate_validity_period_scale": types.StringType,
			"certificate_validity_period_value": types.Int32Type,
			"subject_name_format":               types.StringType,
			"subject_name_format_string":        types.StringType,
			"root_certificate_odata_bind":       types.StringType,
			"key_size":                          types.StringType,
			"key_usage":                         types.SetType{ElemType: types.StringType},
			"custom_subject_alternative_names":  types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"san_type": types.StringType, "name": types.StringType}}},
			"extended_key_usages":               types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType, "object_identifier": types.StringType}}},
			"scep_server_urls":                  types.SetType{ElemType: types.StringType},
			"allow_all_apps_access":             types.BoolType,
		},
	}
}

func PkcsCertificateType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"deployment_channel":                types.StringType,
			"renewal_threshold_percentage":      types.Int32Type,
			"certificate_store":                 types.StringType,
			"certificate_validity_period_scale": types.StringType,
			"certificate_validity_period_value": types.Int32Type,
			"subject_name_format":               types.StringType,
			"subject_name_format_string":        types.StringType,
			"certification_authority":           types.StringType,
			"certification_authority_name":      types.StringType,
			"certificate_template_name":         types.StringType,
			"custom_subject_alternative_names":  types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"san_type": types.StringType, "name": types.StringType}}},
			"allow_all_apps_access":             types.BoolType,
		},
	}
}

func MapRemoteResourceStateToTerraform(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel, remoteResource graphmodels.DeviceConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// Map common properties
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	// Map specific configuration based on type
	switch config := remoteResource.(type) {
	case *graphmodels.MacOSCustomConfiguration:
		mapMacOSCustomConfiguration(ctx, data, config)
	case *graphmodels.MacOSCustomAppConfiguration:
		mapMacOSCustomAppConfiguration(ctx, data, config)
	case *graphmodels.MacOSTrustedRootCertificate:
		mapMacOSTrustedRootCertificate(ctx, data, config)
	case *graphmodels.MacOSScepCertificateProfile:
		mapMacOSScepCertificateProfile(ctx, data, config)
	case *graphmodels.MacOSPkcsCertificateProfile:
		mapMacOSPkcsCertificateProfile(ctx, data, config)
	default:
		tflog.Error(ctx, "Unknown device configuration type", map[string]interface{}{
			"type": fmt.Sprintf("%T", config),
		})
	}

	// Map assignments
	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		data.Assignments = types.SetNull(MacosConfigurationTemplatesAssignmentType())
	} else {
		mapAssignmentsToTerraform(ctx, data, assignments)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

func mapMacOSCustomConfiguration(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel, config *graphmodels.MacOSCustomConfiguration) {
	tflog.Debug(ctx, "Mapping MacOSCustomConfiguration")

	customConfigModel := CustomConfigurationResourceModel{
		DeploymentChannel: convert.GraphToFrameworkEnum(config.GetDeploymentChannel()),
		PayloadFileName:   convert.GraphToFrameworkString(config.GetPayloadFileName()),
		Payload:           convert.GraphToFrameworkBytes(config.GetPayload()),
		PayloadName:       convert.GraphToFrameworkString(config.GetPayloadName()),
	}

	objectValue, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"deployment_channel": types.StringType,
		"payload_file_name":  types.StringType,
		"payload":            types.StringType,
		"payload_name":       types.StringType,
	}, customConfigModel)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to create custom configuration object", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return
	}

	data.CustomConfiguration = objectValue
	// Clear other configuration types
	data.PreferenceFile = types.ObjectNull(PreferenceFileType().AttrTypes)
	data.TrustedCertificate = types.ObjectNull(TrustedCertificateType().AttrTypes)
	data.ScepCertificate = types.ObjectNull(ScepCertificateType().AttrTypes)
	data.PkcsCertificate = types.ObjectNull(PkcsCertificateType().AttrTypes)
}

func mapMacOSCustomAppConfiguration(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel, config *graphmodels.MacOSCustomAppConfiguration) {
	tflog.Debug(ctx, "Mapping MacOSCustomAppConfiguration")

	preferenceModel := PreferenceFileResourceModel{
		FileName:         convert.GraphToFrameworkString(config.GetFileName()),
		ConfigurationXml: convert.GraphToFrameworkBytes(config.GetConfigurationXml()),
		BundleId:         convert.GraphToFrameworkString(config.GetBundleId()),
	}

	objectValue, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"file_name":         types.StringType,
		"configuration_xml": types.StringType,
		"bundle_id":         types.StringType,
	}, preferenceModel)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to create preference file object", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return
	}

	data.PreferenceFile = objectValue
	// Clear other configuration types
	data.CustomConfiguration = types.ObjectNull(CustomConfigurationType().AttrTypes)
	data.TrustedCertificate = types.ObjectNull(TrustedCertificateType().AttrTypes)
	data.ScepCertificate = types.ObjectNull(ScepCertificateType().AttrTypes)
	data.PkcsCertificate = types.ObjectNull(PkcsCertificateType().AttrTypes)
}

func mapMacOSTrustedRootCertificate(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel, config *graphmodels.MacOSTrustedRootCertificate) {
	tflog.Debug(ctx, "Mapping MacOSTrustedRootCertificate")

	// Convert binary certificate data back to base64 for state consistency
	var trustedRootCert types.String
	if certBytes := config.GetTrustedRootCertificate(); certBytes != nil {
		certBase64 := base64.StdEncoding.EncodeToString(certBytes)
		trustedRootCert = types.StringValue(certBase64)
	} else {
		trustedRootCert = types.StringNull()
	}

	certModel := TrustedCertificateResourceModel{
		DeploymentChannel:      convert.GraphToFrameworkEnum(config.GetDeploymentChannel()),
		CertFileName:           convert.GraphToFrameworkString(config.GetCertFileName()),
		TrustedRootCertificate: trustedRootCert,
	}

	objectValue, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"deployment_channel":       types.StringType,
		"cert_file_name":           types.StringType,
		"trusted_root_certificate": types.StringType,
	}, certModel)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to create trusted certificate object", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return
	}

	data.TrustedCertificate = objectValue
	// Clear other configuration types
	data.CustomConfiguration = types.ObjectNull(CustomConfigurationType().AttrTypes)
	data.PreferenceFile = types.ObjectNull(PreferenceFileType().AttrTypes)
	data.ScepCertificate = types.ObjectNull(ScepCertificateType().AttrTypes)
	data.PkcsCertificate = types.ObjectNull(PkcsCertificateType().AttrTypes)
}

func mapMacOSScepCertificateProfile(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel, config *graphmodels.MacOSScepCertificateProfile) {
	tflog.Debug(ctx, "Mapping MacOSScepCertificateProfile")

	// Preserve root certificate odata.bind from existing state. The api doesnt return the additional data reference.
	rootCertRef := types.StringNull()
	if !data.ScepCertificate.IsNull() && !data.ScepCertificate.IsUnknown() {
		var existingScepData ScepCertificateResourceModel
		diags := data.ScepCertificate.As(ctx, &existingScepData, basetypes.ObjectAsOptions{})
		if !diags.HasError() && !existingScepData.RootCertificateOdataBind.IsNull() {
			rootCertRef = existingScepData.RootCertificateOdataBind
		}
	}

	// If we don't have existing state, try to construct from the API response
	if rootCertRef.IsNull() {
		if rootCert := config.GetRootCertificate(); rootCert != nil {
			if id := rootCert.GetId(); id != nil {
				rootCertRef = types.StringValue(fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('%s')", *id))
			}
		}
	}

	scepModel := ScepCertificateResourceModel{
		DeploymentChannel:              convert.GraphToFrameworkEnum(config.GetDeploymentChannel()),
		RenewalThresholdPercentage:     convert.GraphToFrameworkInt32(config.GetRenewalThresholdPercentage()),
		CertificateStore:               convert.GraphToFrameworkEnum(config.GetCertificateStore()),
		CertificateValidityPeriodScale: convert.GraphToFrameworkEnum(config.GetCertificateValidityPeriodScale()),
		CertificateValidityPeriodValue: convert.GraphToFrameworkInt32(config.GetCertificateValidityPeriodValue()),
		SubjectNameFormat:              convert.GraphToFrameworkEnum(config.GetSubjectNameFormat()),
		SubjectNameFormatString:        convert.GraphToFrameworkString(config.GetSubjectNameFormatString()),
		RootCertificateOdataBind:       rootCertRef,
		KeySize:                        convert.GraphToFrameworkEnum(config.GetKeySize()),
		KeyUsage:                       mapKeyUsageToSet(ctx, config.GetKeyUsage()),
		CustomSubjectAlternativeNames:  mapCustomSANsToSet(ctx, config.GetCustomSubjectAlternativeNames()),
		ExtendedKeyUsages:              mapExtendedKeyUsagesToSet(ctx, config.GetExtendedKeyUsages()),
		ScepServerUrls:                 convert.GraphToFrameworkStringSet(ctx, config.GetScepServerUrls()),
		AllowAllAppsAccess:             convert.GraphToFrameworkBool(config.GetAllowAllAppsAccess()),
	}

	// Create the attribute type map for SCEP certificate
	scepAttrTypes := map[string]attr.Type{
		"deployment_channel":                types.StringType,
		"renewal_threshold_percentage":      types.Int32Type,
		"certificate_store":                 types.StringType,
		"certificate_validity_period_scale": types.StringType,
		"certificate_validity_period_value": types.Int32Type,
		"subject_name_format":               types.StringType,
		"subject_name_format_string":        types.StringType,
		"root_certificate_odata_bind":       types.StringType,
		"key_size":                          types.StringType,
		"key_usage":                         types.SetType{ElemType: types.StringType},
		"custom_subject_alternative_names":  types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"san_type": types.StringType, "name": types.StringType}}},
		"extended_key_usages":               types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType, "object_identifier": types.StringType}}},
		"scep_server_urls":                  types.SetType{ElemType: types.StringType},
		"allow_all_apps_access":             types.BoolType,
	}

	objectValue, diags := types.ObjectValueFrom(ctx, scepAttrTypes, scepModel)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to create SCEP certificate object", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return
	}

	data.ScepCertificate = objectValue
	// Clear other configuration types
	data.CustomConfiguration = types.ObjectNull(CustomConfigurationType().AttrTypes)
	data.PreferenceFile = types.ObjectNull(PreferenceFileType().AttrTypes)
	data.TrustedCertificate = types.ObjectNull(TrustedCertificateType().AttrTypes)
	data.PkcsCertificate = types.ObjectNull(PkcsCertificateType().AttrTypes)
}

func mapMacOSPkcsCertificateProfile(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel, config *graphmodels.MacOSPkcsCertificateProfile) {
	tflog.Debug(ctx, "Mapping MacOSPkcsCertificateProfile")

	pkcsModel := PkcsCertificateResourceModel{
		DeploymentChannel:              convert.GraphToFrameworkEnum(config.GetDeploymentChannel()),
		RenewalThresholdPercentage:     convert.GraphToFrameworkInt32(config.GetRenewalThresholdPercentage()),
		CertificateStore:               convert.GraphToFrameworkEnum(config.GetCertificateStore()),
		CertificateValidityPeriodScale: convert.GraphToFrameworkEnum(config.GetCertificateValidityPeriodScale()),
		CertificateValidityPeriodValue: convert.GraphToFrameworkInt32(config.GetCertificateValidityPeriodValue()),
		SubjectNameFormat:              convert.GraphToFrameworkEnum(config.GetSubjectNameFormat()),
		SubjectNameFormatString:        convert.GraphToFrameworkString(config.GetSubjectNameFormatString()),
		CertificationAuthority:         convert.GraphToFrameworkString(config.GetCertificationAuthority()),
		CertificationAuthorityName:     convert.GraphToFrameworkString(config.GetCertificationAuthorityName()),
		CertificateTemplateName:        convert.GraphToFrameworkString(config.GetCertificateTemplateName()),
		CustomSubjectAlternativeNames:  mapCustomSANsToSet(ctx, config.GetCustomSubjectAlternativeNames()),
		AllowAllAppsAccess:             convert.GraphToFrameworkBool(config.GetAllowAllAppsAccess()),
	}

	// Create the attribute type map for PKCS certificate
	pkcsAttrTypes := map[string]attr.Type{
		"deployment_channel":                types.StringType,
		"renewal_threshold_percentage":      types.Int32Type,
		"certificate_store":                 types.StringType,
		"certificate_validity_period_scale": types.StringType,
		"certificate_validity_period_value": types.Int32Type,
		"subject_name_format":               types.StringType,
		"subject_name_format_string":        types.StringType,
		"certification_authority":           types.StringType,
		"certification_authority_name":      types.StringType,
		"certificate_template_name":         types.StringType,
		"custom_subject_alternative_names":  types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"san_type": types.StringType, "name": types.StringType}}},
		"allow_all_apps_access":             types.BoolType,
	}

	objectValue, diags := types.ObjectValueFrom(ctx, pkcsAttrTypes, pkcsModel)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to create PKCS certificate object", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return
	}

	data.PkcsCertificate = objectValue
	// Clear other configuration types
	data.CustomConfiguration = types.ObjectNull(CustomConfigurationType().AttrTypes)
	data.PreferenceFile = types.ObjectNull(PreferenceFileType().AttrTypes)
	data.TrustedCertificate = types.ObjectNull(TrustedCertificateType().AttrTypes)
	data.ScepCertificate = types.ObjectNull(ScepCertificateType().AttrTypes)
}

// Helper functions for complex mappings

func mapKeyUsageToSet(ctx context.Context, keyUsage *graphmodels.KeyUsages) types.Set {
	if keyUsage == nil {
		return types.SetNull(types.StringType)
	}

	// Parse combined key usage back to individual values
	var keyUsageStrings []string

	// Check each possible key usage bit
	if (*keyUsage & graphmodels.DIGITALSIGNATURE_KEYUSAGES) != 0 {
		keyUsageStrings = append(keyUsageStrings, "digitalSignature")
	}
	if (*keyUsage & graphmodels.KEYENCIPHERMENT_KEYUSAGES) != 0 {
		keyUsageStrings = append(keyUsageStrings, "keyEncipherment")
	}

	if len(keyUsageStrings) == 0 {
		return types.SetNull(types.StringType)
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, keyUsageStrings)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create key usage set", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return types.SetNull(types.StringType)
	}

	return setValue
}

func mapCustomSANsToSet(ctx context.Context, sans []graphmodels.CustomSubjectAlternativeNameable) types.Set {
	if len(sans) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
			"san_type": types.StringType,
			"name":     types.StringType,
		}})
	}

	var sanModels []CustomSubjectAlternativeNameResourceModel
	for _, san := range sans {
		if san != nil {
			sanModel := CustomSubjectAlternativeNameResourceModel{
				SanType: convert.GraphToFrameworkEnum(san.GetSanType()),
				Name:    convert.GraphToFrameworkString(san.GetName()),
			}
			sanModels = append(sanModels, sanModel)
		}
	}

	setValue, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"san_type": types.StringType,
		"name":     types.StringType,
	}}, sanModels)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to create custom SANs set", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
			"san_type": types.StringType,
			"name":     types.StringType,
		}})
	}

	return setValue
}

func mapExtendedKeyUsagesToSet(ctx context.Context, ekus []graphmodels.ExtendedKeyUsageable) types.Set {
	if len(ekus) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":              types.StringType,
			"object_identifier": types.StringType,
		}})
	}

	var ekuModels []ExtendedKeyUsageResourceModel
	for _, eku := range ekus {
		if eku != nil {
			ekuModel := ExtendedKeyUsageResourceModel{
				Name:             convert.GraphToFrameworkString(eku.GetName()),
				ObjectIdentifier: convert.GraphToFrameworkString(eku.GetObjectIdentifier()),
			}
			ekuModels = append(ekuModels, ekuModel)
		}
	}

	setValue, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":              types.StringType,
		"object_identifier": types.StringType,
	}}, ekuModels)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to create extended key usages set", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":              types.StringType,
			"object_identifier": types.StringType,
		}})
	}

	return setValue
}
