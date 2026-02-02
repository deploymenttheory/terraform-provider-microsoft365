# Import an api:// identifier URI
terraform import microsoft365_graph_beta_applications_application_identifier_uri.api_uri "00000000-0000-0000-0000-000000000000/api://00000000-0000-0000-0000-000000000000"

# Import an https:// identifier URI (note: URI is NOT URL-encoded in import ID)
terraform import microsoft365_graph_beta_applications_application_identifier_uri.https_uri "00000000-0000-0000-0000-000000000000/https://mycompany.com/my-app"
