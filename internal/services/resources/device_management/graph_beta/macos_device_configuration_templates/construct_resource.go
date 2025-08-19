package graphBetaMacosDeviceConfigurationTemplates

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the macOS configuration template resource for the Terraform provider.
func constructResource(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel) (graphmodels.DeviceConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var requestBody graphmodels.DeviceConfigurationable

	// Determine which configuration type to construct
	if !data.CustomConfiguration.IsNull() && !data.CustomConfiguration.IsUnknown() {
		requestBody = constructMacOSCustomConfiguration(ctx, data)
	} else if !data.PreferenceFile.IsNull() && !data.PreferenceFile.IsUnknown() {
		requestBody = constructMacOSCustomAppConfiguration(ctx, data)
	} else if !data.TrustedCertificate.IsNull() && !data.TrustedCertificate.IsUnknown() {
		requestBody = constructMacOSTrustedRootCertificate(ctx, data)
	} else if !data.ScepCertificate.IsNull() && !data.ScepCertificate.IsUnknown() {
		requestBody = constructMacOSScepCertificateProfile(ctx, data)
	} else if !data.PkcsCertificate.IsNull() && !data.PkcsCertificate.IsUnknown() {
		requestBody = constructMacOSPkcsCertificateProfile(ctx, data)
	} else {
		return nil, fmt.Errorf("no configuration type specified")
	}

	if requestBody == nil {
		return nil, fmt.Errorf("failed to construct configuration")
	}

	// Set common properties
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructMacOSCustomConfiguration constructs a MacOSCustomConfiguration
func constructMacOSCustomConfiguration(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel) graphmodels.DeviceConfigurationable {
	tflog.Debug(ctx, "Constructing MacOSCustomConfiguration")

	customConfig := graphmodels.NewMacOSCustomConfiguration()

	var customConfigData CustomConfigurationResourceModel
	diags := data.CustomConfiguration.As(ctx, &customConfigData, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		tflog.Error(ctx, "Failed to extract custom configuration data")
		return nil
	}

	if err := convert.FrameworkToGraphEnum(customConfigData.DeploymentChannel, graphmodels.ParseAppleDeploymentChannel, customConfig.SetDeploymentChannel); err != nil {
		tflog.Error(ctx, "Failed to set deployment channel", map[string]interface{}{"error": err.Error()})
		return nil
	}

	convert.FrameworkToGraphString(customConfigData.PayloadFileName, customConfig.SetPayloadFileName)
	convert.FrameworkToGraphBytes(customConfigData.Payload, customConfig.SetPayload)
	convert.FrameworkToGraphString(customConfigData.PayloadName, customConfig.SetPayloadName)

	return customConfig
}

// constructMacOSCustomAppConfiguration constructs a MacOSCustomAppConfiguration
func constructMacOSCustomAppConfiguration(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel) graphmodels.DeviceConfigurationable {
	tflog.Debug(ctx, "Constructing MacOSCustomAppConfiguration")

	appConfig := graphmodels.NewMacOSCustomAppConfiguration()

	var preferenceData PreferenceFileResourceModel
	diags := data.PreferenceFile.As(ctx, &preferenceData, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		tflog.Error(ctx, "Failed to extract preference file data")
		return nil
	}

	convert.FrameworkToGraphString(preferenceData.FileName, appConfig.SetFileName)
	convert.FrameworkToGraphBytes(preferenceData.ConfigurationXml, appConfig.SetConfigurationXml)
	convert.FrameworkToGraphString(preferenceData.BundleId, appConfig.SetBundleId)

	return appConfig
}

// constructMacOSTrustedRootCertificate constructs a MacOSTrustedRootCertificate
func constructMacOSTrustedRootCertificate(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel) graphmodels.DeviceConfigurationable {
	tflog.Debug(ctx, "Constructing MacOSTrustedRootCertificate")

	certConfig := graphmodels.NewMacOSTrustedRootCertificate()

	var certData TrustedCertificateResourceModel
	diags := data.TrustedCertificate.As(ctx, &certData, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		tflog.Error(ctx, "Failed to extract trusted certificate data")
		return nil
	}

	if err := convert.FrameworkToGraphEnum(certData.DeploymentChannel, graphmodels.ParseAppleDeploymentChannel, certConfig.SetDeploymentChannel); err != nil {
		tflog.Error(ctx, "Failed to set deployment channel", map[string]interface{}{"error": err.Error()})
		return nil
	}

	convert.FrameworkToGraphString(certData.CertFileName, certConfig.SetCertFileName)
	
	// Handle base64-encoded certificate data from filebase64()
	if !certData.TrustedRootCertificate.IsNull() && !certData.TrustedRootCertificate.IsUnknown() {
		certBase64 := certData.TrustedRootCertificate.ValueString()
		if certBytes, err := base64.StdEncoding.DecodeString(certBase64); err == nil {
			certConfig.SetTrustedRootCertificate(certBytes)
		} else {
			tflog.Error(ctx, "Failed to decode base64 certificate data", map[string]interface{}{"error": err.Error()})
			return nil
		}
	}

	return certConfig
}

// constructMacOSScepCertificateProfile constructs a MacOSScepCertificateProfile
func constructMacOSScepCertificateProfile(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel) graphmodels.DeviceConfigurationable {
	tflog.Debug(ctx, "Constructing MacOSScepCertificateProfile")

	scepConfig := graphmodels.NewMacOSScepCertificateProfile()

	var scepData ScepCertificateResourceModel
	diags := data.ScepCertificate.As(ctx, &scepData, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		tflog.Error(ctx, "Failed to extract SCEP certificate data")
		return nil
	}

	if err := convert.FrameworkToGraphEnum(scepData.DeploymentChannel, graphmodels.ParseAppleDeploymentChannel, scepConfig.SetDeploymentChannel); err != nil {
		tflog.Error(ctx, "Failed to set deployment channel", map[string]interface{}{"error": err.Error()})
		return nil
	}

	convert.FrameworkToGraphInt32(scepData.RenewalThresholdPercentage, scepConfig.SetRenewalThresholdPercentage)

	if err := convert.FrameworkToGraphEnum(scepData.CertificateStore, graphmodels.ParseCertificateStore, scepConfig.SetCertificateStore); err != nil {
		tflog.Error(ctx, "Failed to set certificate store", map[string]interface{}{"error": err.Error()})
		return nil
	}

	if err := convert.FrameworkToGraphEnum(scepData.CertificateValidityPeriodScale, graphmodels.ParseCertificateValidityPeriodScale, scepConfig.SetCertificateValidityPeriodScale); err != nil {
		tflog.Error(ctx, "Failed to set certificate validity period scale", map[string]interface{}{"error": err.Error()})
		return nil
	}

	convert.FrameworkToGraphInt32(scepData.CertificateValidityPeriodValue, scepConfig.SetCertificateValidityPeriodValue)

	if err := convert.FrameworkToGraphEnum(scepData.SubjectNameFormat, graphmodels.ParseAppleSubjectNameFormat, scepConfig.SetSubjectNameFormat); err != nil {
		tflog.Error(ctx, "Failed to set subject name format", map[string]interface{}{"error": err.Error()})
		return nil
	}

	convert.FrameworkToGraphString(scepData.SubjectNameFormatString, scepConfig.SetSubjectNameFormatString)

	// Root certificate needs to be set as an odata bind reference
	if !scepData.RootCertificateOdataBind.IsNull() && !scepData.RootCertificateOdataBind.IsUnknown() {
		// Extract just the ID if the full URL is provided
		rootCertId := scepData.RootCertificateOdataBind.ValueString()

		// Check if it's already a full URL or just an ID
		if !strings.HasPrefix(rootCertId, "https://") {
			// If it's just an ID, construct the full OData binding URL
			rootCertId = fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('%s')", rootCertId)
		}

		// Use additionalData to set the OData bind reference
		additionalData := map[string]interface{}{
			"rootCertificate@odata.bind": rootCertId,
		}
		scepConfig.SetAdditionalData(additionalData)
	}

	if err := convert.FrameworkToGraphEnum(scepData.KeySize, graphmodels.ParseKeySize, scepConfig.SetKeySize); err != nil {
		tflog.Error(ctx, "Failed to set key size", map[string]interface{}{"error": err.Error()})
		return nil
	}

	// Handle key usage set - SCEP needs combined key usage value
	if err := convertKeyUsageSet(ctx, scepData.KeyUsage, func(usages []graphmodels.KeyUsages) {
		if len(usages) > 0 {
			// Combine all key usages using bitwise OR
			var combinedUsage graphmodels.KeyUsages = usages[0]
			for i := 1; i < len(usages); i++ {
				combinedUsage = combinedUsage | usages[i]
			}
			scepConfig.SetKeyUsage(&combinedUsage)
		}
	}); err != nil {
		tflog.Error(ctx, "Failed to set key usage", map[string]interface{}{"error": err.Error()})
		return nil
	}

	// Handle custom subject alternative names
	if err := convertCustomSubjectAlternativeNames(ctx, scepData.CustomSubjectAlternativeNames, scepConfig.SetCustomSubjectAlternativeNames); err != nil {
		tflog.Error(ctx, "Failed to set custom subject alternative names", map[string]interface{}{"error": err.Error()})
		return nil
	}

	// Handle extended key usages
	if err := convertExtendedKeyUsages(ctx, scepData.ExtendedKeyUsages, scepConfig.SetExtendedKeyUsages); err != nil {
		tflog.Error(ctx, "Failed to set extended key usages", map[string]interface{}{"error": err.Error()})
		return nil
	}

	if err := convert.FrameworkToGraphStringSet(ctx, scepData.ScepServerUrls, scepConfig.SetScepServerUrls); err != nil {
		tflog.Error(ctx, "Failed to set SCEP server URLs", map[string]interface{}{"error": err.Error()})
		return nil
	}

	convert.FrameworkToGraphBool(scepData.AllowAllAppsAccess, scepConfig.SetAllowAllAppsAccess)

	return scepConfig
}

// constructMacOSPkcsCertificateProfile constructs a MacOSPkcsCertificateProfile
func constructMacOSPkcsCertificateProfile(ctx context.Context, data *MacosDeviceConfigurationTemplatesResourceModel) graphmodels.DeviceConfigurationable {
	tflog.Debug(ctx, "Constructing MacOSPkcsCertificateProfile")

	pkcsConfig := graphmodels.NewMacOSPkcsCertificateProfile()

	var pkcsData PkcsCertificateResourceModel
	diags := data.PkcsCertificate.As(ctx, &pkcsData, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		tflog.Error(ctx, "Failed to extract PKCS certificate data")
		return nil
	}

	if err := convert.FrameworkToGraphEnum(pkcsData.DeploymentChannel, graphmodels.ParseAppleDeploymentChannel, pkcsConfig.SetDeploymentChannel); err != nil {
		tflog.Error(ctx, "Failed to set deployment channel", map[string]interface{}{"error": err.Error(), "value": pkcsData.DeploymentChannel.ValueString()})
		return nil
	}

	convert.FrameworkToGraphInt32(pkcsData.RenewalThresholdPercentage, pkcsConfig.SetRenewalThresholdPercentage)

	if err := convert.FrameworkToGraphEnum(pkcsData.CertificateStore, graphmodels.ParseCertificateStore, pkcsConfig.SetCertificateStore); err != nil {
		tflog.Error(ctx, "Failed to set certificate store", map[string]interface{}{"error": err.Error(), "value": pkcsData.CertificateStore.ValueString()})
		return nil
	}

	if err := convert.FrameworkToGraphEnum(pkcsData.CertificateValidityPeriodScale, graphmodels.ParseCertificateValidityPeriodScale, pkcsConfig.SetCertificateValidityPeriodScale); err != nil {
		tflog.Error(ctx, "Failed to set certificate validity period scale", map[string]interface{}{"error": err.Error(), "value": pkcsData.CertificateValidityPeriodScale.ValueString()})
		return nil
	}

	convert.FrameworkToGraphInt32(pkcsData.CertificateValidityPeriodValue, pkcsConfig.SetCertificateValidityPeriodValue)

	if err := convert.FrameworkToGraphEnum(pkcsData.SubjectNameFormat, graphmodels.ParseAppleSubjectNameFormat, pkcsConfig.SetSubjectNameFormat); err != nil {
		tflog.Error(ctx, "Failed to set subject name format", map[string]interface{}{"error": err.Error(), "value": pkcsData.SubjectNameFormat.ValueString()})
		return nil
	}

	convert.FrameworkToGraphString(pkcsData.SubjectNameFormatString, pkcsConfig.SetSubjectNameFormatString)
	convert.FrameworkToGraphString(pkcsData.CertificationAuthority, pkcsConfig.SetCertificationAuthority)
	convert.FrameworkToGraphString(pkcsData.CertificationAuthorityName, pkcsConfig.SetCertificationAuthorityName)
	convert.FrameworkToGraphString(pkcsData.CertificateTemplateName, pkcsConfig.SetCertificateTemplateName)

	if err := convertCustomSubjectAlternativeNames(ctx, pkcsData.CustomSubjectAlternativeNames, pkcsConfig.SetCustomSubjectAlternativeNames); err != nil {
		tflog.Error(ctx, "Failed to set custom subject alternative names", map[string]interface{}{"error": err.Error()})
		return nil
	}

	convert.FrameworkToGraphBool(pkcsData.AllowAllAppsAccess, pkcsConfig.SetAllowAllAppsAccess)

	return pkcsConfig
}

// Helper functions for complex conversions

func convertKeyUsageSet(ctx context.Context, keyUsageSet types.Set, setter func([]graphmodels.KeyUsages)) error {
	if keyUsageSet.IsNull() || keyUsageSet.IsUnknown() {
		return nil
	}

	var keyUsageStrings []string
	diags := keyUsageSet.ElementsAs(ctx, &keyUsageStrings, false)
	if diags.HasError() {
		return fmt.Errorf("failed to extract key usage strings: %v", diags.Errors())
	}

	keyUsages := make([]graphmodels.KeyUsages, 0, len(keyUsageStrings))
	for _, usage := range keyUsageStrings {
		if parsedUsage, err := graphmodels.ParseKeyUsages(usage); err == nil {
			if keyUsage, ok := parsedUsage.(*graphmodels.KeyUsages); ok && keyUsage != nil {
				keyUsages = append(keyUsages, *keyUsage)
			}
		} else {
			return fmt.Errorf("invalid key usage: %s", usage)
		}
	}

	setter(keyUsages)
	return nil
}

func convertCustomSubjectAlternativeNames(ctx context.Context, sanSet types.Set, setter func([]graphmodels.CustomSubjectAlternativeNameable)) error {
	if sanSet.IsNull() || sanSet.IsUnknown() {
		return nil
	}

	var sanModels []CustomSubjectAlternativeNameResourceModel
	diags := sanSet.ElementsAs(ctx, &sanModels, false)
	if diags.HasError() {
		return fmt.Errorf("failed to extract custom subject alternative names: %v", diags.Errors())
	}

	sans := make([]graphmodels.CustomSubjectAlternativeNameable, 0, len(sanModels))
	for _, sanModel := range sanModels {
		san := graphmodels.NewCustomSubjectAlternativeName()

		if sanType, err := graphmodels.ParseSubjectAlternativeNameType(sanModel.SanType.ValueString()); err == nil {
			if sanTypePtr, ok := sanType.(*graphmodels.SubjectAlternativeNameType); ok && sanTypePtr != nil {
				san.SetSanType(sanTypePtr)
			}
		} else {
			return fmt.Errorf("invalid SAN type: %s", sanModel.SanType.ValueString())
		}

		convert.FrameworkToGraphString(sanModel.Name, san.SetName)
		sans = append(sans, san)
	}

	setter(sans)
	return nil
}

func convertExtendedKeyUsages(ctx context.Context, ekuSet types.Set, setter func([]graphmodels.ExtendedKeyUsageable)) error {
	if ekuSet.IsNull() || ekuSet.IsUnknown() {
		return nil
	}

	var ekuModels []ExtendedKeyUsageResourceModel
	diags := ekuSet.ElementsAs(ctx, &ekuModels, false)
	if diags.HasError() {
		return fmt.Errorf("failed to extract extended key usages: %v", diags.Errors())
	}

	ekus := make([]graphmodels.ExtendedKeyUsageable, 0, len(ekuModels))
	for _, ekuModel := range ekuModels {
		eku := graphmodels.NewExtendedKeyUsage()
		convert.FrameworkToGraphString(ekuModel.Name, eku.SetName)
		convert.FrameworkToGraphString(ekuModel.ObjectIdentifier, eku.SetObjectIdentifier)
		ekus = append(ekus, eku)
	}

	setter(ekus)
	return nil
}
