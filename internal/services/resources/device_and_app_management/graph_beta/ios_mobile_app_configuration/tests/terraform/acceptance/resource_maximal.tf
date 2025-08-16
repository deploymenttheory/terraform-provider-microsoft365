resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "maximal" {
  display_name        = "Test Maximal iOS Mobile App Configuration - Unique"
  description         = "Maximal iOS mobile app configuration for testing with all features"
  targeted_mobile_apps = ["12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321"]
  role_scope_tag_ids  = ["0", "1"]
  
  encoded_setting_xml = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+dGVzdEtleTwva2V5PgoJPHN0cmluZz50ZXN0VmFsdWU8L3N0cmluZz4KPC9kaWN0Pgo8L3BsaXN0Pgo="
  
  settings = [
    {
      app_config_key       = "testKey1"
      app_config_key_type  = "stringType"
      app_config_key_value = "testValue1"
    },
    {
      app_config_key       = "testKey2"
      app_config_key_type  = "integerType"
      app_config_key_value = "123"
    }
  ]
}