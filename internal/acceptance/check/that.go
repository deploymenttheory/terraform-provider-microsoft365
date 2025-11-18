package check

import (
	"fmt"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type thatType struct {
	// resourceName being the full resource name e.g. microsoft365_graph_beta_device_management_role_scope_tag.test
	resourceName string
}

// That returns a type which can be used for more fluent assertions for a given Resource
func That(resourceName string) thatType {
	return thatType{
		resourceName: resourceName,
	}
}

// DoesNotExistInGraph validates that the specified resource does not exist within Microsoft Graph
func (t thatType) DoesNotExistInGraph(testResource types.TestResource) resource.TestCheckFunc {
	return helpers.DoesNotExistInGraph(testResource, t.resourceName)
}

// ExistsInGraph validates that the specified resource exists within Microsoft Graph
func (t thatType) ExistsInGraph(testResource types.TestResource) resource.TestCheckFunc {
	return helpers.ExistsInGraph(testResource, t.resourceName)
}

// Key returns a type which can be used for more fluent assertions for a given Resource & Key combination
func (t thatType) Key(key string) thatWithKeyType {
	return thatWithKeyType{
		resourceName: t.resourceName,
		key:          key,
	}
}

type thatWithKeyType struct {
	// resourceName being the full resource name e.g. microsoft365_graph_beta_device_management_role_scope_tag.test
	resourceName string
	// key being the specific field we're querying e.g. display_name or a nested object ala assignments.0.group_id
	key string
}

// DoesNotExist returns a TestCheckFunc which validates that the specific key
// does not exist on the resource
func (t thatWithKeyType) DoesNotExist() resource.TestCheckFunc {
	return resource.TestCheckNoResourceAttr(t.resourceName, t.key)
}

// Exists returns a TestCheckFunc which validates that the specific key exists on the resource
func (t thatWithKeyType) Exists() resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(t.resourceName, t.key)
}

// IsEmpty returns a TestCheckFunc which validates that the specific key is empty on the resource
func (t thatWithKeyType) IsEmpty() resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(t.resourceName, t.key, "")
}

// IsNotEmpty returns a TestCheckFunc which validates that the specific key is not empty on the resource
func (t thatWithKeyType) IsNotEmpty() resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(t.resourceName, t.key, func(value string) error {
		if value == "" {
			return fmt.Errorf("value is empty")
		}
		return nil
	})
}

// IsSet returns a TestCheckFunc which validates that the specific key is set on the resource
func (t thatWithKeyType) IsSet() resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(t.resourceName, t.key)
}

// IsUUID returns a TestCheckFunc which validates that the value for the specified key is a UUID
func (t thatWithKeyType) IsUUID() resource.TestCheckFunc {
	uuidRegex := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{3}-[a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return resource.TestMatchResourceAttr(t.resourceName, t.key, uuidRegex)
}

// HasValue returns a TestCheckFunc which validates that the specific key has the
// specified value on the resource
func (t thatWithKeyType) HasValue(value string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(t.resourceName, t.key, value)
}

// MatchesOtherKey returns a TestCheckFunc which validates that the key on this resource
// matches another key on another resource
func (t thatWithKeyType) MatchesOtherKey(other thatWithKeyType) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrPair(t.resourceName, t.key, other.resourceName, other.key)
}

// MatchesRegex returns a TestCheckFunc which validates that the key on this resource matches
// the given regular expression
func (t thatWithKeyType) MatchesRegex(r *regexp.Regexp) resource.TestCheckFunc {
	return resource.TestMatchResourceAttr(t.resourceName, t.key, r)
}

// ValidatesWith returns a TestCheckFunc which runs a custom validation function on the attribute value
func (t thatWithKeyType) ValidatesWith(fn func(value string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, exists := s.RootModule().Resources[t.resourceName]
		if !exists {
			return fmt.Errorf("%q was not found in the state", t.resourceName)
		}

		value, exists := rs.Primary.Attributes[t.key]
		if !exists {
			return fmt.Errorf("the value %q does not exist within %q", t.key, t.resourceName)
		}

		return fn(value)
	}
}
