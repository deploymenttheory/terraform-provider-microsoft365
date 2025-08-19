#!/bin/bash

# Example 1: Import Custom Configuration Template
print_status "Example 1: Importing a Custom Configuration Template"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_config_example \"12345678-1234-1234-1234-123456789012\""
echo ""

# Example 2: Import Preference File Configuration
print_status "Example 2: Importing a Preference File Configuration"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example \"87654321-4321-4321-4321-210987654321\""
echo ""

# Example 3: Import Trusted Certificate Configuration
print_status "Example 3: Importing a Trusted Certificate Configuration"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example \"11111111-2222-3333-4444-555555555555\""
echo ""

# Example 4: Import SCEP Certificate Profile
print_status "Example 4: Importing a SCEP Certificate Profile"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example \"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee\""
echo ""

# Example 5: Import PKCS Certificate Profile
print_status "Example 5: Importing a PKCS Certificate Profile"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example \"ffffffff-eeee-dddd-cccc-bbbbbbbbbbbb\""
echo ""
