# Provider Coverage

This provider offers extensive coverage across Microsoft 365 services including Intune, Microsoft 365, Teams, and Defender. Given the large size of this project, this page provides a comprehensive breakdown of available Terraform components organized by service domain.

*Last updated: 2026-03-18 11:27 UTC*

---

## Summary Statistics

| Terraform Block Type | Count |
|---------------------|-------|
| Resources | 144 |
| Data Sources | 51 |
| List Resources | 4 |
| Ephemerals | 2 |
| Actions | 43 |
| **Total Components** | **244** |

---

## Agents

**11 Resources**

<details>
<summary><b>Resources (11)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_agents_agent_collection` | v0.38.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_collection_assignment` | v0.38.0 | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_identity` | v0.38.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_identity_blueprint` | v0.38.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_identity_blueprint_certificate_credential` | v0.38.0 | — | Experimental | 3 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential` | v0.38.0 | — | Experimental | 3 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri` | v0.38.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential` | v0.38.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal` | v0.38.0 | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_instance` | v0.38.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_agents_agent_user` | v0.38.0 | — | Experimental | 3 | ✅ | ✅ |

</details>

---

## Applications

**10 Resources • 2 Data Sources**

<details>
<summary><b>Resources (10)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_applications_application` | v0.43.0 | — | Experimental | 6 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_application_certificate_credential` | v0.43.0 | — | Experimental | 3 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_application_federated_identity_credential` | v0.43.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_application_identifier_uri` | v0.43.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_application_owner` | v0.43.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_application_password_credential` | v0.43.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_on_premises_ip_application_segment` | v0.33.0 | v0.42.0 | Experimental | 5 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_service_principal` | v0.43.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_service_principal_app_role_assigned_to` | v0.38.0 | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_service_principal_owner` | v0.43.0 | — | Experimental | 2 | ✅ | ✅ |

</details>

<details>
<summary><b>Data Sources (2)</b></summary>

| Data Source Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_applications_application` | v0.43.0-alpha | — | Experimental | 10 | ✅ | ✅ |
| `microsoft365_graph_beta_applications_service_principal` | v0.31.0-alpha | v0.43.0-alpha | Experimental | 6 | ✅ | ✅ |

</details>

---

## Device Management

**55 Resources • 17 Data Sources • 2 List Resources • 86 Actions**

<details>
<summary><b>Resources (55)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_device_management_android_device_owner_compliance_policy` | v0.23.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_aosp_device_owner_compliance_policy` | v0.23.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls` | v0.28.0-alpha | — | Testing | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_app_control_for_business_policy` | v0.28.0-alpha | — | Testing | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy` | v0.32.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_assignment_filter` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_autopatch_groups` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ❌ |
| `microsoft365_graph_beta_device_management_device_category` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_device_compliance_notification_template` | v0.25.0-alpha | v0.27.0-alpha | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_device_enrollment_limit_configuration` | v0.16.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_device_enrollment_notification` | v0.27.0-alpha | v0.28-alpha | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_endpoint_privilege_management_json` | v0.14.1-alpha | v0.25.0-alpha | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_group_policy_configuration` | v0.29.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_group_policy_definition` | v0.40.0-alpha | — | Experimental | 5 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files` | v0.29.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_intune_branding_profile` | v0.23.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_ios_device_compliance_policy` | v0.23.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_linux_device_compliance_policy` | v0.25.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_linux_device_compliance_script` | v0.23.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_linux_platform_script` | v0.14.1-alpha | v0.25.0-alpha | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_macos_custom_attribute_script` | v0.21.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_macos_device_compliance_policy` | v0.23.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_macos_device_configuration_templates` | v0.27.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_macos_platform_script` | v0.14.1-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_macos_software_update_configuration` | v0.16.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_managed_device_cleanup_rule` | v0.14.1-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_operation_approval_policy` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_policy_set` | v0.33.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_rbac_resource_operation` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_reuseable_policy_setting` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_role_assignment` | v0.14.1-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_role_definition` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_role_scope_tag` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_settings_catalog_configuration_policy` | v0.14.1-alpha | v0.25.0-alpha | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json` | v0.14.1-alpha | v0.25.0-alpha | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_settings_catalog_template_json` | v0.14.1-alpha | v0.25.0-alpha | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_terms_and_conditions` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_autopilot_device_identity` | v0.21.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy` | v0.16.0-alpha | v0.47.0-alpha | Experimental | 3 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_backup_and_restore` | v0.31.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_device_compliance_policy` | v0.23.0-alpha | — | Experimental | 6 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_device_compliance_script` | v0.23.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_driver_update_inventory` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_driver_update_profile` | v0.14.1-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_enrollment_status_page` | v0.27.3-alpha | v0.28.0-alpha | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_feature_update_policy` | v0.14.1-alpha | — | Experimental | 4 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_platform_script` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy` | v0.14.1-alpha | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_quality_update_policy` | v0.14.1-alpha | v0.39.0-alpha | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_remediation_script` | v0.14.1-alpha | v0.39.0-alpha | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_update_deployment` | — | — | — | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_update_ring` | v0.23.0 | v0.39.0 | Experimental | 7 | ✅ | ✅ |
| `microsoft365_graph_device_management_device_configuration_assignment` | v0.17.0-alpha | — | Experimental | 1 | ✅ | ❌ |

</details>

<details>
<summary><b>Data Sources (17)</b></summary>

| Data Source Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_device_management_assignment_filter` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_device_category` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_group_policy_category` | v0.29.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_group_policy_value_reference` | v0.41.0-alpha | — | Experimental | 6 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_linux_platform_script` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_managed_device` | v0.18.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_reuseable_policy_setting` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_role_scope_tag` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_driver_update_inventory` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_driver_update_profile` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_feature_update_policy` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_platform_script` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_quality_update_policy` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_remediation_script` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_update_catalog_item` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_management_windows_update_ring` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |

</details>

<details>
<summary><b>List Resources (2)</b></summary>

| List Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|--------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_device_management_settings_catalog_configuration_policy` | v0.40.0-alpha | — | Experimental | 14 | ✅ | ✅ |
| `microsoft365_graph_beta_device_management_windows_platform_script` | v0.45.0-alpha | — | Experimental | 11 | ✅ | ✅ |

</details>

<details>
<summary><b>Actions (86)</b></summary>

Device management actions for managed devices including lifecycle operations, security actions, and maintenance tasks.

</details>

---

## Device and App Management

**18 Resources • 4 Data Sources**

<details>
<summary><b>Resources (18)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy` | v0.36.0-alpha | — | Experimental | 13 | ✅ | ✅ |
| `microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app` | v0.23.0 | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_and_app_management_application_category` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_ios_ipados_web_clip` | v0.21.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy` | v0.36.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app` | v1.0.0 | — | Stable | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_and_app_management_ios_store_app` | v0.21.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_macos_dmg_app` | v0.15.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_macos_lob_app` | v0.15.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_macos_pkg_app` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_macos_vpp_app` | v0.21.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_mobile_app_assignment` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence` | v0.14.1-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_office_suite_app` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_win32_app` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_win_get_app` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app` | v0.23.0 | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_and_app_management_windows_web_app` | v0.21.0-alpha | — | Experimental | 1 | ❌ | ❌ |

</details>

<details>
<summary><b>Data Sources (4)</b></summary>

| Data Source Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_device_and_app_management_application_category` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_and_app_management_mobile_app` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package` | v0.32.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_device_and_app_management_mobile_app_relationship` | v0.14.1-alpha | — | Experimental | 1 | ✅ | ✅ |

</details>

---

## Groups

**7 Resources • 1 Data Source**

<details>
<summary><b>Resources (7)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_groups_group` | v0.15.0-alpha | — | Experimental | 6 | ✅ | ✅ |
| `microsoft365_graph_beta_groups_group_app_role_assignment` | v0.39.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_groups_group_lifecycle_expiration_policy` | v0.18.0-alpha | v0.37.0-alpha | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_groups_group_lifecycle_expiration_policy_assignment` | v0.37.0-alpha | — | Preview | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_groups_group_member_assignment` | v0.15.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_groups_group_owner_assignment` | v0.15.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_groups_license_assignment` | v0.15.0-alpha | v0.37.0-alpha | Experimental | 1 | ✅ | ✅ |

</details>

<details>
<summary><b>Data Sources (1)</b></summary>

| Data Source Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_groups_group` | v0.42.0-alpha | — | Experimental | 5 | ✅ | ✅ |

</details>

---

## Identity and Access

**13 Resources • 5 Data Sources • 1 List Resource**

<details>
<summary><b>Resources (13)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_identity_and_access_attribute_set` | v0.33.0-alpha | — | Experimental | 1 | ✅ | ❌ |
| `microsoft365_graph_beta_identity_and_access_authentication_context` | v0.31.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_identity_and_access_authentication_strength_policy` | v0.36.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_identity_and_access_conditional_access_policy` | v0.19.0-alpha | v0.34.0-alpha | Experimental | 48 | ✅ | ✅ |
| `microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings` | v0.49.0 | — | Experimental | 8 | ✅ | ✅ |
| `microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings` | v0.49.0 | — | Experimental | 1 | ✅ | ❌ |
| `microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_user_sync_settings` | v0.49.0 | — | Experimental | 1 | ✅ | ❌ |
| `microsoft365_graph_beta_identity_and_access_cross_tenant_access_policy` | v0.49.0 | — | Experimental | 3 | ✅ | ✅ |
| `microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value` | v0.33.0-alpha | — | Experimental | 1 | ✅ | ❌ |
| `microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition` | v0.33.0-alpha | — | Experimental | 1 | ✅ | ❌ |
| `microsoft365_graph_beta_identity_and_access_named_location` | v0.28.0-alpha | v0.38.0-alpha | Experimental | 10 | ✅ | ✅ |
| `microsoft365_graph_beta_identity_and_access_network_filtering_policy` | v0.36.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_identity_and_access_conditional_access_terms_of_use` | v0.34.0-alpha | — | Experimental | 1 | ✅ | ✅ |

</details>

<details>
<summary><b>Data Sources (5)</b></summary>

| Data Source Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_identity_and_access_conditional_access_template` | v0.41.0-alpha | — | Experimental | 19 | ✅ | ✅ |
| `microsoft365_graph_beta_identity_and_access_directory_setting_templates` | v0.15.0-alpha | v0.33.0-alpha | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_identity_and_access_role_definitions` | v0.33.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_identity_and_access_subscribed_skus` | v0.15.0-alpha | v0.35.0-alpha | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_identity_and_access_tenant_information` | v0.39.0-alpha | — | Experimental | 2 | ✅ | ✅ |

</details>

<details>
<summary><b>List Resources (1)</b></summary>

| List Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|--------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_identity_and_access_conditional_access_policy` | v0.45.0-alpha | — | Experimental | 12 | ✅ | ✅ |

</details>

---

## M365 Admin

**3 Resources • 2 Data Sources**

<details>
<summary><b>Resources (3)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_m365_admin_browser_site` | v0.15.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_m365_admin_browser_site_list` | v0.15.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_m365_admin_m365_apps_installation_options` | v0.14.1-alpha | — | Experimental | 1 | ✅ | ✅ |

</details>

<details>
<summary><b>Data Sources (2)</b></summary>

| Data Source Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_m365_admin_browser_site` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_m365_admin_browser_site_list` | v0.42.0-alpha | — | Experimental | 1 | ❌ | ❌ |

</details>

---

## Microsoft Teams

**2 Resources**

<details>
<summary><b>Resources (2)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_powershell_microsoft_teams_teams_calling_policy` | v0.21.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_powershell_microsoft_teams_teams_meeting_policy` | v0.21.0-alpha | — | Experimental | 1 | ❌ | ❌ |

</details>

---

## Multitenant Management

**2 Ephemerals**

<details>
<summary><b>Ephemerals (2)</b></summary>

| Ephemeral Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|----------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_audit_events` | — | — | — | — | ❌ | ❌ |
| `microsoft365_graph_beta_windows_autopilot_device_c_s_v_import` | — | — | — | — | ❌ | ❌ |

</details>

---

## Users

**4 Resources • 1 List Resource**

<details>
<summary><b>Resources (4)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_users_user` | v0.16.0-alpha | v0.36.0-alpha | Experimental | 3 | ✅ | ✅ |
| `microsoft365_graph_beta_users_user_license_assignment` | v0.15.0-alpha | v0.37.0-alpha | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_users_user_mailbox_settings` | v0.40.0-alpha | — | Experimental | 3 | ✅ | ❌ |
| `microsoft365_graph_beta_users_user_manager` | v0.38.0-alpha | — | Experimental | 2 | ✅ | ✅ |

</details>

<details>
<summary><b>List Resources (1)</b></summary>

| List Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|--------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_users_user` | v0.45.0-alpha | — | Experimental | 13 | ✅ | ✅ |

</details>

---

## Utility

**10 Data Sources**

<details>
<summary><b>Data Sources (10)</b></summary>

| Data Source Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_utility_deployment_scheduler` | v0.41.0-alpha | — | Experimental | 8 | ✅ | ❌ |
| `microsoft365_utility_entra_id_sid_converter` | v0.35.0-alpha | — | Experimental | — | ✅ | ✅ |
| `microsoft365_utility_guid_list_sharder` | v0.42.0-alpha | v0.43.0-alpha | Experimental | 9 | ✅ | ✅ |
| `microsoft365_utility_itunes_app_metadata` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_utility_licensing_service_plan_reference` | v0.37.0-alpha | — | Experimental | 4 | ✅ | ✅ |
| `microsoft365_utility_macos_pkg_app_metadata` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ❌ |
| `microsoft365_utility_microsoft_365_endpoint_reference` | v0.36.0-alpha | — | Experimental | 4 | ✅ | ✅ |
| `microsoft365_utility_microsoft_store_package_manifest_metadata` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_utility_windows_msi_app_metadata` | v0.42.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_utility_windows_remediation_script_registry_key_generator` | v0.35.0-alpha | — | Experimental | 1 | ✅ | ✅ |

</details>

---

## Windows 365

**7 Resources • 6 Data Sources**

<details>
<summary><b>Resources (7)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_windows_365_azure_network_connection` | v0.19.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_365_cloud_pc_alert_rule` | v0.20.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_365_cloud_pc_device_image` | v0.19.1-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_365_cloud_pc_organization_settings` | v0.19.1-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy` | v0.19.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_365_cloud_pc_role_definition` | v0.25.0-alpha | — | Experimental | 1 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_365_cloud_pc_user_setting` | v0.19.0-alpha | — | Experimental | 1 | ✅ | ✅ |

</details>

<details>
<summary><b>Data Sources (6)</b></summary>

| Data Source Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_windows_365_cloud_pc_audit_event` | v0.18.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_windows_365_cloud_pc_device_image` | v0.18.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_windows_365_cloud_pc_frontline_service_plan` | v0.18.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_windows_365_cloud_pc_gallery_image` | v0.18.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_windows_365_cloud_pc_source_device_image` | v0.18.0-alpha | — | Experimental | 1 | ❌ | ❌ |
| `microsoft365_graph_beta_windows_365_cloud_pcs` | v0.18.0-alpha | — | Experimental | 1 | ❌ | ❌ |

</details>

---

## Windows Updates

**7 Resources • 2 Data Sources**

<details>
<summary><b>Resources (7)</b></summary>

| Resource Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|---------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_windows_updates_autopatch_deployment` | v0.50.0-alpha | — | Experimental | 6 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_updates_autopatch_deployment_state` | v0.50.0-alpha | — | Experimental | 3 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_updates_autopatch_operational_insights_connection` | v0.50.0-alpha | — | Experimental | 1 | ✅ | ❌ |
| `microsoft365_graph_beta_windows_updates_autopatch_policy_approval` | v0.50.0-alpha | — | Experimental | 2 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_updates_autopatch_ring` | v0.50.0-alpha | — | Experimental | 3 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group` | v0.50.0-alpha | — | Experimental | 1 | ✅ | ❌ |
| `microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group_assignment` | v0.50.0-alpha | — | Experimental | 1 | ✅ | ✅ |

</details>

<details>
<summary><b>Data Sources (2)</b></summary>

| Data Source Name | Version Introduced | Last Updated | Status | Examples | Unit Tests | Acceptance Tests |
|------------------|-------------------|--------------|--------|----------|------------|------------------|
| `microsoft365_graph_beta_windows_updates_catalog_enteries` | v0.50.0-alpha | — | Experimental | 5 | ✅ | ✅ |
| `microsoft365_graph_beta_windows_updates_product` | v0.50.0-alpha | — | Experimental | 2 | ❌ | ✅ |

</details>

---
