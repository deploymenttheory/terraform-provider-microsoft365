# example of inline XML policy
resource "microsoft365_graph_beta_device_management_app_control_for_business_policy" "inline" {
  name        = "unit-test-app-control-policy-inline"
  description = "unit-test-app-control-policy-inline"

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
  <!--EKUS-->
  <EKUs />
  <!--File Rules-->
  <FileRules>
    <Allow ID="ID_ALLOW_A_1" FileName="*" />
    <Allow ID="ID_ALLOW_A_2" FileName="*" />
  </FileRules>
  <!--Signers-->
  <Signers />
  <!--Driver Signing Scenarios-->
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

  # Role scope tags for testing
  role_scope_tag_ids = ["0", "1", "2"]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_id   = "44444444-4444-4444-4444-444444444444"
      filter_type = "include"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "55555555-5555-5555-5555-555555555555"
      filter_type = "exclude"
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}

# App Control for Business XML Policy configuration referencing a file
resource "microsoft365_graph_beta_device_management_app_control_for_business_policy" "file" {
  name        = "app-control-policy-file"
  description = "app-control-policy-file"

  # XML policy referencing a file
  policy_xml = file("${path.module}/AllowAll_EnableHVCI.xml")

  role_scope_tag_ids = ["0", "1", "2"]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.group_1.id
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.assignment_filter_1.id
      filter_type = "include"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.assignment_filter_2.id
      filter_type = "exclude"
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}