# https://www.base64decode.org/

# admx/adml for google
resource "microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files" "google" {
  file_name = "google.admx"
  content   = <<-EOT
<?xml version="1.0" ?>
<policyDefinitions revision="1.0" schemaVersion="1.0">
  <policyNamespaces>
    <target namespace="Google.Policies" prefix="Google"/>
  </policyNamespaces>
  <resources minRequiredRevision="1.0" />
  <categories>
    <category displayName="$(string.google)" name="Cat_Google"/>
  </categories>
</policyDefinitions>
EOT

  // or: "PD94bWwgdmVyc2lvbj0iMS4wIiA/Pg0KPHBvbGljeURlZmluaXRpb25zIHJldmlzaW9uPSIxLjAiIHNjaGVtYVZlcnNpb249IjEuMCI+DQogIDxwb2xpY3lOYW1lc3BhY2VzPg0KICAgIDx0YXJnZXQgbmFtZXNwYWNlPSJHb29nbGUuUG9saWNpZXMiIHByZWZpeD0iR29vZ2xlIi8+DQogIDwvcG9saWN5TmFtZXNwYWNlcz4NCiAgPHJlc291cmNlcyBtaW5SZXF1aXJlZFJldmlzaW9uPSIxLjAiIC8+DQogIDxjYXRlZ29yaWVzPg0KICAgIDxjYXRlZ29yeSBkaXNwbGF5TmFtZT0iJChzdHJpbmcuZ29vZ2xlKSIgbmFtZT0iQ2F0X0dvb2dsZSIvPg0KICA8L2NhdGVnb3JpZXM+DQo8L3BvbGljeURlZmluaXRpb25zPg0K"

  group_policy_uploaded_language_files = [
    {
      file_name     = "google.adml"
      language_code = "en-US"
      content       = <<-EOT
<?xml version="1.0" ?>
<policyDefinitionResources revision="1.0" schemaVersion="1.0">
  <displayName/>
  <description/>
  <resources>
    <stringTable>
      <string id="google">Google</string>
    </stringTable>
  </resources>
</policyDefinitionResources>
EOT
      // or: "PD94bWwgdmVyc2lvbj0iMS4wIiA/Pg0KPHBvbGljeURlZmluaXRpb25SZXNvdXJjZXMgcmV2aXNpb249IjEuMCIgc2NoZW1hVmVyc2lvbj0iMS4wIj4NCiAgPGRpc3BsYXlOYW1lLz4NCiAgPGRlc2NyaXB0aW9uLz4NCiAgPHJlc291cmNlcz4NCiAgICA8c3RyaW5nVGFibGU+DQogICAgICA8c3RyaW5nIGlkPSJnb29nbGUiPkdvb2dsZTwvc3RyaW5nPg0KICAgIDwvc3RyaW5nVGFibGU+DQogIDwvcmVzb3VyY2VzPg0KPC9wb2xpY3lEZWZpbml0aW9uUmVzb3VyY2VzPg0K"
    }
  ]

  timeouts = {
    create = "10m" // allow a generous amount of time for the upload to complete
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}