# admx/adml for mozilla. which is a pre-req for the firefox files of the same
resource "microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files" "moziila" {
  file_name = "mozilla.admx"
  content   = <<-EOT
    <?xml version="1.0" ?>
    <policyDefinitions revision="4.8" schemaVersion="1.0">
      <policyNamespaces>
        <target namespace="Mozilla.Policies" prefix="Mozilla"/>
      </policyNamespaces>
      <resources minRequiredRevision="4.8" />
      <categories>
        <category displayName="$(string.mozilla)" name="Cat_Mozilla"/>
      </categories>
    </policyDefinitions>
  EOT

  group_policy_uploaded_language_files = [
    {
      file_name     = "mozilla.adml"
      language_code = "en-US"
      content       = <<-EOT
        <?xml version="1.0" ?>
        <policyDefinitionResources revision="4.8" schemaVersion="1.0">
          <displayName/>
          <description/>
          <resources>
            <stringTable>
              <string id="mozilla">Mozilla</string>
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