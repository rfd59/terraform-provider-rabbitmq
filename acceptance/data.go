// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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

func (TestData) RandomInteger() int {
	return randTimeInt()
}

func (TestData) RandomIntegerString() string {
	return strconv.Itoa(randTimeInt())
}

func (TestData) RandomString() string {
	return randString(5)
}

func (TestData) BuildArrayString(array []string) string {
	var str = ""

	for i := 0; i < len(array); i++ {
		str += fmt.Sprintf("\"%s\"", array[i])
		if i < len(array)-1 {
			str += ","
		}
	}

	return fmt.Sprintf("[%s]", str)
}

// RandomIntOfLength is a random 8 to 18 digit integer which is unique to this test case
func (td *TestData) RandomIntOfLength(len int) int {
	// len should not be
	//  - greater then 18, longest a int can represent
	//  - less then 8, as that gives us YYMMDDRR
	if 8 > len || len > 18 {
		panic("Invalid Test: RandomIntOfLength: len is not between 8 or 18 inclusive")
	}

	// 18 - just return the int
	if len >= 18 {
		return td.RandomInteger()
	}

	// 16-17 just strip off the last 1-2 digits
	if len >= 16 {
		return td.RandomInteger() / int(math.Pow10(18-len))
	}

	// 8-15 keep len - 2 digits and add 2 characters of randomness on
	s := strconv.Itoa(td.RandomInteger())
	r := s[16:18]
	v := s[0 : len-2]
	i, _ := strconv.Atoi(v + r)

	return i
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

func randTimeInt() int {
	// acctest.RantInt() returns a value of size:
	// 000000000000000000
	// YYMMddHHmmsshhRRRR

	// go format: 2006-01-02 15:04:05.00

	timeStr := strings.Replace(time.Now().Local().Format("060102150405.00"), ".", "", 1) // no way to not have a .?
	postfix := randStringFromCharSet(4, "0123456789")

	i, err := strconv.Atoi(timeStr + postfix)
	if err != nil {
		panic(err)
	}

	return i
}
