# Import an intent by its ID. Settings are not tracked until they are added to configuration.
terraform import microsoft365_graph_beta_device_management_endpoint_security_intent.bitlocker_minimal 00000000-0000-0000-0000-000000000000

# Import an intent and the settings to manage, as a comma-separated list of setting definition IDs after a slash.
terraform import microsoft365_graph_beta_device_management_endpoint_security_intent.bitlocker_maximal 00000000-0000-0000-0000-000000000000/deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerEncryptDevice,deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerRecoveryPasswordRotation
