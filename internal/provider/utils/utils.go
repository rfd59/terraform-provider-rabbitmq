package utils

import (
	"strings"
)

// Because slashes are used to separate different components when constructing binding IDs,
// we need a way to ensure any components that include slashes can survive the round trip.
// Percent-encoding is a straightforward way of doing so.
// (reference: https://developer.mozilla.org/en-US/docs/Glossary/percent-encoding)

func PercentEncodeSlashes(s string) string {
	// Encode any percent signs, then encode any forward slashes.
	return strings.ReplaceAll(strings.ReplaceAll(s, "%", "%25"), "/", "%2F")
}

func PercentDecodeSlashes(s string) string {
	// Decode any forward slashes, then decode any percent signs.
	return strings.ReplaceAll(strings.ReplaceAll(s, "%2F", "/"), "%25", "%")
}
