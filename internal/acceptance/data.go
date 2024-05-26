// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

const (
	// charSetAlphaNum is the alphanumeric character set for use with randStringFromCharSet
	charSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"
)

func init() {
	// unit testing
	if os.Getenv("TF_ACC") == "" {
		return
	}
}

type TestData struct {
	// ResourceName is the fully qualified resource name, comprising the
	// resource type and then the resource label
	// e.g. `azurerm_resource_group.test`
	ResourceName string

	// ResourceType is the Terraform Resource Type - `azurerm_resource_group`
	ResourceType string

	// resourceLabel is the local used for the resource - generally "test""
	ResourceLabel string
}

// BuildTestData generates some test data for the given resource
func BuildTestData(resourceType string, resourceLabel string) TestData {
	testData := TestData{
		ResourceType:  resourceType,
		ResourceLabel: resourceLabel,
		ResourceName:  fmt.Sprintf("%s.%s", resourceType, resourceLabel),
	}

	return testData
}

func (td *TestData) RandomInteger() int {
	return td.RandomIntOfLength(999999)
}

func (td *TestData) RandomIntegerString() string {
	return strconv.Itoa(td.RandomInteger())
}

func (td *TestData) RandomString() string {
	return randString(5)
}

func (td *TestData) BuildArrayString(array []string) string {
	var str = ""

	for i := 0; i < len(array); i++ {
		str += fmt.Sprintf("\"%s\"", array[i])
		if i < len(array)-1 {
			str += ","
		}
	}

	return fmt.Sprintf("[%s]", str)
}

func (td *TestData) RandomIntOfLength(len int) int {
	return rand.Intn(len)
}

// RandomStringOfLength is a random 1 to 1024 character string which is unique to this test case
func (td *TestData) RandomStringOfLength(len int) string {
	// len should not be less then 1 or greater than 1024
	if 1 > len || len > 1024 {
		panic("Invalid Test: RandomStringOfLength: length argument must be between 1 and 1024 characters")
	}

	return randString(len)
}

// randString generates a random alphanumeric string of the length specified
func randString(strlen int) string {
	return randStringFromCharSet(strlen, charSetAlphaNum)
}

// randStringFromCharSet generates a random string by selecting characters from
// the charset provided
func randStringFromCharSet(strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}
