package graphBetaDevice

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructDeviceItems constructs a list of DeviceItemModel from a list of Deviceable
func ConstructDeviceItems(devices []graphmodels.Deviceable) []DeviceItemModel {
	if devices == nil {
		return []DeviceItemModel{}
	}

	items := make([]DeviceItemModel, 0, len(devices))
	for _, device := range devices {
		if device != nil {
			items = append(items, ConstructDeviceItem(device))
		}
	}

	return items
}

// ConstructDeviceItem constructs a DeviceItemModel from a Deviceable
func ConstructDeviceItem(device graphmodels.Deviceable) DeviceItemModel {
	item := DeviceItemModel{
		ID:                            types.StringPointerValue(device.GetId()),
		AccountEnabled:                types.BoolPointerValue(device.GetAccountEnabled()),
		ApproximateLastSignInDateTime: types.StringNull(),
		ComplianceExpirationDateTime:  types.StringNull(),
		DeviceCategory:                types.StringPointerValue(device.GetDeviceCategory()),
		DeviceId:                      types.StringPointerValue(device.GetDeviceId()),
		DeviceMetadata:                types.StringPointerValue(device.GetDeviceMetadata()),
		DeviceOwnership:               types.StringPointerValue(device.GetDeviceOwnership()),
		DeviceVersion:                 types.Int64Null(),
		DisplayName:                   types.StringPointerValue(device.GetDisplayName()),
		DomainName:                    types.StringPointerValue(device.GetDomainName()),
		EnrollmentProfileName:         types.StringPointerValue(device.GetEnrollmentProfileName()),
		EnrollmentType:                types.StringPointerValue(device.GetEnrollmentType()),
		IsCompliant:                   types.BoolPointerValue(device.GetIsCompliant()),
		IsManaged:                     types.BoolPointerValue(device.GetIsManaged()),
		IsManagementRestricted:        types.BoolPointerValue(device.GetIsManagementRestricted()),
		IsRooted:                      types.BoolPointerValue(device.GetIsRooted()),
		ManagementType:                types.StringPointerValue(device.GetManagementType()),
		Manufacturer:                  types.StringPointerValue(device.GetManufacturer()),
		MdmAppId:                      types.StringPointerValue(device.GetMdmAppId()),
		Model:                         types.StringPointerValue(device.GetModel()),
		OnPremisesLastSyncDateTime:    types.StringNull(),
		OnPremisesSecurityIdentifier:  types.StringPointerValue(device.GetOnPremisesSecurityIdentifier()),
		OnPremisesSyncEnabled:         types.BoolPointerValue(device.GetOnPremisesSyncEnabled()),
		OperatingSystem:               types.StringPointerValue(device.GetOperatingSystem()),
		OperatingSystemVersion:        types.StringPointerValue(device.GetOperatingSystemVersion()),
		ProfileType:                   types.StringPointerValue(device.GetProfileType()),
		RegistrationDateTime:          types.StringNull(),
		TrustType:                     types.StringPointerValue(device.GetTrustType()),
	}

	// Handle time fields
	if lastSignIn := device.GetApproximateLastSignInDateTime(); lastSignIn != nil {
		item.ApproximateLastSignInDateTime = types.StringValue(lastSignIn.Format("2006-01-02T15:04:05Z"))
	}

	if complianceExpiration := device.GetComplianceExpirationDateTime(); complianceExpiration != nil {
		item.ComplianceExpirationDateTime = types.StringValue(complianceExpiration.Format("2006-01-02T15:04:05Z"))
	}

	if lastSync := device.GetOnPremisesLastSyncDateTime(); lastSync != nil {
		item.OnPremisesLastSyncDateTime = types.StringValue(lastSync.Format("2006-01-02T15:04:05Z"))
	}

	if registrationTime := device.GetRegistrationDateTime(); registrationTime != nil {
		item.RegistrationDateTime = types.StringValue(registrationTime.Format("2006-01-02T15:04:05Z"))
	}

	// Handle device version
	if deviceVersion := device.GetDeviceVersion(); deviceVersion != nil {
		item.DeviceVersion = types.Int64Value(int64(*deviceVersion))
	}

	// Handle alternative security IDs
	if altSecIds := device.GetAlternativeSecurityIds(); len(altSecIds) > 0 {
		item.AlternativeSecurityIds = make([]AlternativeSecurityId, 0, len(altSecIds))
		for _, altSecId := range altSecIds {
			if altSecId != nil {
				item.AlternativeSecurityIds = append(item.AlternativeSecurityIds, AlternativeSecurityId{
					Type:             types.Int64Null(),
					IdentityProvider: types.StringPointerValue(altSecId.GetIdentityProvider()),
					Key:              types.StringNull(),
				})

				// Handle type
				if altSecIdType := altSecId.GetTypeEscaped(); altSecIdType != nil {
					item.AlternativeSecurityIds[len(item.AlternativeSecurityIds)-1].Type = types.Int64Value(int64(*altSecIdType))
				}

				// Handle key (base64 encoded)
				if key := altSecId.GetKey(); len(key) > 0 {
					item.AlternativeSecurityIds[len(item.AlternativeSecurityIds)-1].Key = types.StringValue(string(key))
				}
			}
		}
	}

	// Handle physical IDs
	if physicalIds := device.GetPhysicalIds(); len(physicalIds) > 0 {
		item.PhysicalIds = make([]types.String, 0, len(physicalIds))
		for _, physicalId := range physicalIds {
			item.PhysicalIds = append(item.PhysicalIds, types.StringValue(physicalId))
		}
	}

	// Handle system labels
	if systemLabels := device.GetSystemLabels(); len(systemLabels) > 0 {
		item.SystemLabels = make([]types.String, 0, len(systemLabels))
		for _, label := range systemLabels {
			item.SystemLabels = append(item.SystemLabels, types.StringValue(label))
		}
	}

	// Handle extension attributes
	if extAttrs := device.GetExtensionAttributes(); extAttrs != nil {
		item.ExtensionAttributes = &OnPremisesExtensionAttributes{
			ExtensionAttribute1:  types.StringPointerValue(extAttrs.GetExtensionAttribute1()),
			ExtensionAttribute2:  types.StringPointerValue(extAttrs.GetExtensionAttribute2()),
			ExtensionAttribute3:  types.StringPointerValue(extAttrs.GetExtensionAttribute3()),
			ExtensionAttribute4:  types.StringPointerValue(extAttrs.GetExtensionAttribute4()),
			ExtensionAttribute5:  types.StringPointerValue(extAttrs.GetExtensionAttribute5()),
			ExtensionAttribute6:  types.StringPointerValue(extAttrs.GetExtensionAttribute6()),
			ExtensionAttribute7:  types.StringPointerValue(extAttrs.GetExtensionAttribute7()),
			ExtensionAttribute8:  types.StringPointerValue(extAttrs.GetExtensionAttribute8()),
			ExtensionAttribute9:  types.StringPointerValue(extAttrs.GetExtensionAttribute9()),
			ExtensionAttribute10: types.StringPointerValue(extAttrs.GetExtensionAttribute10()),
			ExtensionAttribute11: types.StringPointerValue(extAttrs.GetExtensionAttribute11()),
			ExtensionAttribute12: types.StringPointerValue(extAttrs.GetExtensionAttribute12()),
			ExtensionAttribute13: types.StringPointerValue(extAttrs.GetExtensionAttribute13()),
			ExtensionAttribute14: types.StringPointerValue(extAttrs.GetExtensionAttribute14()),
			ExtensionAttribute15: types.StringPointerValue(extAttrs.GetExtensionAttribute15()),
		}
	}

	return item
}

// ConstructDirectoryObjectItems constructs a list of DirectoryObjectItem from a list of DirectoryObjectable
func ConstructDirectoryObjectItems(objects []graphmodels.DirectoryObjectable) []DirectoryObjectItem {
	if objects == nil {
		return []DirectoryObjectItem{}
	}

	items := make([]DirectoryObjectItem, 0, len(objects))
	for _, obj := range objects {
		if obj != nil {
			items = append(items, ConstructDirectoryObjectItem(obj))
		}
	}

	return items
}

// ConstructDirectoryObjectItem constructs a DirectoryObjectItem from a DirectoryObjectable
func ConstructDirectoryObjectItem(obj graphmodels.DirectoryObjectable) DirectoryObjectItem {
	item := DirectoryObjectItem{
		ID:        types.StringPointerValue(obj.GetId()),
		ODataType: types.StringPointerValue(obj.GetOdataType()),
	}

	// Try to get displayName from common directory object types
	// The DirectoryObject base type doesn't have displayName, but Group and User do
	if additionalData := obj.GetAdditionalData(); additionalData != nil {
		if displayName, ok := additionalData["displayName"].(string); ok {
			item.DisplayName = types.StringValue(displayName)
		}
	}

	return item
}
