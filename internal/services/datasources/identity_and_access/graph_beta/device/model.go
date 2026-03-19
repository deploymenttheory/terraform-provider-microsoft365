// REF: https://learn.microsoft.com/en-us/graph/api/resources/device?view=graph-rest-beta
package graphBetaDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceDataSourceModel represents the Terraform data source model for Microsoft Entra Devices
type DeviceDataSourceModel struct {
	ID                       types.String          `tfsdk:"id"`
	ObjectId                 types.String          `tfsdk:"object_id"`
	DisplayName              types.String          `tfsdk:"display_name"`
	DeviceId                 types.String          `tfsdk:"device_id"`
	ListAll                  types.Bool            `tfsdk:"list_all"`
	ListMemberOf             types.Bool            `tfsdk:"list_member_of"`
	ListRegisteredOwners     types.Bool            `tfsdk:"list_registered_owners"`
	ListRegisteredUsers      types.Bool            `tfsdk:"list_registered_users"`
	ODataQuery               types.String          `tfsdk:"odata_query"`
	Items                    []DeviceItemModel     `tfsdk:"items"`
	MemberOf                 []DirectoryObjectItem `tfsdk:"member_of"`
	RegisteredOwners         []DirectoryObjectItem `tfsdk:"registered_owners"`
	RegisteredUsers          []DirectoryObjectItem `tfsdk:"registered_users"`
	Timeouts                 timeouts.Value        `tfsdk:"timeouts"`
}

// DeviceItemModel represents an individual Microsoft Entra Device
type DeviceItemModel struct {
	ID                                types.String                         `tfsdk:"id"`
	AccountEnabled                    types.Bool                           `tfsdk:"account_enabled"`
	AlternativeSecurityIds            []AlternativeSecurityId              `tfsdk:"alternative_security_ids"`
	ApproximateLastSignInDateTime     types.String                         `tfsdk:"approximate_last_sign_in_date_time"`
	ComplianceExpirationDateTime      types.String                         `tfsdk:"compliance_expiration_date_time"`
	DeviceCategory                    types.String                         `tfsdk:"device_category"`
	DeviceId                          types.String                         `tfsdk:"device_id"`
	DeviceMetadata                    types.String                         `tfsdk:"device_metadata"`
	DeviceOwnership                   types.String                         `tfsdk:"device_ownership"`
	DeviceVersion                     types.Int64                          `tfsdk:"device_version"`
	DisplayName                       types.String                         `tfsdk:"display_name"`
	DomainName                        types.String                         `tfsdk:"domain_name"`
	EnrollmentProfileName             types.String                         `tfsdk:"enrollment_profile_name"`
	EnrollmentType                    types.String                         `tfsdk:"enrollment_type"`
	ExtensionAttributes               *OnPremisesExtensionAttributes       `tfsdk:"extension_attributes"`
	IsCompliant                       types.Bool                           `tfsdk:"is_compliant"`
	IsManaged                         types.Bool                           `tfsdk:"is_managed"`
	IsManagementRestricted            types.Bool                           `tfsdk:"is_management_restricted"`
	IsRooted                          types.Bool                           `tfsdk:"is_rooted"`
	ManagementType                    types.String                         `tfsdk:"management_type"`
	Manufacturer                      types.String                         `tfsdk:"manufacturer"`
	MdmAppId                          types.String                         `tfsdk:"mdm_app_id"`
	Model                             types.String                         `tfsdk:"model"`
	OnPremisesLastSyncDateTime        types.String                         `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesSecurityIdentifier      types.String                         `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled             types.Bool                           `tfsdk:"on_premises_sync_enabled"`
	OperatingSystem                   types.String                         `tfsdk:"operating_system"`
	OperatingSystemVersion            types.String                         `tfsdk:"operating_system_version"`
	PhysicalIds                       []types.String                       `tfsdk:"physical_ids"`
	ProfileType                       types.String                         `tfsdk:"profile_type"`
	RegistrationDateTime              types.String                         `tfsdk:"registration_date_time"`
	SystemLabels                      []types.String                       `tfsdk:"system_labels"`
	TrustType                         types.String                         `tfsdk:"trust_type"`
}

// AlternativeSecurityId represents an alternative security identifier
type AlternativeSecurityId struct {
	Type             types.Int64  `tfsdk:"type"`
	IdentityProvider types.String `tfsdk:"identity_provider"`
	Key              types.String `tfsdk:"key"`
}

// OnPremisesExtensionAttributes represents extension attributes for on-premises sync
type OnPremisesExtensionAttributes struct {
	ExtensionAttribute1  types.String `tfsdk:"extension_attribute1"`
	ExtensionAttribute2  types.String `tfsdk:"extension_attribute2"`
	ExtensionAttribute3  types.String `tfsdk:"extension_attribute3"`
	ExtensionAttribute4  types.String `tfsdk:"extension_attribute4"`
	ExtensionAttribute5  types.String `tfsdk:"extension_attribute5"`
	ExtensionAttribute6  types.String `tfsdk:"extension_attribute6"`
	ExtensionAttribute7  types.String `tfsdk:"extension_attribute7"`
	ExtensionAttribute8  types.String `tfsdk:"extension_attribute8"`
	ExtensionAttribute9  types.String `tfsdk:"extension_attribute9"`
	ExtensionAttribute10 types.String `tfsdk:"extension_attribute10"`
	ExtensionAttribute11 types.String `tfsdk:"extension_attribute11"`
	ExtensionAttribute12 types.String `tfsdk:"extension_attribute12"`
	ExtensionAttribute13 types.String `tfsdk:"extension_attribute13"`
	ExtensionAttribute14 types.String `tfsdk:"extension_attribute14"`
	ExtensionAttribute15 types.String `tfsdk:"extension_attribute15"`
}

// DirectoryObjectItem represents a directory object (for memberOf, registeredOwners, registeredUsers)
type DirectoryObjectItem struct {
	ID          types.String `tfsdk:"id"`
	ODataType   types.String `tfsdk:"odata_type"`
	DisplayName types.String `tfsdk:"display_name"`
}
