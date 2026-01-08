# Maximal Group Policy Uploaded Definition File configuration
resource "microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files" "maximal" {
  file_name             = "unit-test-msedge.admx"
  content               = "UEQ5NGJXd2dkbVZ5YzJsdmJqMGlNUzR3SWlCbGJtTnZaR2x1WnowaWRYUm1MVGdpUHo0S1BIQnZiR2xqZVVSbFptbHVhWFJwYjI1eklISmxkbWx6YVc5dVBTSTNMaklpSUhOamFHVnRZVlpsY25OcGIyNDlJakV1TUNJKwogIDx3b2xpY3lPYW1lc3BhY2VzPgogICAgPHRhcmdldCBwcmVmaXg9Im1zZWRnZSIgbmFtZXNwYWNlPSJNaWNyb3NvZnQuUG9saWNpZXMuRWRnZSIvPgogICAgPHVzaW5nIHByZWZpeD0iTWljcm9zb2Z0IiBuYW1lc3BhY2U9Ik1pY3Jvc29mdC5Qb2xpY2llcyIvPgogIDwvcG9saWN5TmFtZXNwYWNlcz4KICA8cmVzb3VyY2VzIG1pblJlPgogICAgICAgICAgICAgIDxzdHJpbmc+dG9vbGJhcjwvc3RyaW5nPgogICAgICAgIDxkZWNpbWFsIHZhbHVlPSIxIi8+CiAgICAgIDwvZW5hYmxlZFZhbHVlPgogICAgICA8ZGlzYWJsZWRWYWx1ZT4KICAgICAgICA8ZGVjaW1hbCB2YWx1ZT0iMCIvPgogICAgICA8L2Rpc2FibGVkVmFsdWU+CiAgICA8L3BvbGljeT4KICA8L3BvbGljaWVzPgo8L3BvbGljeURlZmluaXRpb25zPgo="
  //default_language_code = "en-US"

  # Multiple language files
  group_policy_uploaded_language_files = [
    {
      file_name     = "msedge.adml"
      language_code = "en-US"
      content       = "UEQ5NGJXd2dkbVZ5YzJsdmJqMGlNUzR3SWlCbGJtTnZaR2x1WnowaWRYUm1MVGdpUHo0S1BIQnZiR2xqZVVSbFptbHVhWFJwYjI1U1pYTnZkWEpqWlhNZ2NtVjJhWE5wYjI0OUlqY3VNaUlnYzJOb1pXMWhWbVZ5YzJsdmJqMGlNUzR3SWo0S0lDQmtjMkJ2YVhOd2JHRjVUbUZ0WlM4K0NpQWdaR1Z6WTNKcGNIUnBiMjR2UGdvZ0lDQWdJQ0E4TDNCeVpYTmxiblJoZEdsdmJqNEtJQ0FnIDwvY0hKbGMyVnVkR0YwYVc5dVZHRmliR1UrQ2lBZ1BDOXlaWE52ZFhKalpYTStDanhpYjJSNUNqd3ZZbTlrZVQ0S1BDOXdiMnhwWTNsRVpXWnBibWwwYVc5dVVtVnpiM1Z5WTJWUD0="
    },
    {
      file_name     = "msedge.adml"
      language_code = "fr-FR"
      content       = "UEQ5NGJXd2dkbVZ5YzJsdmJqMGlNUzR3SWlCbGJtTnZaR2x1WnowaWRYUm1MVGdpUHo0S1BIQnZiR2xqZVVSbFptbHVhWFJwYjI1U1pYTnZkWEpqWlhNZ2NtVjJhWE5wYjI0OUlqY3VNaUlnYzJOb1pXMWhWbVZ5YzJsdmJqMGlNUzR3SWo0S0lDQmtjMkJ2YVhOd2JHRjVUbUZ0WlM4K0NpQWdaR1Z6WTNKcGNIUnBiMjR2UGdvZ0lDQWdJQ0E4TDNCeVpYTmxiblJoZEdsdmJqNEtJQ0FnIDwvY0hKbGMyVnVkR0YwYVc5dVZHRmliR1UrQ2lBZ1BDOXlaWE52ZFhKalpYTStDanhpYjJSNUNqd3ZZbTlrZVQ0S1BDOXdiMnhwWTNsRVpXWnBibWwwYVc5dVVtVnpiM1Z5WTJWUD0="
    }
  ]

  # Optional timeouts
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
