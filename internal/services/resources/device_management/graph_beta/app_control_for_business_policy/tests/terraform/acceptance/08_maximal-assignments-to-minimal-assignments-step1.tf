# Test 08 Step 1: Maximal Assignments to Minimal Assignments - Initial State (Maximal Assignments)
# Purpose: Deploy minimal policy with maximal assignments (3 targets), will be reduced in step 2

# Dependencies
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Group Dependency
resource "microsoft365_graph_beta_groups_group" "test_group_1" {
  display_name     = "acc-test-acfb-group-1-${random_string.suffix.result}"
  description      = "Test group 1 for app control policy assignment testing"
  mail_nickname    = "acc-test-acfb-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Wait for group to be ready
resource "time_sleep" "wait_for_groups" {
  depends_on = [microsoft365_graph_beta_groups_group.test_group_1]

  create_duration = "30s"
}

# Start with Maximal Assignments
resource "microsoft365_graph_beta_device_management_app_control_for_business_policy" "assignment_step_test" {
  depends_on = [time_sleep.wait_for_groups]

  name        = "acc-test-app-control-policy-assignment-step-test"
  description = "Assignment step test policy - starts with maximal assignments"

  policy_xml = <<-EOT
<?xml version="1.0" encoding="utf-8"?>
<SiPolicy xmlns="urn:schemas-microsoft-com:sipolicy" PolicyType="Base Policy">
  <VersionEx>1.0.3.0</VersionEx>
  <PolicyID>{264C0644-19BE-418F-BAED-29E5E36250AD}</PolicyID>
  <BasePolicyID>{264C0644-19BE-418F-BAED-29E5E36250AD}</BasePolicyID>
  <PlatformID>{2E07F7E4-194C-4D20-B7C9-6F44A6C5A234}</PlatformID>
  <Rules>
    <Rule>
      <Option>Enabled:Unsigned System Integrity Policy</Option>
    </Rule>
    <Rule>
      <Option>Enabled:Advanced Boot Options Menu</Option>
    </Rule>
	<Rule>
      <Option>Enabled:UMCI</Option>
    </Rule>
	<Rule>
      <Option>Enabled:Update Policy No Reboot</Option>
    </Rule>
    <Rule>
      <Option>Enabled:Inherit Default Policy</Option>
    </Rule>
    <Rule>
      <Option>Enabled:Revoked Expired As Unsigned</Option>
    </Rule>
    <Rule>
      <Option>Disabled:Script Enforcement</Option>
    </Rule>
  </Rules>
  <EKUs />
  <FileRules>
    <Allow ID="ID_ALLOW_A_1" FileName="*" />
    <Allow ID="ID_ALLOW_A_2" FileName="*" />
  </FileRules>
  <Signers />
  <SigningScenarios>
    <SigningScenario Value="131" ID="ID_SIGNINGSCENARIO_DRIVERS_1" FriendlyName="Auto generated policy on 08-17-2015">
      <ProductSigners>
        <FileRulesRef>
          <FileRuleRef RuleID="ID_ALLOW_A_1" />
        </FileRulesRef>
      </ProductSigners>
    </SigningScenario>
    <SigningScenario Value="12" ID="ID_SIGNINGSCENARIO_WINDOWS" FriendlyName="Auto generated policy on 08-17-2015">
      <ProductSigners>
        <FileRulesRef>
          <FileRuleRef RuleID="ID_ALLOW_A_2" />
        </FileRulesRef>
      </ProductSigners>
    </SigningScenario>
  </SigningScenarios>
  <UpdatePolicySigners />
  <CiSigners />
  <HvciOptions>1</HvciOptions>
  <Settings>
	<Setting Provider="AllHostIds" Key="AllKeys" ValueName="EnterpriseDefinedClsId">
	  <Value>
        <Boolean>true</Boolean>
      </Value>
	</Setting>
    <Setting Provider="PolicyInfo" Key="Information" ValueName="Name">
      <Value>
        <String>AllowAllEnableHVCI</String>
      </Value>
    </Setting>
    <Setting Provider="PolicyInfo" Key="Information" ValueName="Id">
      <Value>
        <String>022422</String>
      </Value>
    </Setting>
  </Settings>
</SiPolicy>
  EOT

  role_scope_tag_ids = ["0"]

  # Maximal Assignments - 3 targets
  assignments = [
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.test_group_1.id
      filter_type = "none"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
