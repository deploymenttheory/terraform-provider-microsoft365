package sentinels

import "errors"

// Validation-related sentinel errors used across resource validation functions.
// These errors are used when validating Terraform plan data before API requests.
var (
	// ErrInvalidTenantID indicates an invalid Microsoft Entra organization tenant ID
	ErrInvalidTenantID = errors.New("invalid Microsoft Entra organization tenant ID")

	// ErrInvalidRoleGUIDs indicates invalid role GUIDs that don't exist in the tenant
	ErrInvalidRoleGUIDs = errors.New("invalid role GUIDs found")

	// ErrInvalidUserGUIDs indicates invalid user GUIDs that don't exist in the tenant
	ErrInvalidUserGUIDs = errors.New("invalid user GUIDs found")

	// ErrInvalidGroupGUIDs indicates invalid group GUIDs that don't exist in the tenant
	ErrInvalidGroupGUIDs = errors.New("group IDs not found in tenant")

	// ErrInvalidLocationGUIDs indicates invalid named location GUIDs that don't exist in the tenant
	ErrInvalidLocationGUIDs = errors.New("invalid location GUIDs found")

	// ErrUserNotFound indicates a user was not found in the tenant
	ErrUserNotFound = errors.New("user not found")

	// ErrGroupNotFound indicates a group was not found in the tenant
	ErrGroupNotFound = errors.New("group not found")

	// ErrEmptyDisplayName indicates display_name field is empty
	ErrEmptyDisplayName = errors.New("display_name cannot be empty")

	// ErrDuplicateDisplayName indicates a resource with the same display_name already exists
	ErrDuplicateDisplayName = errors.New("resource with this display_name already exists")

	// ErrInvalidResourceOperations indicates invalid resource operation IDs
	ErrInvalidResourceOperations = errors.New("invalid resource operation ID(s)")

	// ErrNoLifecyclePolicy indicates no lifecycle policy exists in the tenant
	ErrNoLifecyclePolicy = errors.New("no lifecycle policy found")

	// ErrLifecyclePolicyIDNull indicates the lifecycle policy ID is null
	ErrLifecyclePolicyIDNull = errors.New("policy ID is null")

	// ErrLifecyclePolicyNotSelected indicates lifecycle policy managedGroupTypes is not 'Selected'
	ErrLifecyclePolicyNotSelected = errors.New("lifecycle policy managedGroupTypes is not 'Selected'")

	// ErrGroupNotM365 indicates a group is not a Microsoft 365 group
	ErrGroupNotM365 = errors.New("group is not a Microsoft 365 group")

	// ErrInvalidPolicySetItemType indicates an unsupported policy set item type
	ErrInvalidPolicySetItemType = errors.New("unsupported policy set item type")

	// ErrInvalidMobileAppID indicates an invalid Intune mobile app ID
	ErrInvalidMobileAppID = errors.New("not a valid Intune mobile app ID")

	// ErrInvalidManagedAppPolicyID indicates an invalid Intune managed app policy ID
	ErrInvalidManagedAppPolicyID = errors.New("not a valid Intune managed app policy ID")

	// ErrInvalidDeviceConfigurationID indicates an invalid Intune device configuration ID
	ErrInvalidDeviceConfigurationID = errors.New("not a valid Intune device configuration ID")

	// ErrInvalidAutopilotProfileID indicates an invalid Intune Windows Autopilot deployment profile ID
	ErrInvalidAutopilotProfileID = errors.New("not a valid Intune deployment profile ID")

	// ErrInvalidCompliancePolicyID indicates an invalid Intune device compliance policy ID
	ErrInvalidCompliancePolicyID = errors.New("not a valid Intune compliance policy ID")

	// Conditional Access Policy specific errors

	// ErrEmptyUserInclusions indicates user inclusion assignments are empty when they should contain All or None
	ErrEmptyUserInclusions = errors.New("when conditional access policy user inclusion assignments are empty for 'include_users', 'include_groups', 'include_roles', and 'include_guests_or_external_users', then 'include_users' must be either 'All' or 'None'")

	// ErrInvalidUserInclusionValue indicates user inclusions contain invalid values when empty
	ErrInvalidUserInclusionValue = errors.New("when conditional access policy user inclusion assignments are empty for 'include_users', 'include_groups', 'include_roles', and 'include_guests_or_external_users', then 'include_users' must contain only 'All' or 'None'")

	// ErrUserInclusionWithSpecificAssignments indicates user inclusions cannot contain All/None when specific assignments exist
	ErrUserInclusionWithSpecificAssignments = errors.New("when conditional access policy has specific user inclusion assignments configured for 'include_groups', 'include_roles', or 'include_guests_or_external_users', then 'include_users' cannot contain 'All' or 'None'")

	// ErrEmptyApplicationFields indicates all application fields are empty when include_applications should be None
	ErrEmptyApplicationFields = errors.New("when conditional access policy application fields 'include_applications', 'exclude_applications', 'include_user_actions', and 'include_authentication_context_class_references' are all empty, then 'include_applications' must be set to 'None'")

	// ErrApplicationFilterWithSpecialValues indicates application_filter cannot be used with All/None/AllAgentIdResources
	ErrApplicationFilterWithSpecialValues = errors.New("conditional access policy 'application_filter' cannot be used when 'include_applications' contains 'All', 'None', or 'AllAgentIdResources' values. It can be used with GUID values or 'Office365'")

	// ErrApplicationsAndUserActions indicates both include_applications and include_user_actions are configured
	ErrApplicationsAndUserActions = errors.New("conditional access policy cannot have both 'include_applications' and 'include_user_actions' configured at the same time")

	// ErrApplicationsAndAuthContext indicates both include_applications and include_authentication_context_class_references are configured
	ErrApplicationsAndAuthContext = errors.New("conditional access policy cannot have both 'include_applications' and 'include_authentication_context_class_references' configured at the same time")

	// ErrFrequencyIntervalEveryTimeWithType indicates type field is set when frequency_interval is everyTime
	ErrFrequencyIntervalEveryTimeWithType = errors.New("when 'frequency_interval' is set to 'everyTime', the 'type' field must not be set in the configuration")

	// ErrFrequencyIntervalEveryTimeWithValue indicates value field is set when frequency_interval is everyTime
	ErrFrequencyIntervalEveryTimeWithValue = errors.New("when 'frequency_interval' is set to 'everyTime', the 'value' field must not be set in the configuration")

	// Agent-related errors

	// ErrSponsorUserIDsNullOrUnknown indicates sponsor_user_ids is null or unknown
	ErrSponsorUserIDsNullOrUnknown = errors.New("sponsor_user_ids cannot be null or unknown")

	// ErrAtLeastOneSponsorRequired indicates at least one sponsor is required
	ErrAtLeastOneSponsorRequired = errors.New("at least one sponsor is required")

	// ErrSponsorUserObjectNull indicates a sponsor user object is null
	ErrSponsorUserObjectNull = errors.New("sponsor user object is null")

	// ErrOwnerUserIDsNullOrUnknown indicates owner_user_ids is null or unknown
	ErrOwnerUserIDsNullOrUnknown = errors.New("owner_user_ids cannot be null or unknown")

	// ErrAtLeastOneOwnerRequired indicates at least one owner is required
	ErrAtLeastOneOwnerRequired = errors.New("at least one owner is required")

	// ErrOwnerUserObjectNull indicates an owner user object is null
	ErrOwnerUserObjectNull = errors.New("owner user object is null")

	// ErrAgentInstanceIDEmpty indicates agent_instance_id is empty
	ErrAgentInstanceIDEmpty = errors.New("agent_instance_id cannot be empty")

	// ErrNoAgentInstancesFound indicates no agent instances were found
	ErrNoAgentInstancesFound = errors.New("no agent instances found")

	// ErrAgentInstanceNotFound indicates an agent instance was not found
	ErrAgentInstanceNotFound = errors.New("agent instance does not exist")

	// Credential validation errors

	// ErrStartDateInPast indicates credential start_date_time is in the past
	ErrStartDateInPast = errors.New("start_date_time is in the past")

	// ErrEndDateBeforeStartDate indicates end_date_time is before start_date_time
	ErrEndDateBeforeStartDate = errors.New("end_date_time must be after start_date_time")

	// Authentication context errors

	// ErrAuthContextIDExists indicates an authentication context class reference ID already exists
	ErrAuthContextIDExists = errors.New("authentication context class reference with this ID already exists")

	// Group policy definition errors

	// ErrNoPresentationsFound indicates no presentations found for a policy
	ErrNoPresentationsFound = errors.New("no presentations found for policy")

	// ErrLabelNotFound indicates a label was not found in policy presentations
	ErrLabelNotFound = errors.New("label not found in policy")

	// ErrCheckboxRequiresBooleanValue indicates a checkbox requires a boolean value
	ErrCheckboxRequiresBooleanValue = errors.New("checkbox requires a boolean value ('true' or 'false')")

	// ErrDropdownRequiresNonEmptyValue indicates a dropdown/combobox requires a non-empty value
	ErrDropdownRequiresNonEmptyValue = errors.New("dropdown/combobox requires a non-empty value")

	// ErrPresentationReadOnly indicates a presentation is read-only and cannot have a value set
	ErrPresentationReadOnly = errors.New("presentation is read-only and cannot have a value set")

	// ErrUnsupportedPresentationType indicates an unsupported presentation type
	ErrUnsupportedPresentationType = errors.New("unsupported presentation type")

	// Cloud PC role definition errors

	// ErrInvalidCloudPCResourceOperations indicates invalid Cloud PC resource operation IDs
	ErrInvalidCloudPCResourceOperations = errors.New("invalid Cloud PC resource operation(s)")

	// ErrRoleDefinitionNameExists indicates a role definition with the same display name already exists
	ErrRoleDefinitionNameExists = errors.New("a role definition with this display name already exists")

	// ErrEntraDeviceValidationFailed indicates Entra ID device validation failed
	ErrEntraDeviceValidationFailed = errors.New("entra device validation failed")

	// ErrEntraDeviceNotFound indicates an Entra ID device was not found
	ErrEntraDeviceNotFound = errors.New("entra device not found")

	// ErrRetrieveEntraDevice indicates failure to retrieve an Entra ID device
	ErrRetrieveEntraDevice = errors.New("failed to retrieve Entra ID device")

	// ErrNoDevicesEnrolled indicates no devices are enrolled for an update category
	ErrNoDevicesEnrolled = errors.New("no devices enrolled for update category")

	// ErrUpdatableAssetGroupNotFound indicates a Windows Updates updatable asset group was not found
	ErrUpdatableAssetGroupNotFound = errors.New("updatable asset group not found")

	// ErrUpdatableAssetGroupValidationFailed indicates validation of a Windows Updates updatable asset group failed
	ErrUpdatableAssetGroupValidationFailed = errors.New("updatable asset group validation failed")

	// ErrInvalidEntraDeviceIDs indicates one or more Entra ID device IDs do not exist
	ErrInvalidEntraDeviceIDs = errors.New("one or more Entra ID device IDs not found")
)
