# Using an explicit project ID, the import ID is:
# {project_id}:{cluster_id}
terraform import microsoft365_graph_beta_device_and_app_management_settings_catalog.example f709ec73-55d4-46d8-897d-816ebba28778:settings-catalog
# Using the provider-default project ID, the import ID is:
# {cluster_id}
terraform import microsoft365_graph_beta_device_and_app_management_settings_catalog.example settings-catalog