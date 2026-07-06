# Enrollment time grouping (ETG) adds devices to a static Microsoft Entra security group as they
# enroll, via device_security_group. This requires two prerequisites, both shown below:
#
#  1. A static (Assigned) Microsoft Entra security group.
#  2. The "Intune Provisioning Client" service principal (AppId f1346770-5b25-470b-88bd-d5744ab7952c,
#     sometimes shown as "Intune Autopilot ConfidentialClient") set as an OWNER of that group, so
#     Intune can add enrolling devices to it.
#
# ~> Known Microsoft Graph limitation: setting/clearing device_security_group calls the
# setEnrollmentTimeDeviceMembershipTarget / clearEnrollmentTimeDeviceMembershipTarget actions,
# which currently return a 500 error from the Intune backend when called with application
# permissions (client credentials) - the auth flow this provider always uses. These calls succeed
# with delegated (signed-in user) permissions only. See the resource's "Known Issues" section.

data "microsoft365_graph_beta_applications_service_principal" "intune_provisioning_client" {
  app_id = "f1346770-5b25-470b-88bd-d5744ab7952c"
}

resource "microsoft365_graph_beta_groups_group" "ios_ade_enrollment_group" {
  display_name     = "ios-ade-enrollment-time-grouping"
  mail_nickname    = "iosadeenrollmenttimegrouping"
  mail_enabled     = false
  security_enabled = true

  # Membership type: Assigned (static) - omitting "DynamicMembership" from group_types keeps
  # this a static group. Microsoft Entra roles can be assigned to the group: No.
  is_assignable_to_role = false

  description = "Static security group used for iOS/iPadOS ADE enrollment time grouping."

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Give Microsoft Entra ID time to propagate the new group before granting ownership.
resource "time_sleep" "wait_for_group_propagation" {
  depends_on = [
    microsoft365_graph_beta_groups_group.ios_ade_enrollment_group,
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "ios_ade_enrollment_group_owner" {
  group_id          = microsoft365_graph_beta_groups_group.ios_ade_enrollment_group.id
  owner_id          = data.microsoft365_graph_beta_applications_service_principal.intune_provisioning_client.id
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_group_propagation]
}

# Give Microsoft Entra ID time to propagate the new ownership before the policy references the
# group (Intune validates ownership synchronously against Graph, which can lag after the write).
resource "time_sleep" "wait_for_owner_propagation" {
  depends_on = [
    microsoft365_graph_beta_groups_group_owner_assignment.ios_ade_enrollment_group_owner,
  ]

  create_duration = "60s"
}

resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "with_enrollment_time_grouping" {
  name = "iOS ADE - Enrollment Time Grouping"

  requires_user_authentication = false

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  # Only valid when the provider is configured with delegated (user) credentials - see the
  # Known Issues note above.
  device_security_group = microsoft365_graph_beta_groups_group.ios_ade_enrollment_group.id

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_owner_propagation]
}
