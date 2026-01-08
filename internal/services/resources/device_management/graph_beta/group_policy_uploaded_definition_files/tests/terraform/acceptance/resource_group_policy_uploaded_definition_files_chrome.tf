# admx/adml for chrome
resource "microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files" "test_update" {
  file_name = "chrome.admx"
  content   = <<-EOT
<?xml version="1.0" ?>
<policyDefinitions revision="1.0" schemaVersion="1.0">
  <policyNamespaces>
    <target namespace="Google.Policies.Chrome" prefix="Chrome"/>
  </policyNamespaces>
  <resources minRequiredRevision="1.0" />
  <categories>
    <category displayName="$(string.chrome)" name="Cat_Chrome"/>
  </categories>
</policyDefinitions>
EOT

  group_policy_uploaded_language_files = [
    {
      file_name     = "chrome.adml"
      language_code = "en-US"
      content       = <<-EOT
<?xml version="1.0" ?>
<policyDefinitionResources revision="1.0" schemaVersion="1.0">
  <displayName/>
  <description/>
  <resources>
    <stringTable>
      <string id="chrome">Google Chrome</string>
    </stringTable>
  </resources>
</policyDefinitionResources>
EOT
    }
  ]

  timeouts = {
    create = "10m" // allow a generous amount of time for the upload to complete
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

