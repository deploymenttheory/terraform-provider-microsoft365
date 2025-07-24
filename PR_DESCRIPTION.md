# Pull Request Description

## Summary

Add support for iOS mobile app configuration management in Microsoft Intune through new Terraform resource and data source. This implementation enables Infrastructure as Code management of app-specific settings for managed iOS applications.

### Issue Reference

Fixes #[Issue Number]

### Motivation and Context

- **Why is this change needed?** Organizations managing iOS devices through Microsoft Intune need to configure app-specific settings programmatically through Terraform
- **What problem does it solve?** Enables automated deployment and management of iOS app configurations without manual portal intervention
- **Key capabilities:**
  - Create, update, and delete iOS mobile app configurations
  - Configure app settings using key-value pairs or XML
  - Assign configurations to user/device groups
  - Query existing configurations by ID or display name

### Dependencies

- Microsoft Graph API v1.0
- Existing Kiota SDK client
- No new external dependencies required

## Type of Change

Please mark the relevant option with an `x`:

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [x] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [x] 📝 Documentation update (Wiki/README/Code comments)
- [ ] ♻️ Refactor (code improvement without functional changes)
- [ ] 🎨 Style update (formatting, renaming)
- [ ] 🔧 Configuration change
- [ ] 📦 Dependency update

## Testing

- [x] I have added unit tests that prove my fix is effective or that my feature works
- [x] New and existing unit tests pass locally with my changes
- [x] I have tested this code in the following browsers/environments: Local development environment with mocked Graph API responses

### Test Coverage
- **Unit Tests**: Complete coverage for CRUD operations, state mapping, and error handling
- **Acceptance Tests**: Full integration tests for create, read, update, delete, and import operations
- **Mock Tests**: Comprehensive HTTP response mocks for reliable testing

## Quality Checklist

- [x] I have reviewed my own code before requesting review
- [x] I have verified there are no other open Pull Requests for the same update/change
- [x] All CI/CD pipelines pass without errors or warnings
- [x] My code follows the established style guidelines of this project
- [x] I have added necessary documentation (if appropriate)
- [x] I have commented my code, particularly in complex areas
- [x] I have made corresponding changes to the README and other relevant documentation
- [x] My changes generate no new warnings
- [x] I have performed a self-review of my own code
- [x] My code is properly formatted according to project standards

## Implementation Details

### Resource: `microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration`
- Full CRUD operations
- Support for both key-value settings and XML configuration
- Sensitive data handling for XML content
- Assignment management for targeting user/device groups
- Import functionality

### Data Source: `microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration`
- Lookup by ID or display name
- Read-only access to all configuration attributes
- Consistent schema with resource

### Key Features
- **Flexible Configuration**: Support both structured key-value pairs and raw XML
- **Secure Handling**: Sensitive XML data properly marked and handled
- **Assignment Support**: Target configurations to specific groups
- **Import Capability**: Import existing configurations into Terraform state
- **Comprehensive Testing**: Unit and acceptance tests with full mocking

## Additional Notes

- Follows existing provider patterns for consistency
- Uses Graph API v1.0 (stable) endpoints
- Implements proper error handling with descriptive messages
- Includes retry logic for eventual consistency
- Documentation auto-generated for Terraform Registry

### Files Changed
- **Resource Implementation**: `internal/services/resources/device_and_app_management/graph_v1.0/ios_mobile_app_configuration/`
- **Data Source Implementation**: `internal/services/datasources/device_and_app_management/graph_v1.0/ios_mobile_app_configuration/`
- **Provider Registration**: Updated to include new resource and data source
- **Documentation**: Auto-generated docs for both resource and data source