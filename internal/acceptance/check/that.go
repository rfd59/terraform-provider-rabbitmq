// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
)

type thatType struct {
	// resourceName being the full resource name e.g. azurerm_foo.bar
	resourceName string
}

// Key returns a type which can be used for more fluent assertions for a given Resource
func That(resourceName string) thatType {
	return thatType{
		resourceName: resourceName,
	}
}

func (t thatType) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[t.resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", t.resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource id not set")
		}

		return nil
	}
}

// // ExistsInAzure validates that the specified resource exists within Azure
func (t thatType) ExistsInRabbitMQ(any interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if strings.HasPrefix(t.resourceName, "rabbitmq_user.") {
			return any.(acceptance.UserResource).ExistsInRabbitMQ()
		}
		if strings.HasPrefix(t.resourceName, "rabbitmq_vhost.") {
			return any.(acceptance.VhostResource).ExistsInRabbitMQ()
		}
		return fmt.Errorf("'ExistsInRabbitMQ' method not found for this resource!!!")
	}
}

// Key returns a type which can be used for more fluent assertions for a given Resource & Key combination
func (t thatType) Key(key string) thatWithKeyType {
	return thatWithKeyType{
		resourceName: t.resourceName,
		key:          key,
	}
}

type thatWithKeyType struct {
	// resourceName being the full resource name e.g. azurerm_foo.bar
	resourceName string

	// key being the specific field we're querying e.g. bar or a nested object ala foo.0.bar
	key string

	// Skip the test
	skip bool
}

// JsonAssertionFunc is a function which takes a deserialized JSON object and asserts on it
type JsonAssertionFunc func(input []interface{}) (*bool, error)

// ContainsKeyValue returns a TestCheckFunc which asserts upon a given JSON string set into
// the State by deserializing it and then asserting on it via the JsonAssertionFunc
func (t thatWithKeyType) ContainsJsonValue(assertion JsonAssertionFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, exists := s.RootModule().Resources[t.resourceName]
		if !exists {
			return fmt.Errorf("%q was not found in the state", t.resourceName)
		}

		value, exists := rs.Primary.Attributes[t.key]
		if !exists {
			return fmt.Errorf("the value %q does not exist within %q", t.key, t.resourceName)
		}

		if value == "" {
			return fmt.Errorf("the value for %q was empty", t.key)
		}

		var out []interface{}
		if err := json.Unmarshal([]byte(value), &out); err != nil {
			return fmt.Errorf("deserializing the value for %q (%q) to json: %+v", t.key, value, err)
		}

		ok, err := assertion(out)
		if err != nil {
			return fmt.Errorf("asserting value for %q: %+v", t.key, err)
		}

		if ok == nil || !*ok {
			return fmt.Errorf("assertion failed for %q: %+v", t.key, err)
		}

		return nil
	}
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
	if t.skip {
		return skipTest()
	} else {
		return resource.TestCheckResourceAttr(t.resourceName, t.key, "")
	}
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
func (t thatWithKeyType) Count(length int) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(t.resourceName, t.key+".#", strconv.Itoa(length))
}

// IsUUID returns a TestCheckFunc which validates that the value for the specified key
// is a UUID.
func (t thatWithKeyType) IsUUID() resource.TestCheckFunc {
	var uuidRegex = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return resource.TestMatchResourceAttr(t.resourceName, t.key, uuidRegex)
}

// HasValue returns a TestCheckFunc which validates that the specific key has the
// specified value on the resource
func (t thatWithKeyType) HasValue(value string) resource.TestCheckFunc {
	if t.skip {
		return skipTest()
	} else {
		return resource.TestCheckResourceAttr(t.resourceName, t.key, value)
	}
}

// MatchesOtherKey returns a TestCheckFunc which validates that the key on this resource
// matches another other key on another resource
func (t thatWithKeyType) MatchesOtherKey(other string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrPair(t.resourceName, t.key, t.resourceName, other)
}

// MatchesRegex returns a TestCheckFunc which validates that the key on this resource matches
// the given regular expression
func (t thatWithKeyType) MatchesRegex(r *regexp.Regexp) resource.TestCheckFunc {
	return resource.TestMatchResourceAttr(t.resourceName, t.key, r)
}

// Load skip state
func (t thatWithKeyType) RmqFeature(feature bool) thatWithKeyType {
	t.skip = !feature
	return t
}

// Skip a Test
func skipTest() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}
